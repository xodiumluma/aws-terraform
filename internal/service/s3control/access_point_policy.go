// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package s3control

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3control"
	"github.com/aws/aws-sdk-go-v2/service/s3control/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_s3control_access_point_policy")
func resourceAccessPointPolicy() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAccessPointPolicyCreate,
		ReadWithoutTimeout:   resourceAccessPointPolicyRead,
		UpdateWithoutTimeout: resourceAccessPointPolicyUpdate,
		DeleteWithoutTimeout: resourceAccessPointPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAccessPointPolicyImport,
		},

		Schema: map[string]*schema.Schema{
			"access_point_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"has_public_access_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"policy": {
				Type:                  schema.TypeString,
				Required:              true,
				ValidateFunc:          validation.StringIsJSON,
				DiffSuppressFunc:      verify.SuppressEquivalentPolicyDiffs,
				DiffSuppressOnRefresh: true,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
		},
	}
}

func resourceAccessPointPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).S3ControlClient(ctx)

	resourceID, err := AccessPointCreateResourceID(d.Get("access_point_arn").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	accountID, name, err := AccessPointParseResourceID(resourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input := &s3control.PutAccessPointPolicyInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
		Policy:    aws.String(policy),
	}

	_, err = conn.PutAccessPointPolicy(ctx, input)

	if err != nil {
		return diag.Errorf("creating S3 Access Point (%s) Policy: %s", resourceID, err)
	}

	d.SetId(resourceID)

	return resourceAccessPointPolicyRead(ctx, d, meta)
}

func resourceAccessPointPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).S3ControlClient(ctx)

	accountID, name, err := AccessPointParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policy, status, err := findAccessPointPolicyAndStatusByTwoPartKey(ctx, conn, accountID, name)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] S3 Access Point Policy (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading S3 Access Point Policy (%s): %s", d.Id(), err)
	}

	d.Set("has_public_access_policy", status.IsPublic)

	if policy != "" {
		policyToSet, err := verify.PolicyToSet(d.Get("policy").(string), policy)
		if err != nil {
			return diag.FromErr(err)
		}

		d.Set("policy", policyToSet)
	} else {
		d.Set("policy", "")
	}

	return nil
}

func resourceAccessPointPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).S3ControlClient(ctx)

	accountID, name, err := AccessPointParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := structure.NormalizeJsonString(d.Get("policy").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input := &s3control.PutAccessPointPolicyInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
		Policy:    aws.String(policy),
	}

	_, err = conn.PutAccessPointPolicy(ctx, input)

	if err != nil {
		return diag.Errorf("updating S3 Access Point Policy (%s): %s", d.Id(), err)
	}

	return resourceAccessPointPolicyRead(ctx, d, meta)
}

func resourceAccessPointPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).S3ControlClient(ctx)

	accountID, name, err := AccessPointParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting S3 Access Point Policy: %s", d.Id())
	_, err = conn.DeleteAccessPointPolicy(ctx, &s3control.DeleteAccessPointPolicyInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
	})

	if tfawserr.ErrCodeEquals(err, errCodeNoSuchAccessPoint, errCodeNoSuchAccessPointPolicy) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting S3 Access Point Policy (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAccessPointPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	resourceID, err := AccessPointCreateResourceID(d.Id())

	if err != nil {
		return nil, err
	}

	d.Set("access_point_arn", d.Id())
	d.SetId(resourceID)

	return []*schema.ResourceData{d}, nil
}

func findAccessPointPolicyAndStatusByTwoPartKey(ctx context.Context, conn *s3control.Client, accountID, name string) (string, *types.PolicyStatus, error) {
	inputGAPP := &s3control.GetAccessPointPolicyInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
	}

	outputGAPP, err := conn.GetAccessPointPolicy(ctx, inputGAPP)

	if tfawserr.ErrCodeEquals(err, errCodeNoSuchAccessPoint, errCodeNoSuchAccessPointPolicy) {
		return "", nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: inputGAPP,
		}
	}

	if err != nil {
		return "", nil, err
	}

	if outputGAPP == nil {
		return "", nil, tfresource.NewEmptyResultError(inputGAPP)
	}

	policy := aws.ToString(outputGAPP.Policy)

	if policy == "" {
		return "", nil, tfresource.NewEmptyResultError(inputGAPP)
	}

	inputGAPPS := &s3control.GetAccessPointPolicyStatusInput{
		AccountId: aws.String(accountID),
		Name:      aws.String(name),
	}

	outputGAPPS, err := conn.GetAccessPointPolicyStatus(ctx, inputGAPPS)

	if tfawserr.ErrCodeEquals(err, errCodeNoSuchAccessPoint, errCodeNoSuchAccessPointPolicy) {
		return "", nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: inputGAPPS,
		}
	}

	if err != nil {
		return "", nil, err
	}

	if outputGAPPS == nil || outputGAPPS.PolicyStatus == nil {
		return "", nil, tfresource.NewEmptyResultError(inputGAPPS)
	}

	return policy, outputGAPPS.PolicyStatus, nil
}
