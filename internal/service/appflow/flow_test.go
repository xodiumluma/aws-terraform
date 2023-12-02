// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appflow_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/service/appflow/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfappflow "github.com/hashicorp/terraform-provider-aws/internal/service/appflow"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccAppFlowFlow_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"
	scheduleStartTime := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_basic(rSourceName, rDestinationName, rFlowName, scheduleStartTime),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "appflow", regexache.MustCompile(`flow/.+`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rFlowName),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.0.connector_type"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.0.destination_connector_properties.#"),
					resource.TestCheckResourceAttrSet(resourceName, "source_flow_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "source_flow_config.0.connector_type"),
					resource.TestCheckResourceAttrSet(resourceName, "source_flow_config.0.source_connector_properties.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.0.source_fields.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.0.task_type"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_type", "Scheduled"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.0.data_pull_mode", "Incremental"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.0.schedule_expression", "rate(3hours)"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.0.schedule_start_time", scheduleStartTime),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "0"),
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

func TestAccAppFlowFlow_S3_outputFormatConfig_ParquetFileType(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"
	scheduleStartTime := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_S3_OutputFormatConfig_ParquetFileType(rSourceName, rDestinationName, rFlowName, scheduleStartTime, "PARQUET", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.0.connector_type"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.0.destination_connector_properties.#"),
					resource.TestCheckResourceAttr(resourceName, "destination_flow_config.0.destination_connector_properties.0.s3.0.s3_output_format_config.0.preserve_source_data_typing", "true"),
					resource.TestCheckResourceAttr(resourceName, "destination_flow_config.0.destination_connector_properties.0.s3.0.s3_output_format_config.0.file_type", "PARQUET"),
					resource.TestCheckResourceAttrSet(resourceName, "task.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.0.source_fields.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.0.task_type"),
				),
			},
			{
				Config: testAccFlowConfig_S3_OutputFormatConfig_ParquetFileType(rSourceName, rDestinationName, rFlowName, scheduleStartTime, "PARQUET", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.0.connector_type"),
					resource.TestCheckResourceAttrSet(resourceName, "destination_flow_config.0.destination_connector_properties.#"),
					resource.TestCheckResourceAttr(resourceName, "destination_flow_config.0.destination_connector_properties.0.s3.0.s3_output_format_config.0.preserve_source_data_typing", "false"),
					resource.TestCheckResourceAttr(resourceName, "destination_flow_config.0.destination_connector_properties.0.s3.0.s3_output_format_config.0.file_type", "PARQUET"),
					resource.TestCheckResourceAttrSet(resourceName, "task.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.0.source_fields.#"),
					resource.TestCheckResourceAttrSet(resourceName, "task.0.task_type"),
				),
			},
		},
	})
}

func TestAccAppFlowFlow_update(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"
	description := "test description"
	scheduleStartTime := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_basic(rSourceName, rDestinationName, rFlowName, scheduleStartTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
				),
			},
			{
				Config: testAccFlowConfig_update(rSourceName, rDestinationName, rFlowName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, description),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_type", "Scheduled"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.0.data_pull_mode", "Complete"),
					resource.TestCheckResourceAttr(resourceName, "trigger_config.0.trigger_properties.0.scheduled.0.schedule_expression", "rate(6hours)"),
				),
			},
		},
	})
}

func TestAccAppFlowFlow_taskProperties(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_taskProperties(rSourceName, rDestinationName, rFlowName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttr(resourceName, "task.0.task_properties.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "task.0.task_properties.SOURCE_DATA_TYPE", "CSV"),
					resource.TestCheckResourceAttr(resourceName, "task.0.task_properties.DESTINATION_DATA_TYPE", "CSV"),
				),
			},
		},
	})
}

func TestAccAppFlowFlow_taskUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_multipleTasks(rSourceName, rDestinationName, rFlowName, "aThirdTestField"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttr(resourceName, "task.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "task.*", map[string]string{
						"destination_field": "",
						"source_fields.#":   "2",
						"source_fields.0":   "testField",
						"source_fields.1":   "anotherTestField",
						"task_type":         "Filter",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "task.*", map[string]string{
						"destination_field":                     "testField",
						"source_fields.#":                       "1",
						"source_fields.0":                       "testField",
						"task_properties.%":                     "2",
						"task_properties.DESTINATION_DATA_TYPE": "string",
						"task_properties.SOURCE_DATA_TYPE":      "string",
						"task_type":                             "Map",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "task.*", map[string]string{
						"destination_field":                     "aThirdTestField",
						"source_fields.#":                       "1",
						"source_fields.0":                       "aThirdTestField",
						"task_properties.%":                     "2",
						"task_properties.DESTINATION_DATA_TYPE": "id",
						"task_properties.SOURCE_DATA_TYPE":      "id",
						"task_type":                             "Map",
					}),
				),
			},
			{
				Config: testAccFlowConfig_multipleTasks(rSourceName, rDestinationName, rFlowName, "anotherTestField"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttr(resourceName, "task.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "task.*", map[string]string{
						"destination_field": "",
						"source_fields.#":   "2",
						"source_fields.0":   "testField",
						"source_fields.1":   "anotherTestField",
						"task_type":         "Filter",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "task.*", map[string]string{
						"destination_field":                     "testField",
						"source_fields.#":                       "1",
						"source_fields.0":                       "testField",
						"task_properties.%":                     "2",
						"task_properties.DESTINATION_DATA_TYPE": "string",
						"task_properties.SOURCE_DATA_TYPE":      "string",
						"task_type":                             "Map",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "task.*", map[string]string{
						"destination_field":                     "anotherTestField",
						"source_fields.#":                       "1",
						"source_fields.0":                       "anotherTestField",
						"task_properties.%":                     "2",
						"task_properties.DESTINATION_DATA_TYPE": "id",
						"task_properties.SOURCE_DATA_TYPE":      "id",
						"task_type":                             "Map",
					}),
				),
			},
		},
	})
}

func TestAccAppFlowFlow_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_tags1(rSourceName, rDestinationName, rFlowName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
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
				Config: testAccFlowConfig_tags2(rSourceName, rDestinationName, rFlowName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccFlowConfig_tags1(rSourceName, rDestinationName, rFlowName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccAppFlowFlow_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var flowOutput types.FlowDefinition
	rSourceName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rDestinationName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rFlowName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appflow_flow.test"
	scheduleStartTime := time.Now().UTC().AddDate(0, 0, 1).Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.AppFlowEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckFlowDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccFlowConfig_basic(rSourceName, rDestinationName, rFlowName, scheduleStartTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFlowExists(ctx, resourceName, &flowOutput),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfappflow.ResourceFlow(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFlowConfig_base(rSourceName, rDestinationName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_s3_bucket" "test_source" {
  bucket = %[1]q
}

resource "aws_s3_bucket_policy" "test_source" {
  bucket = aws_s3_bucket.test_source.id
  policy = <<EOF
{
    "Statement": [
        {
            "Effect": "Allow",
            "Sid": "AllowAppFlowSourceActions",
            "Principal": {
                "Service": "appflow.amazonaws.com"
            },
            "Action": [
                "s3:ListBucket",
                "s3:GetObject"
            ],
            "Resource": [
                "arn:${data.aws_partition.current.partition}:s3:::%[1]s",
                "arn:${data.aws_partition.current.partition}:s3:::%[1]s/*"
            ]
        }
    ],
	"Version": "2012-10-17"
}
EOF
}

resource "aws_s3_object" "test" {
  bucket = aws_s3_bucket.test_source.id
  key    = "flow_source.csv"
  source = "test-fixtures/flow_source.csv"
}

resource "aws_s3_bucket" "test_destination" {
  bucket = %[2]q
}

resource "aws_s3_bucket_policy" "test_destination" {
  bucket = aws_s3_bucket.test_destination.id
  policy = <<EOF

{
    "Statement": [
        {
            "Effect": "Allow",
            "Sid": "AllowAppFlowDestinationActions",
            "Principal": {
                "Service": "appflow.amazonaws.com"
            },
            "Action": [
                "s3:PutObject",
                "s3:AbortMultipartUpload",
                "s3:ListMultipartUploadParts",
                "s3:ListBucketMultipartUploads",
                "s3:GetBucketAcl",
                "s3:PutObjectAcl"
            ],
            "Resource": [
                "arn:${data.aws_partition.current.partition}:s3:::%[2]s",
                "arn:${data.aws_partition.current.partition}:s3:::%[2]s/*"
            ]
        }
    ],
	"Version": "2012-10-17"
}
EOF
}
`, rSourceName, rDestinationName)
}

func testAccFlowConfig_basic(rSourceName, rDestinationName, rFlowName, scheduleStartTime string) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name = %[1]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }
        }
      }
    }
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    connector_operator {
      s3 = "NO_OP"
    }
  }

  trigger_config {
    trigger_type = "Scheduled"

    trigger_properties {
      scheduled {
        data_pull_mode      = "Incremental"
        schedule_expression = "rate(3hours)"
        schedule_start_time = %[2]q
      }
    }
  }
}
`, rFlowName, scheduleStartTime),
	)
}

func testAccFlowConfig_S3_OutputFormatConfig_ParquetFileType(rSourceName, rDestinationName, rFlowName, scheduleStartTime, fileType string, preserveSourceDataTyping bool) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name = %[1]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }

          file_type                   = %[3]q
          preserve_source_data_typing = %[4]t

          aggregation_config {
            aggregation_type = "None"
          }
        }
      }
    }
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    task_properties = {
      "DESTINATION_DATA_TYPE" = "string"
      "SOURCE_DATA_TYPE"      = "string"
    }

    connector_operator {
      s3 = "NO_OP"
    }
  }

  trigger_config {
    trigger_type = "Scheduled"

    trigger_properties {
      scheduled {
        data_pull_mode      = "Incremental"
        schedule_expression = "rate(3hours)"
        schedule_start_time = %[2]q
      }
    }
  }
}
`, rFlowName, scheduleStartTime, fileType, preserveSourceDataTyping),
	)
}

func testAccFlowConfig_update(rSourceName, rDestinationName, rFlowName, description string) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name        = %[1]q
  description = %[2]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }
        }
      }
    }
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    connector_operator {
      s3 = "NO_OP"
    }
  }

  trigger_config {
    trigger_type = "Scheduled"

    trigger_properties {
      scheduled {
        data_pull_mode      = "Complete"
        schedule_expression = "rate(6hours)"
      }
    }
  }
}
`, rFlowName, description),
	)
}

func testAccFlowConfig_taskProperties(rSourceName, rDestinationName, rFlowName string) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name = %[1]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }
        }
      }
    }
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    task_properties = {
      SOURCE_DATA_TYPE      = "CSV"
      DESTINATION_DATA_TYPE = "CSV"
    }

    connector_operator {
      s3 = "NO_OP"
    }
  }

  trigger_config {
    trigger_type = "OnDemand"
  }
}
`, rFlowName),
	)
}

func testAccFlowConfig_multipleTasks(rSourceName, rDestinationName, rFlowName, fieldName string) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name = %[1]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }
        }
      }
    }
  }

  task {
    source_fields = [
      "testField",
      "anotherTestField"
    ]
    connector_operator {
      s3 = "PROJECTION"
    }
    task_type         = "Filter"
    destination_field = ""
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    connector_operator {
      s3 = "NO_OP"
    }

    task_properties = {
      "DESTINATION_DATA_TYPE" = "string"
      "SOURCE_DATA_TYPE"      = "string"
    }
  }

  task {
    source_fields     = [%[2]q]
    destination_field = %[2]q
    task_type         = "Map"

    connector_operator {
      s3 = "NO_OP"
    }

    task_properties = {
      "DESTINATION_DATA_TYPE" = "id"
      "SOURCE_DATA_TYPE"      = "id"
    }
  }

  trigger_config {
    trigger_type = "OnDemand"
  }
}
`, rFlowName, fieldName),
	)
}

func testAccFlowConfig_tags1(rSourceName, rDestinationName, rFlowName string, tagKey1 string, tagValue1 string) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name = %[1]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }
        }
      }
    }
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    connector_operator {
      s3 = "NO_OP"
    }
  }

  trigger_config {
    trigger_type = "OnDemand"
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rFlowName, tagKey1, tagValue1),
	)
}

func testAccFlowConfig_tags2(rSourceName, rDestinationName, rFlowName string, tagKey1 string, tagValue1 string, tagKey2 string, tagValue2 string) string {
	return acctest.ConfigCompose(
		testAccFlowConfig_base(rSourceName, rDestinationName),
		fmt.Sprintf(`
resource "aws_appflow_flow" "test" {
  name = %[1]q

  source_flow_config {
    connector_type = "S3"
    source_connector_properties {
      s3 {
        bucket_name   = aws_s3_bucket_policy.test_source.bucket
        bucket_prefix = "flow"
      }
    }
  }

  destination_flow_config {
    connector_type = "S3"
    destination_connector_properties {
      s3 {
        bucket_name = aws_s3_bucket_policy.test_destination.bucket

        s3_output_format_config {
          prefix_config {
            prefix_type = "PATH"
          }
        }
      }
    }
  }

  task {
    source_fields     = ["testField"]
    destination_field = "testField"
    task_type         = "Map"

    connector_operator {
      s3 = "NO_OP"
    }
  }

  trigger_config {
    trigger_type = "OnDemand"
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rFlowName, tagKey1, tagValue1, tagKey2, tagValue2),
	)
}

func testAccCheckFlowExists(ctx context.Context, n string, v *types.FlowDefinition) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppFlowClient(ctx)

		output, err := tfappflow.FindFlowByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckFlowDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).AppFlowClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_appflow_flow" {
				continue
			}

			_, err := tfappflow.FindFlowByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("AppFlow Flow %s still exists", rs.Primary.ID)
		}

		return nil
	}
}
