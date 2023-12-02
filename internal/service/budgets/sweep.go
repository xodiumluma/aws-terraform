// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package budgets

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/budgets"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_budgets_budget_action", &resource.Sweeper{
		Name: "aws_budgets_budget_action",
		F:    sweepBudgetActions,
	})

	resource.AddTestSweepers("aws_budgets_budget", &resource.Sweeper{
		Name: "aws_budgets_budget",
		F:    sweepBudgets,
		Dependencies: []string{
			"aws_budgets_budget_action",
		},
	})
}

func sweepBudgetActions(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.BudgetsConn(ctx)
	accountID := client.AccountID
	input := &budgets.DescribeBudgetActionsForAccountInput{
		AccountId: aws.String(accountID),
	}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeBudgetActionsForAccountPagesWithContext(ctx, input, func(page *budgets.DescribeBudgetActionsForAccountOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Actions {
			r := ResourceBudgetAction()
			d := r.Data(nil)
			d.SetId(BudgetActionCreateResourceID(accountID, aws.StringValue(v.ActionId), aws.StringValue(v.BudgetName)))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Budget Action sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Budget Actions (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Budget Actions (%s): %w", region, err)
	}

	return nil
}

func sweepBudgets(region string) error { // nosemgrep:ci.budgets-in-func-name
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.BudgetsConn(ctx)
	accountID := client.AccountID
	input := &budgets.DescribeBudgetsInput{
		AccountId: aws.String(accountID),
	}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.DescribeBudgetsPagesWithContext(ctx, input, func(page *budgets.DescribeBudgetsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Budgets {
			// skip budgets we have configured to track our spend
			budgetName := aws.StringValue(v.BudgetName)
			if !strings.HasPrefix(budgetName, "tf-acc") {
				continue
			}

			r := ResourceBudget()
			d := r.Data(nil)
			d.SetId(BudgetCreateResourceID(accountID, budgetName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Budget sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Budgets (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Budgets (%s): %w", region, err)
	}

	return nil
}
