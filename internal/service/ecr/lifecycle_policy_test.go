// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ecr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccECRLifecyclePolicy_basic(t *testing.T) {
	ctx := acctest.Context(t)
	randString := sdkacctest.RandString(10)
	rName := fmt.Sprintf("tf-acc-test-lifecycle-%s", randString)
	resourceName := "aws_ecr_lifecycle_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ecr.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLifecyclePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecyclePolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecyclePolicyExists(ctx, resourceName),
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

func TestAccECRLifecyclePolicy_ignoreEquivalent(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_ecr_lifecycle_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ecr.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLifecyclePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecyclePolicyConfig_order(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecyclePolicyExists(ctx, resourceName),
				),
			},
			{
				Config:   testAccLifecyclePolicyConfig_newOrder(rName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccECRLifecyclePolicy_detectDiff(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_ecr_lifecycle_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, ecr.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLifecyclePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLifecyclePolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecyclePolicyExists(ctx, resourceName),
				),
			},
			{
				Config:             testAccLifecyclePolicyConfig_changed(rName),
				ExpectNonEmptyPlan: true,
				PlanOnly:           true,
			},
		},
	})
}

func testAccCheckLifecyclePolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ECRConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_ecr_lifecycle_policy" {
				continue
			}

			input := &ecr.GetLifecyclePolicyInput{
				RepositoryName: aws.String(rs.Primary.ID),
			}

			_, err := conn.GetLifecyclePolicyWithContext(ctx, input)
			if err != nil {
				if tfawserr.ErrCodeEquals(err, ecr.ErrCodeRepositoryNotFoundException) {
					return nil
				}
				if tfawserr.ErrCodeEquals(err, ecr.ErrCodeLifecyclePolicyNotFoundException) {
					return nil
				}
				return err
			}
		}

		return nil
	}
}

func testAccCheckLifecyclePolicyExists(ctx context.Context, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ECRConn(ctx)

		input := &ecr.GetLifecyclePolicyInput{
			RepositoryName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetLifecyclePolicyWithContext(ctx, input)
		return err
	}
}

func testAccLifecyclePolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecr_repository" "test" {
  name = "%s"
}

resource "aws_ecr_lifecycle_policy" "test" {
  repository = aws_ecr_repository.test.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire images older than 14 days",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 14
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}
`, rName)
}

func testAccLifecyclePolicyConfig_changed(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecr_repository" "test" {
  name = "%s"
}

resource "aws_ecr_lifecycle_policy" "test" {
  repository = aws_ecr_repository.test.name

  policy = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "Expire images older than 14 days",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 7
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}
`, rName)
}

func testAccLifecyclePolicyConfig_order(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecr_repository" "test" {
  name = %[1]q
}

resource "aws_ecr_lifecycle_policy" "test" {
  repository = aws_ecr_repository.test.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Expire images older than 14 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 14
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 2
        description  = "Expire tagged images older than 14 days"
        selection = {
          tagStatus = "tagged"
          tagPrefixList = [
            "first",
            "second",
            "third",
          ]
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 14
        }
        action = {
          type = "expire"
        }
      },
    ]
  })
}
`, rName)
}

func testAccLifecyclePolicyConfig_newOrder(rName string) string {
	return fmt.Sprintf(`
resource "aws_ecr_repository" "test" {
  name = "%s"
}

resource "aws_ecr_lifecycle_policy" "test" {
  repository = aws_ecr_repository.test.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 2
        description  = "Expire tagged images older than 14 days"
        selection = {
          tagStatus = "tagged"
          tagPrefixList = [
            "third",
            "first",
            "second",
          ]
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 14
        }
        action = {
          type = "expire"
        }
      },
      {
        rulePriority = 1
        description  = "Expire images older than 14 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 14
        }
        action = {
          type = "expire"
        }
      },
    ]
  })
}
`, rName)
}
