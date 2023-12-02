// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sagemaker_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsagemaker "github.com/hashicorp/terraform-provider-aws/internal/service/sagemaker"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccSageMakerPipeline_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var pipeline sagemaker.DescribePipelineOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rNameUpdated := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_pipeline.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPipelineDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPipelinePipelineConfig_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "pipeline_name", rName),
					resource.TestCheckResourceAttr(resourceName, "pipeline_display_name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexache.MustCompile(`pipeline/.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "parallelism_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPipelinePipelineConfig_basic(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "pipeline_name", rName),
					resource.TestCheckResourceAttr(resourceName, "pipeline_display_name", rNameUpdated),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "sagemaker", regexache.MustCompile(`pipeline/.+`)),
					resource.TestCheckResourceAttrPair(resourceName, "role_arn", "aws_iam_role.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "parallelism_configuration.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccSageMakerPipeline_parallelism(t *testing.T) {
	ctx := acctest.Context(t)
	var pipeline sagemaker.DescribePipelineOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_pipeline.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPipelineDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPipelinePipelineConfig_parallelism(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "pipeline_name", rName),
					resource.TestCheckResourceAttr(resourceName, "parallelism_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "parallelism_configuration.0.max_parallel_execution_steps", "1"),
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

func TestAccSageMakerPipeline_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var pipeline sagemaker.DescribePipelineOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_pipeline.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPipelineDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPipelinePipelineConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
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
				Config: testAccPipelinePipelineConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccPipelinePipelineConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccSageMakerPipeline_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var pipeline sagemaker.DescribePipelineOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_sagemaker_pipeline.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, sagemaker.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPipelineDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPipelinePipelineConfig_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPipelineExists(ctx, resourceName, &pipeline),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfsagemaker.ResourcePipeline(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckPipelineDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_sagemaker_pipeline" {
				continue
			}

			_, err := tfsagemaker.FindPipelineByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("SageMaker Pipeline %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckPipelineExists(ctx context.Context, n string, pipeline *sagemaker.DescribePipelineOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SageMaker Pipeline ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SageMakerConn(ctx)

		output, err := tfsagemaker.FindPipelineByName(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*pipeline = *output

		return nil
	}
}

func testAccPipelinePipelineConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name               = %[1]q
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.test.json
}

data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["sagemaker.amazonaws.com"]
    }
  }
}
`, rName)
}

func testAccPipelinePipelineConfig_basic(rName, dispName string) string {
	return acctest.ConfigCompose(testAccPipelinePipelineConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_pipeline" "test" {
  pipeline_name         = %[1]q
  pipeline_display_name = %[2]q
  role_arn              = aws_iam_role.test.arn

  pipeline_definition = jsonencode({
    Version = "2020-12-01"
    Steps = [{
      Name = "Test"
      Type = "Fail"
      Arguments = {
        ErrorMessage = "test"
      }
    }]
  })
}
`, rName, dispName))
}

func testAccPipelinePipelineConfig_parallelism(rName string) string {
	return acctest.ConfigCompose(testAccPipelinePipelineConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_pipeline" "test" {
  pipeline_name         = %[1]q
  pipeline_display_name = %[1]q
  role_arn              = aws_iam_role.test.arn

  pipeline_definition = jsonencode({
    Version = "2020-12-01"
    Steps = [{
      Name = "Test"
      Type = "Fail"
      Arguments = {
        ErrorMessage = "test"
      }
    }]
  })

  parallelism_configuration {
    max_parallel_execution_steps = 1
  }
}
`, rName))
}

func testAccPipelinePipelineConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(testAccPipelinePipelineConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_pipeline" "test" {
  pipeline_name         = %[1]q
  pipeline_display_name = %[1]q
  role_arn              = aws_iam_role.test.arn

  pipeline_definition = jsonencode({
    Version = "2020-12-01"
    Steps = [{
      Name = "Test"
      Type = "Fail"
      Arguments = {
        ErrorMessage = "test"
      }
    }]
  })

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccPipelinePipelineConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(testAccPipelinePipelineConfig_base(rName), fmt.Sprintf(`
resource "aws_sagemaker_pipeline" "test" {
  pipeline_name         = %[1]q
  pipeline_display_name = %[1]q
  role_arn              = aws_iam_role.test.arn

  pipeline_definition = jsonencode({
    Version = "2020-12-01"
    Steps = [{
      Name = "Test"
      Type = "Fail"
      Arguments = {
        ErrorMessage = "test"
      }
    }]
  })

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
