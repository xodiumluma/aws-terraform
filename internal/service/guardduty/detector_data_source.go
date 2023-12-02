// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package guardduty

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

// @SDKDataSource("aws_guardduty_detector")
func DataSourceDetector() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceDetectorRead,

		Schema: map[string]*schema.Schema{
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_configuration": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"finding_publishing_frequency": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_role_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDetectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).GuardDutyConn(ctx)

	detectorID := d.Get("id").(string)

	if detectorID == "" {
		output, err := FindDetector(ctx, conn)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "reading this account's single GuardDuty Detector: %s", err)
		}

		detectorID = aws.StringValue(output)
	}

	gdo, err := FindDetectorByID(ctx, conn, detectorID)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading GuardDuty Detector (%s): %s", detectorID, err)
	}

	d.SetId(detectorID)
	if gdo.Features != nil {
		if err := d.Set("features", flattenDetectorFeatureConfigurationResults(gdo.Features)); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting features: %s", err)
		}
	} else {
		d.Set("features", nil)
	}
	d.Set("finding_publishing_frequency", gdo.FindingPublishingFrequency)
	d.Set("service_role_arn", gdo.ServiceRole)
	d.Set("status", gdo.Status)

	return diags
}
