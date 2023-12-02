// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kms_test

import (
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/kms"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccKMSAliasDataSource_service(t *testing.T) {
	ctx := acctest.Context(t)
	rName := "alias/aws/s3"
	resourceName := "data.aws_kms_alias.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAliasDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "kms", rName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "target_key_arn", "kms", regexache.MustCompile(`key/[0-9a-z]{8}-([0-9a-z]{4}-){3}[0-9a-z]{12}`)),
					resource.TestMatchResourceAttr(resourceName, "target_key_id", regexache.MustCompile("^[0-9a-z]{8}-([0-9a-z]{4}-){3}[0-9a-z]{12}$")),
				),
			},
		},
	})
}

func TestAccKMSAliasDataSource_cmk(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	aliasResourceName := "aws_kms_alias.test"
	datasourceAliasResourceName := "data.aws_kms_alias.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAliasDataSourceConfig_cmk(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceAliasResourceName, "arn", aliasResourceName, "arn"),
					resource.TestCheckResourceAttrPair(datasourceAliasResourceName, "target_key_arn", aliasResourceName, "target_key_arn"),
					resource.TestCheckResourceAttrPair(datasourceAliasResourceName, "target_key_id", aliasResourceName, "target_key_id"),
				),
			},
		},
	})
}

func testAccAliasDataSourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "aws_kms_alias" "test" {
  name = %[1]q
}
`, rName)
}

func testAccAliasDataSourceConfig_cmk(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  name          = "alias/%[1]s"
  target_key_id = aws_kms_key.test.key_id
}

data "aws_kms_alias" "test" {
  name = aws_kms_alias.test.name
}
`, rName)
}
