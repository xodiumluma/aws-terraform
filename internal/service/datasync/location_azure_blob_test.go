// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasync_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/datasync"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdatasync "github.com/hashicorp/terraform-provider-aws/internal/service/datasync"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccDataSyncLocationAzureBlob_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationAzureBlobOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datasync_location_azure_blob.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationAzureBlobDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationAzureBlobConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "access_tier", "HOT"),
					resource.TestCheckResourceAttr(resourceName, "agent_arns.#", "1"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "datasync", regexache.MustCompile(`location/loc-.+`)),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "blob_type", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "container_url", "https://example.com/path"),
					resource.TestCheckResourceAttr(resourceName, "sas_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "sas_configuration.0.token"),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/path/"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "uri", regexache.MustCompile(`^azure-blob://.+/`)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"container_url", "sas_configuration"},
			},
		},
	})
}

func TestAccDataSyncLocationAzureBlob_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationAzureBlobOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datasync_location_azure_blob.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationAzureBlobDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationAzureBlobConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdatasync.ResourceLocationAzureBlob(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSyncLocationAzureBlob_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationAzureBlobOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datasync_location_azure_blob.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationAzureBlobDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationAzureBlobConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"container_url", "sas_configuration"},
			},
			{
				Config: testAccLocationAzureBlobConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLocationAzureBlobConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func TestAccDataSyncLocationAzureBlob_update(t *testing.T) {
	ctx := acctest.Context(t)
	var v datasync.DescribeLocationAzureBlobOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_datasync_location_azure_blob.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationAzureBlobDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationAzureBlobConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "access_tier", "HOT"),
					resource.TestCheckResourceAttr(resourceName, "agent_arns.#", "1"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "datasync", regexache.MustCompile(`location/loc-.+`)),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "blob_type", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "container_url", "https://example.com/path"),
					resource.TestCheckResourceAttr(resourceName, "sas_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "sas_configuration.0.token"),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/path/"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "uri", regexache.MustCompile(`^azure-blob://.+/`)),
				),
			},
			{
				Config: testAccLocationAzureBlobConfig_updated(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckLocationAzureBlobExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "access_tier", "COOL"),
					resource.TestCheckResourceAttr(resourceName, "agent_arns.#", "1"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "datasync", regexache.MustCompile(`location/loc-.+`)),
					resource.TestCheckResourceAttr(resourceName, "authentication_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "blob_type", "BLOCK"),
					resource.TestCheckResourceAttr(resourceName, "container_url", "https://example.com/path"),
					resource.TestCheckResourceAttr(resourceName, "sas_configuration.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "sas_configuration.0.token"),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/path/"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "uri", regexache.MustCompile(`^azure-blob://.+/`)),
				),
			},
		},
	})
}

func testAccCheckLocationAzureBlobDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DataSyncConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_datasync_location_azure_blob" {
				continue
			}

			_, err := tfdatasync.FindLocationAzureBlobByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DataSync Location Microsoft Azure Blob Storage %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckLocationAzureBlobExists(ctx context.Context, n string, v *datasync.DescribeLocationAzureBlobOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DataSyncConn(ctx)

		output, err := tfdatasync.FindLocationAzureBlobByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccLocationAzureBlobConfig_base(rName string) string {
	return acctest.ConfigCompose(testAccAgentAgentConfig_base(rName), fmt.Sprintf(`
resource "aws_datasync_agent" "test" {
  ip_address = aws_instance.test.public_ip
  name       = %[1]q
}
`, rName))
}

func testAccLocationAzureBlobConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccLocationAzureBlobConfig_base(rName), `
resource "aws_datasync_location_azure_blob" "test" {
  agent_arns          = [aws_datasync_agent.test.arn]
  authentication_type = "SAS"
  container_url       = "https://example.com/path"

  sas_configuration {
    token = "sp=r&st=2023-12-20T14:54:52Z&se=2023-12-20T22:54:52Z&spr=https&sv=2021-06-08&sr=c&sig=aBBKDWQvyuVcTPH9EBp%%2FXTI9E%%2F%%2Fmq171%%2BZU178wcwqU%%3D"
  }
}
`)
}

func testAccLocationAzureBlobConfig_tags1(rName, key1, value1 string) string {
	return acctest.ConfigCompose(testAccLocationAzureBlobConfig_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_azure_blob" "test" {
  agent_arns          = [aws_datasync_agent.test.arn]
  authentication_type = "SAS"
  container_url       = "https://example.com/path"

  sas_configuration {
    token = "sp=r&st=2023-12-20T14:54:52Z&se=2023-12-20T22:54:52Z&spr=https&sv=2021-06-08&sr=c&sig=aBBKDWQvyuVcTPH9EBp%%2FXTI9E%%2F%%2Fmq171%%2BZU178wcwqU%%3D"
  }

  tags = {
    %[1]q = %[2]q
  }
}
`, key1, value1))
}

func testAccLocationAzureBlobConfig_tags2(rName, key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(testAccLocationAzureBlobConfig_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_azure_blob" "test" {
  agent_arns          = [aws_datasync_agent.test.arn]
  authentication_type = "SAS"
  container_url       = "https://example.com/path"

  sas_configuration {
    token = "sp=r&st=2023-12-20T14:54:52Z&se=2023-12-20T22:54:52Z&spr=https&sv=2021-06-08&sr=c&sig=aBBKDWQvyuVcTPH9EBp%%2FXTI9E%%2F%%2Fmq171%%2BZU178wcwqU%%3D"
  }

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, key1, value1, key2, value2))
}

func testAccLocationAzureBlobConfig_updated(rName string) string {
	return acctest.ConfigCompose(testAccLocationAzureBlobConfig_base(rName), `
resource "aws_datasync_location_azure_blob" "test" {
  access_tier         = "COOL"
  agent_arns          = [aws_datasync_agent.test.arn]
  authentication_type = "SAS"
  container_url       = "https://example.com/path"

  sas_configuration {
    token = "sp=r&st=2023-12-20T14:54:52Z&se=2023-12-20T22:54:52Z&spr=https&sv=2021-06-08&sr=c&sig=aBBKDWQvyuVcTPH9EBp%%2FXTI9E%%2F%%2Fmq171%%2BZU178wcwqU%%3D"
  }
}
`)
}
