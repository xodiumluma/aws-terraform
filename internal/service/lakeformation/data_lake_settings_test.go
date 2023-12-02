// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package lakeformation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lakeformation"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tflakeformation "github.com/hashicorp/terraform-provider-aws/internal/service/lakeformation"
)

func testAccDataLakeSettings_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_lakeformation_data_lake_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, lakeformation.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, lakeformation.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDataLakeSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeSettingsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataLakeSettingsExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "catalog_id", "data.aws_caller_identity.current", "account_id"),
					resource.TestCheckResourceAttr(resourceName, "admins.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "admins.0", "data.aws_iam_session_context.current", "issuer_arn"),
					resource.TestCheckResourceAttr(resourceName, "create_database_default_permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "create_database_default_permissions.0.principal", "IAM_ALLOWED_PRINCIPALS"),
					resource.TestCheckResourceAttr(resourceName, "create_database_default_permissions.0.permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "create_database_default_permissions.0.permissions.0", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "create_table_default_permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "create_table_default_permissions.0.principal", "IAM_ALLOWED_PRINCIPALS"),
					resource.TestCheckResourceAttr(resourceName, "create_table_default_permissions.0.permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "create_table_default_permissions.0.permissions.0", "ALL"),
					resource.TestCheckResourceAttr(resourceName, "allow_external_data_filtering", "true"),
					resource.TestCheckResourceAttr(resourceName, "external_data_filtering_allow_list.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "external_data_filtering_allow_list.0", "data.aws_caller_identity.current", "account_id"),
					resource.TestCheckResourceAttr(resourceName, "authorized_session_tag_value_list.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "authorized_session_tag_value_list.0", "engine1"),
				),
			},
		},
	})
}

func testAccDataLakeSettings_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_lakeformation_data_lake_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, lakeformation.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, lakeformation.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDataLakeSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeSettingsConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataLakeSettingsExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tflakeformation.ResourceDataLakeSettings(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccDataLakeSettings_withoutCatalogID(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_lakeformation_data_lake_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, lakeformation.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDataLakeSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeSettingsConfig_withoutCatalogID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataLakeSettingsExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "admins.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "admins.0", "data.aws_iam_session_context.current", "issuer_arn"),
				),
			},
		},
	})
}

func testAccDataLakeSettings_readOnlyAdmins(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_lakeformation_data_lake_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, lakeformation.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDataLakeSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDataLakeSettingsConfig_readOnlyAdmins,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataLakeSettingsExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "read_only_admins.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "read_only_admins.0", "data.aws_iam_session_context.current", "issuer_arn"),
				),
			},
		},
	})
}

func testAccCheckDataLakeSettingsDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).LakeFormationConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_lakeformation_data_lake_settings" {
				continue
			}

			input := &lakeformation.GetDataLakeSettingsInput{}

			if rs.Primary.Attributes["catalog_id"] != "" {
				input.CatalogId = aws.String(rs.Primary.Attributes["catalog_id"])
			}

			output, err := conn.GetDataLakeSettingsWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, lakeformation.ErrCodeEntityNotFoundException) {
				continue
			}

			if err != nil {
				return fmt.Errorf("error getting Lake Formation data lake settings (%s): %w", rs.Primary.ID, err)
			}

			if output != nil && output.DataLakeSettings != nil && len(output.DataLakeSettings.DataLakeAdmins) > 0 {
				return fmt.Errorf("Lake Formation data lake admin(s) (%s) still exist", rs.Primary.ID)
			}

			if output != nil && output.DataLakeSettings != nil && len(output.DataLakeSettings.ReadOnlyAdmins) > 0 {
				return fmt.Errorf("Lake Formation data lake read only admin(s) (%s) still exist", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckDataLakeSettingsExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).LakeFormationConn(ctx)

		input := &lakeformation.GetDataLakeSettingsInput{}

		if rs.Primary.Attributes["catalog_id"] != "" {
			input.CatalogId = aws.String(rs.Primary.Attributes["catalog_id"])
		}

		_, err := conn.GetDataLakeSettingsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("error getting Lake Formation data lake settings (%s): %w", rs.Primary.ID, err)
		}

		return nil
	}
}

const testAccDataLakeSettingsConfig_basic = `
data "aws_caller_identity" "current" {}

data "aws_iam_session_context" "current" {
  arn = data.aws_caller_identity.current.arn
}

resource "aws_lakeformation_data_lake_settings" "test" {
  catalog_id = data.aws_caller_identity.current.account_id

  create_database_default_permissions {
    principal   = "IAM_ALLOWED_PRINCIPALS"
    permissions = ["ALL"]
  }

  create_table_default_permissions {
    principal   = "IAM_ALLOWED_PRINCIPALS"
    permissions = ["ALL"]
  }

  admins                             = [data.aws_iam_session_context.current.issuer_arn]
  trusted_resource_owners            = [data.aws_caller_identity.current.account_id]
  allow_external_data_filtering      = true
  external_data_filtering_allow_list = [data.aws_caller_identity.current.account_id]
  authorized_session_tag_value_list  = ["engine1"]
}
`

const testAccDataLakeSettingsConfig_withoutCatalogID = `
data "aws_caller_identity" "current" {}

data "aws_iam_session_context" "current" {
  arn = data.aws_caller_identity.current.arn
}

resource "aws_lakeformation_data_lake_settings" "test" {
  admins = [data.aws_iam_session_context.current.issuer_arn]
}
`

const testAccDataLakeSettingsConfig_readOnlyAdmins = `
data "aws_caller_identity" "current" {}

data "aws_iam_session_context" "current" {
  arn = data.aws_caller_identity.current.arn
}

resource "aws_lakeformation_data_lake_settings" "test" {
  catalog_id = data.aws_caller_identity.current.account_id

  read_only_admins = [data.aws_iam_session_context.current.issuer_arn]
}
`
