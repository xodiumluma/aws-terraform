package s3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfs3 "github.com/hashicorp/terraform-provider-aws/internal/service/s3"
)

func TestAccS3BucketAccelerateConfiguration_basic(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_accelerate_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccelerateConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketAccelerateConfigurationConfig_basic(bucketName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketAccelerateConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", s3.BucketAccelerateStatusEnabled),
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

func TestAccS3BucketAccelerateConfiguration_update(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_accelerate_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccelerateConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketAccelerateConfigurationConfig_basic(bucketName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketAccelerateConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", s3.BucketAccelerateStatusEnabled),
				),
			},
			{
				Config: testAccBucketAccelerateConfigurationConfig_basic(bucketName, s3.BucketAccelerateStatusSuspended),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketAccelerateConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "aws_s3_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", s3.BucketAccelerateStatusSuspended),
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

func TestAccS3BucketAccelerateConfiguration_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	bucketName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_accelerate_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, cloudfront.EndpointsID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccelerateConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketAccelerateConfigurationConfig_basic(bucketName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketAccelerateConfigurationExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfs3.ResourceBucketAccelerateConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3BucketAccelerateConfiguration_migrate_noChange(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_accelerate_configuration.test"
	bucketResourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccelerateConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_acceleration(rName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketExists(ctx, bucketResourceName),
					resource.TestCheckResourceAttr(bucketResourceName, "acceleration_status", s3.BucketAccelerateStatusEnabled),
				),
			},
			{
				Config: testAccBucketAccelerateConfigurationConfig_basic(rName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketAccelerateConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status", s3.BucketAccelerateStatusEnabled),
				),
			},
		},
	})
}

func TestAccS3BucketAccelerateConfiguration_migrate_withChange(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_accelerate_configuration.test"
	bucketResourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketAccelerateConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketConfig_acceleration(rName, s3.BucketAccelerateStatusEnabled),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketExists(ctx, bucketResourceName),
					resource.TestCheckResourceAttr(bucketResourceName, "acceleration_status", s3.BucketAccelerateStatusEnabled),
				),
			},
			{
				Config: testAccBucketAccelerateConfigurationConfig_basic(rName, s3.BucketAccelerateStatusSuspended),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckBucketAccelerateConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", bucketResourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "status", s3.BucketAccelerateStatusSuspended),
				),
			},
		},
	})
}

func testAccCheckBucketAccelerateConfigurationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_s3_bucket_accelerate_configuration" {
				continue
			}

			bucket, expectedBucketOwner, err := tfs3.ParseResourceID(rs.Primary.ID)
			if err != nil {
				return err
			}

			input := &s3.GetBucketAccelerateConfigurationInput{
				Bucket: aws.String(bucket),
			}

			if expectedBucketOwner != "" {
				input.ExpectedBucketOwner = aws.String(expectedBucketOwner)
			}

			output, err := conn.GetBucketAccelerateConfigurationWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, s3.ErrCodeNoSuchBucket) {
				continue
			}

			if err != nil {
				return fmt.Errorf("error getting S3 Bucket accelerate configuration (%s): %w", rs.Primary.ID, err)
			}

			if output != nil {
				return fmt.Errorf("S3 Bucket accelerate configuration (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckBucketAccelerateConfigurationExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Conn()

		bucket, expectedBucketOwner, err := tfs3.ParseResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		input := &s3.GetBucketAccelerateConfigurationInput{
			Bucket: aws.String(bucket),
		}

		if expectedBucketOwner != "" {
			input.ExpectedBucketOwner = aws.String(expectedBucketOwner)
		}

		output, err := conn.GetBucketAccelerateConfigurationWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("error getting S3 Bucket accelerate configuration (%s): %w", rs.Primary.ID, err)
		}

		if output == nil {
			return fmt.Errorf("S3 Bucket accelerate configuration (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccBucketAccelerateConfigurationConfig_basic(bucketName, status string) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_accelerate_configuration" "test" {
  bucket = aws_s3_bucket.test.id
  status = %[2]q
}
`, bucketName, status)
}
