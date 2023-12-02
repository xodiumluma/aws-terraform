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

func TestAccSiteVPNGatewayAttachment_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.VpcAttachment
	resourceName := "aws_vpn_gateway_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPNGatewayAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVPNGatewayAttachmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayAttachmentExists(ctx, resourceName, &v),
				),
			},
		},
	})
}

func TestAccSiteVPNGatewayAttachment_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v ec2.VpcAttachment
	resourceName := "aws_vpn_gateway_attachment.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckVPNGatewayAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVPNGatewayAttachmentConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPNGatewayAttachmentExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceVPNGatewayAttachment(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckVPNGatewayAttachmentExists(ctx context.Context, n string, v *ec2.VpcAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No EC2 VPN Gateway Attachment ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		output, err := tfec2.FindVPNGatewayVPCAttachment(ctx, conn, rs.Primary.Attributes["vpn_gateway_id"], rs.Primary.Attributes["vpc_id"])

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckVPNGatewayAttachmentDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_vpn_gateway_attachment" {
				continue
			}

			_, err := tfec2.FindVPNGatewayVPCAttachment(ctx, conn, rs.Primary.Attributes["vpn_gateway_id"], rs.Primary.Attributes["vpc_id"])

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("EC2 VPN Gateway Attachment %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccSiteVPNGatewayAttachmentConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = %[1]q
  }
}

resource "aws_vpn_gateway" "test" {
  tags = {
    Name = %[1]q
  }
}

resource "aws_vpn_gateway_attachment" "test" {
  vpc_id         = aws_vpc.test.id
  vpn_gateway_id = aws_vpn_gateway.test.id
}
`, rName)
}
