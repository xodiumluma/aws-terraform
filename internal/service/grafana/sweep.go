// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package grafana

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/managedgrafana"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_grafana_workspace", &resource.Sweeper{
		Name: "aws_grafana_workspace",
		F:    sweepWorkSpaces,
	})
}

func sweepWorkSpaces(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	conn := client.GrafanaConn(ctx)

	sweepResources := make([]sweep.Sweepable, 0)
	var errs *multierror.Error

	input := &managedgrafana.ListWorkspacesInput{}

	err = conn.ListWorkspacesPagesWithContext(ctx, input, func(page *managedgrafana.ListWorkspacesOutput, lastPage bool) bool {
		if len(page.Workspaces) == 0 {
			log.Printf("[INFO] No Grafana Workspaces to sweep")
			return false
		}
		for _, workspace := range page.Workspaces {
			id := aws.StringValue(workspace.Id)
			log.Printf("[INFO] Deleting Grafana Workspace: %s", id)
			r := ResourceWorkspace()
			d := r.Data(nil)
			d.SetId(id)

			if err != nil {
				err := fmt.Errorf("reading Grafana Workspace (%s): %w", id, err)
				errs = multierror.Append(errs, err)
				continue
			}

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}
		return !lastPage
	})

	if err != nil {
		errs = multierror.Append(errs, fmt.Errorf("listing Grafana Workspace for %s: %w", region, err))
	}

	if err := sweep.SweepOrchestrator(ctx, sweepResources); err != nil {
		errs = multierror.Append(errs, fmt.Errorf("sweeping Grafana Workspace for %s: %w", region, err))
	}

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Grafana Workspace sweep for %s: %s", region, errs)
		return nil
	}

	return errs.ErrorOrNil()
}
