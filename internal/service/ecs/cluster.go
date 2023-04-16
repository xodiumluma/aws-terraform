package ecs

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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

// @SDKResource("aws_ecs_cluster", name="Cluster")
// @Tags(identifierAttribute="id")
func ResourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceClusterCreate,
		ReadWithoutTimeout:   resourceClusterRead,
		UpdateWithoutTimeout: resourceClusterUpdate,
		DeleteWithoutTimeout: resourceClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterImport,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"capacity_providers": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Use the aws_ecs_cluster_capacity_providers resource instead",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"execute_command_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kms_key_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"log_configuration": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cloud_watch_encryption_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"cloud_watch_log_group_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"s3_bucket_encryption_enabled": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"s3_bucket_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"s3_key_prefix": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"logging": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(ecs.ExecuteCommandLogging_Values(), false),
									},
								},
							},
						},
					},
				},
			},
			"default_capacity_provider_strategy": {
				Type:       schema.TypeSet,
				Optional:   true,
				Computed:   true,
				Deprecated: "Use the aws_ecs_cluster_capacity_providers resource instead",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 100000),
						},
						"capacity_provider": {
							Type:     schema.TypeString,
							Required: true,
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 1000),
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateClusterName,
			},
			"service_connect_defaults": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
			},
			"setting": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(ecs.ClusterSettingName_Values(), false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},
	}
}

func resourceClusterImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	d.SetId(arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		AccountID: meta.(*conns.AWSClient).AccountID,
		Service:   "ecs",
		Resource:  fmt.Sprintf("cluster/%s", d.Id()),
	}.String())
	return []*schema.ResourceData{d}, nil
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECSConn()

	clusterName := d.Get("name").(string)
	log.Printf("[DEBUG] Creating ECS cluster %s", clusterName)

	input := &ecs.CreateClusterInput{
		ClusterName:                     aws.String(clusterName),
		DefaultCapacityProviderStrategy: expandCapacityProviderStrategy(d.Get("default_capacity_provider_strategy").(*schema.Set)),
		Tags:                            GetTagsIn(ctx),
	}

	if v, ok := d.GetOk("capacity_providers"); ok {
		input.CapacityProviders = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("service_connect_defaults"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.ServiceConnectDefaults = expandClusterServiceConnectDefaultsRequest(v.([]interface{})[0].(map[string]interface{}))
	}

	if v, ok := d.GetOk("setting"); ok {
		input.Settings = expandClusterSettings(v.(*schema.Set))
	}

	if v, ok := d.GetOk("configuration"); ok && len(v.([]interface{})) > 0 {
		input.Configuration = expandClusterConfiguration(v.([]interface{}))
	}

	// CreateCluster will create the ECS IAM Service Linked Role on first ECS provision
	// This process does not complete before the initial API call finishes.
	out, err := retryClusterCreate(ctx, conn, input)

	// Some partitions (i.e., ISO) may not support tag-on-create
	if input.Tags != nil && verify.ErrorISOUnsupported(conn.PartitionID, err) {
		log.Printf("[WARN] ECS tagging failed creating Cluster (%s) with tags: %s. Trying create without tags.", clusterName, err)
		input.Tags = nil

		out, err = retryClusterCreate(ctx, conn, input)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating ECS Cluster (%s): %s", clusterName, err)
	}

	log.Printf("[DEBUG] ECS cluster %s created", aws.StringValue(out.Cluster.ClusterArn))

	d.SetId(aws.StringValue(out.Cluster.ClusterArn))

	if _, err := waitClusterAvailable(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for ECS Cluster (%s) to become Available while creating: %s", d.Id(), err)
	}

	// Some partitions (i.e., ISO) may not support tag-on-create, attempt tag after create
	if tags := KeyValueTags(ctx, GetTagsIn(ctx)); input.Tags == nil && len(tags) > 0 {
		err := UpdateTags(ctx, conn, d.Id(), nil, tags)

		if v, ok := d.GetOk("tags"); (!ok || len(v.(map[string]interface{})) == 0) && verify.ErrorISOUnsupported(conn.PartitionID, err) {
			// If default tags only, log and continue. Otherwise, error.
			log.Printf("[WARN] ECS tagging failed adding tags after create for Cluster (%s): %s", d.Id(), err)
			return append(diags, resourceClusterRead(ctx, d, meta)...)
		}

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "ECS tagging failed adding tags after create for Cluster (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceClusterRead(ctx, d, meta)...)
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECSConn()

	outputRaw, err := tfresource.RetryWhenNewResourceNotFound(ctx, clusterReadTimeout, func() (interface{}, error) {
		return FindClusterByNameOrARN(ctx, conn, d.Id())
	}, d.IsNewResource())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] ECS Cluster (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading ECS Cluster (%s): %s", d.Id(), err)
	}

	cluster := outputRaw.(*ecs.Cluster)
	d.Set("arn", cluster.ClusterArn)
	d.Set("name", cluster.ClusterName)

	if err := d.Set("capacity_providers", aws.StringValueSlice(cluster.CapacityProviders)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting capacity_providers: %s", err)
	}
	if err := d.Set("default_capacity_provider_strategy", flattenCapacityProviderStrategy(cluster.DefaultCapacityProviderStrategy)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting default_capacity_provider_strategy: %s", err)
	}

	if cluster.ServiceConnectDefaults != nil {
		if err := d.Set("service_connect_defaults", []interface{}{flattenClusterServiceConnectDefaults(cluster.ServiceConnectDefaults)}); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting service_connect_defaults: %s", err)
		}
	} else {
		d.Set("service_connect_defaults", nil)
	}

	if err := d.Set("setting", flattenClusterSettings(cluster.Settings)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting setting: %s", err)
	}

	if cluster.Configuration != nil {
		if err := d.Set("configuration", flattenClusterConfiguration(cluster.Configuration)); err != nil {
			return sdkdiag.AppendErrorf(diags, "setting configuration: %s", err)
		}
	}

	SetTagsOut(ctx, cluster.Tags)

	return diags
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ECSConn()

	if d.HasChanges("setting", "configuration", "service_connect_defaults") {
		input := &ecs.UpdateClusterInput{
			Cluster: aws.String(d.Id()),
		}

		if v, ok := d.GetOk("setting"); ok {
			input.Settings = expandClusterSettings(v.(*schema.Set))
		}

		if v, ok := d.GetOk("configuration"); ok && len(v.([]interface{})) > 0 {
			input.Configuration = expandClusterConfiguration(v.([]interface{}))
		}

		if v, ok := d.GetOk("service_connect_defaults"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
			input.ServiceConnectDefaults = expandClusterServiceConnectDefaultsRequest(v.([]interface{})[0].(map[string]interface{}))
		}

		_, err := conn.UpdateClusterWithContext(ctx, input)

		if err != nil {
			return diag.Errorf("updating ECS Cluster (%s): %s", d.Id(), err)
		}

		if _, err := waitClusterAvailable(ctx, conn, d.Id()); err != nil {
			return diag.Errorf("waiting for ECS Cluster (%s) update: %s", d.Id(), err)
		}
	}

	if d.HasChanges("capacity_providers", "default_capacity_provider_strategy") {
		input := ecs.PutClusterCapacityProvidersInput{
			Cluster:                         aws.String(d.Id()),
			CapacityProviders:               flex.ExpandStringSet(d.Get("capacity_providers").(*schema.Set)),
			DefaultCapacityProviderStrategy: expandCapacityProviderStrategy(d.Get("default_capacity_provider_strategy").(*schema.Set)),
		}

		err := retryClusterCapacityProvidersPut(ctx, conn, &input)

		if err != nil {
			return diag.Errorf("updating ECS Cluster (%s) capacity providers: %s", d.Id(), err)
		}

		if _, err := waitClusterAvailable(ctx, conn, d.Id()); err != nil {
			return diag.Errorf("waiting for ECS Cluster (%s) capacity providers update: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ECSConn()

	log.Printf("[DEBUG] Deleting ECS cluster %s", d.Id())
	input := &ecs.DeleteClusterInput{
		Cluster: aws.String(d.Id()),
	}
	err := retry.RetryContext(ctx, clusterDeleteTimeout, func() *retry.RetryError {
		_, err := conn.DeleteClusterWithContext(ctx, input)

		if err == nil {
			log.Printf("[DEBUG] ECS cluster %s deleted", d.Id())
			return nil
		}

		if tfawserr.ErrCodeEquals(err, "ClusterContainsContainerInstancesException") {
			log.Printf("[TRACE] Retrying ECS cluster %q deletion after %s", d.Id(), err)
			return retry.RetryableError(err)
		}
		if tfawserr.ErrCodeEquals(err, "ClusterContainsServicesException") {
			log.Printf("[TRACE] Retrying ECS cluster %q deletion after %s", d.Id(), err)
			return retry.RetryableError(err)
		}
		if tfawserr.ErrCodeEquals(err, "ClusterContainsTasksException") {
			log.Printf("[TRACE] Retrying ECS cluster %q deletion after %s", d.Id(), err)
			return retry.RetryableError(err)
		}
		if tfawserr.ErrCodeEquals(err, ecs.ErrCodeUpdateInProgressException) {
			log.Printf("[TRACE] Retrying ECS cluster %q deletion after %s", d.Id(), err)
			return retry.RetryableError(err)
		}
		return retry.NonRetryableError(err)
	})
	if tfresource.TimedOut(err) {
		_, err = conn.DeleteClusterWithContext(ctx, input)
	}
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting ECS cluster: %s", err)
	}

	if _, err := waitClusterDeleted(ctx, conn, d.Id()); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for ECS Cluster (%s) to become Deleted: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] ECS cluster %q deleted", d.Id())
	return diags
}

func retryClusterCreate(ctx context.Context, conn *ecs.ECS, input *ecs.CreateClusterInput) (*ecs.CreateClusterOutput, error) {
	var output *ecs.CreateClusterOutput
	err := retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		var err error
		output, err = conn.CreateClusterWithContext(ctx, input)

		if tfawserr.ErrMessageContains(err, ecs.ErrCodeInvalidParameterException, "Unable to assume the service linked role") {
			return retry.RetryableError(err)
		}

		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if tfresource.TimedOut(err) {
		output, err = conn.CreateClusterWithContext(ctx, input)
	}

	return output, err
}

func expandClusterSettings(configured *schema.Set) []*ecs.ClusterSetting {
	list := configured.List()
	if len(list) == 0 {
		return nil
	}

	settings := make([]*ecs.ClusterSetting, 0, len(list))

	for _, raw := range list {
		data := raw.(map[string]interface{})

		setting := &ecs.ClusterSetting{
			Name:  aws.String(data["name"].(string)),
			Value: aws.String(data["value"].(string)),
		}

		settings = append(settings, setting)
	}

	return settings
}

func expandClusterServiceConnectDefaultsRequest(tfMap map[string]interface{}) *ecs.ClusterServiceConnectDefaultsRequest {
	if tfMap == nil {
		return nil
	}

	apiObject := &ecs.ClusterServiceConnectDefaultsRequest{}

	if v, ok := tfMap["namespace"].(string); ok && v != "" {
		apiObject.Namespace = aws.String(v)
	}

	return apiObject
}

func flattenClusterServiceConnectDefaults(apiObject *ecs.ClusterServiceConnectDefaults) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Namespace; v != nil {
		tfMap["namespace"] = aws.StringValue(v)
	}

	return tfMap
}

func flattenClusterSettings(list []*ecs.ClusterSetting) []map[string]interface{} {
	if len(list) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(list))
	for _, setting := range list {
		l := map[string]interface{}{
			"name":  aws.StringValue(setting.Name),
			"value": aws.StringValue(setting.Value),
		}

		result = append(result, l)
	}
	return result
}

func flattenClusterConfiguration(apiObject *ecs.ClusterConfiguration) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if apiObject.ExecuteCommandConfiguration != nil {
		tfMap["execute_command_configuration"] = flattenClusterConfigurationExecuteCommandConfiguration(apiObject.ExecuteCommandConfiguration)
	}
	return []interface{}{tfMap}
}

func flattenClusterConfigurationExecuteCommandConfiguration(apiObject *ecs.ExecuteCommandConfiguration) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if apiObject.KmsKeyId != nil {
		tfMap["kms_key_id"] = aws.StringValue(apiObject.KmsKeyId)
	}

	if apiObject.LogConfiguration != nil {
		tfMap["log_configuration"] = flattenClusterConfigurationExecuteCommandConfigurationLogConfiguration(apiObject.LogConfiguration)
	}

	if apiObject.Logging != nil {
		tfMap["logging"] = aws.StringValue(apiObject.Logging)
	}

	return []interface{}{tfMap}
}

func flattenClusterConfigurationExecuteCommandConfigurationLogConfiguration(apiObject *ecs.ExecuteCommandLogConfiguration) []interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	tfMap["cloud_watch_encryption_enabled"] = aws.BoolValue(apiObject.CloudWatchEncryptionEnabled)
	tfMap["s3_bucket_encryption_enabled"] = aws.BoolValue(apiObject.S3EncryptionEnabled)

	if apiObject.CloudWatchLogGroupName != nil {
		tfMap["cloud_watch_log_group_name"] = aws.StringValue(apiObject.CloudWatchLogGroupName)
	}

	if apiObject.S3BucketName != nil {
		tfMap["s3_bucket_name"] = aws.StringValue(apiObject.S3BucketName)
	}

	if apiObject.S3KeyPrefix != nil {
		tfMap["s3_key_prefix"] = aws.StringValue(apiObject.S3KeyPrefix)
	}

	return []interface{}{tfMap}
}

func expandClusterConfiguration(nc []interface{}) *ecs.ClusterConfiguration {
	if len(nc) == 0 {
		return &ecs.ClusterConfiguration{}
	}
	raw := nc[0].(map[string]interface{})

	config := &ecs.ClusterConfiguration{}
	if v, ok := raw["execute_command_configuration"].([]interface{}); ok && len(v) > 0 {
		config.ExecuteCommandConfiguration = expandClusterConfigurationExecuteCommandConfiguration(v)
	}

	return config
}

func expandClusterConfigurationExecuteCommandConfiguration(nc []interface{}) *ecs.ExecuteCommandConfiguration {
	if len(nc) == 0 {
		return &ecs.ExecuteCommandConfiguration{}
	}
	raw := nc[0].(map[string]interface{})

	config := &ecs.ExecuteCommandConfiguration{}
	if v, ok := raw["log_configuration"].([]interface{}); ok && len(v) > 0 {
		config.LogConfiguration = expandClusterConfigurationExecuteCommandLogConfiguration(v)
	}

	if v, ok := raw["kms_key_id"].(string); ok && v != "" {
		config.KmsKeyId = aws.String(v)
	}

	if v, ok := raw["logging"].(string); ok && v != "" {
		config.Logging = aws.String(v)
	}

	return config
}

func expandClusterConfigurationExecuteCommandLogConfiguration(nc []interface{}) *ecs.ExecuteCommandLogConfiguration {
	if len(nc) == 0 {
		return &ecs.ExecuteCommandLogConfiguration{}
	}
	raw := nc[0].(map[string]interface{})

	config := &ecs.ExecuteCommandLogConfiguration{}

	if v, ok := raw["cloud_watch_log_group_name"].(string); ok && v != "" {
		config.CloudWatchLogGroupName = aws.String(v)
	}

	if v, ok := raw["s3_bucket_name"].(string); ok && v != "" {
		config.S3BucketName = aws.String(v)
	}

	if v, ok := raw["s3_key_prefix"].(string); ok && v != "" {
		config.S3KeyPrefix = aws.String(v)
	}

	if v, ok := raw["cloud_watch_encryption_enabled"].(bool); ok {
		config.CloudWatchEncryptionEnabled = aws.Bool(v)
	}

	if v, ok := raw["s3_bucket_encryption_enabled"].(bool); ok {
		config.S3EncryptionEnabled = aws.Bool(v)
	}

	return config
}
