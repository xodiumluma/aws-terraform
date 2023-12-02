// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kinesisanalyticsv2

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesisanalyticsv2"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_kinesisanalyticsv2_application", &resource.Sweeper{
		Name: "aws_kinesisanalyticsv2_application",
		F:    sweepApplication,
	})
}

func sweepApplication(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.KinesisAnalyticsV2Conn(ctx)

	sweepResources := make([]sweep.Sweepable, 0)
	var sweeperErrs *multierror.Error

	input := &kinesisanalyticsv2.ListApplicationsInput{}
	err = listApplicationsPages(ctx, conn, input, func(page *kinesisanalyticsv2.ListApplicationsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, applicationSummary := range page.ApplicationSummaries {
			arn := aws.StringValue(applicationSummary.ApplicationARN)
			name := aws.StringValue(applicationSummary.ApplicationName)

			application, err := FindApplicationDetailByName(ctx, conn, name)

			if err != nil {
				sweeperErr := fmt.Errorf("error reading Kinesis Analytics v2 Application (%s): %w", arn, err)
				log.Printf("[ERROR] %s", err)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}

			r := ResourceApplication()
			d := r.Data(nil)
			d.SetId(arn)
			d.Set("create_timestamp", aws.TimeValue(application.CreateTimestamp).Format(time.RFC3339))
			d.Set("name", name)

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Kinesis Analytics v2 Application sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Kinesis Analytics v2 Applications: %w", err))
	}

	if err := sweep.SweepOrchestrator(ctx, sweepResources); err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping Kinesis Analytics v2 Applications: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}
