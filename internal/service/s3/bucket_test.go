// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package s3_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/YakDriver/regexache"
	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	s3_sdkv2 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfcloudformation "github.com/hashicorp/terraform-provider-aws/internal/service/cloudformation"
	tfs3 "github.com/hashicorp/terraform-provider-aws/internal/service/s3"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func init() {
	acctest.RegisterServiceErrorCheckFunc(s3.EndpointsID, testAccErrorCheckSkip)
}

// testAccErrorCheckSkip skips tests that have error messages indicating unsupported features
func testAccErrorCheckSkip(t *testing.T) resource.ErrorCheckFunc {
	return acctest.ErrorCheckSkipMessagesContaining(t,
		"Number of distinct destination bucket ARNs cannot exceed",
		"destination is not allowed",
	)
}

func TestAccS3Bucket_Basic_basic(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	region := acctest.Region()
	hostedZoneID, _ := tfs3.HostedZoneIDForRegion(region)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "hosted_zone_id", hostedZoneID),
					resource.TestCheckResourceAttr(resourceName, "region", region),
					resource.TestCheckNoResourceAttr(resourceName, "website_endpoint"),
					acctest.CheckResourceAttrGlobalARNNoAccount(resourceName, "arn", "s3", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket", bucketName),
					testAccCheckBucketDomainName(resourceName, "bucket_domain_name", bucketName),
					resource.TestCheckResourceAttr(resourceName, "bucket_regional_domain_name", testAccBucketRegionalDomainName(bucketName, region)),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.mfa_delete", "false"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

// Support for common Terraform 0.11 pattern
// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/7868
func TestAccS3Bucket_Basic_emptyString(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_emptyString,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					acctest.CheckResourceAttrNameGenerated(resourceName, "bucket"),
					resource.TestCheckResourceAttr(resourceName, "bucket_prefix", id.UniqueIdPrefix),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Basic_nameGenerated(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_nameGenerated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					acctest.CheckResourceAttrNameGenerated(resourceName, "bucket"),
					resource.TestCheckResourceAttr(resourceName, "bucket_prefix", id.UniqueIdPrefix),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Basic_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_namePrefix("tf-test-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "bucket", "tf-test-"),
					resource.TestCheckResourceAttr(resourceName, "bucket_prefix", "tf-test-"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Basic_forceDestroy(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_forceDestroy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					testAccCheckBucketAddObjects(ctx, resourceName, "data.txt", "prefix/more_data.txt"),
				),
			},
		},
	})
}

// By default, the AWS Go SDK cleans up URIs by removing extra slashes
// when the service API requests use the URI as part of making a request.
// While the aws_s3_object resource automatically cleans the key
// to not contain these extra slashes, out-of-band handling and other AWS
// services may create keys with extra slashes (empty "directory" prefixes).
func TestAccS3Bucket_Basic_forceDestroyWithEmptyPrefixes(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_forceDestroy(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					testAccCheckBucketAddObjects(ctx, resourceName, "data.txt", "/extraleadingslash.txt"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Basic_forceDestroyWithObjectLockEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_forceDestroyObjectLockEnabled(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					testAccCheckBucketAddObjectsWithLegalHold(ctx, resourceName, "data.txt", "prefix/more_data.txt"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Basic_acceleration(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_acceleration(bucketName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "acceleration_status", s3.BucketAccelerateStatusEnabled),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccBucketConfig_acceleration(bucketName, s3.BucketAccelerateStatusSuspended),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "acceleration_status", s3.BucketAccelerateStatusSuspended),
				),
			},
		},
	})
}

func TestAccS3Bucket_Basic_keyEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_defaultEncryptionKeyEnabledKMSMasterKey(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.sse_algorithm", "aws:kms"),
					resource.TestMatchResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.kms_master_key_id", regexache.MustCompile("^arn")),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.bucket_key_enabled", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Basic_requestPayer(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_requestPayer(bucketName, s3.PayerBucketOwner),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "request_payer", s3.PayerBucketOwner),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				Config: testAccBucketConfig_requestPayer(bucketName, s3.PayerRequester),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "request_payer", s3.PayerRequester),
				),
			},
		},
	})
}

// Test TestAccS3Bucket_disappears is designed to fail with a "plan
// not empty" error in Terraform, to check against regressions.
// See https://github.com/hashicorp/terraform/pull/2925
func TestAccS3Bucket_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfs3.ResourceBucket(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3Bucket_Duplicate_basic(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	region := acctest.Region()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckRegionNot(t, endpoints.UsEast1RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketConfig_duplicate(region, bucketName),
				ExpectError: regexache.MustCompile(s3.ErrCodeBucketAlreadyOwnedByYou),
			},
		},
	})
}

func TestAccS3Bucket_Duplicate_UsEast1(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartition(t, endpoints.AwsPartitionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketConfig_duplicate(endpoints.UsEast1RegionID, bucketName),
				ExpectError: regexache.MustCompile(tfs3.ErrMessageBucketAlreadyExists),
			},
		},
	})
}

func TestAccS3Bucket_Duplicate_UsEast1AltAccount(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartition(t, endpoints.AwsPartitionID)
			acctest.PreCheckAlternateAccount(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketConfig_duplicateAltAccount(endpoints.UsEast1RegionID, bucketName),
				ExpectError: regexache.MustCompile(s3.ErrCodeBucketAlreadyExists),
			},
		},
	})
}

func TestAccS3Bucket_Tags_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rInt := sdkacctest.RandInt()
	resourceName := "aws_s3_bucket.bucket1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_multiTags(rInt),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Tags_withNoSystemTags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_tags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "AAA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccBucketConfig_updatedTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "XXX"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key5", "EEE"),
				),
			},
			{
				Config: testAccBucketConfig_noTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			// Verify update from 0 tags.
			{
				Config: testAccBucketConfig_tags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "AAA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Tags_withSystemTags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")

	var stackID string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			testAccCheckBucketDestroy(ctx),
			func(s *terraform.State) error {
				// Tear down CF stack.
				conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFormationConn(ctx)

				requestToken := id.UniqueId()
				req := &cloudformation.DeleteStackInput{
					StackName:          aws.String(stackID),
					ClientRequestToken: aws.String(requestToken),
				}

				log.Printf("[DEBUG] Deleting CloudFormation stack: %s", req)
				if _, err := conn.DeleteStackWithContext(ctx, req); err != nil {
					return fmt.Errorf("error deleting CloudFormation stack: %w", err)
				}

				if _, err := tfcloudformation.WaitStackDeleted(ctx, conn, stackID, requestToken, 10*time.Minute); err != nil {
					return fmt.Errorf("Error waiting for CloudFormation stack deletion: %s", err)
				}

				return nil
			},
		),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_noTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfs3.ResourceBucket(), resourceName),
					testAccCheckBucketCreateViaCloudFormation(ctx, bucketName, &stackID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccBucketConfig_tags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "AAA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
					testAccCheckBucketTagKeys(ctx, resourceName, "aws:cloudformation:stack-name", "aws:cloudformation:stack-id", "aws:cloudformation:logical-id"),
				),
			},
			{
				Config: testAccBucketConfig_updatedTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "XXX"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key4", "DDD"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key5", "EEE"),
					testAccCheckBucketTagKeys(ctx, resourceName, "aws:cloudformation:stack-name", "aws:cloudformation:stack-id", "aws:cloudformation:logical-id"),
				),
			},
			{
				Config: testAccBucketConfig_noTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					testAccCheckBucketTagKeys(ctx, resourceName, "aws:cloudformation:stack-name", "aws:cloudformation:stack-id", "aws:cloudformation:logical-id"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Tags_ignoreTags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.test"
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigIgnoreTagsKeyPrefixes1("ignorekey"),
					testAccBucketConfig_noTags(bucketName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					testAccCheckBucketUpdateTags(ctx, resourceName, nil, map[string]string{"ignorekey1": "ignorevalue1"}),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					testAccCheckBucketCheckTags(ctx, resourceName, map[string]string{
						"ignorekey1": "ignorevalue1",
					}),
				),
			},
			{
				Config: acctest.ConfigCompose(
					acctest.ConfigIgnoreTagsKeyPrefixes1("ignorekey"),
					testAccBucketConfig_tags(bucketName)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key1", "AAA"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key2", "BBB"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key3", "CCC"),
					testAccCheckBucketCheckTags(ctx, resourceName, map[string]string{
						"ignorekey1": "ignorevalue1",
						"Key1":       "AAA",
						"Key2":       "BBB",
						"Key3":       "CCC",
					}),
				),
			},
		},
	})
}

func TestAccS3Bucket_Manage_lifecycleBasic(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_lifecycle(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.0.days", "365"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.0.date", ""),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.0.expired_object_delete_marker", "false"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.0.transition.*", map[string]string{
						"date":          "",
						"days":          "30",
						"storage_class": "STANDARD_IA",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.0.transition.*", map[string]string{
						"date":          "",
						"days":          "60",
						"storage_class": "INTELLIGENT_TIERING",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.0.transition.*", map[string]string{
						"date":          "",
						"days":          "90",
						"storage_class": "ONEZONE_IA",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.0.transition.*", map[string]string{
						"date":          "",
						"days":          "120",
						"storage_class": "GLACIER",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.0.transition.*", map[string]string{
						"date":          "",
						"days":          "210",
						"storage_class": "DEEP_ARCHIVE",
					}),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.expiration.0.date", "2016-01-12"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.expiration.0.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.1.expiration.0.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.2.transition.*", map[string]string{
						"days": "0",
					}),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.3.id", "id4"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.3.prefix", "path4/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.3.tags.tagKey", "tagValue"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.3.tags.terraform", "hashicorp"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.4.id", "id5"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.4.tags.tagKey", "tagValue"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.4.tags.terraform", "hashicorp"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.4.transition.*", map[string]string{
						"days":          "0",
						"storage_class": "GLACIER",
					}),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.5.id", "id6"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.5.tags.tagKey", "tagValue"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.5.transition.*", map[string]string{
						"days":          "0",
						"storage_class": "GLACIER",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
				),
			},
		},
	})
}

func TestAccS3Bucket_Manage_lifecycleExpireMarkerOnly(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_lifecycleExpireMarker(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.0.days", "0"),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.0.date", ""),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.0.expiration.0.expired_object_delete_marker", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/11420
func TestAccS3Bucket_Manage_lifecycleRuleExpirationEmptyBlock(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_lifecycleRuleExpirationEmptyBlock(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/15138
func TestAccS3Bucket_Manage_lifecycleRuleAbortIncompleteMultipartUploadDaysNoExpiration(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_lifecycleRuleAbortIncompleteMultipartUploadDays(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_lifecycleRemove(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_lifecycle(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.#", "6"),
				),
			},
			{
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					// As Lifecycle Rule is a Computed field, removing them from terraform will not
					// trigger an update to remove them from the S3 bucket.
					resource.TestCheckResourceAttr(resourceName, "lifecycle_rule.#", "6"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Manage_objectLock(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_objectLockEnabledNoDefaultRetention(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.object_lock_enabled", s3.ObjectLockEnabledEnabled),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
			{
				Config: testAccBucketConfig_objectLockEnabledDefaultRetention(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.object_lock_enabled", "Enabled"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.rule.0.default_retention.0.mode", "COMPLIANCE"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.rule.0.default_retention.0.days", "3"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Manage_objectLock_deprecatedEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_objectLockEnabledNoDefaultRetentionDeprecatedEnabled(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.object_lock_enabled", s3.ObjectLockEnabledEnabled),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.rule.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_objectLock_migrate(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_objectLockEnabledNoDefaultRetentionDeprecatedEnabled(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.object_lock_enabled", s3.ObjectLockEnabledEnabled),
				),
			},
			{
				Config:   testAccBucketConfig_objectLockEnabledNoDefaultRetention(bucketName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccS3Bucket_Manage_objectLockWithVersioning(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_objectLockEnabledVersioning(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.object_lock_enabled", s3.ObjectLockEnabledEnabled),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_objectLockWithVersioning_deprecatedEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_objectLockEnabledVersioningDeprecatedEnabled(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_lock_configuration.0.object_lock_enabled", s3.ObjectLockEnabledEnabled),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_versioning(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_versioning(bucketName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.mfa_delete", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				Config: testAccBucketConfig_versioning(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.mfa_delete", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_versioningDisabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_versioning(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.mfa_delete", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_MFADeleteDisabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_versioningMFADelete(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.mfa_delete", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Manage_versioningAndMFADeleteDisabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_versioningDisabledAndMFADelete(bucketName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "versioning.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "versioning.0.mfa_delete", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Replication_basic(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	iamRoleResourceName := "aws_iam_role.role"
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replication(bucketName, s3.StorageClassStandard),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replication(bucketName, s3.StorageClassGlacier),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationSSEKMSEncryptedObjects(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Replication_multipleDestinationsEmptyFilter(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationMultipleDestinationsEmptyFilter(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination2", acctest.RegionProviderFunc(alternateRegion, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination3", acctest.RegionProviderFunc(alternateRegion, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule1",
						"priority":                    "1",
						"status":                      "Enabled",
						"filter.#":                    "1",
						"filter.0.prefix":             "",
						"destination.#":               "1",
						"destination.0.storage_class": "STANDARD",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule2",
						"priority":                    "2",
						"status":                      "Enabled",
						"filter.#":                    "1",
						"filter.0.prefix":             "",
						"destination.#":               "1",
						"destination.0.storage_class": "STANDARD_IA",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule3",
						"priority":                    "3",
						"status":                      "Disabled",
						"filter.#":                    "1",
						"filter.0.prefix":             "",
						"destination.#":               "1",
						"destination.0.storage_class": "ONEZONE_IA",
					}),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationMultipleDestinationsEmptyFilter(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Replication_multipleDestinationsNonEmptyFilter(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationMultipleDestinationsNonEmptyFilter(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination2", acctest.RegionProviderFunc(alternateRegion, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination3", acctest.RegionProviderFunc(alternateRegion, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule1",
						"priority":                    "1",
						"status":                      "Enabled",
						"filter.#":                    "1",
						"filter.0.prefix":             "prefix1",
						"destination.#":               "1",
						"destination.0.storage_class": "STANDARD",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule2",
						"priority":                    "2",
						"status":                      "Enabled",
						"filter.#":                    "1",
						"filter.0.tags.%":             "1",
						"filter.0.tags.Key2":          "Value2",
						"destination.#":               "1",
						"destination.0.storage_class": "STANDARD_IA",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule3",
						"priority":                    "3",
						"status":                      "Disabled",
						"filter.#":                    "1",
						"filter.0.prefix":             "prefix3",
						"filter.0.tags.%":             "1",
						"filter.0.tags.Key3":          "Value3",
						"destination.#":               "1",
						"destination.0.storage_class": "ONEZONE_IA",
					}),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationMultipleDestinationsNonEmptyFilter(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Replication_twoDestination(t *testing.T) {
	ctx := acctest.Context(t)

	// This tests 2 destinations since GovCloud and possibly other non-standard partitions allow a max of 2
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationMultipleDestinationsTwoDestination(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination2", acctest.RegionProviderFunc(alternateRegion, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule1",
						"priority":                    "1",
						"status":                      "Enabled",
						"filter.#":                    "1",
						"filter.0.prefix":             "prefix1",
						"destination.#":               "1",
						"destination.0.storage_class": "STANDARD",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_configuration.0.rules.*", map[string]string{
						"id":                          "rule2",
						"priority":                    "2",
						"status":                      "Enabled",
						"filter.#":                    "1",
						"filter.0.tags.%":             "1",
						"filter.0.tags.Key2":          "Value2",
						"destination.#":               "1",
						"destination.0.storage_class": "STANDARD_IA",
					}),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationMultipleDestinationsTwoDestination(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Replication_ruleDestinationAccessControlTranslation(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := acctest.Region()
	iamRoleResourceName := "aws_iam_role.role"
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationAccessControlTranslation(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationAccessControlTranslation(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl", "versioning"},
			},
			{
				Config: testAccBucketConfig_replicationSSEKMSEncryptedObjectsAndAccessControlTranslation(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
				),
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/12480
func TestAccS3Bucket_Replication_ruleDestinationAddAccessControlTranslation(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := acctest.Region()
	iamRoleResourceName := "aws_iam_role.role"
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationRulesDestination(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationAccessControlTranslation(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl", "versioning"},
			},
			{
				Config: testAccBucketConfig_replicationAccessControlTranslation(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
				),
			},
		},
	})
}

// StorageClass issue: https://github.com/hashicorp/terraform/issues/10909
func TestAccS3Bucket_Replication_withoutStorageClass(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationNoStorageClass(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationNoStorageClass(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Replication_expectVersioningValidationError(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketConfig_replicationNoVersioning(bucketName),
				ExpectError: regexache.MustCompile(`versioning must be enabled to allow S3 bucket replication`),
			},
		},
	})
}

// Prefix issue: https://github.com/hashicorp/terraform-provider-aws/issues/6340
func TestAccS3Bucket_Replication_withoutPrefix(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationNoPrefix(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationNoPrefix(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Replication_schemaV2(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	iamRoleResourceName := "aws_iam_role.role"
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationV2DeleteMarkerReplicationDisabled(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationV2NoTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config:                  testAccBucketConfig_replicationV2NoTags(bucketName),
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				Config: testAccBucketConfig_replicationV2OnlyOneTag(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationV2PrefixAndTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationV2MultipleTags(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
		},
	})
}

func TestAccS3Bucket_Replication_schemaV2SameRegion(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_s3_bucket.source"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	destinationResourceName := "aws_s3_bucket.destination"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationV2SameRegionNoTags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					acctest.CheckResourceAttrGlobalARN(resourceName, "replication_configuration.0.role", "iam", fmt.Sprintf("role/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExists(ctx, destinationResourceName),
				),
			},
			{
				Config:            testAccBucketConfig_replicationV2SameRegionNoTags(rName),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
					"acl",
				},
			},
		},
	})
}

func TestAccS3Bucket_Replication_RTC_valid(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	alternateRegion := acctest.AlternateRegion()
	region := acctest.Region()
	iamRoleResourceName := "aws_iam_role.role"
	resourceName := "aws_s3_bucket.source"

	// record the initialized providers so that we can use them to check for the instances in each region
	var providers []*schema.Provider

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesPlusProvidersAlternate(ctx, t, &providers),
		CheckDestroy:             acctest.CheckWithProviders(testAccCheckBucketDestroyWithProvider(ctx), &providers),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_replicationV2RTC(bucketName, 15),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationV2RTCNoMinutes(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationV2RTCNoStatus(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
			{
				Config: testAccBucketConfig_replicationV2RTCNotConfigured(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExistsWithProvider(ctx, resourceName, acctest.RegionProviderFunc(region, &providers)),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "replication_configuration.0.role", iamRoleResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.0.destination.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.0.destination.0.replication_time.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_configuration.0.rules.0.destination.0.metrics.#", "1"),
					testAccCheckBucketExistsWithProvider(ctx, "aws_s3_bucket.destination", acctest.RegionProviderFunc(alternateRegion, &providers)),
				),
			},
		},
	})
}

func TestAccS3Bucket_Security_corsUpdate(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	updateBucketCors := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("Not found: %s", n)
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)
			_, err := conn.PutBucketCorsWithContext(ctx, &s3.PutBucketCorsInput{
				Bucket: aws.String(rs.Primary.ID),
				CORSConfiguration: &s3.CORSConfiguration{
					CORSRules: []*s3.CORSRule{
						{
							AllowedHeaders: []*string{aws.String("*")},
							AllowedMethods: []*string{aws.String("GET")},
							AllowedOrigins: []*string{aws.String("https://www.example.com")},
						},
					},
				},
			})
			if err != nil && !tfawserr.ErrCodeEquals(err, tfs3.ErrCodeNoSuchCORSConfiguration) {
				return err
			}
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_cors(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.0", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.1", "POST"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_origins.0", "https://www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.0", "x-amz-server-side-encryption"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.1", "ETag"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.max_age_seconds", "3000"),
					updateBucketCors(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				Config: testAccBucketConfig_cors(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.0", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.1", "POST"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_origins.0", "https://www.example.com"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.0", "x-amz-server-side-encryption"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.1", "ETag"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.max_age_seconds", "3000"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Security_corsDelete(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	deleteBucketCors := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("Not found: %s", n)
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)
			_, err := conn.DeleteBucketCorsWithContext(ctx, &s3.DeleteBucketCorsInput{
				Bucket: aws.String(rs.Primary.ID),
			})
			if err != nil && !tfawserr.ErrCodeEquals(err, tfs3.ErrCodeNoSuchCORSConfiguration) {
				return err
			}
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_cors(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					deleteBucketCors(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3Bucket_Security_corsEmptyOrigin(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_corsEmptyOrigin(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.0", "PUT"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_methods.1", "POST"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.allowed_origins.0", ""),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.0", "x-amz-server-side-encryption"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.expose_headers.1", "ETag"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.0.max_age_seconds", "3000"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Security_corsSingleMethodAndEmptyOrigin(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_corsSingleMethodAndEmptyOrigin(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy"},
			},
		},
	})
}

func TestAccS3Bucket_Security_logging(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_logging(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "logging.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "logging.0.target_bucket", "aws_s3_bucket.log_bucket", "id"),
					resource.TestCheckResourceAttr(resourceName, "logging.0.target_prefix", "log/"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Security_enableDefaultEncryptionWhenTypical(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_defaultEncryptionKMSMasterKey(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.sse_algorithm", s3.ServerSideEncryptionAwsKms),
					resource.TestMatchResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.kms_master_key_id", regexache.MustCompile("^arn")),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Security_enableDefaultEncryptionWhenAES256IsUsed(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_defaultEncryptionDefaultKey(bucketName, s3.ServerSideEncryptionAes256),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.sse_algorithm", s3.ServerSideEncryptionAes256),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.0.rule.0.apply_server_side_encryption_by_default.0.kms_master_key_id", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
		},
	})
}

func TestAccS3Bucket_Security_disableDefaultEncryptionWhenDefaultEncryptionIsEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix("tf-test-bucket")
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_defaultEncryptionDefaultKey(bucketName, s3.ServerSideEncryptionAwsKms),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl"},
			},
			{
				// As ServerSide Encryption Configuration is a Computed field, removing them from terraform will not
				// trigger an update to remove it from the S3 bucket.
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "server_side_encryption_configuration.#", "1"),
				),
			},
		},
	})
}

func TestAccS3Bucket_Web_simple(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := acctest.Region()
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_website(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.index_document", "index.html"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl", "grant"},
			},
			{
				Config: testAccBucketConfig_websiteAndError(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr(resourceName, "website.0.error_document", "error.html"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
			{
				// As Website is a Computed field, removing them from terraform will not
				// trigger an update to remove them from the S3 bucket.
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr(resourceName, "website.0.error_document", "error.html"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
		},
	})
}

func TestAccS3Bucket_Web_redirect(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := acctest.Region()
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_websiteAndRedirect(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.redirect_all_requests_to", "hashicorp.com?my=query"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl", "grant"},
			},
			{
				Config: testAccBucketConfig_websiteAndHTTPSRedirect(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.redirect_all_requests_to", "https://hashicorp.com?my=query"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
			{
				// As Website is a Computed field, removing them from terraform will not
				// trigger an update to remove them from the S3 bucket.
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.redirect_all_requests_to", "https://hashicorp.com?my=query"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
		},
	})
}

func TestAccS3Bucket_Web_routingRules(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	region := acctest.Region()
	resourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_websiteAndRoutingRules(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.error_document", "error.html"),
					resource.TestCheckResourceAttr(resourceName, "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttrSet(resourceName, "website.0.routing_rules"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_destroy", "acl", "grant"},
			},
			{
				// As Website is a Computed field, removing them from terraform will not
				// trigger an update to remove them from the S3 bucket.
				Config: testAccBucketConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "website.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "website.0.error_document", "error.html"),
					resource.TestCheckResourceAttr(resourceName, "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttrSet(resourceName, "website.0.routing_rules"),
					testAccCheckBucketWebsiteEndpoint(resourceName, "website_endpoint", bucketName, region),
				),
			},
		},
	})
}

func TestBucketName(t *testing.T) {
	t.Parallel()

	validDnsNames := []string{
		"foobar",
		"foo.bar",
		"foo.bar.baz",
		"1234",
		"foo-bar",
		strings.Repeat("x", 63),
	}

	for _, v := range validDnsNames {
		if err := tfs3.ValidBucketName(v, endpoints.UsWest2RegionID); err != nil {
			t.Fatalf("%q should be a valid S3 bucket name", v)
		}
	}

	invalidDnsNames := []string{
		"foo..bar",
		"Foo.Bar",
		"192.168.0.1",
		"127.0.0.1",
		".foo",
		"bar.",
		"foo_bar",
		strings.Repeat("x", 64),
	}

	for _, v := range invalidDnsNames {
		if err := tfs3.ValidBucketName(v, endpoints.UsWest2RegionID); err == nil {
			t.Fatalf("%q should not be a valid S3 bucket name", v)
		}
	}

	validEastNames := []string{
		"foobar",
		"foo_bar",
		"127.0.0.1",
		"foo..bar",
		"foo_bar_baz",
		"foo.bar.baz",
		"Foo.Bar",
		strings.Repeat("x", 255),
	}

	for _, v := range validEastNames {
		if err := tfs3.ValidBucketName(v, endpoints.UsEast1RegionID); err != nil {
			t.Fatalf("%q should be a valid S3 bucket name", v)
		}
	}

	invalidEastNames := []string{
		"foo;bar",
		strings.Repeat("x", 256),
	}

	for _, v := range invalidEastNames {
		if err := tfs3.ValidBucketName(v, endpoints.UsEast1RegionID); err == nil {
			t.Fatalf("%q should not be a valid S3 bucket name", v)
		}
	}
}

func TestBucketRegionalDomainName(t *testing.T) {
	t.Parallel()

	const bucket = "bucket-name"

	var testCases = []struct {
		ExpectedErrCount int
		ExpectedOutput   string
		Region           string
	}{
		{
			Region:           "",
			ExpectedErrCount: 0,
			ExpectedOutput:   bucket + ".s3.amazonaws.com",
		},
		{
			Region:           "custom",
			ExpectedErrCount: 0,
			ExpectedOutput:   bucket + ".s3.custom.amazonaws.com",
		},
		{
			Region:           endpoints.UsEast1RegionID,
			ExpectedErrCount: 0,
			ExpectedOutput:   bucket + fmt.Sprintf(".s3.%s.%s", endpoints.UsEast1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			Region:           endpoints.UsWest2RegionID,
			ExpectedErrCount: 0,
			ExpectedOutput:   bucket + fmt.Sprintf(".s3.%s.%s", endpoints.UsWest2RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			Region:           endpoints.UsGovWest1RegionID,
			ExpectedErrCount: 0,
			ExpectedOutput:   bucket + fmt.Sprintf(".s3.%s.%s", endpoints.UsGovWest1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			Region:           endpoints.CnNorth1RegionID,
			ExpectedErrCount: 0,
			ExpectedOutput:   bucket + fmt.Sprintf(".s3.%s.amazonaws.com.cn", endpoints.CnNorth1RegionID),
		},
	}

	for _, tc := range testCases {
		output, err := tfs3.BucketRegionalDomainName(bucket, tc.Region)
		if tc.ExpectedErrCount == 0 && err != nil {
			t.Fatalf("expected %q not to trigger an error, received: %s", tc.Region, err)
		}
		if tc.ExpectedErrCount > 0 && err == nil {
			t.Fatalf("expected %q to trigger an error", tc.Region)
		}
		if output != tc.ExpectedOutput {
			t.Fatalf("expected %q, received %q", tc.ExpectedOutput, output)
		}
	}
}

func TestWebsiteEndpoint(t *testing.T) {
	t.Parallel()

	// https://docs.aws.amazon.com/AmazonS3/latest/dev/WebsiteEndpoints.html
	testCases := []struct {
		TestingClient      *conns.AWSClient
		LocationConstraint string
		Expected           string
	}{
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.UsEast1RegionID,
			},
			LocationConstraint: "",
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.UsEast1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.UsWest2RegionID,
			},
			LocationConstraint: endpoints.UsWest2RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.UsWest2RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.UsWest1RegionID,
			},
			LocationConstraint: endpoints.UsWest1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.UsWest1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.EuWest1RegionID,
			},
			LocationConstraint: endpoints.EuWest1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.EuWest1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.EuWest3RegionID,
			},
			LocationConstraint: endpoints.EuWest3RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.%s", endpoints.EuWest3RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.EuCentral1RegionID,
			},
			LocationConstraint: endpoints.EuCentral1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.%s", endpoints.EuCentral1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.ApSouth1RegionID,
			},
			LocationConstraint: endpoints.ApSouth1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.%s", endpoints.ApSouth1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.ApSoutheast1RegionID,
			},
			LocationConstraint: endpoints.ApSoutheast1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.ApSoutheast1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.ApNortheast1RegionID,
			},
			LocationConstraint: endpoints.ApNortheast1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.ApNortheast1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.ApSoutheast2RegionID,
			},
			LocationConstraint: endpoints.ApSoutheast2RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.ApSoutheast2RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.ApNortheast2RegionID,
			},
			LocationConstraint: endpoints.ApNortheast2RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.%s", endpoints.ApNortheast2RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.SaEast1RegionID,
			},
			LocationConstraint: endpoints.SaEast1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.SaEast1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.UsGovEast1RegionID,
			},
			LocationConstraint: endpoints.UsGovEast1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.%s", endpoints.UsGovEast1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com",
				Region:    endpoints.UsGovWest1RegionID,
			},
			LocationConstraint: endpoints.UsGovWest1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website-%s.%s", endpoints.UsGovWest1RegionID, acctest.PartitionDNSSuffix()),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "c2s.ic.gov",
				Region:    endpoints.UsIsoEast1RegionID,
			},
			LocationConstraint: endpoints.UsIsoEast1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.c2s.ic.gov", endpoints.UsIsoEast1RegionID),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "sc2s.sgov.gov",
				Region:    endpoints.UsIsobEast1RegionID,
			},
			LocationConstraint: endpoints.UsIsobEast1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.sc2s.sgov.gov", endpoints.UsIsobEast1RegionID),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com.cn",
				Region:    endpoints.CnNorthwest1RegionID,
			},
			LocationConstraint: endpoints.CnNorthwest1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.amazonaws.com.cn", endpoints.CnNorthwest1RegionID),
		},
		{
			TestingClient: &conns.AWSClient{
				DNSSuffix: "amazonaws.com.cn",
				Region:    endpoints.CnNorth1RegionID,
			},
			LocationConstraint: endpoints.CnNorth1RegionID,
			Expected:           fmt.Sprintf("bucket-name.s3-website.%s.amazonaws.com.cn", endpoints.CnNorth1RegionID),
		},
	}

	for _, testCase := range testCases {
		got := tfs3.WebsiteEndpoint(testCase.TestingClient, "bucket-name", testCase.LocationConstraint)
		if got.Endpoint != testCase.Expected {
			t.Errorf("WebsiteEndpointUrl(\"bucket-name\", %q) => %q, want %q", testCase.LocationConstraint, got.Endpoint, testCase.Expected)
		}
	}
}

func testAccCheckBucketDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error { return testAccCheckBucketDestroyWithProvider(ctx)(s, acctest.Provider) }
}

func testAccCheckBucketDestroyWithProvider(ctx context.Context) acctest.TestCheckWithProviderFunc {
	return func(s *terraform.State, provider *schema.Provider) error {
		conn := provider.Meta().(*conns.AWSClient).S3Client(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_s3_bucket" {
				continue
			}

			err := tfs3.FindBucket(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("S3 Bucket %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckBucketExists(ctx context.Context, n string) resource.TestCheckFunc {
	return testAccCheckBucketExistsWithProvider(ctx, n, func() *schema.Provider { return acctest.Provider })
}

func testAccCheckBucketExistsWithProvider(ctx context.Context, n string, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket ID is set")
		}

		conn := providerF().Meta().(*conns.AWSClient).S3Client(ctx)

		return tfs3.FindBucket(ctx, conn, rs.Primary.ID)
	}
}

func testAccCheckBucketAddObjects(ctx context.Context, n string, keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Client(ctx)

		for _, key := range keys {
			_, err := conn.PutObject(ctx, &s3_sdkv2.PutObjectInput{
				Bucket: aws_sdkv2.String(rs.Primary.ID),
				Key:    aws_sdkv2.String(key),
			})

			if err != nil {
				return fmt.Errorf("PutObject error: %s", err)
			}
		}

		return nil
	}
}

func testAccCheckBucketAddObjectsWithLegalHold(ctx context.Context, n string, keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)

		for _, key := range keys {
			_, err := conn.PutObjectWithContext(ctx, &s3.PutObjectInput{
				Bucket:                    aws.String(rs.Primary.ID),
				Key:                       aws.String(key),
				ObjectLockLegalHoldStatus: aws.String(s3.ObjectLockLegalHoldStatusOn),
			})

			if err != nil {
				return fmt.Errorf("PutObject error: %s", err)
			}
		}

		return nil
	}
}

func testAccCheckBucketAddObjectWithMetadata(ctx context.Context, n string, key string, metadata map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)

		_, err := conn.PutObjectWithContext(ctx, &s3.PutObjectInput{
			Bucket:   aws.String(rs.Primary.ID),
			Key:      aws.String(key),
			Metadata: aws.StringMap(metadata),
		})

		if err != nil {
			return fmt.Errorf("PutObject error: %s", err)
		}

		return nil
	}
}

// Create an S3 bucket via a CF stack so that it has system tags.
func testAccCheckBucketCreateViaCloudFormation(ctx context.Context, n string, stackID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).CloudFormationConn(ctx)
		stackName := sdkacctest.RandomWithPrefix("tf-acc-test-s3tags")
		templateBody := fmt.Sprintf(`{
  "Resources": {
    "TfTestBucket": {
      "Type": "AWS::S3::Bucket",
      "Properties": {
        "BucketName": "%s"
      }
    }
  }
}`, n)

		requestToken := id.UniqueId()
		req := &cloudformation.CreateStackInput{
			StackName:          aws.String(stackName),
			TemplateBody:       aws.String(templateBody),
			ClientRequestToken: aws.String(requestToken),
		}

		log.Printf("[DEBUG] Creating CloudFormation stack: %s", req)
		resp, err := conn.CreateStackWithContext(ctx, req)
		if err != nil {
			return fmt.Errorf("error creating CloudFormation stack: %w", err)
		}

		stack, err := tfcloudformation.WaitStackCreated(ctx, conn, aws.StringValue(resp.StackId), requestToken, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("Error waiting for CloudFormation stack creation: %w", err)
		}
		status := aws.StringValue(stack.StackStatus)
		if status != cloudformation.StackStatusCreateComplete {
			return fmt.Errorf("Invalid CloudFormation stack creation status: %s", status)
		}

		*stackID = aws.StringValue(resp.StackId)
		return nil
	}
}

func testAccCheckBucketTagKeys(ctx context.Context, n string, keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)

		got, err := tfs3.BucketListTags(ctx, conn, rs.Primary.Attributes["bucket"])
		if err != nil {
			return err
		}

		for _, want := range keys {
			ok := false
			for _, key := range got.Keys() {
				if want == key {
					ok = true
					break
				}
			}
			if !ok {
				return fmt.Errorf("Key %s not found in bucket's tag set", want)
			}
		}

		return nil
	}
}

func testAccCheckBucketDomainName(resourceName string, attributeName string, bucketName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expectedValue := acctest.Provider.Meta().(*conns.AWSClient).PartitionHostname(fmt.Sprintf("%s.s3", bucketName))

		return resource.TestCheckResourceAttr(resourceName, attributeName, expectedValue)(s)
	}
}

func testAccBucketRegionalDomainName(bucket, region string) string {
	regionalEndpoint, err := tfs3.BucketRegionalDomainName(bucket, region)
	if err != nil {
		return fmt.Sprintf("Regional endpoint not found for bucket %s", bucket)
	}
	return regionalEndpoint
}

func testAccCheckBucketWebsiteEndpoint(resourceName string, attributeName string, bucketName string, region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		website := tfs3.WebsiteEndpoint(acctest.Provider.Meta().(*conns.AWSClient), bucketName, region)
		expectedValue := website.Endpoint

		return resource.TestCheckResourceAttr(resourceName, attributeName, expectedValue)(s)
	}
}

func testAccCheckBucketUpdateTags(ctx context.Context, n string, oldTags, newTags map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)

		return tfs3.BucketUpdateTags(ctx, conn, rs.Primary.Attributes["bucket"], oldTags, newTags)
	}
}

func testAccCheckBucketCheckTags(ctx context.Context, n string, expectedTags map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn(ctx)

		got, err := tfs3.BucketListTags(ctx, conn, rs.Primary.Attributes["bucket"])
		if err != nil {
			return err
		}

		want := tftags.New(ctx, expectedTags)
		if !reflect.DeepEqual(want, got) {
			return fmt.Errorf("Incorrect tags, want: %v got: %v", want, got)
		}

		return nil
	}
}

func testAccBucketConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}
`, bucketName)
}

func testAccBucketConfig_acceleration(bucketName, acceleration string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket              = %[1]q
  acceleration_status = %[2]q
}
`, bucketName, acceleration)
}

func testAccBucketConfig_acl(bucketName, acl string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  acl    = %[2]q
}
`, bucketName, acl)
}

func testAccBucketConfig_cors(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["x-amz-server-side-encryption", "ETag"]
    max_age_seconds = 3000
  }
}
`, bucketName)
}

func testAccBucketConfig_corsSingleMethodAndEmptyOrigin(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = [""]
  }
}
`, bucketName)
}

func testAccBucketConfig_corsEmptyOrigin(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = [""]
    expose_headers  = ["x-amz-server-side-encryption", "ETag"]
    max_age_seconds = 3000
  }
}
`, bucketName)
}

func testAccBucketConfig_defaultEncryptionDefaultKey(bucketName, sseAlgorithm string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = %[2]q
      }
    }
  }
}
`, bucketName, sseAlgorithm)
}

func testAccBucketConfig_defaultEncryptionKMSMasterKey(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = "KMS Key for Bucket %[1]s"
  deletion_window_in_days = 10
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.test.arn
        sse_algorithm     = "aws:kms"
      }
    }
  }
}
`, bucketName)
}

func testAccBucketConfig_defaultEncryptionKeyEnabledKMSMasterKey(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = "KMS Key for Bucket %[1]s"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        kms_master_key_id = aws_kms_key.test.arn
        sse_algorithm     = "aws:kms"
      }
      bucket_key_enabled = true
    }
  }
}
`, bucketName)
}

func testAccBucketConfig_lifecycle(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }

    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }

    transition {
      days          = 60
      storage_class = "INTELLIGENT_TIERING"
    }

    transition {
      days          = 90
      storage_class = "ONEZONE_IA"
    }

    transition {
      days          = 120
      storage_class = "GLACIER"
    }

    transition {
      days          = 210
      storage_class = "DEEP_ARCHIVE"
    }
  }

  lifecycle_rule {
    id      = "id2"
    prefix  = "path2/"
    enabled = true

    expiration {
      date = "2016-01-12"
    }
  }

  lifecycle_rule {
    id      = "id3"
    prefix  = "path3/"
    enabled = true

    transition {
      days          = 0
      storage_class = "GLACIER"
    }
  }

  lifecycle_rule {
    id      = "id4"
    prefix  = "path4/"
    enabled = true

    tags = {
      "tagKey"    = "tagValue"
      "terraform" = "hashicorp"
    }

    expiration {
      date = "2016-01-12"
    }
  }

  lifecycle_rule {
    id      = "id5"
    enabled = true

    tags = {
      "tagKey"    = "tagValue"
      "terraform" = "hashicorp"
    }

    transition {
      days          = 0
      storage_class = "GLACIER"
    }
  }

  lifecycle_rule {
    id      = "id6"
    enabled = true

    tags = {
      "tagKey" = "tagValue"
    }

    transition {
      days          = 0
      storage_class = "GLACIER"
    }
  }
}
`, bucketName)
}

func testAccBucketConfig_lifecycleExpireMarker(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    expiration {
      expired_object_delete_marker = "true"
    }
  }
}
`, bucketName)
}

func testAccBucketConfig_lifecycleRuleExpirationEmptyBlock(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  lifecycle_rule {
    enabled = true
    id      = "id1"

    expiration {}
  }
}
`, rName)
}

func testAccBucketConfig_lifecycleRuleAbortIncompleteMultipartUploadDays(rName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  lifecycle_rule {
    abort_incomplete_multipart_upload_days = 7
    enabled                                = true
    id                                     = "id1"
  }
}
`, rName)
}

func testAccBucketConfig_logging(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "log_bucket" {
  bucket = "%[1]s-log"
}

resource "aws_s3_bucket_ownership_controls" "log_bucket_ownership" {
  bucket = aws_s3_bucket.log_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "log_bucket_acl" {
  depends_on = [aws_s3_bucket_ownership_controls.log_bucket_ownership]

  bucket = aws_s3_bucket.log_bucket.id
  acl    = "log-delivery-write"
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  logging {
    target_bucket = aws_s3_bucket.log_bucket.id
    target_prefix = "log/"
  }
}
`, bucketName)
}

func testAccBucketConfig_policy(bucketName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}
data "aws_caller_identity" "current" {}

data "aws_iam_policy_document" "policy" {
  statement {
    effect = "Allow"

    actions = [
      "s3:*",
    ]

    resources = [
      "arn:${data.aws_partition.current.partition}:s3:::%[1]s",
      "arn:${data.aws_partition.current.partition}:s3:::%[1]s/*",
    ]

    principals {
      type        = "AWS"
      identifiers = ["arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"]
    }
  }
}

resource "aws_s3_bucket" "test" {
  bucket = %[1]q
  policy = data.aws_iam_policy_document.policy.json
}
`, bucketName)
}

func testAccBucketConfig_ReplicationBase(bucketName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigMultipleRegionProvider(2),
		fmt.Sprintf(`
data "aws_partition" "current" {}

resource "aws_iam_role" "role" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "s3.${data.aws_partition.current.dns_suffix}"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

resource "aws_s3_bucket" "destination" {
  provider = "awsalternate"
  bucket   = "%[1]s-destination"

  versioning {
    enabled = true
  }
}
`, bucketName))
}

func testAccBucketConfig_replication(bucketName, storageClass string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = %[2]q
      }
    }
  }
}
`, bucketName, storageClass))
}

func testAccBucketConfig_replicationAccessControlTranslation(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        account_id    = data.aws_caller_identity.current.account_id
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"

        access_control_translation {
          owner = "Destination"
        }
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationMultipleDestinationsEmptyFilter(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "destination2" {
  provider = "awsalternate"
  bucket   = "%[1]s-destination2"

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "destination3" {
  provider = "awsalternate"
  bucket   = "%[1]s-destination3"

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id       = "rule1"
      priority = 1
      status   = "Enabled"

      filter {}

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }

    rules {
      id       = "rule2"
      priority = 2
      status   = "Enabled"

      filter {}

      destination {
        bucket        = aws_s3_bucket.destination2.arn
        storage_class = "STANDARD_IA"
      }
    }

    rules {
      id       = "rule3"
      priority = 3
      status   = "Disabled"

      filter {}

      destination {
        bucket        = aws_s3_bucket.destination3.arn
        storage_class = "ONEZONE_IA"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationMultipleDestinationsNonEmptyFilter(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "destination2" {
  provider = "awsalternate"
  bucket   = "%[1]s-destination2"

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "destination3" {
  provider = "awsalternate"
  bucket   = "%[1]s-destination3"

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id       = "rule1"
      priority = 1
      status   = "Enabled"

      filter {
        prefix = "prefix1"
      }

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }

    rules {
      id       = "rule2"
      priority = 2
      status   = "Enabled"

      filter {
        tags = {
          Key2 = "Value2"
        }
      }

      destination {
        bucket        = aws_s3_bucket.destination2.arn
        storage_class = "STANDARD_IA"
      }
    }

    rules {
      id       = "rule3"
      priority = 3
      status   = "Disabled"

      filter {
        prefix = "prefix3"

        tags = {
          Key3 = "Value3"
        }
      }

      destination {
        bucket        = aws_s3_bucket.destination3.arn
        storage_class = "ONEZONE_IA"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationMultipleDestinationsTwoDestination(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "destination2" {
  provider = "awsalternate"
  bucket   = "%[1]s-destination2"

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id       = "rule1"
      priority = 1
      status   = "Enabled"

      filter {
        prefix = "prefix1"
      }

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }

    rules {
      id       = "rule2"
      priority = 2
      status   = "Enabled"

      filter {
        tags = {
          Key2 = "Value2"
        }
      }

      destination {
        bucket        = aws_s3_bucket.destination2.arn
        storage_class = "STANDARD_IA"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationNoVersioning(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s"

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationRulesDestination(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        account_id    = data.aws_caller_identity.current.account_id
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }

  versioning {
    enabled = true
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationSSEKMSEncryptedObjects(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_kms_key" "replica" {
  provider                = "awsalternate"
  description             = "TF Acceptance Test S3 repl KMS key"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        bucket             = aws_s3_bucket.destination.arn
        storage_class      = "STANDARD"
        replica_kms_key_id = aws_kms_key.replica.arn
      }

      source_selection_criteria {
        sse_kms_encrypted_objects {
          enabled = true
        }
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationSSEKMSEncryptedObjectsAndAccessControlTranslation(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
data "aws_caller_identity" "current" {}

resource "aws_kms_key" "replica" {
  provider                = "awsalternate"
  description             = "TF Acceptance Test S3 repl KMS key"
  deletion_window_in_days = 7
}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        account_id         = data.aws_caller_identity.current.account_id
        bucket             = aws_s3_bucket.destination.arn
        storage_class      = "STANDARD"
        replica_kms_key_id = aws_kms_key.replica.arn

        access_control_translation {
          owner = "Destination"
        }
      }

      source_selection_criteria {
        sse_kms_encrypted_objects {
          enabled = true
        }
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationNoPrefix(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      status = "Enabled"

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationNoStorageClass(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      prefix = "foo"
      status = "Enabled"

      destination {
        bucket = aws_s3_bucket.destination.arn
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2SameRegionNoTags(rName string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "s3.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.test.arn

    rules {
      id     = "testid"
      status = "Enabled"

      filter {
        prefix = "testprefix"
      }

      delete_marker_replication_status = "Enabled"

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}

resource "aws_s3_bucket" "destination" {
  bucket = "%[1]s-destination"

  versioning {
    enabled = true
  }
}
`, rName)
}

func testAccBucketConfig_replicationV2DeleteMarkerReplicationDisabled(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      status = "Enabled"

      filter {
        prefix = "foo"
      }

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2NoTags(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      status = "Enabled"

      filter {
        prefix = "foo"
      }

      delete_marker_replication_status = "Enabled"

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2OnlyOneTag(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      status = "Enabled"

      priority = 42

      filter {
        tags = {
          ReplicateMe = "Yes"
        }
      }

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2PrefixAndTags(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      status = "Enabled"

      priority = 41

      filter {
        prefix = "foo"

        tags = {
          AnotherTag  = "OK"
          ReplicateMe = "Yes"
        }
      }

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2MultipleTags(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn

    rules {
      id     = "foobar"
      status = "Enabled"

      filter {
        tags = {
          AnotherTag  = "OK"
          Foo         = "Bar"
          ReplicateMe = "Yes"
        }
      }

      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2RTC(bucketName string, minutes int) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn
    rules {
      id     = "rtc"
      status = "Enabled"
      filter {
        tags = {}
      }
      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
        metrics {
          status  = "Enabled"
          minutes = %[2]d
        }
        replication_time {
          status  = "Enabled"
          minutes = %[2]d
        }
      }
    }
  }
}
`, bucketName, minutes))
}

func testAccBucketConfig_replicationV2RTCNoMinutes(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn
    rules {
      id     = "rtc-no-minutes"
      status = "Enabled"
      filter {
        tags = {}
      }
      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
        metrics {}
        replication_time {
          status = "Enabled"
        }
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2RTCNoStatus(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"

  versioning {
    enabled = true
  }

  replication_configuration {
    role = aws_iam_role.role.arn
    rules {
      id     = "rtc-no-status"
      status = "Enabled"
      filter {
        prefix = "foo"
      }
      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
        metrics {}
        replication_time {
          minutes = 15
        }
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_replicationV2RTCNotConfigured(bucketName string) string {
	return acctest.ConfigCompose(
		testAccBucketConfig_ReplicationBase(bucketName),
		fmt.Sprintf(`
resource "aws_s3_bucket" "source" {
  bucket = "%[1]s-source"
  versioning {
    enabled = true
  }
  replication_configuration {
    role = aws_iam_role.role.arn
    rules {
      id     = "rtc-no-config"
      status = "Enabled"
      filter {
        prefix = "foo"
      }
      destination {
        bucket        = aws_s3_bucket.destination.arn
        storage_class = "STANDARD"
        metrics {}
        replication_time {}
      }
    }
  }
}
`, bucketName))
}

func testAccBucketConfig_requestPayer(bucketName, requestPayer string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  request_payer = %[2]q
}
`, bucketName, requestPayer)
}

func testAccBucketConfig_versioning(bucketName string, enabled bool) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  versioning {
    enabled = %[2]t
  }
}
`, bucketName, enabled)
}

func testAccBucketConfig_versioningMFADelete(bucketName string, mfaDelete bool) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  versioning {
    mfa_delete = %[2]t
  }
}
`, bucketName, mfaDelete)
}

func testAccBucketConfig_versioningDisabledAndMFADelete(bucketName string, mfaDelete bool) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  versioning {
    enabled    = false
    mfa_delete = %[2]t
  }
}
`, bucketName, mfaDelete)
}

func testAccBucketConfig_website(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  website {
    index_document = "index.html"
  }
}
`, bucketName)
}

func testAccBucketConfig_websiteAndError(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
`, bucketName)
}

func testAccBucketConfig_websiteAndRedirect(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  website {
    redirect_all_requests_to = "hashicorp.com?my=query"
  }
}
`, bucketName)
}

func testAccBucketConfig_websiteAndHTTPSRedirect(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  website {
    redirect_all_requests_to = "https://hashicorp.com?my=query"
  }
}
`, bucketName)
}

func testAccBucketConfig_websiteAndRoutingRules(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  website {
    index_document = "index.html"
    error_document = "error.html"

    routing_rules = <<EOF
[
  {
    "Condition": {
      "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
      "ReplaceKeyPrefixWith": "documents/"
    }
  }
]
EOF

  }
}
`, bucketName)
}

func testAccBucketConfig_noTags(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = false
}
`, bucketName)
}

func testAccBucketConfig_tags(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = false

  tags = {
    Key1 = "AAA"
    Key2 = "BBB"
    Key3 = "CCC"
  }
}
`, bucketName)
}

func testAccBucketConfig_updatedTags(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = false

  tags = {
    Key2 = "BBB"
    Key3 = "XXX"
    Key4 = "DDD"
    Key5 = "EEE"
  }
}
`, bucketName)
}

func testAccBucketConfig_multiTags(randInt int) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "bucket1" {
  bucket        = "tf-test-bucket-1-%[1]d"
  force_destroy = true

  tags = {
    Name        = "tf-test-bucket-1-%[1]d"
    Environment = "%[1]d"
  }
}

resource "aws_s3_bucket" "bucket2" {
  bucket        = "tf-test-bucket-2-%[1]d"
  force_destroy = true

  tags = {
    Name        = "tf-test-bucket-2-%[1]d"
    Environment = "%[1]d"
  }
}

resource "aws_s3_bucket" "bucket3" {
  bucket        = "tf-test-bucket-3-%[1]d"
  force_destroy = true

  tags = {
    Name        = "tf-test-bucket-3-%[1]d"
    Environment = "%[1]d"
  }
}

resource "aws_s3_bucket" "bucket4" {
  bucket        = "tf-test-bucket-4-%[1]d"
  force_destroy = true

  tags = {
    Name        = "tf-test-bucket-4-%[1]d"
    Environment = "%[1]d"
  }
}

resource "aws_s3_bucket" "bucket5" {
  bucket        = "tf-test-bucket-5-%[1]d"
  force_destroy = true

  tags = {
    Name        = "tf-test-bucket-5-%[1]d"
    Environment = "%[1]d"
  }
}

resource "aws_s3_bucket" "bucket6" {
  bucket        = "tf-test-bucket-6-%[1]d"
  force_destroy = true

  tags = {
    Name        = "tf-test-bucket-6-%[1]d"
    Environment = "%[1]d"
  }
}
`, randInt)
}

func testAccBucketConfig_objectLockEnabledNoDefaultRetention(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_enabled = true
}
`, bucketName)
}

func testAccBucketConfig_objectLockEnabledDefaultRetention(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_configuration {
    object_lock_enabled = "Enabled"

    rule {
      default_retention {
        mode = "COMPLIANCE"
        days = 3
      }
    }
  }
}
`, bucketName)
}

func testAccBucketConfig_objectLockEnabledNoDefaultRetentionDeprecatedEnabled(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  object_lock_configuration {
    object_lock_enabled = "Enabled"
  }
}
`, bucketName)
}

func testAccBucketConfig_objectLockEnabledVersioning(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true

  object_lock_enabled = true
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}
`, bucketName)
}

func testAccBucketConfig_objectLockEnabledVersioningDeprecatedEnabled(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = %[1]q
  force_destroy = true

  object_lock_configuration {
    object_lock_enabled = "Enabled"
  }
}

resource "aws_s3_bucket_versioning" "test" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}
`, bucketName)
}

func testAccBucketConfig_forceDestroy(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = "%s"
  force_destroy = true
}
`, bucketName)
}

func testAccBucketConfig_forceDestroyObjectLockEnabled(bucketName string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket        = "%s"
  force_destroy = true

  object_lock_enabled = true
}

resource "aws_s3_bucket_versioning" "bucket" {
  bucket = aws_s3_bucket.test.id
  versioning_configuration {
    status = "Enabled"
  }
}
`, bucketName)
}

const testAccBucketConfig_emptyString = `
resource "aws_s3_bucket" "test" {
  bucket = ""
}
`

func testAccBucketConfig_namePrefix(namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket_prefix = %[1]q
}
`, namePrefix)
}

const testAccBucketConfig_nameGenerated = `
resource "aws_s3_bucket" "test" {}
`

func testAccBucketConfig_duplicate(region, bucketName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigRegionalProvider(region),
		fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  depends_on = [aws_s3_bucket.duplicate]
}

resource "aws_s3_bucket" "duplicate" {
  bucket = %[1]q
}
  `, bucketName),
	)
}

func testAccBucketConfig_duplicateAltAccount(region, bucketName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigRegionalProvider(region),
		acctest.ConfigAlternateAccountProvider(),
		fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q

  depends_on = [aws_s3_bucket.duplicate]
}

resource "aws_s3_bucket" "duplicate" {
  provider = "awsalternate"
  bucket   = %[1]q
}
  `, bucketName),
	)
}
