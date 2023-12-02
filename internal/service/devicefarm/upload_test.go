// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package devicefarm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/devicefarm"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdevicefarm "github.com/hashicorp/terraform-provider-aws/internal/service/devicefarm"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccDeviceFarmUpload_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var proj devicefarm.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rNameUpdated := sdkacctest.RandomWithPrefix("tf-acc-test-updated")
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, devicefarm.EndpointsID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, devicefarm.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "devicefarm", regexache.MustCompile(`upload:.+`)),
					resource.TestCheckResourceAttr(resourceName, "type", "APPIUM_JAVA_TESTNG_TEST_SPEC"),
					resource.TestCheckResourceAttr(resourceName, "category", "PRIVATE"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"url"},
			},
			{
				Config: testAccUploadConfig_basic(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "devicefarm", regexache.MustCompile(`upload:.+`)),
					resource.TestCheckResourceAttr(resourceName, "type", "APPIUM_JAVA_TESTNG_TEST_SPEC"),
					resource.TestCheckResourceAttr(resourceName, "category", "PRIVATE"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
		},
	})
}

func TestAccDeviceFarmUpload_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var proj devicefarm.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, devicefarm.EndpointsID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, devicefarm.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceUpload(), resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceUpload(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDeviceFarmUpload_disappears_project(t *testing.T) {
	ctx := acctest.Context(t)
	var proj devicefarm.Upload
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_upload.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, devicefarm.EndpointsID)
			// Currently, DeviceFarm is only supported in us-west-2
			// https://docs.aws.amazon.com/general/latest/gr/devicefarm.html
			acctest.PreCheckRegion(t, endpoints.UsWest2RegionID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, devicefarm.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckUploadDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccUploadConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUploadExists(ctx, resourceName, &proj),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceProject(), "aws_devicefarm_project.test"),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceUpload(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckUploadExists(ctx context.Context, n string, v *devicefarm.Upload) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DeviceFarmConn(ctx)
		resp, err := tfdevicefarm.FindUploadByARN(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("DeviceFarm Upload not found")
		}

		*v = *resp

		return nil
	}
}

func testAccCheckUploadDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DeviceFarmConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_devicefarm_upload" {
				continue
			}

			// Try to find the resource
			_, err := tfdevicefarm.FindUploadByARN(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DeviceFarm Upload %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccUploadConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_devicefarm_project" "test" {
  name = %[1]q
}

resource "aws_devicefarm_upload" "test" {
  name        = %[1]q
  project_arn = aws_devicefarm_project.test.arn
  type        = "APPIUM_JAVA_TESTNG_TEST_SPEC"
}
`, rName)
}
