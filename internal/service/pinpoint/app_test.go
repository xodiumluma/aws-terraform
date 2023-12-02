// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pinpoint_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfpinpoint "github.com/hashicorp/terraform-provider-aws/internal/service/pinpoint"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccPinpointApp_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "campaign_hook.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "limits.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					resource.TestCheckResourceAttr(resourceName, "quiet_time.#", "1"),
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

func TestAccPinpointApp_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfpinpoint.ResourceApp(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccPinpointApp_nameGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_nameGenerated(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", id.UniqueIdPrefix),
				),
			},
		},
	})
}

func TestAccPinpointApp_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_namePrefix("tf-acc-test-prefix-"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
				),
			},
		},
	})
}

func TestAccPinpointApp_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
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
				Config: testAccAppConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAppConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccPinpointApp_campaignHookLambda(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_campaignHookLambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "campaign_hook.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "campaign_hook.0.mode", "DELIVERY"),
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

func TestAccPinpointApp_limits(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_limits(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "limits.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "limits.0.total", "100"),
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

func TestAccPinpointApp_quietTime(t *testing.T) {
	ctx := acctest.Context(t)
	var application pinpoint.ApplicationResponse
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_pinpoint_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, pinpoint.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAppDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAppConfig_quietTime(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAppExists(ctx, resourceName, &application),
					resource.TestCheckResourceAttr(resourceName, "quiet_time.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "quiet_time.0.start", "00:00"),
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

func testAccPreCheckApp(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointConn(ctx)

	input := &pinpoint.GetAppsInput{}

	_, err := conn.GetAppsWithContext(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckAppDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_pinpoint_app" {
				continue
			}

			_, err := tfpinpoint.FindAppByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Pinpoint App %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAppExists(ctx context.Context, n string, v *pinpoint.ApplicationResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointConn(ctx)

		output, err := tfpinpoint.FindAppByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccAppConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name = %[1]q
}
`, rName)
}

func testAccAppConfig_nameGenerated() string {
	return `
resource "aws_pinpoint_app" "test" {}
`
}

func testAccAppConfig_namePrefix(namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}

func testAccAppConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccAppConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccAppConfig_campaignHookLambda(rName string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name = %[1]q

  campaign_hook {
    lambda_function_name = aws_lambda_function.test.arn
    mode                 = "DELIVERY"
  }

  depends_on = [aws_lambda_permission.test]
}

resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdapinpoint.zip"
  function_name = %[1]q
  role          = aws_iam_role.test.arn
  handler       = "lambdapinpoint.handler"
  runtime       = "nodejs16.x"
  publish       = true
}

data "aws_partition" "current" {}

resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.${data.aws_partition.current.dns_suffix}"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

data "aws_caller_identity" "aws" {}

data "aws_region" "current" {}

resource "aws_lambda_permission" "test" {
  statement_id  = "AllowExecutionFromPinpoint"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test.function_name
  principal     = "pinpoint.${data.aws_region.current.name}.${data.aws_partition.current.dns_suffix}"
  source_arn    = "arn:${data.aws_partition.current.partition}:mobiletargeting:${data.aws_region.current.name}:${data.aws_caller_identity.aws.account_id}:/apps/*"
}
`, rName)
}

func testAccAppConfig_limits(rName string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name = %[1]q

  limits {
    daily               = 3
    maximum_duration    = 600
    messages_per_second = 50
    total               = 100
  }
}
`, rName)
}

func testAccAppConfig_quietTime(rName string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test" {
  name = %[1]q

  quiet_time {
    start = "00:00"
    end   = "03:00"
  }
}
`, rName)
}
