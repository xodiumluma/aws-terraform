// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkfirewall_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfnetworkfirewall "github.com/hashicorp/terraform-provider-aws/internal/service/networkfirewall"
)

func init() {
	acctest.RegisterServiceErrorCheckFunc(networkfirewall.EndpointsID, testAccErrorCheckSkip)
}

func testAccErrorCheckSkip(t *testing.T) resource.ErrorCheckFunc {
	return acctest.ErrorCheckSkipMessagesContaining(t,
		"The supplied policy does not match RAM managed permissions",
	)
}

func TestAccNetworkFirewallResourcePolicy_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_networkfirewall_resource_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkfirewall.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourcePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "resource_arn", "aws_networkfirewall_firewall_policy.test", "arn"),
					resource.TestMatchResourceAttr(resourceName, "policy", regexache.MustCompile(`"Action":\["network-firewall:CreateFirewall","network-firewall:UpdateFirewall","network-firewall:AssociateFirewallPolicy","network-firewall:ListFirewallPolicies"\]`)),
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

func TestAccNetworkFirewallResourcePolicy_ignoreEquivalent(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_networkfirewall_resource_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkfirewall.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourcePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "resource_arn", "aws_networkfirewall_firewall_policy.test", "arn"),
					resource.TestMatchResourceAttr(resourceName, "policy", regexache.MustCompile(`"Action":\["network-firewall:CreateFirewall","network-firewall:UpdateFirewall","network-firewall:AssociateFirewallPolicy","network-firewall:ListFirewallPolicies"\]`)),
				),
			},
			{
				Config: testAccResourcePolicyConfig_equivalent(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					resource.TestMatchResourceAttr(resourceName, "policy", regexache.MustCompile(`"Action":\["network-firewall:CreateFirewall","network-firewall:UpdateFirewall","network-firewall:AssociateFirewallPolicy","network-firewall:ListFirewallPolicies"\]`)),
				),
			},
		},
	})
}

func TestAccNetworkFirewallResourcePolicy_ruleGroup(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_networkfirewall_resource_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkfirewall.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourcePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyConfig_ruleGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "resource_arn", "aws_networkfirewall_rule_group.test", "arn"),
					resource.TestMatchResourceAttr(resourceName, "policy", regexache.MustCompile(`"Action":\["network-firewall:CreateFirewallPolicy","network-firewall:UpdateFirewallPolicy","network-firewall:ListRuleGroups"\]`)),
				),
			},
			{
				// Update the policy's Actions
				Config: testAccResourcePolicyConfig_ruleGroupUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					resource.TestMatchResourceAttr(resourceName, "policy", regexache.MustCompile(`"Action":\["network-firewall:CreateFirewallPolicy","network-firewall:UpdateFirewallPolicy","network-firewall:ListRuleGroups"\]`)),
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

func TestAccNetworkFirewallResourcePolicy_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_networkfirewall_resource_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkfirewall.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourcePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfnetworkfirewall.ResourceResourcePolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNetworkFirewallResourcePolicy_Disappears_firewallPolicy(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_networkfirewall_resource_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkfirewall.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourcePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfnetworkfirewall.ResourceFirewallPolicy(), "aws_networkfirewall_firewall_policy.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccNetworkFirewallResourcePolicy_Disappears_ruleGroup(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_networkfirewall_resource_policy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, networkfirewall.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckResourcePolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePolicyConfig_ruleGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourcePolicyExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfnetworkfirewall.ResourceRuleGroup(), "aws_networkfirewall_rule_group.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckResourcePolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_networkfirewall_resource_policy" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkFirewallConn(ctx)
			policy, err := tfnetworkfirewall.FindResourcePolicy(ctx, conn, rs.Primary.ID)
			if tfawserr.ErrCodeEquals(err, networkfirewall.ErrCodeResourceNotFoundException) {
				continue
			}
			if err != nil {
				return err
			}
			if policy != nil {
				return fmt.Errorf("NetworkFirewall Resource Policy (for resource: %s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckResourcePolicyExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No NetworkFirewall Resource Policy ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).NetworkFirewallConn(ctx)
		policy, err := tfnetworkfirewall.FindResourcePolicy(ctx, conn, rs.Primary.ID)
		if err != nil {
			return err
		}

		if policy == nil {
			return fmt.Errorf("NetworkFirewall Resource Policy (for resource: %s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccResourcePolicyFirewallPolicyBaseConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_networkfirewall_firewall_policy" "test" {
  name = %[1]q
  firewall_policy {
    stateless_fragment_default_actions = ["aws:drop"]
    stateless_default_actions          = ["aws:pass"]
  }
}
`, rName)
}

func testAccResourcePolicyConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccResourcePolicyFirewallPolicyBaseConfig(rName), `
resource "aws_networkfirewall_resource_policy" "test" {
  resource_arn = aws_networkfirewall_firewall_policy.test.arn
  # policy's Action element must include all of the following operations
  policy = jsonencode({
    Statement = [{
      Action = [
        "network-firewall:CreateFirewall",
        "network-firewall:UpdateFirewall",
        "network-firewall:AssociateFirewallPolicy",
        "network-firewall:ListFirewallPolicies",
      ]
      Effect   = "Allow"
      Resource = aws_networkfirewall_firewall_policy.test.arn
      Principal = {
        AWS = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"
      }
    }]
    Version = "2012-10-17"
  })
}
`)
}

func testAccResourcePolicyConfig_equivalent(rName string) string {
	return acctest.ConfigCompose(
		testAccResourcePolicyFirewallPolicyBaseConfig(rName), `
resource "aws_networkfirewall_resource_policy" "test" {
  resource_arn = aws_networkfirewall_firewall_policy.test.arn
  # policy's Action element must include all of the following operations
  policy = jsonencode({
    Statement = [{
      Action = [
        "network-firewall:UpdateFirewall",
        "network-firewall:AssociateFirewallPolicy",
        "network-firewall:ListFirewallPolicies",
        "network-firewall:CreateFirewall",
      ]
      Effect   = "Allow"
      Resource = aws_networkfirewall_firewall_policy.test.arn
      Principal = {
        AWS = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"
      }
    }]
    Version = "2012-10-17"
  })
}
`)
}

func testAccResourcePolicyRuleGroupBaseConfig(rName string) string {
	return fmt.Sprintf(`
data "aws_partition" "current" {}

data "aws_caller_identity" "current" {}

resource "aws_networkfirewall_rule_group" "test" {
  capacity = 100
  name     = %q
  type     = "STATEFUL"
  rule_group {
    rules_source {
      rules_source_list {
        generated_rules_type = "ALLOWLIST"
        target_types         = ["HTTP_HOST"]
        targets              = ["test.example.com"]
      }
    }
  }
}
`, rName)
}

func testAccResourcePolicyConfig_ruleGroup(rName string) string {
	return acctest.ConfigCompose(
		testAccResourcePolicyRuleGroupBaseConfig(rName), `
resource "aws_networkfirewall_resource_policy" "test" {
  resource_arn = aws_networkfirewall_rule_group.test.arn
  # policy's Action element must include all of the following operations
  policy = jsonencode({
    Statement = [{
      Action = [
        "network-firewall:CreateFirewallPolicy",
        "network-firewall:UpdateFirewallPolicy",
        "network-firewall:ListRuleGroups"
      ]
      Effect   = "Allow"
      Resource = aws_networkfirewall_rule_group.test.arn
      Principal = {
        AWS = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"
      }
    }]
    Version = "2012-10-17"
  })
}
`)
}

func testAccResourcePolicyConfig_ruleGroupUpdate(rName string) string {
	return acctest.ConfigCompose(
		testAccResourcePolicyRuleGroupBaseConfig(rName), `
resource "aws_networkfirewall_resource_policy" "test" {
  resource_arn = aws_networkfirewall_rule_group.test.arn
  # policy's Action element must include all of the following operations
  policy = jsonencode({
    Statement = [{
      Action = [
        "network-firewall:UpdateFirewallPolicy",
        "network-firewall:ListRuleGroups",
        "network-firewall:CreateFirewallPolicy",
      ]
      Effect   = "Allow"
      Resource = aws_networkfirewall_rule_group.test.arn
      Principal = {
        AWS = "arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"
      }
    }]
    Version = "2012-10-17"
  })
}
`)
}
