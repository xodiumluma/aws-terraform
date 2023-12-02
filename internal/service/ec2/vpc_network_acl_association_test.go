// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccVPCNetworkACLAssociation_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.NetworkAclAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_network_acl_association.test"
	naclResourceName := "aws_network_acl.test"
	subnetResourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkACLAssociationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkACLAssociationExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "network_acl_id", naclResourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", subnetResourceName, "id"),
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

func TestAccVPCNetworkACLAssociation_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.NetworkAclAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_network_acl_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkACLAssociationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkACLAssociationExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceNetworkACLAssociation(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNetworkACLAssociation_disappears_NACL(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.NetworkAclAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_network_acl_association.test"
	naclResourceName := "aws_network_acl.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkACLAssociationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkACLAssociationExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceNetworkACL(), naclResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNetworkACLAssociation_disappears_Subnet(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.NetworkAclAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_network_acl_association.test"
	subnetResourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkACLAssociationConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkACLAssociationExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceSubnet(), subnetResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccVPCNetworkACLAssociation_twoAssociations(t *testing.T) {
	ctx := acctest.Context(t)
	var v1, v2 ec2.NetworkAclAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resource1Name := "aws_network_acl_association.test1"
	resource2Name := "aws_network_acl_association.test2"
	naclResourceName := "aws_network_acl.test"
	subnet1ResourceName := "aws_subnet.test1"
	subnet2ResourceName := "aws_subnet.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkACLAssociationConfig_twoAssociations(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkACLAssociationExists(ctx, resource1Name, &v1),
					testAccCheckNetworkACLAssociationExists(ctx, resource1Name, &v2),
					resource.TestCheckResourceAttrPair(resource1Name, "network_acl_id", naclResourceName, "id"),
					resource.TestCheckResourceAttrPair(resource1Name, "subnet_id", subnet1ResourceName, "id"),
					resource.TestCheckResourceAttrPair(resource2Name, "network_acl_id", naclResourceName, "id"),
					resource.TestCheckResourceAttrPair(resource2Name, "subnet_id", subnet2ResourceName, "id"),
				),
			},
			{
				ResourceName:      resource1Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      resource2Name,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVPCNetworkACLAssociation_associateWithDefaultNACL(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.NetworkAclAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_network_acl_association.test"
	subnetResourceName := "aws_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckNetworkACLAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCNetworkACLAssociationConfig_default(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckNetworkACLAssociationExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id", subnetResourceName, "id"),
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

func testAccCheckNetworkACLAssociationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_network_acl_association" {
				continue
			}

			_, err := tfec2.FindNetworkACLAssociationByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("EC2 Network ACL Association %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNetworkACLAssociationExists(ctx context.Context, n string, v *ec2.NetworkAclAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 Network ACL Association ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		output, err := tfec2.FindNetworkACLAssociationByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccVPCNetworkACLAssociationConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_acl" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id = aws_vpc.test.id

  cidr_block = "10.1.33.0/24"

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_acl_association" "test" {
  network_acl_id = aws_network_acl.test.id
  subnet_id      = aws_subnet.test.id
}
`, rName)
}

func testAccVPCNetworkACLAssociationConfig_twoAssociations(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_acl" "test" {
  vpc_id = aws_vpc.test.id

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test1" {
  vpc_id = aws_vpc.test.id

  cidr_block = "10.1.33.0/24"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test2" {
  vpc_id = aws_vpc.test.id

  cidr_block = "10.1.34.0/24"

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_acl_association" "test1" {
  network_acl_id = aws_network_acl.test.id
  subnet_id      = aws_subnet.test1.id
}

resource "aws_network_acl_association" "test2" {
  network_acl_id = aws_network_acl.test.id
  subnet_id      = aws_subnet.test2.id
}
`, rName)
}

func testAccVPCNetworkACLAssociationConfig_default(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.1.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_subnet" "test" {
  vpc_id = aws_vpc.test.id

  cidr_block = "10.1.33.0/24"

  tags = {
    Name = %[1]q
  }
}

resource "aws_network_acl_association" "test" {
  network_acl_id = aws_vpc.test.default_network_acl_id
  subnet_id      = aws_subnet.test.id
}
`, rName)
}
