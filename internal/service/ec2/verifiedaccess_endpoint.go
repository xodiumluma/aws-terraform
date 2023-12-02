// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_verifiedaccess_endpoint", name="Verified Access Endpoint")
// @Tags(identifierAttribute="id")
func ResourceVerifiedAccessEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceVerifiedAccessEndpointCreate,
		ReadWithoutTimeout:   resourceVerifiedAccessEndpointRead,
		UpdateWithoutTimeout: resourceVerifiedAccessEndpointUpdate,
		DeleteWithoutTimeout: resourceVerifiedAccessEndpointDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"application_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(verifiedAccessAttachmentType_Values(), false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"device_validation_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_certificate_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"endpoint_domain_prefix": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(verifiedAccessEndpointType_Values(), false),
			},
			"load_balancer_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: verify.ValidARN,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IsPortNumber,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(verifiedAccessEndpointProtocol_Values(), false),
						},
						"subnet_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"network_interface_options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IsPortNumber,
						},
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(verifiedAccessEndpointProtocol_Values(), false),
						},
					},
				},
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sse_specification": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_managed_key_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"kms_key_arn": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"verified_access_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"verified_access_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceVerifiedAccessEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	input := &ec2.CreateVerifiedAccessEndpointInput{
		ApplicationDomain:     aws.String(d.Get("application_domain").(string)),
		AttachmentType:        types.VerifiedAccessEndpointAttachmentType(d.Get("attachment_type").(string)),
		DomainCertificateArn:  aws.String(d.Get("domain_certificate_arn").(string)),
		EndpointDomainPrefix:  aws.String(d.Get("endpoint_domain_prefix").(string)),
		EndpointType:          types.VerifiedAccessEndpointType(d.Get("endpoint_type").(string)),
		TagSpecifications:     getTagSpecificationsInV2(ctx, types.ResourceTypeVerifiedAccessEndpoint),
		VerifiedAccessGroupId: aws.String(d.Get("verified_access_group_id").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("load_balancer_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.LoadBalancerOptions = expandCreateVerifiedAccessEndpointLoadBalancerOptions(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("network_interface_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.NetworkInterfaceOptions = expandCreateVerifiedAccessEndpointEniOptions(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("security_group_ids"); ok && v.(*schema.Set).Len() > 0 {
		input.SecurityGroupIds = flex.ExpandStringValueSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("sse_specification"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.SseSpecification = expandCreateVerifiedAccessEndpointSseSpecification(v.([]interface{})[0].(map[string]interface{}))
	}

	output, err := conn.CreateVerifiedAccessEndpoint(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Verified Access Endpoint: %s", err)
	}

	d.SetId(aws.ToString(output.VerifiedAccessEndpoint.VerifiedAccessEndpointId))

	if _, err := WaitVerifiedAccessEndpointCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Verified Access Endpoint (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceVerifiedAccessEndpointRead(ctx, d, meta)...)
}

func resourceVerifiedAccessEndpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	ep, err := FindVerifiedAccessEndpointByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] EC2 Verified Access Endpoint (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Verified Access Endpoint (%s): %s", d.Id(), err)
	}

	d.Set("application_domain", ep.ApplicationDomain)
	d.Set("attachment_type", ep.AttachmentType)
	d.Set("description", ep.Description)
	d.Set("device_validation_domain", ep.DeviceValidationDomain)
	d.Set("domain_certificate_arn", ep.DomainCertificateArn)
	d.Set("endpoint_domain", ep.EndpointDomain)
	d.Set("endpoint_type", ep.EndpointType)
	if err := d.Set("load_balancer_options", flattenVerifiedAccessEndpointLoadBalancerOptions(ep.LoadBalancerOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting load_balancer_options: %s", err)
	}
	if err := d.Set("network_interface_options", flattenVerifiedAccessEndpointEniOptions(ep.NetworkInterfaceOptions)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting network_interface_options: %s", err)
	}
	d.Set("security_group_ids", aws.StringSlice(ep.SecurityGroupIds))
	if err := d.Set("sse_specification", flattenVerifiedAccessSseSpecificationRequest(ep.SseSpecification)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting sse_specification: %s", err)
	}
	d.Set("verified_access_group_id", ep.VerifiedAccessGroupId)
	d.Set("verified_access_instance_id", ep.VerifiedAccessInstanceId)

	return diags
}

func resourceVerifiedAccessEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	if d.HasChangesExcept("policy_document", "tags", "tags_all") {
		input := &ec2.ModifyVerifiedAccessEndpointInput{
			VerifiedAccessEndpointId: aws.String(d.Id()),
		}

		if d.HasChanges("description") {
			input.Description = aws.String(d.Get("description").(string))
		}

		if d.HasChanges("load_balancer_options") {
			if v, ok := d.GetOk("load_balancer_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.LoadBalancerOptions = expandModifyVerifiedAccessEndpointLoadBalancerOptions(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChanges("network_interface_options") {
			if v, ok := d.GetOk("network_interface_options"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
				input.NetworkInterfaceOptions = expandModifyVerifiedAccessEndpointEniOptions(v.([]interface{})[0].(map[string]interface{}))
			}
		}

		if d.HasChanges("verified_access_group_id") {
			input.VerifiedAccessGroupId = aws.String(d.Get("verified_access_group_id").(string))
		}

		_, err := conn.ModifyVerifiedAccessEndpoint(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Verified Access Endpoint (%s): %s", d.Id(), err)
		}

		if _, err := WaitVerifiedAccessEndpointUpdated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for Verified Access Endpoint (%s) update: %s", d.Id(), err)
		}
	}

	return append(diags, resourceVerifiedAccessEndpointRead(ctx, d, meta)...)
}

func resourceVerifiedAccessEndpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	log.Printf("[INFO] Deleting Verified Access Endpoint: %s", d.Id())
	_, err := conn.DeleteVerifiedAccessEndpoint(ctx, &ec2.DeleteVerifiedAccessEndpointInput{
		VerifiedAccessEndpointId: aws.String(d.Id()),
	})

	if tfawserr.ErrCodeEquals(err, errCodeInvalidVerifiedAccessEndpointIdNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Verified Access Endpoint (%s): %s", d.Id(), err)
	}

	if _, err := WaitVerifiedAccessEndpointDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Verified Access Endpoint (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func flattenVerifiedAccessEndpointLoadBalancerOptions(apiObject *types.VerifiedAccessEndpointLoadBalancerOptions) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfmap := map[string]interface{}{}

	if v := apiObject.LoadBalancerArn; v != nil {
		tfmap["load_balancer_arn"] = aws.ToString(v)
	}

	if v := apiObject.Port; v != nil {
		tfmap["port"] = aws.ToInt32(v)
	}

	if v := apiObject.Protocol; v != "" {
		tfmap["protocol"] = v
	}

	if v := apiObject.SubnetIds; v != nil {
		tfmap["subnet_ids"] = aws.StringSlice(v)
	}

	return []interface{}{tfmap}
}

func flattenVerifiedAccessEndpointEniOptions(apiObject *types.VerifiedAccessEndpointEniOptions) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfmap := map[string]interface{}{}

	if v := apiObject.NetworkInterfaceId; v != nil {
		tfmap["network_interface_id"] = aws.ToString(v)
	}

	if v := apiObject.Port; v != nil {
		tfmap["port"] = aws.ToInt32(v)
	}

	if v := apiObject.Protocol; v != "" {
		tfmap["protocol"] = v
	}

	return []interface{}{tfmap}
}

func flattenVerifiedAccessSseSpecificationRequest(apiObject *types.VerifiedAccessSseSpecificationResponse) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfmap := map[string]interface{}{}

	if v := apiObject.CustomerManagedKeyEnabled; v != nil {
		tfmap["customer_managed_key_enabled"] = aws.ToBool(v)
	}

	if v := apiObject.KmsKeyArn; v != nil {
		tfmap["kms_key_arn"] = aws.ToString(v)
	}

	return []interface{}{tfmap}
}

func expandCreateVerifiedAccessEndpointLoadBalancerOptions(tfMap map[string]interface{}) *types.CreateVerifiedAccessEndpointLoadBalancerOptions {
	if tfMap == nil {
		return nil
	}

	apiobject := &types.CreateVerifiedAccessEndpointLoadBalancerOptions{}

	if v, ok := tfMap["load_balancer_arn"].(string); ok && v != "" {
		apiobject.LoadBalancerArn = aws.String(v)
	}

	if v, ok := tfMap["port"].(int); ok {
		apiobject.Port = aws.Int32(int32(v))
	}

	if v, ok := tfMap["protocol"].(string); ok && v != "" {
		apiobject.Protocol = types.VerifiedAccessEndpointProtocol(v)
	}

	if v, ok := tfMap["subnet_ids"].(*schema.Set); ok && v.Len() > 0 {
		apiobject.SubnetIds = flex.ExpandStringValueSet(v)
	}

	return apiobject
}

func expandCreateVerifiedAccessEndpointEniOptions(tfMap map[string]interface{}) *types.CreateVerifiedAccessEndpointEniOptions {
	if tfMap == nil {
		return nil
	}

	apiobject := &types.CreateVerifiedAccessEndpointEniOptions{}

	if v, ok := tfMap["network_interface_id"].(string); ok && v != "" {
		apiobject.NetworkInterfaceId = aws.String(v)
	}

	if v, ok := tfMap["port"].(int); ok {
		apiobject.Port = aws.Int32(int32(v))
	}

	if v, ok := tfMap["protocol"].(string); ok && v != "" {
		apiobject.Protocol = types.VerifiedAccessEndpointProtocol(v)
	}
	return apiobject
}

func expandModifyVerifiedAccessEndpointLoadBalancerOptions(tfMap map[string]interface{}) *types.ModifyVerifiedAccessEndpointLoadBalancerOptions {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.ModifyVerifiedAccessEndpointLoadBalancerOptions{}

	if v, ok := tfMap["port"].(int); ok {
		apiObject.Port = aws.Int32(int32(v))
	}

	if v, ok := tfMap["protocol"].(string); ok && v != "" {
		apiObject.Protocol = types.VerifiedAccessEndpointProtocol(v)
	}

	if v, ok := tfMap["subnet_ids"]; ok {
		apiObject.SubnetIds = flex.ExpandStringValueList(v.([]interface{}))
	}

	return apiObject
}

func expandModifyVerifiedAccessEndpointEniOptions(tfMap map[string]interface{}) *types.ModifyVerifiedAccessEndpointEniOptions {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.ModifyVerifiedAccessEndpointEniOptions{}

	if v, ok := tfMap["port"].(int); ok {
		apiObject.Port = aws.Int32(int32(v))
	}

	if v, ok := tfMap["protocol"].(string); ok && v != "" {
		apiObject.Protocol = types.VerifiedAccessEndpointProtocol(v)
	}
	return apiObject
}

func expandCreateVerifiedAccessEndpointSseSpecification(tfMap map[string]interface{}) *types.VerifiedAccessSseSpecificationRequest {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.VerifiedAccessSseSpecificationRequest{}

	if v, ok := tfMap["customer_managed_key_enabled"].(bool); ok {
		apiObject.CustomerManagedKeyEnabled = aws.Bool(v)
	}

	if v, ok := tfMap["kms_key_arn"].(string); ok && v != "" {
		apiObject.KmsKeyArn = aws.String(v)
	}
	return apiObject
}
