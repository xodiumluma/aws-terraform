// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kms_test

import (
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/kms"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	tfkms "github.com/hashicorp/terraform-provider-aws/internal/service/kms"
)

func TestAccKMSReplicaKey_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var key kms.KeyMetadata
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	primaryKeyResourceName := "aws_kms_key.test"
	resourceName := "aws_kms_replica_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaKeyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "kms", regexache.MustCompile(`key/.+`)),
					resource.TestCheckResourceAttr(resourceName, "bypass_policy_lockout_safety_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "deletion_window_in_days", "30"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "key_rotation_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "key_spec", "SYMMETRIC_DEFAULT"),
					resource.TestMatchResourceAttr(resourceName, "policy", regexache.MustCompile(`Enable IAM User Permissions`)),
					resource.TestCheckResourceAttrPair(resourceName, "primary_key_arn", primaryKeyResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_window_in_days", "bypass_policy_lockout_safety_check"},
			},
		},
	})
}

func TestAccKMSReplicaKey_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var key kms.KeyMetadata
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_replica_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaKeyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfkms.ResourceReplicaKey(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccKMSReplicaKey_descriptionAndEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	var key kms.KeyMetadata
	rName1 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName3 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName4 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_replica_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaKeyConfig_descriptionAndEnabled(rName1, rName2, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "description", rName2),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_window_in_days", "bypass_policy_lockout_safety_check"},
			},
			{
				Config: testAccReplicaKeyConfig_descriptionAndEnabled(rName1, rName3, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "description", rName3),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccReplicaKeyConfig_descriptionAndEnabled(rName1, rName4, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "description", rName4),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAccKMSReplicaKey_policy(t *testing.T) {
	ctx := acctest.Context(t)
	var key kms.KeyMetadata
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_replica_key.test"
	policy1 := `{"Id":"kms-tf-1","Statement":[{"Action":"kms:*","Effect":"Allow","Principal":{"AWS":"*"},"Resource":"*","Sid":"Enable IAM User Permissions 1"}],"Version":"2012-10-17"}`
	policy2 := `{"Id":"kms-tf-1","Statement":[{"Action":"kms:*","Effect":"Allow","Principal":{"AWS":"*"},"Resource":"*","Sid":"Enable IAM User Permissions 2"}],"Version":"2012-10-17"}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaKeyConfig_policy(rName, policy1, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "bypass_policy_lockout_safety_check", "false"),
					testAccCheckKeyHasPolicy(ctx, resourceName, policy1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_window_in_days", "bypass_policy_lockout_safety_check"},
			},
			{
				Config: testAccReplicaKeyConfig_policy(rName, policy2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "bypass_policy_lockout_safety_check", "true"),
					testAccCheckExternalKeyHasPolicy(ctx, resourceName, policy2),
				),
			},
		},
	})
}

func TestAccKMSReplicaKey_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var key kms.KeyMetadata
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_replica_key.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 2)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesAlternate(ctx, t),
		CheckDestroy:             testAccCheckKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaKeyConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"deletion_window_in_days",
					"bypass_policy_lockout_safety_check",
				},
			},
			{
				Config: testAccReplicaKeyConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccReplicaKeyConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccReplicaKeyConfig_tags0(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"deletion_window_in_days",
					"bypass_policy_lockout_safety_check",
				},
			},
		},
	})
}

func TestAccKMSReplicaKey_twoReplicas(t *testing.T) {
	ctx := acctest.Context(t)
	var key kms.KeyMetadata
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_kms_replica_key.test1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckMultipleRegion(t, 3)
		},
		ErrorCheck:               acctest.ErrorCheck(t, kms.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5FactoriesMultipleRegions(ctx, t, 3),
		CheckDestroy:             testAccCheckKeyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccReplicaKeyConfig_two(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(ctx, resourceName, &key),
				),
			},
		},
	})
}

func testAccReplicaKeyConfig_basic(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true
}

resource "aws_kms_replica_key" "test" {
  primary_key_arn = aws_kms_key.test.arn
}
`, rName))
}

func testAccReplicaKeyConfig_descriptionAndEnabled(rName, description string, enabled bool) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test" {
  description     = %[2]q
  enabled         = %[3]t
  primary_key_arn = aws_kms_key.test.arn

  deletion_window_in_days = 7
}
`, rName, description, enabled))
}

func testAccReplicaKeyConfig_policy(rName, policy string, bypassLockoutCheck bool) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test" {
  description     = %[1]q
  primary_key_arn = aws_kms_key.test.arn

  deletion_window_in_days = 7

  bypass_policy_lockout_safety_check = %[3]t

  policy = %[2]q
}
`, rName, policy, bypassLockoutCheck))
}

func testAccReplicaKeyConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true

  tags = {
    Name = %[1]q
  }

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test" {
  description     = %[1]q
  primary_key_arn = aws_kms_key.test.arn

  tags = {
    %[2]q = %[3]q
  }

  deletion_window_in_days = 7
}
`, rName, tagKey1, tagValue1))
}

func testAccReplicaKeyConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true

  tags = {
    Name = %[1]q
  }

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test" {
  description     = %[1]q
  primary_key_arn = aws_kms_key.test.arn

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }

  deletion_window_in_days = 7
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccReplicaKeyConfig_tags0(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigAlternateRegionProvider(), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test" {
  description     = %[1]q
  primary_key_arn = aws_kms_key.test.arn

  deletion_window_in_days = 7
}
`, rName))
}

func testAccReplicaKeyConfig_two(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigMultipleRegionProvider(3), fmt.Sprintf(`
resource "aws_kms_key" "test" {
  provider = awsalternate

  description  = %[1]q
  multi_region = true

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test1" {
  description     = %[1]q
  primary_key_arn = aws_kms_key.test.arn

  deletion_window_in_days = 7
}

resource "aws_kms_replica_key" "test2" {
  provider = awsthird

  description     = %[1]q
  primary_key_arn = aws_kms_key.test.arn

  deletion_window_in_days = 7
}
`, rName))
}
