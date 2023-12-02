// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wafv2_test

import (
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/wafv2"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccWAFV2IPSetDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	name := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_wafv2_ip_set.test"
	datasourceName := "data.aws_wafv2_ip_set.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckScopeRegional(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, wafv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccIPSetDataSourceConfig_nonExistent(name),
				ExpectError: regexache.MustCompile(`WAFv2 IPSet not found`),
			},
			{
				Config: testAccIPSetDataSourceConfig_name(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "addresses", resourceName, "addresses"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
					acctest.MatchResourceAttrRegionalARN(datasourceName, "arn", "wafv2", regexache.MustCompile(fmt.Sprintf("regional/ipset/%v/.+$", name))),
					resource.TestCheckResourceAttrPair(datasourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(datasourceName, "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "ip_address_version", resourceName, "ip_address_version"),
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "scope", resourceName, "scope"),
				),
			},
		},
	})
}

func testAccIPSetDataSourceConfig_name(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_ip_set" "test" {
  name               = "%s"
  scope              = "REGIONAL"
  ip_address_version = "IPV4"
}

data "aws_wafv2_ip_set" "test" {
  name  = aws_wafv2_ip_set.test.name
  scope = "REGIONAL"
}
`, name)
}

func testAccIPSetDataSourceConfig_nonExistent(name string) string {
	return fmt.Sprintf(`
resource "aws_wafv2_ip_set" "test" {
  name               = "%s"
  scope              = "REGIONAL"
  ip_address_version = "IPV4"
}

data "aws_wafv2_ip_set" "test" {
  name  = "tf-acc-test-does-not-exist"
  scope = "REGIONAL"
}
`, name)
}
