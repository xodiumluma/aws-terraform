// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationinsights"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapplicationinsights "github.com/hashicorp/terraform-provider-aws/internal/service/applicationinsights"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccApplicationInsightsApplication_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var app applicationinsights.ApplicationInfo
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_applicationinsights_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, applicationinsights.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "resource_group_name", rName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "applicationinsights", fmt.Sprintf("application/resource-group/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "auto_config_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "cwe_monitor_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ops_center_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccApplicationConfig_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "resource_group_name", rName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "applicationinsights", fmt.Sprintf("application/resource-group/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "auto_config_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "cwe_monitor_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ops_center_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccApplicationInsightsApplication_autoConfig(t *testing.T) {
	ctx := acctest.Context(t)
	var app applicationinsights.ApplicationInfo
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_applicationinsights_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, applicationinsights.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "resource_group_name", rName),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "applicationinsights", fmt.Sprintf("application/resource-group/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "auto_config_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "cwe_monitor_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "ops_center_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccApplicationInsightsApplication_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var app applicationinsights.ApplicationInfo
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_applicationinsights_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, applicationinsights.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
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
				Config: testAccApplicationConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccApplicationConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccApplicationInsightsApplication_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var app applicationinsights.ApplicationInfo
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_applicationinsights_application.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, applicationinsights.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckApplicationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplicationExists(ctx, resourceName, &app),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfapplicationinsights.ResourceApplication(), resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfapplicationinsights.ResourceApplication(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckApplicationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ApplicationInsightsConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_applicationinsights_application" {
				continue
			}

			app, err := tfapplicationinsights.FindApplicationByName(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			if aws.StringValue(app.ResourceGroupName) == rs.Primary.ID {
				return fmt.Errorf("applicationinsights Application %q still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckApplicationExists(ctx context.Context, n string, app *applicationinsights.ApplicationInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No applicationinsights Application ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ApplicationInsightsConn(ctx)
		resp, err := tfapplicationinsights.FindApplicationByName(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		*app = *resp

		return nil
	}
}

func testAccApplicationConfigBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_resourcegroups_group" "test" {
  name = %[1]q

  resource_query {
    query = <<JSON
	{
		"ResourceTypeFilters": [
		  "AWS::EC2::Instance"
		],
		"TagFilters": [
		  {
			"Key": "Stage",
			"Values": [
			  "Test"
			]
		  }
		]
	  }
JSON
  }
}
`, rName)
}

func testAccApplicationConfig_basic(rName string) string {
	return testAccApplicationConfigBase(rName) + `
resource "aws_applicationinsights_application" "test" {
  resource_group_name = aws_resourcegroups_group.test.name
}
`
}

func testAccApplicationConfig_updated(rName string) string {
	return testAccApplicationConfigBase(rName) + `
resource "aws_applicationinsights_application" "test" {
  resource_group_name = aws_resourcegroups_group.test.name
  auto_config_enabled = true
}
`
}

func testAccApplicationConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return testAccApplicationConfigBase(rName) + fmt.Sprintf(`
resource "aws_applicationinsights_application" "test" {
  resource_group_name = aws_resourcegroups_group.test.name

  tags = {
    %[1]q = %[2]q
  }
}
`, tagKey1, tagValue1)
}

func testAccApplicationConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return testAccApplicationConfigBase(rName) + fmt.Sprintf(`
resource "aws_applicationinsights_application" "test" {
  resource_group_name = aws_resourcegroups_group.test.name

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tagKey1, tagValue1, tagKey2, tagValue2)
}
