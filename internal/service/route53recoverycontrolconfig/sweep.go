// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package route53recoverycontrolconfig

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	r53rcc "github.com/aws/aws-sdk-go/service/route53recoverycontrolconfig"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_route53recoverycontrolconfig_cluster", &resource.Sweeper{
		Name: "aws_route53recoverycontrolconfig_cluster",
		F:    sweepClusters,
		Dependencies: []string{
			"aws_route53recoverycontrolconfig_control_panel",
		},
	})

	resource.AddTestSweepers("aws_route53recoverycontrolconfig_control_panel", &resource.Sweeper{
		Name: "aws_route53recoverycontrolconfig_control_panel",
		F:    sweepControlPanels,
		Dependencies: []string{
			"aws_route53recoverycontrolconfig_routing_control",
			"aws_route53recoverycontrolconfig_safety_rule",
		},
	})

	resource.AddTestSweepers("aws_route53recoverycontrolconfig_routing_control", &resource.Sweeper{
		Name: "aws_route53recoverycontrolconfig_routing_control",
		F:    sweepRoutingControls,
	})

	resource.AddTestSweepers("aws_route53recoverycontrolconfig_safety_rule", &resource.Sweeper{
		Name: "aws_route53recoverycontrolconfig_safety_rule",
		F:    sweepSafetyRules,
	})
}

func sweepClusters(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.Route53RecoveryControlConfigConn(ctx)
	input := &r53rcc.ListClustersInput{}
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListClustersPagesWithContext(ctx, input, func(page *r53rcc.ListClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Clusters {
			r := ResourceCluster()
			d := r.Data(nil)
			d.SetId(aws.StringValue(v.ClusterArn))

			sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Route53 Recovery Control Config Cluster sweep for %s: %s", region, err)
		return nil
	}

	if err != nil {
		return fmt.Errorf("error listing Route53 Recovery Control Config Clusters (%s): %w", region, err)
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		return fmt.Errorf("error sweeping Route53 Recovery Control Config Clusters (%s): %w", region, err)
	}

	return nil
}

func sweepControlPanels(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.Route53RecoveryControlConfigConn(ctx)
	input := &r53rcc.ListClustersInput{}
	var sweeperErrs *multierror.Error
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListClustersPagesWithContext(ctx, input, func(page *r53rcc.ListClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Clusters {
			input := &r53rcc.ListControlPanelsInput{
				ClusterArn: v.ClusterArn,
			}

			err := conn.ListControlPanelsPagesWithContext(ctx, input, func(page *r53rcc.ListControlPanelsOutput, lastPage bool) bool {
				if page == nil {
					return !lastPage
				}

				for _, v := range page.ControlPanels {
					if aws.BoolValue(v.DefaultControlPanel) {
						continue
					}

					r := ResourceControlPanel()
					d := r.Data(nil)
					d.SetId(aws.StringValue(v.ControlPanelArn))

					sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
				}

				return !lastPage
			})

			if awsv1.SkipSweepError(err) {
				continue
			}

			if err != nil {
				sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Control Panels (%s): %w", region, err))
			}
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Route53 Recovery Control Config Control Panel sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Clusters (%s): %w", region, err))
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping Route53 Recovery Control Config Control Panels (%s): %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepRoutingControls(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.Route53RecoveryControlConfigConn(ctx)
	input := &r53rcc.ListClustersInput{}
	var sweeperErrs *multierror.Error
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListClustersPagesWithContext(ctx, input, func(page *r53rcc.ListClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Clusters {
			input := &r53rcc.ListControlPanelsInput{
				ClusterArn: v.ClusterArn,
			}

			err := conn.ListControlPanelsPagesWithContext(ctx, input, func(page *r53rcc.ListControlPanelsOutput, lastPage bool) bool {
				if page == nil {
					return !lastPage
				}

				for _, v := range page.ControlPanels {
					input := &r53rcc.ListRoutingControlsInput{
						ControlPanelArn: v.ControlPanelArn,
					}

					err := conn.ListRoutingControlsPagesWithContext(ctx, input, func(page *r53rcc.ListRoutingControlsOutput, lastPage bool) bool {
						if page == nil {
							return !lastPage
						}

						for _, v := range page.RoutingControls {
							r := ResourceRoutingControl()
							d := r.Data(nil)
							d.SetId(aws.StringValue(v.RoutingControlArn))

							sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
						}

						return !lastPage
					})

					if awsv1.SkipSweepError(err) {
						continue
					}

					if err != nil {
						sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Routing Controls (%s): %w", region, err))
					}
				}

				return !lastPage
			})

			if awsv1.SkipSweepError(err) {
				continue
			}

			if err != nil {
				sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Control Panels (%s): %w", region, err))
			}
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Route53 Recovery Control Config Routing Control sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Clusters (%s): %w", region, err))
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping Route53 Recovery Control Config Routing Controls (%s): %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepSafetyRules(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.Route53RecoveryControlConfigConn(ctx)
	input := &r53rcc.ListClustersInput{}
	var sweeperErrs *multierror.Error
	sweepResources := make([]sweep.Sweepable, 0)

	err = conn.ListClustersPagesWithContext(ctx, input, func(page *r53rcc.ListClustersOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Clusters {
			input := &r53rcc.ListControlPanelsInput{
				ClusterArn: v.ClusterArn,
			}

			err := conn.ListControlPanelsPagesWithContext(ctx, input, func(page *r53rcc.ListControlPanelsOutput, lastPage bool) bool {
				if page == nil {
					return !lastPage
				}

				for _, v := range page.ControlPanels {
					input := &r53rcc.ListSafetyRulesInput{
						ControlPanelArn: v.ControlPanelArn,
					}

					err := conn.ListSafetyRulesPagesWithContext(ctx, input, func(page *r53rcc.ListSafetyRulesOutput, lastPage bool) bool {
						if page == nil {
							return !lastPage
						}

						for _, v := range page.SafetyRules {
							r := ResourceSafetyRule()
							d := r.Data(nil)
							if v.ASSERTION != nil {
								d.SetId(aws.StringValue(v.ASSERTION.SafetyRuleArn))
							} else if v.GATING != nil {
								d.SetId(aws.StringValue(v.GATING.SafetyRuleArn))
							} else {
								continue
							}

							sweepResources = append(sweepResources, sweep.NewSweepResource(r, d, client))
						}

						return !lastPage
					})

					if awsv1.SkipSweepError(err) {
						continue
					}

					if err != nil {
						sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Safety Rules (%s): %w", region, err))
					}
				}

				return !lastPage
			})

			if awsv1.SkipSweepError(err) {
				continue
			}

			if err != nil {
				sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Control Panels (%s): %w", region, err))
			}
		}

		return !lastPage
	})

	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping Route53 Recovery Control Config Safety Rule sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error listing Route53 Recovery Control Config Clusters (%s): %w", region, err))
	}

	err = sweep.SweepOrchestrator(ctx, sweepResources)

	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("error sweeping Route53 Recovery Control Config Safety Rules (%s): %w", region, err))
	}

	return sweeperErrs.ErrorOrNil()
}
