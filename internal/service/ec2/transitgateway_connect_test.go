// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func testAccTransitGatewayConnect_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.TransitGatewayConnect
	resourceName := "aws_ec2_transit_gateway_connect.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	transitGatewayVpcAttachmentResourceName := "aws_ec2_transit_gateway_vpc_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckTransitGatewayVPCAttachment(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayConnectDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConnectConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "protocol", "gre"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "true"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "transit_gateway_id", transitGatewayResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "transport_attachment_id", transitGatewayVpcAttachmentResourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayConnect_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.TransitGatewayConnect
	resourceName := "aws_ec2_transit_gateway_connect.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckTransitGatewayVPCAttachment(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayConnectDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConnectConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceTransitGatewayConnect(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccTransitGatewayConnect_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.TransitGatewayConnect
	resourceName := "aws_ec2_transit_gateway_connect.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckTransitGatewayVPCAttachment(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayConnectDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConnectConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayConnectConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccTransitGatewayConnectConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccTransitGatewayConnect_TransitGatewayDefaultRouteTableAssociationAndPropagationDisabled(t *testing.T) {
	ctx := acctest.Context(t)
	var transitGateway1 ec2.TransitGateway
	var transitGatewayConnect1 ec2.TransitGatewayConnect
	resourceName := "aws_ec2_transit_gateway_connect.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckTransitGatewayVPCAttachment(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayConnectDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTableAssociationAndPropagationDisabled(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect1),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentNotAssociated(ctx, &transitGateway1, &transitGatewayConnect1),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentNotPropagated(ctx, &transitGateway1, &transitGatewayConnect1),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "false"),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTransitGatewayConnect_TransitGatewayDefaultRouteTableAssociation(t *testing.T) {
	ctx := acctest.Context(t)
	var transitGateway1, transitGateway2, transitGateway3 ec2.TransitGateway
	var transitGatewayConnect1, transitGatewayConnect2, transitGatewayConnect3 ec2.TransitGatewayConnect
	resourceName := "aws_ec2_transit_gateway_connect.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckTransitGatewayVPCAttachment(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayConnectDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTableAssociation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect1),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentNotAssociated(ctx, &transitGateway1, &transitGatewayConnect1),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTableAssociation(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway2),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect2),
					testAccCheckTransitGatewayConnectNotRecreated(&transitGatewayConnect1, &transitGatewayConnect2),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentAssociated(ctx, &transitGateway2, &transitGatewayConnect2),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "true"),
				),
			},
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTableAssociation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway3),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect3),
					testAccCheckTransitGatewayConnectNotRecreated(&transitGatewayConnect2, &transitGatewayConnect3),
					testAccCheckTransitGatewayAssociationDefaultRouteTableAttachmentNotAssociated(ctx, &transitGateway3, &transitGatewayConnect3),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_association", "false"),
				),
			},
		},
	})
}

func testAccTransitGatewayConnect_TransitGatewayDefaultRouteTablePropagation(t *testing.T) {
	ctx := acctest.Context(t)
	var transitGateway1, transitGateway2, transitGateway3 ec2.TransitGateway
	var transitGatewayConnect1, transitGatewayConnect2, transitGatewayConnect3 ec2.TransitGatewayConnect
	resourceName := "aws_ec2_transit_gateway_connect.test"
	transitGatewayResourceName := "aws_ec2_transit_gateway.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckTransitGatewayVPCAttachment(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTransitGatewayConnectDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTablePropagation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway1),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect1),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentNotPropagated(ctx, &transitGateway1, &transitGatewayConnect1),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTablePropagation(rName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway2),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect2),
					testAccCheckTransitGatewayConnectNotRecreated(&transitGatewayConnect1, &transitGatewayConnect2),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentPropagated(ctx, &transitGateway2, &transitGatewayConnect2),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "true"),
				),
			},
			{
				Config: testAccTransitGatewayConnectConfig_defaultRouteTablePropagation(rName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTransitGatewayExists(ctx, transitGatewayResourceName, &transitGateway3),
					testAccCheckTransitGatewayConnectExists(ctx, resourceName, &transitGatewayConnect3),
					testAccCheckTransitGatewayConnectNotRecreated(&transitGatewayConnect2, &transitGatewayConnect3),
					testAccCheckTransitGatewayPropagationDefaultRouteTableAttachmentNotPropagated(ctx, &transitGateway3, &transitGatewayConnect3),
					resource.TestCheckResourceAttr(resourceName, "transit_gateway_default_route_table_propagation", "false"),
				),
			},
		},
	})
}

func testAccCheckTransitGatewayConnectExists(ctx context.Context, n string, v *ec2.TransitGatewayConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Transit Gateway Connect ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		output, err := tfec2.FindTransitGatewayConnectByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckTransitGatewayConnectDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_ec2_transit_gateway_connect" {
				continue
			}

			_, err := tfec2.FindTransitGatewayConnectByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("EC2 Transit Gateway Connect %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckTransitGatewayConnectNotRecreated(i, j *ec2.TransitGatewayConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.StringValue(i.TransitGatewayAttachmentId) != aws.StringValue(j.TransitGatewayAttachmentId) {
			return errors.New("EC2 Transit Gateway Connect was recreated")
		}

		return nil
	}
}

func testAccTransitGatewayConnectConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_connect" "test" {
  transit_gateway_id      = aws_ec2_transit_gateway.test.id
  transport_attachment_id = aws_ec2_transit_gateway_vpc_attachment.test.id
}
`, rName))
}

func testAccTransitGatewayConnectConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_connect" "test" {
  transit_gateway_id      = aws_ec2_transit_gateway.test.id
  transport_attachment_id = aws_ec2_transit_gateway_vpc_attachment.test.id

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccTransitGatewayConnectConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids         = [aws_subnet.test.id]
  transit_gateway_id = aws_ec2_transit_gateway.test.id
  vpc_id             = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_connect" "test" {
  transit_gateway_id      = aws_ec2_transit_gateway.test.id
  transport_attachment_id = aws_ec2_transit_gateway_vpc_attachment.test.id

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccTransitGatewayConnectConfig_defaultRouteTableAssociationAndPropagationDisabled(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  default_route_table_association = "disable"
  default_route_table_propagation = "disable"

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids                                      = [aws_subnet.test.id]
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  vpc_id                                          = aws_vpc.test.id
  transit_gateway_default_route_table_association = false
  transit_gateway_default_route_table_propagation = false

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_connect" "test" {
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  transport_attachment_id                         = aws_ec2_transit_gateway_vpc_attachment.test.id
  transit_gateway_default_route_table_association = false
  transit_gateway_default_route_table_propagation = false

  tags = {
    Name = %[1]q
  }
}
`, rName))
}

func testAccTransitGatewayConnectConfig_defaultRouteTableAssociation(rName string, transitGatewayDefaultRouteTableAssociation bool) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids                                      = [aws_subnet.test.id]
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  vpc_id                                          = aws_vpc.test.id
  transit_gateway_default_route_table_association = %[2]t

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_connect" "test" {
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  transport_attachment_id                         = aws_ec2_transit_gateway_vpc_attachment.test.id
  transit_gateway_default_route_table_association = %[2]t

  tags = {
    Name = %[1]q
  }
}
`, rName, transitGatewayDefaultRouteTableAssociation))
}

func testAccTransitGatewayConnectConfig_defaultRouteTablePropagation(rName string, transitGatewayDefaultRouteTablePropagation bool) string {
	return acctest.ConfigCompose(acctest.ConfigAvailableAZsNoOptInDefaultExclude(), fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_vpc_attachment" "test" {
  subnet_ids                                      = [aws_subnet.test.id]
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  vpc_id                                          = aws_vpc.test.id
  transit_gateway_default_route_table_propagation = %[2]t

  tags = {
    Name = %[1]q
  }
}

resource "aws_ec2_transit_gateway_connect" "test" {
  transit_gateway_id                              = aws_ec2_transit_gateway.test.id
  transport_attachment_id                         = aws_ec2_transit_gateway_vpc_attachment.test.id
  transit_gateway_default_route_table_propagation = %[2]t

  tags = {
    Name = %[1]q
  }
}
`, rName, transitGatewayDefaultRouteTablePropagation))
}
