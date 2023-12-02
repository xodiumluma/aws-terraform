// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dms

import (
	"context"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKDataSource("aws_dms_certificate")
func DataSourceCertificate() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceCertificateRead,

		Schema: map[string]*schema.Schema{
			"certificate_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 255),
					validation.StringMatch(regexache.MustCompile("^[A-Za-z][0-9A-Za-z-]+$"), "must start with a letter, only contain alphanumeric characters and hyphens"),
					validation.StringDoesNotMatch(regexache.MustCompile(`--`), "cannot contain two consecutive hyphens"),
					validation.StringDoesNotMatch(regexache.MustCompile(`-$`), "cannot end in a hyphen"),
				),
			},
			"certificate_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_pem": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"certificate_wallet": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"key_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"signing_algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"valid_from_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_to_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCertificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DMSConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	certificateID := d.Get("certificate_id").(string)
	out, err := FindCertificateByID(ctx, conn, certificateID)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DMS Certificate (%s): %s", certificateID, err)
	}

	d.SetId(aws.StringValue(out.CertificateIdentifier))
	arn := aws.StringValue(out.CertificateArn)
	d.Set("certificate_arn", arn)
	d.Set("certificate_id", out.CertificateIdentifier)
	d.Set("certificate_pem", out.CertificatePem)
	if len(out.CertificateWallet) != 0 {
		d.Set("certificate_wallet", verify.Base64Encode(out.CertificateWallet))
	}
	d.Set("key_length", out.KeyLength)
	d.Set("signing_algorithm", out.SigningAlgorithm)
	d.Set("valid_from_date", out.ValidFromDate.String())
	d.Set("valid_to_date", out.ValidToDate.String())

	tags, err := listTags(ctx, conn, arn)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "listing tags for DMS Certificate (%s): %s", arn, err)
	}

	tags = tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return nil
}
