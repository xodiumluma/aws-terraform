// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package applicationinsights

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationinsights"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_applicationinsights_application", &resource.Sweeper{
		Name: "aws_applicationinsights_application",
		F:    sweepApplications,
	})
}

func sweepApplications(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)

	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	conn := client.ApplicationInsightsConn(ctx)
	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	err = conn.ListApplicationsPagesWithContext(ctx, &applicationinsights.ListApplicationsInput{}, func(resp *applicationinsights.ListApplicationsOutput, lastPage bool) bool {
		if len(resp.ApplicationInfoList) == 0 {
			log.Print("[DEBUG] No ApplicationInsights Applications to sweep")
			return !lastPage
		}

		for _, c := range resp.ApplicationInfoList {
			r := ResourceApplication()
			d := r.Data(nil)
			d.SetId(aws.StringValue(c.ResourceGroupName))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error describing ApplicationInsights Applications: %w", err))
		// in case work can be done, don't jump out yet
	}

	if err = sweep.SweepOrchestrator(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("error sweeping ApplicationInsights Applications for %s: %w", region, err))
	}

	if awsv1.SkipSweepError(errs.ErrorOrNil()) {
		log.Printf("[WARN] Skipping ApplicationInsights Application sweep for %s: %s", region, err)
		return nil
	}

	return errs.ErrorOrNil()
}
