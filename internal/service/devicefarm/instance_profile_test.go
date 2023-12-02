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

func TestAccDeviceFarmInstanceProfile_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var profile devicefarm.InstanceProfile
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rNameUpdated := sdkacctest.RandomWithPrefix("tf-acc-test-updated")
	resourceName := "aws_devicefarm_instance_profile.test"

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
		CheckDestroy:             testAccCheckInstanceProfileDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceProfileConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceProfileExists(ctx, resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "reboot_after_use", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "devicefarm", regexache.MustCompile(`instanceprofile:.+`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccInstanceProfileConfig_basic(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceProfileExists(ctx, resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "reboot_after_use", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "devicefarm", regexache.MustCompile(`instanceprofile:.+`)),
				),
			},
		},
	})
}

func TestAccDeviceFarmInstanceProfile_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var profile devicefarm.InstanceProfile
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_instance_profile.test"

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
		CheckDestroy:             testAccCheckInstanceProfileDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceProfileConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceProfileExists(ctx, resourceName, &profile),
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
				Config: testAccInstanceProfileConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceProfileExists(ctx, resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccInstanceProfileConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceProfileExists(ctx, resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccDeviceFarmInstanceProfile_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var profile devicefarm.InstanceProfile
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_devicefarm_instance_profile.test"

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
		CheckDestroy:             testAccCheckInstanceProfileDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceProfileConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceProfileExists(ctx, resourceName, &profile),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceInstanceProfile(), resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdevicefarm.ResourceInstanceProfile(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckInstanceProfileExists(ctx context.Context, n string, v *devicefarm.InstanceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DeviceFarmConn(ctx)
		resp, err := tfdevicefarm.FindInstanceProfileByARN(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("DeviceFarm Instance Profile not found")
		}

		*v = *resp

		return nil
	}
}

func testAccCheckInstanceProfileDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DeviceFarmConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_devicefarm_instance_profile" {
				continue
			}

			// Try to find the resource
			_, err := tfdevicefarm.FindInstanceProfileByARN(ctx, conn, rs.Primary.ID)
			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DeviceFarm Instance Profile %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccInstanceProfileConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_devicefarm_instance_profile" "test" {
  name = %[1]q
}
`, rName)
}

func testAccInstanceProfileConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_devicefarm_instance_profile" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccInstanceProfileConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_devicefarm_instance_profile" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
