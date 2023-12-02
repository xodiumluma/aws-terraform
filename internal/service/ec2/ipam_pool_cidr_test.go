// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccIPAMPoolCIDR_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var cidr ec2.IpamPoolCidr
	resourceName := "aws_vpc_ipam_pool_cidr.test"
	cidrBlock := "10.0.0.0/24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIPAMPoolCIDRDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccIPAMPoolCIDRConfig_provisionedIPv4(cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMPoolCIDRExists(ctx, resourceName, &cidr),
					resource.TestCheckResourceAttr(resourceName, "cidr", cidrBlock),
					resource.TestCheckResourceAttrPair(resourceName, "ipam_pool_id", "aws_vpc_ipam_pool.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"netmask_length",
				},
			},
		},
	})
}

func TestAccIPAMPoolCIDR_basicNetmaskLength(t *testing.T) {
	ctx := acctest.Context(t)
	var cidr ec2.IpamPoolCidr
	resourceName := "aws_vpc_ipam_pool_cidr.test"
	netmaskLength := "24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIPAMPoolCIDRDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccIPAMPoolCIDRConfig_provisionedIPv4NetmaskLength(netmaskLength),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMPoolCIDRExists(ctx, resourceName, &cidr),
					resource.TestCheckResourceAttr(resourceName, "netmask_length", netmaskLength),
					testAccCheckIPAMPoolCIDRPrefix(&cidr, netmaskLength),
					resource.TestCheckResourceAttrPair(resourceName, "ipam_pool_id", "aws_vpc_ipam_pool.testchild", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"netmask_length",
				},
			},
		},
	})
}

func TestAccIPAMPoolCIDR_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var cidr ec2.IpamPoolCidr
	resourceName := "aws_vpc_ipam_pool_cidr.test"
	cidrBlock := "10.0.0.0/24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIPAMPoolCIDRDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccIPAMPoolCIDRConfig_provisionedIPv4(cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMPoolCIDRExists(ctx, resourceName, &cidr),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceIPAMPoolCIDR(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccIPAMPoolCIDR_Disappears_ipam(t *testing.T) {
	ctx := acctest.Context(t)
	var cidr ec2.IpamPoolCidr
	resourceName := "aws_vpc_ipam_pool_cidr.test"
	ipamResourceName := "aws_vpc_ipam.test"
	cidrBlock := "10.0.0.0/24"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ec2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckIPAMPoolCIDRDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccIPAMPoolCIDRConfig_provisionedIPv4(cidrBlock),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAMPoolCIDRExists(ctx, resourceName, &cidr),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfec2.ResourceIPAM(), ipamResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckIPAMPoolCIDRExists(ctx context.Context, n string, v *ec2.IpamPoolCidr) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IPAM Pool CIDR ID is set")
		}

		cidrBlock, poolID, err := tfec2.IPAMPoolCIDRParseResourceID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		output, err := tfec2.FindIPAMPoolCIDRByTwoPartKey(ctx, conn, cidrBlock, poolID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckIPAMPoolCIDRDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Conn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_vpc_ipam_pool_cidr" {
				continue
			}

			cidrBlock, poolID, err := tfec2.IPAMPoolCIDRParseResourceID(rs.Primary.ID)

			if err != nil {
				return err
			}

			_, err = tfec2.FindIPAMPoolCIDRByTwoPartKey(ctx, conn, cidrBlock, poolID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("IPAM Pool CIDR still exists: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckIPAMPoolCIDRPrefix(cidr *ec2.IpamPoolCidr, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if strings.Split(aws.StringValue(cidr.Cidr), "/")[1] != expected {
			return fmt.Errorf("Bad cidr prefix: got %s, expected %s", aws.StringValue(cidr.Cidr), expected)
		}

		return nil
	}
}

const TestAccIPAMPoolCIDRConfig_base = `
data "aws_region" "current" {}

resource "aws_vpc_ipam" "test" {
  description = "test"

  operating_regions {
    region_name = data.aws_region.current.name
  }

  cascade = true
}
`

const TestAccIPAMPoolCIDRConfig_privatePool = `
resource "aws_vpc_ipam_pool" "test" {
  address_family = "ipv4"
  ipam_scope_id  = aws_vpc_ipam.test.private_default_scope_id
  locale         = data.aws_region.current.name
}
`

const TestAccIPAMPoolCIDRConfig_privatePoolWithCIDR = `
resource "aws_vpc_ipam_pool" "test" {
  address_family = "ipv4"
  ipam_scope_id  = aws_vpc_ipam.test.private_default_scope_id
  locale         = data.aws_region.current.name
}

resource "aws_vpc_ipam_pool_cidr" "testparent" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = "10.0.0.0/16"
}

resource "aws_vpc_ipam_pool" "testchild" {
  address_family      = "ipv4"
  ipam_scope_id       = aws_vpc_ipam.test.private_default_scope_id
  locale              = data.aws_region.current.name
  source_ipam_pool_id = aws_vpc_ipam_pool.test.id
}
`

func testAccIPAMPoolCIDRConfig_provisionedIPv4(cidr string) string {
	return acctest.ConfigCompose(TestAccIPAMPoolCIDRConfig_base, TestAccIPAMPoolCIDRConfig_privatePool, fmt.Sprintf(`
resource "aws_vpc_ipam_pool_cidr" "test" {
  ipam_pool_id = aws_vpc_ipam_pool.test.id
  cidr         = %[1]q
}
`, cidr))
}

func testAccIPAMPoolCIDRConfig_provisionedIPv4NetmaskLength(netmaskLength string) string {
	return acctest.ConfigCompose(TestAccIPAMPoolCIDRConfig_base, TestAccIPAMPoolCIDRConfig_privatePoolWithCIDR, fmt.Sprintf(`
resource "aws_vpc_ipam_pool_cidr" "test" {
  ipam_pool_id   = aws_vpc_ipam_pool.testchild.id
  netmask_length = %[1]q
  depends_on     = [aws_vpc_ipam_pool_cidr.testparent]
}
`, netmaskLength))
}
