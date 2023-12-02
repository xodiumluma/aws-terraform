// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package s3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

// @SDKDataSource("aws_s3_bucket_policy")
func DataSourceBucketPolicy() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceBucketPolicyRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBucketPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).S3Client(ctx)

	name := d.Get("bucket").(string)

	policy, err := findBucketPolicy(ctx, conn, name)

	if err != nil {
		return diag.Errorf("reading S3 Bucket (%s) Policy: %s", name, err)
	}

	policy, err = structure.NormalizeJsonString(policy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	d.Set("policy", policy)

	return nil
}
