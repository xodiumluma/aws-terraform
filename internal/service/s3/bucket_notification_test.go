package s3_test

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccS3BucketNotification_eventbridge(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_eventBridge(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketEventBridgeNotification(ctx, "aws_s3_bucket.test")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccS3BucketNotification_lambdaFunction(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_lambdaFunction(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLambdaFunctionConfiguration(ctx, "aws_s3_bucket.test",
						"notification-lambda",
						"aws_lambda_function.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						&s3.KeyFilter{
							FilterRules: []*s3.FilterRule{
								{
									Name:  aws.String("Prefix"),
									Value: aws.String("tf-acc-test/"),
								},
								{
									Name:  aws.String("Suffix"),
									Value: aws.String(".png"),
								},
							},
						},
					),
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

func TestAccS3BucketNotification_LambdaFunctionLambdaFunctionARN_alias(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_lambdaFunctionLambdaFunctionARNAlias(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketLambdaFunctionConfiguration(ctx, "aws_s3_bucket.test",
						"test",
						"aws_lambda_alias.test",
						[]string{"s3:ObjectCreated:*"},
						nil,
					),
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

func TestAccS3BucketNotification_queue(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_queue(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketQueueNotification(ctx, "aws_s3_bucket.test",
						"notification-sqs",
						"aws_sqs_queue.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						&s3.KeyFilter{
							FilterRules: []*s3.FilterRule{
								{
									Name:  aws.String("Prefix"),
									Value: aws.String("tf-acc-test/"),
								},
								{
									Name:  aws.String("Suffix"),
									Value: aws.String(".mp4"),
								},
							},
						},
					),
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

func TestAccS3BucketNotification_topic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_topic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketTopicNotification(ctx, "aws_s3_bucket.test",
						"notification-sns1",
						"aws_sns_topic.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						nil,
					),
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

func TestAccS3BucketNotification_Topic_multiple(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_notification.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_topicMultiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketTopicNotification(ctx, "aws_s3_bucket.test",
						"notification-sns1",
						"aws_sns_topic.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						&s3.KeyFilter{
							FilterRules: []*s3.FilterRule{
								{
									Name:  aws.String("Prefix"),
									Value: aws.String("tf-acc-test/"),
								},
								{
									Name:  aws.String("Suffix"),
									Value: aws.String(".txt"),
								},
							},
						},
					),
					testAccCheckBucketTopicNotification(ctx, "aws_s3_bucket.test",
						"notification-sns2",
						"aws_sns_topic.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						&s3.KeyFilter{
							FilterRules: []*s3.FilterRule{
								{
									Name:  aws.String("Suffix"),
									Value: aws.String(".log"),
								},
							},
						},
					),
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

func TestAccS3BucketNotification_update(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketNotificationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketNotificationConfig_topic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketTopicNotification(ctx, "aws_s3_bucket.test",
						"notification-sns1",
						"aws_sns_topic.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						nil,
					),
				),
			},
			{
				Config: testAccBucketNotificationConfig_queue(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketQueueNotification(ctx, "aws_s3_bucket.test",
						"notification-sqs",
						"aws_sqs_queue.test",
						[]string{"s3:ObjectCreated:*", "s3:ObjectRemoved:Delete"},
						&s3.KeyFilter{
							FilterRules: []*s3.FilterRule{
								{
									Name:  aws.String("Prefix"),
									Value: aws.String("tf-acc-test/"),
								},
								{
									Name:  aws.String("Suffix"),
									Value: aws.String(".mp4"),
								},
							},
						},
					),
				),
			},
		},
	})
}

func testAccCheckBucketNotificationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_s3_bucket_notification" {
				continue
			}
			err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
				out, err := conn.GetBucketNotificationConfigurationWithContext(ctx, &s3.GetBucketNotificationConfigurationRequest{
					Bucket: aws.String(rs.Primary.ID),
				})

				if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket) {
					return nil
				}

				if err != nil {
					return retry.NonRetryableError(err)
				}

				if len(out.TopicConfigurations) > 0 {
					return retry.RetryableError(fmt.Errorf("TopicConfigurations is exists: %v", out))
				}
				if len(out.LambdaFunctionConfigurations) > 0 {
					return retry.RetryableError(fmt.Errorf("LambdaFunctionConfigurations is exists: %v", out))
				}
				if len(out.QueueConfigurations) > 0 {
					return retry.RetryableError(fmt.Errorf("QueueConfigurations is exists: %v", out))
				}

				return nil
			})

			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccCheckBucketTopicNotification(ctx context.Context, n, i, t string, events []string, filters *s3.KeyFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		topicArn := s.RootModule().Resources[t].Primary.ID
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
			out, err := conn.GetBucketNotificationConfigurationWithContext(ctx, &s3.GetBucketNotificationConfigurationRequest{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("GetBucketNotification error: %v", err))
			}

			eventSlice := sort.StringSlice(events)
			eventSlice.Sort()

			outputTopics := out.TopicConfigurations
			matched := false
			for _, outputTopic := range outputTopics {
				if *outputTopic.Id == i {
					matched = true

					if *outputTopic.TopicArn != topicArn {
						return retry.RetryableError(fmt.Errorf("bad topic arn, expected: %s, got %#v", topicArn, *outputTopic.TopicArn))
					}

					if filters != nil {
						if !reflect.DeepEqual(filters, outputTopic.Filter.Key) {
							return retry.RetryableError(fmt.Errorf("bad notification filters, expected: %#v, got %#v", filters, outputTopic.Filter.Key))
						}
					} else {
						if outputTopic.Filter != nil {
							return retry.RetryableError(fmt.Errorf("bad notification filters, expected: nil, got %#v", outputTopic.Filter))
						}
					}

					outputEventSlice := sort.StringSlice(aws.StringValueSlice(outputTopic.Events))
					outputEventSlice.Sort()
					if !reflect.DeepEqual(eventSlice, outputEventSlice) {
						return retry.RetryableError(fmt.Errorf("bad notification events, expected: %#v, got %#v", events, outputEventSlice))
					}
				}
			}

			if !matched {
				return retry.RetryableError(fmt.Errorf("No match topic configurations: %#v", out))
			}

			return nil
		})

		return err
	}
}

func testAccCheckBucketEventBridgeNotification(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
			out, err := conn.GetBucketNotificationConfigurationWithContext(ctx, &s3.GetBucketNotificationConfigurationRequest{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("GetBucketNotification error: %v", err))
			}

			if out.EventBridgeConfiguration == nil {
				return retry.RetryableError(fmt.Errorf("No EventBridge configuration: %#v", out))
			} else {
				return nil
			}
		})

		return err
	}
}

func testAccCheckBucketQueueNotification(ctx context.Context, n, i, t string, events []string, filters *s3.KeyFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		queueArn := s.RootModule().Resources[t].Primary.Attributes["arn"]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
			out, err := conn.GetBucketNotificationConfigurationWithContext(ctx, &s3.GetBucketNotificationConfigurationRequest{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("GetBucketNotification error: %v", err))
			}

			eventSlice := sort.StringSlice(events)
			eventSlice.Sort()

			outputQueues := out.QueueConfigurations
			matched := false
			for _, outputQueue := range outputQueues {
				if *outputQueue.Id == i {
					matched = true

					if *outputQueue.QueueArn != queueArn {
						return retry.RetryableError(fmt.Errorf("bad queue arn, expected: %s, got %#v", queueArn, *outputQueue.QueueArn))
					}

					if filters != nil {
						if !reflect.DeepEqual(filters, outputQueue.Filter.Key) {
							return retry.RetryableError(fmt.Errorf("bad notification filters, expected: %#v, got %#v", filters, outputQueue.Filter.Key))
						}
					} else {
						if outputQueue.Filter != nil {
							return retry.RetryableError(fmt.Errorf("bad notification filters, expected: nil, got %#v", outputQueue.Filter))
						}
					}

					outputEventSlice := sort.StringSlice(aws.StringValueSlice(outputQueue.Events))
					outputEventSlice.Sort()
					if !reflect.DeepEqual(eventSlice, outputEventSlice) {
						return retry.RetryableError(fmt.Errorf("bad notification events, expected: %#v, got %#v", events, outputEventSlice))
					}
				}
			}

			if !matched {
				return retry.RetryableError(fmt.Errorf("No match queue configurations: %#v", out))
			}

			return nil
		})

		return err
	}
}

func testAccCheckBucketLambdaFunctionConfiguration(ctx context.Context, n, i, t string, events []string, filters *s3.KeyFilter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		funcArn := s.RootModule().Resources[t].Primary.Attributes["arn"]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		err := retry.RetryContext(ctx, 1*time.Minute, func() *retry.RetryError {
			out, err := conn.GetBucketNotificationConfigurationWithContext(ctx, &s3.GetBucketNotificationConfigurationRequest{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("GetBucketNotification error: %v", err))
			}

			eventSlice := sort.StringSlice(events)
			eventSlice.Sort()

			outputFunctions := out.LambdaFunctionConfigurations
			matched := false
			for _, outputFunc := range outputFunctions {
				if *outputFunc.Id == i {
					matched = true

					if *outputFunc.LambdaFunctionArn != funcArn {
						return retry.RetryableError(fmt.Errorf("bad lambda function arn, expected: %s, got %#v", funcArn, *outputFunc.LambdaFunctionArn))
					}

					if filters != nil {
						if !reflect.DeepEqual(filters, outputFunc.Filter.Key) {
							return retry.RetryableError(fmt.Errorf("bad notification filters, expected: %#v, got %#v", filters, outputFunc.Filter.Key))
						}
					} else {
						if outputFunc.Filter != nil {
							return retry.RetryableError(fmt.Errorf("bad notification filters, expected: nil, got %#v", outputFunc.Filter))
						}
					}

					outputEventSlice := sort.StringSlice(aws.StringValueSlice(outputFunc.Events))
					outputEventSlice.Sort()
					if !reflect.DeepEqual(eventSlice, outputEventSlice) {
						return retry.RetryableError(fmt.Errorf("bad notification events, expected: %#v, got %#v", events, outputEventSlice))
					}
				}
			}

			if !matched {
				return retry.RetryableError(fmt.Errorf("No match lambda function configurations: %#v", out))
			}

			return nil
		})

		return err
	}
}

func testAccBucketNotificationConfig_eventBridge(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "test" {
  depends_on = [
    aws_s3_bucket_public_access_block.test,
    aws_s3_bucket_ownership_controls.test,
  ]

  bucket = aws_s3_bucket.test.id
  acl    = "public-read"
}

resource "aws_s3_bucket_notification" "test" {
  bucket = aws_s3_bucket.test.id

  eventbridge = true
}
`, rName)
}

func testAccBucketNotificationConfig_topicMultiple(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_sns_topic" "test" {
  name = %[1]q

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": { "Service": "s3.${data.aws_partition.current.dns_suffix}" },
      "Action": "SNS:Publish",
      "Resource": "arn:${data.aws_partition.current.partition}:sns:*:*:%[1]s",
      "Condition": {
        "ArnLike": {
          "aws:SourceArn": "${aws_s3_bucket.test.arn}"
        }
      }
    }
  ]
}
POLICY
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "test" {
  depends_on = [
    aws_s3_bucket_public_access_block.test,
    aws_s3_bucket_ownership_controls.test,
  ]

  bucket = aws_s3_bucket.test.id
  acl    = "public-read"
}

resource "aws_s3_bucket_notification" "test" {
  bucket = aws_s3_bucket.test.id

  topic {
    id        = "notification-sns1"
    topic_arn = aws_sns_topic.test.arn

    events = [
      "s3:ObjectCreated:*",
      "s3:ObjectRemoved:Delete",
    ]

    filter_prefix = "tf-acc-test/"
    filter_suffix = ".txt"
  }

  topic {
    id        = "notification-sns2"
    topic_arn = aws_sns_topic.test.arn

    events = [
      "s3:ObjectCreated:*",
      "s3:ObjectRemoved:Delete",
    ]

    filter_suffix = ".log"
  }
}
`, rName)
}

func testAccBucketNotificationConfig_queue(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_sqs_queue" "test" {
  name = %[1]q

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": "*",
      "Action": "sqs:SendMessage",
      "Resource": "arn:${data.aws_partition.current.partition}:sqs:*:*:%[1]s",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_s3_bucket.test.arn}"
        }
      }
    }
  ]
}
POLICY
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "test" {
  depends_on = [
    aws_s3_bucket_public_access_block.test,
    aws_s3_bucket_ownership_controls.test,
  ]

  bucket = aws_s3_bucket.test.id
  acl    = "public-read"
}

resource "aws_s3_bucket_notification" "test" {
  bucket = aws_s3_bucket.test.id

  queue {
    id        = "notification-sqs"
    queue_arn = aws_sqs_queue.test.arn

    events = [
      "s3:ObjectCreated:*",
      "s3:ObjectRemoved:Delete",
    ]

    filter_prefix = "tf-acc-test/"
    filter_suffix = ".mp4"
  }
}
`, rName)
}

func testAccBucketNotificationConfig_lambdaFunction(rName string) string {
	return fmt.Sprintf(`
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

resource "aws_lambda_permission" "test" {
  statement_id  = "AllowExecutionFromS3Bucket"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test.arn
  principal     = "s3.${data.aws_partition.current.dns_suffix}"
  source_arn    = aws_s3_bucket.test.arn
}

resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = %[1]q
  role          = aws_iam_role.test.arn
  handler       = "exports.example"
  runtime       = "nodejs16.x"
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "test" {
  depends_on = [
    aws_s3_bucket_public_access_block.test,
    aws_s3_bucket_ownership_controls.test,
  ]

  bucket = aws_s3_bucket.test.id
  acl    = "public-read"
}

resource "aws_s3_bucket_notification" "test" {
  bucket = aws_s3_bucket.test.id

  lambda_function {
    id                  = "notification-lambda"
    lambda_function_arn = aws_lambda_function.test.arn

    events = [
      "s3:ObjectCreated:*",
      "s3:ObjectRemoved:Delete",
    ]

    filter_prefix = "tf-acc-test/"
    filter_suffix = ".png"
  }
}
`, rName)
}

func testAccBucketNotificationConfig_lambdaFunctionLambdaFunctionARNAlias(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

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
      "Effect": "Allow"
    }
  ]
}
EOF

  name = %[1]q
}

resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = %[1]q
  handler       = "exports.example"
  role          = aws_iam_role.test.arn
  runtime       = "nodejs16.x"
}

resource "aws_lambda_alias" "test" {
  function_name    = aws_lambda_function.test.arn
  function_version = "$LATEST"
  name             = "testalias"
}

resource "aws_lambda_permission" "test" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.test.arn
  principal     = "s3.${data.aws_partition.current.dns_suffix}"
  qualifier     = aws_lambda_alias.test.name
  source_arn    = aws_s3_bucket.test.arn
  statement_id  = "AllowExecutionFromS3Bucket"
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_notification" "test" {
  bucket = aws_s3_bucket.test.id

  lambda_function {
    events              = ["s3:ObjectCreated:*"]
    id                  = "test"
    lambda_function_arn = aws_lambda_alias.test.arn
  }
}
`, rName)
}

func testAccBucketNotificationConfig_topic(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_sns_topic" "test" {
  name = %[1]q

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "s3.${data.aws_partition.current.dns_suffix}"
      },
      "Action": "SNS:Publish",
      "Resource": "arn:${data.aws_partition.current.partition}:sns:*:*:%[1]s",
      "Condition": {
        "ArnLike": {
          "aws:SourceArn": "${aws_s3_bucket.test.arn}"
        }
      }
    }
  ]
}
POLICY
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_ownership_controls" "test" {
  bucket = aws_s3_bucket.test.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "test" {
  depends_on = [
    aws_s3_bucket_public_access_block.test,
    aws_s3_bucket_ownership_controls.test,
  ]

  bucket = aws_s3_bucket.test.id
  acl    = "public-read"
}

resource "aws_s3_bucket_notification" "test" {
  bucket = aws_s3_bucket.test.id

  topic {
    id        = "notification-sns1"
    topic_arn = aws_sns_topic.test.arn

    events = [
      "s3:ObjectCreated:*",
      "s3:ObjectRemoved:Delete",
    ]
  }
}
`, rName)
}
