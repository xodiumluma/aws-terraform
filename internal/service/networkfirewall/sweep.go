// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkfirewall

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_networkfirewall_firewall_policy", &resource.Sweeper{
		Name: "aws_networkfirewall_firewall_policy",
		F:    sweepFirewallPolicies,
		Dependencies: []string{
			"aws_networkfirewall_firewall",
		},
	})

	resource.AddTestSweepers("aws_networkfirewall_firewall", &resource.Sweeper{
		Name: "aws_networkfirewall_firewall",
		F:    sweepFirewalls,
		Dependencies: []string{
			"aws_networkfirewall_logging_configuration",
		},
	})

	resource.AddTestSweepers("aws_networkfirewall_logging_configuration", &resource.Sweeper{
		Name: "aws_networkfirewall_logging_configuration",
		F:    sweepLoggingConfigurations,
	})

	resource.AddTestSweepers("aws_networkfirewall_rule_group", &resource.Sweeper{
		Name: "aws_networkfirewall_rule_group",
		F:    sweepRuleGroups,
		Dependencies: []string{
			"aws_networkfirewall_firewall_policy",
		},
	})
}

func sweepFirewallPolicies(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.NetworkFirewallConn(ctx)
	input := &networkfirewall.ListFirewallPoliciesInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListFirewallPoliciesPagesWithContext(ctx, input, func(page *networkfirewall.ListFirewallPoliciesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.FirewallPolicies {
			r := ResourceFirewallPolicy()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.Arn))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping NetworkFirewall Firewall Policy sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing NetworkFirewall Firewall Policies (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping NetworkFirewall Firewall Policies (%s): %w", region, err)
	}

	return nil
}

func sweepFirewalls(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.NetworkFirewallConn(ctx)
	input := &networkfirewall.ListFirewallsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListFirewallsPagesWithContext(ctx, input, func(page *networkfirewall.ListFirewallsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Firewalls {
			r := ResourceFirewall()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.FirewallArn))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping NetworkFirewall Firewall sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing NetworkFirewall Firewalls (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping NetworkFirewall Firewalls (%s): %w", region, err)
	}

	return nil
}

func sweepLoggingConfigurations(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.NetworkFirewallConn(ctx)
	input := &networkfirewall.ListFirewallsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListFirewallsPagesWithContext(ctx, input, func(page *networkfirewall.ListFirewallsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Firewalls {
			r := ResourceLoggingConfiguration()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.FirewallArn))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping NetworkFirewall Logging Configuration sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing NetworkFirewall Firewalls (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping NetworkFirewall Logging Configurations (%s): %w", region, err)
	}

	return nil
}

func sweepRuleGroups(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.NetworkFirewallConn(ctx)
	input := &networkfirewall.ListRuleGroupsInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListRuleGroupsPagesWithContext(ctx, input, func(page *networkfirewall.ListRuleGroupsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.RuleGroups {
			r := ResourceRuleGroup()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.Arn))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping NetworkFirewall Rule Group sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing NetworkFirewall Rule Groups (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping NetworkFirewall Rule Groups (%s): %w", region, err)
	}

	return nil
}
