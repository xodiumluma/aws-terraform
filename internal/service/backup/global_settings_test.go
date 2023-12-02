// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backup_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccBackupGlobalSettings_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var settings backup.DescribeGlobalSettingsOutput

	resourceName := "aws_backup_global_settings.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckOrganizationManagementAccount(ctx, t)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, backup.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGlobalSettingsConfig_basic("true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(ctx, &settings),
					resource.TestCheckResourceAttr(resourceName, "global_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "global_settings.isCrossAccountBackupEnabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGlobalSettingsConfig_basic("false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(ctx, &settings),
					resource.TestCheckResourceAttr(resourceName, "global_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "global_settings.isCrossAccountBackupEnabled", "false"),
				),
			},
			{
				Config: testAccGlobalSettingsConfig_basic("true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGlobalSettingsExists(ctx, &settings),
					resource.TestCheckResourceAttr(resourceName, "global_settings.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "global_settings.isCrossAccountBackupEnabled", "true"),
				),
			},
		},
	})
}

func testAccCheckGlobalSettingsExists(ctx context.Context, settings *backup.DescribeGlobalSettingsOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).BackupConn(ctx)
		resp, err := conn.DescribeGlobalSettingsWithContext(ctx, &backup.DescribeGlobalSettingsInput{})
		if err != nil {
			return err
		}

		*settings = *resp

		return nil
	}
}

func testAccGlobalSettingsConfig_basic(setting string) string {
	return fmt.Sprintf(`
resource "aws_backup_global_settings" "test" {
  global_settings = {
    "isCrossAccountBackupEnabled" = %[1]q
  }
}
`, setting)
}
