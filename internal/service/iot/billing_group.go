// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iot

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_iot_billing_group", name="Billing Group")
// @Tags(identifierAttribute="arn")
func ResourceBillingGroup() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceBillingGroupCreate,
		ReadWithoutTimeout:   resourceBillingGroupRead,
		UpdateWithoutTimeout: resourceBillingGroupUpdate,
		DeleteWithoutTimeout: resourceBillingGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceBillingGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTConn(ctx)

	name := d.Get("name").(string)
	input := &iot.CreateBillingGroupInput{
		BillingGroupName: aws.String(name),
		Tags:             getTagsIn(ctx),
	}

	if v, ok := d.GetOk("properties"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.BillingGroupProperties = expandBillingGroupProperties(v.([]interface{})[0].(map[string]interface{}))
	}

	output, err := conn.CreateBillingGroupWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating IoT Billing Group (%s): %s", name, err)
	}

	d.SetId(aws.StringValue(output.BillingGroupName))

	return append(diags, resourceBillingGroupRead(ctx, d, meta)...)
}

func resourceBillingGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTConn(ctx)

	output, err := FindBillingGroupByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] IoT Billing Group (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading IoT Billing Group (%s): %s", d.Id(), err)
	}

	d.Set("arn", output.BillingGroupArn)
	d.Set("name", output.BillingGroupName)

	if output.BillingGroupMetadata != nil {
		if err := d.Set("metadata", []interface{}{flattenBillingGroupMetadata(output.BillingGroupMetadata)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting metadata: %s", err)
		}
	} else {
		d.Set("metadata", nil)
	}
	if v := flattenBillingGroupProperties(output.BillingGroupProperties); len(v) > 0 {
		if err := d.Set("properties", []interface{}{v}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting properties: %s", err)
		}
	} else {
		d.Set("properties", nil)
	}
	d.Set("version", output.Version)

	return diags
}

func resourceBillingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTConn(ctx)

	if d.HasChangesExcept("tags", "tags_all") {
		input := &iot.UpdateBillingGroupInput{
			BillingGroupName: aws.String(d.Id()),
			ExpectedVersion:  aws.Int64(int64(d.Get("version").(int))),
		}

		if v, ok := d.GetOk("properties"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.BillingGroupProperties = expandBillingGroupProperties(v.([]interface{})[0].(map[string]interface{}))
		} else {
			input.BillingGroupProperties = &iot.BillingGroupProperties{}
		}

		_, err := conn.UpdateBillingGroupWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating IoT Billing Group (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceBillingGroupRead(ctx, d, meta)...)
}

func resourceBillingGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).IoTConn(ctx)

	log.Printf("[DEBUG] Deleting IoT Billing Group: %s", d.Id())
	_, err := conn.DeleteBillingGroupWithContext(ctx, &iot.DeleteBillingGroupInput{
		BillingGroupName: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, iot.ErrCodeResourceNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting IoT Billing Group (%s): %s", d.Id(), err)
	}

	return diags
}

func FindBillingGroupByName(ctx context.Context, conn *iot.IoT, name string) (*iot.DescribeBillingGroupOutput, error) {
	input := &iot.DescribeBillingGroupInput{
		BillingGroupName: aws.String(name),
	}

	output, err := conn.DescribeBillingGroupWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, iot.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func expandBillingGroupProperties(tfMap map[string]interface{}) *iot.BillingGroupProperties {
	if tfMap == nil {
		return nil
	}

	apiObject := &iot.BillingGroupProperties{}

	if v, ok := tfMap["description"].(string); ok && v != "" {
		apiObject.BillingGroupDescription = aws.String(v)
	}

	return apiObject
}

func flattenBillingGroupMetadata(apiObject *iot.BillingGroupMetadata) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.CreationDate; v != nil {
		tfMap["creation_date"] = aws.TimeValue(v).Format(time.RFC3339)
	}

	return tfMap
}

func flattenBillingGroupProperties(apiObject *iot.BillingGroupProperties) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BillingGroupDescription; v != nil {
		tfMap["description"] = aws.StringValue(v)
	}

	return tfMap
}
