// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package elbv2_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccELBV2TargetGroupAttachment_basic(t *testing.T) {
	ctx := acctest.Context(t)
	targetGroupName := fmt.Sprintf("test-target-group-%s", sdkacctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupAttachmentConfig_idInstance(targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupAttachmentExists(ctx, "aws_lb_target_group_attachment.test"),
				),
			},
		},
	})
}

func TestAccELBV2TargetGroupAttachment_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	targetGroupName := fmt.Sprintf("test-target-group-%s", sdkacctest.RandString(10))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupAttachmentConfig_idInstance(targetGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTargetGroupAttachmentExists(ctx, "aws_lb_target_group_attachment.test"),
					testAccCheckTargetGroupAttachmentDisappears(ctx, "aws_lb_target_group_attachment.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccELBV2TargetGroupAttachment_backwardsCompatibility(t *testing.T) {
	ctx := acctest.Context(t)
	targetGroupName := fmt.Sprintf("test-target-group-%s", sdkacctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupAttachmentConfig_backwardsCompatibility(targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupAttachmentExists(ctx, "aws_alb_target_group_attachment.test"),
				),
			},
		},
	})
}

func TestAccELBV2TargetGroupAttachment_port(t *testing.T) {
	ctx := acctest.Context(t)
	targetGroupName := fmt.Sprintf("test-target-group-%s", sdkacctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupAttachmentConfig_port(targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupAttachmentExists(ctx, "aws_lb_target_group_attachment.test"),
				),
			},
		},
	})
}

func TestAccELBV2TargetGroupAttachment_ipAddress(t *testing.T) {
	ctx := acctest.Context(t)
	targetGroupName := fmt.Sprintf("test-target-group-%s", sdkacctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupAttachmentConfig_idIPAddress(targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupAttachmentExists(ctx, "aws_lb_target_group_attachment.test"),
				),
			},
		},
	})
}

func TestAccELBV2TargetGroupAttachment_lambda(t *testing.T) {
	ctx := acctest.Context(t)
	targetGroupName := fmt.Sprintf("test-target-group-%s", sdkacctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elbv2.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckTargetGroupAttachmentDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccTargetGroupAttachmentConfig_idLambda(targetGroupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTargetGroupAttachmentExists(ctx, "aws_lb_target_group_attachment.test"),
				),
			},
		},
	})
}

func testAccCheckTargetGroupAttachmentDisappears(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Attachment not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBV2Conn(ctx)
		targetGroupArn := rs.Primary.Attributes["target_group_arn"]

		target := &elbv2.TargetDescription{
			Id: aws.String(rs.Primary.Attributes["target_id"]),
		}

		_, hasPort := rs.Primary.Attributes["port"]
		if hasPort {
			port, _ := strconv.Atoi(rs.Primary.Attributes["port"])
			target.Port = aws.Int64(int64(port))
		}

		params := &elbv2.DeregisterTargetsInput{
			TargetGroupArn: aws.String(targetGroupArn),
			Targets:        []*elbv2.TargetDescription{target},
		}

		_, err := conn.DeregisterTargetsWithContext(ctx, params)
		if err != nil && !tfawserr.ErrCodeEquals(err, elbv2.ErrCodeTargetGroupNotFoundException) {
			return fmt.Errorf("Error deregistering Targets: %s", err)
		}

		return err
	}
}

func testAccCheckTargetGroupAttachmentExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Target Group Attachment ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBV2Conn(ctx)

		_, hasPort := rs.Primary.Attributes["port"]
		targetGroupArn := rs.Primary.Attributes["target_group_arn"]

		target := &elbv2.TargetDescription{
			Id: aws.String(rs.Primary.Attributes["target_id"]),
		}
		if hasPort {
			port, _ := strconv.Atoi(rs.Primary.Attributes["port"])
			target.Port = aws.Int64(int64(port))
		}

		describe, err := conn.DescribeTargetHealthWithContext(ctx, &elbv2.DescribeTargetHealthInput{
			TargetGroupArn: aws.String(targetGroupArn),
			Targets:        []*elbv2.TargetDescription{target},
		})

		if err != nil {
			return err
		}

		if len(describe.TargetHealthDescriptions) != 1 {
			return errors.New("Target Group Attachment not found")
		}

		return nil
	}
}

func testAccCheckTargetGroupAttachmentDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ELBV2Conn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_lb_target_group_attachment" && rs.Type != "aws_alb_target_group_attachment" {
				continue
			}

			_, hasPort := rs.Primary.Attributes["port"]
			targetGroupArn := rs.Primary.Attributes["target_group_arn"]

			target := &elbv2.TargetDescription{
				Id: aws.String(rs.Primary.Attributes["target_id"]),
			}
			if hasPort {
				port, _ := strconv.Atoi(rs.Primary.Attributes["port"])
				target.Port = aws.Int64(int64(port))
			}

			describe, err := conn.DescribeTargetHealthWithContext(ctx, &elbv2.DescribeTargetHealthInput{
				TargetGroupArn: aws.String(targetGroupArn),
				Targets:        []*elbv2.TargetDescription{target},
			})
			if err == nil {
				if len(describe.TargetHealthDescriptions) != 0 {
					return fmt.Errorf("Target Group Attachment %q still exists", rs.Primary.ID)
				}
			}

			// Verify the error
			if tfawserr.ErrCodeEquals(err, elbv2.ErrCodeTargetGroupNotFoundException) || tfawserr.ErrCodeEquals(err, elbv2.ErrCodeInvalidTargetException) {
				return nil
			} else {
				return fmt.Errorf("Unexpected error checking LB destroyed: %s", err)
			}
		}

		return nil
	}
}

func testAccTargetGroupAttachmentInstanceBaseConfig() string {
	return `
data "aws_availability_zones" "available" {
  # t2.micro instance type is not available in these Availability Zones
  exclude_zone_ids = ["usw2-az4"]
  state            = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

data "aws_ami" "amzn-ami-minimal-hvm-ebs" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn-ami-minimal-hvm-*"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

resource "aws_instance" "test" {
  ami           = data.aws_ami.amzn-ami-minimal-hvm-ebs.id
  instance_type = "t2.micro"
  subnet_id     = aws_subnet.test.id
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.1.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "tf-acc-test-lb-target-group-attachment"
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "tf-acc-test-lb-target-group-attachment"
  }
}
`
}

func testAccTargetGroupAttachmentConfig_idInstance(rName string) string {
	return testAccTargetGroupAttachmentInstanceBaseConfig() + fmt.Sprintf(`
resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 443
  protocol = "HTTPS"
  vpc_id   = aws_vpc.test.id
}

resource "aws_lb_target_group_attachment" "test" {
  target_group_arn = aws_lb_target_group.test.arn
  target_id        = aws_instance.test.id
}
`, rName)
}

func testAccTargetGroupAttachmentConfig_port(rName string) string {
	return testAccTargetGroupAttachmentInstanceBaseConfig() + fmt.Sprintf(`
resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 443
  protocol = "HTTPS"
  vpc_id   = aws_vpc.test.id
}

resource "aws_lb_target_group_attachment" "test" {
  target_group_arn = aws_lb_target_group.test.arn
  target_id        = aws_instance.test.id
  port             = 80
}
`, rName)
}

func testAccTargetGroupAttachmentConfig_backwardsCompatibility(rName string) string {
	return testAccTargetGroupAttachmentInstanceBaseConfig() + fmt.Sprintf(`
resource "aws_lb_target_group" "test" {
  name     = %[1]q
  port     = 443
  protocol = "HTTPS"
  vpc_id   = aws_vpc.test.id
}

resource "aws_alb_target_group_attachment" "test" {
  target_group_arn = aws_lb_target_group.test.arn
  target_id        = aws_instance.test.id
  port             = 80
}
`, rName)
}

func testAccTargetGroupAttachmentConfig_idIPAddress(rName string) string {
	return testAccTargetGroupAttachmentInstanceBaseConfig() + fmt.Sprintf(`
resource "aws_lb_target_group" "test" {
  name        = %[1]q
  port        = 443
  protocol    = "HTTPS"
  target_type = "ip"
  vpc_id      = aws_vpc.test.id
}

resource "aws_lb_target_group_attachment" "test" {
  availability_zone = aws_instance.test.availability_zone
  target_group_arn  = aws_lb_target_group.test.arn
  target_id         = aws_instance.test.private_ip
}
`, rName)
}

func testAccTargetGroupAttachmentConfig_idLambda(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_lambda_permission" "test" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test.arn
  principal     = "elasticloadbalancing.${data.aws_partition.current.dns_suffix}"
  qualifier     = aws_lambda_alias.test.name
  source_arn    = aws_lb_target_group.test.arn
  statement_id  = "AllowExecutionFromlb"
}

resource "aws_lb_target_group" "test" {
  name        = %[1]q
  target_type = "lambda"
}

resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambda_elb.zip"
  function_name = %[1]q
  role          = aws_iam_role.test.arn
  handler       = "lambda_elb.lambda_handler"
  runtime       = "python3.7"
}

resource "aws_lambda_alias" "test" {
  name             = "test"
  description      = "a sample description"
  function_name    = aws_lambda_function.test.function_name
  function_version = "$LATEST"
}

resource "aws_iam_role" "test" {
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

resource "aws_lb_target_group_attachment" "test" {
  depends_on = [aws_lambda_permission.test]

  target_group_arn = aws_lb_target_group.test.arn
  target_id        = aws_lambda_alias.test.arn
}
`, rName)
}
