// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfkms "github.com/hashicorp/terraform-provider-aws/internal/service/kms"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccKMSAlias_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"
	keyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "kms", regexache.MustCompile(`alias/.+`)),
					resource.TestCheckResourceAttr(resourceName, "name", tfkms.AliasNamePrefix+rName),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_arn", keyResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_id", keyResourceName, "id"),
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

func TestAccKMSAlias_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfkms.ResourceAlias(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKMSAlias_Name_generated(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_nameGenerated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					resource.TestMatchResourceAttr(resourceName, "name", regexache.MustCompile(fmt.Sprintf("%s[[:xdigit:]]{%d}", tfkms.AliasNamePrefix, id.UniqueIDSuffixLength))),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", tfkms.AliasNamePrefix),
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

func TestAccKMSAlias_namePrefix(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_namePrefix(rName, tfkms.AliasNamePrefix+"tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", tfkms.AliasNamePrefix+"tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", tfkms.AliasNamePrefix+"tf-acc-test-prefix-"),
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

func TestAccKMSAlias_updateKeyID(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"
	key1ResourceName := "aws_kms_key.test"
	key2ResourceName := "aws_kms_key.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_name(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_arn", key1ResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_id", key1ResourceName, "id"),
				),
			},
			{
				Config: testAccAliasConfig_updatedKeyID(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_arn", key2ResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_id", key2ResourceName, "id"),
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

func TestAccKMSAlias_multipleAliasesForSameKey(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"
	alias2ResourceName := "aws_kms_alias.test2"
	keyResourceName := "aws_kms_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_multiple(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_arn", keyResourceName, "arn"),
					resource.TestCheckResourceAttrPair(resourceName, "target_key_id", keyResourceName, "id"),
					testAccCheckAliasExists(ctx, alias2ResourceName, &alias),
					resource.TestCheckResourceAttrPair(alias2ResourceName, "target_key_arn", keyResourceName, "arn"),
					resource.TestCheckResourceAttrPair(alias2ResourceName, "target_key_id", keyResourceName, "id"),
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

func TestAccKMSAlias_arnDiffSuppress(t *testing.T) {
	ctx := acctest.Context(t)
	var alias kms.AliasListEntry
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_alias.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAliasDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasConfig_diffSuppress(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAliasExists(ctx, resourceName, &alias),
					resource.TestCheckResourceAttrSet(resourceName, "target_key_arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				Config:             testAccAliasConfig_diffSuppress(rName),
			},
		},
	})
}

func testAccCheckAliasDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).KMSConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_kms_alias" {
				continue
			}

			_, err := tfkms.FindAliasByName(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("KMS Alias %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckAliasExists(ctx context.Context, name string, v *kms.AliasListEntry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No KMS Alias ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).KMSConn(ctx)

		output, err := tfkms.FindAliasByName(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccAliasConfig_name(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  name          = "alias/%[1]s"
  target_key_id = aws_kms_key.test.id
}
`, rName)
}

func testAccAliasConfig_nameGenerated(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  target_key_id = aws_kms_key.test.id
}
`, rName)
}

func testAccAliasConfig_namePrefix(rName, namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  name_prefix   = %[2]q
  target_key_id = aws_kms_key.test.id
}
`, rName, namePrefix)
}

func testAccAliasConfig_updatedKeyID(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_key" "test2" {
  description             = "%[1]s-2"
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  name          = "alias/%[1]s"
  target_key_id = aws_kms_key.test2.id
}
`, rName)
}

func testAccAliasConfig_multiple(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  name          = "alias/%[1]s-1"
  target_key_id = aws_kms_key.test.key_id
}

resource "aws_kms_alias" "test2" {
  name          = "alias/%[1]s-2"
  target_key_id = aws_kms_key.test.key_id
}
`, rName)
}

func testAccAliasConfig_diffSuppress(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  description             = %[1]q
  deletion_window_in_days = 7
}

resource "aws_kms_alias" "test" {
  name          = "alias/%[1]s"
  target_key_id = aws_kms_key.test.arn
}
`, rName)
}
