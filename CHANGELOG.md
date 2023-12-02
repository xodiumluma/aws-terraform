## 5.30.0 (Unreleased)

FEATURES:

* **New Data Source:** `aws_codeguruprofiler_profiling_group` ([#34672](https://github.com/hashicorp/terraform-provider-aws/issues/34672))
* **New Data Source:** `aws_lb_trust_store` ([#34584](https://github.com/hashicorp/terraform-provider-aws/issues/34584))
* **New Resource:** `aws_codeguruprofiler_profiling_group` ([#34672](https://github.com/hashicorp/terraform-provider-aws/issues/34672))
* **New Resource:** `aws_lb_trust_store` ([#34584](https://github.com/hashicorp/terraform-provider-aws/issues/34584))
* **New Resource:** `aws_lb_trust_store_revocation` ([#34584](https://github.com/hashicorp/terraform-provider-aws/issues/34584))

ENHANCEMENTS:

* data-source/aws_lb_listener: Add `mutual_authentication` attribute ([#34584](https://github.com/hashicorp/terraform-provider-aws/issues/34584))
* resource/aws_cloudwatch_log_group: Add `log_group_class` argument ([#34679](https://github.com/hashicorp/terraform-provider-aws/issues/34679))
* resource/aws_lb_listener: Add `mutual_authentication` configuration block ([#34584](https://github.com/hashicorp/terraform-provider-aws/issues/34584))

## 5.29.0 (November 30, 2023)

FEATURES:

* **New Resource:** `aws_docdbelastic_cluster` ([#31033](https://github.com/hashicorp/terraform-provider-aws/issues/31033))
* **New Resource:** `aws_eks_pod_identity_association` ([#34566](https://github.com/hashicorp/terraform-provider-aws/issues/34566))

ENHANCEMENTS:

* resource/aws_docdb_cluster: Add `storage_type` argument ([#34637](https://github.com/hashicorp/terraform-provider-aws/issues/34637))
* resource/aws_neptune_parameter_group: Add `name_prefix` argument ([#34500](https://github.com/hashicorp/terraform-provider-aws/issues/34500))

BUG FIXES:

* resource/aws_networkmanager_attachment_accepter: Now revokes attachment on deletion for VPC Attachments ([#34547](https://github.com/hashicorp/terraform-provider-aws/issues/34547))
* resource/aws_networkmanager_vpc_attachment: Fixes error when modifying `options` fields while waiting for acceptance ([#34547](https://github.com/hashicorp/terraform-provider-aws/issues/34547))
* resource/aws_networkmanager_vpc_attachment: Fixes error where VPC Attachments waiting for acceptance could not be deleted ([#34547](https://github.com/hashicorp/terraform-provider-aws/issues/34547))
* resource/aws_s3_directory_bucket: Fix `NotImplemented: This bucket does not support Object Versioning` errors on resource Delete when `force_destroy` is `true` ([#34647](https://github.com/hashicorp/terraform-provider-aws/issues/34647))

## 5.28.0 (November 29, 2023)

FEATURES:

* **New Data Source:** `aws_s3_directory_buckets` ([#34612](https://github.com/hashicorp/terraform-provider-aws/issues/34612))
* **New Resource:** `aws_s3_directory_bucket` ([#34612](https://github.com/hashicorp/terraform-provider-aws/issues/34612))

ENHANCEMENTS:

* resource/aws_s3control_access_grants_instance: Add `identity_center_arn` argument and `identity_center_application_arn` attribute ([#34582](https://github.com/hashicorp/terraform-provider-aws/issues/34582))

BUG FIXES:

* resource/aws_elaticache_replication_group: Fix regression caused by the introduction of the `auth_token_update_strategy` argument with a default value ([#34600](https://github.com/hashicorp/terraform-provider-aws/issues/34600))

## 5.27.0 (November 27, 2023)

NOTES:

* provider: This release includes an update to the AWS SDK for Go v2 with breaking type changes to several services: `internetmonitor`, `ivschat`, `pipes`, and `s3`. These changes primarily affect how arguments with default values are serialized for outbound requests, changing scalar types to pointers. See [this AWS SDK for Go V2 issue](https://github.com/aws/aws-sdk-go-v2/issues/2162) for additional context. The corresponding provider changes should make this breakfix transparent to users, but as with any breaking change there is the potential for missed edge cases. If errors are observed in the impacted resources, please link to this dependency update pull request in the bug report ([#34476](https://github.com/hashicorp/terraform-provider-aws/issues/34476))

FEATURES:

* **New Data Source:** `aws_emr_supported_instance_types` ([#34481](https://github.com/hashicorp/terraform-provider-aws/issues/34481))
* **New Resource:** `aws_apprunner_default_auto_scaling_configuration_version` ([#34292](https://github.com/hashicorp/terraform-provider-aws/issues/34292))
* **New Resource:** `aws_lexv2models_bot_version` ([#33858](https://github.com/hashicorp/terraform-provider-aws/issues/33858))
* **New Resource:** `aws_s3control_access_grant` ([#34564](https://github.com/hashicorp/terraform-provider-aws/issues/34564))
* **New Resource:** `aws_s3control_access_grants_instance` ([#34564](https://github.com/hashicorp/terraform-provider-aws/issues/34564))
* **New Resource:** `aws_s3control_access_grants_instance_resource_policy` ([#34564](https://github.com/hashicorp/terraform-provider-aws/issues/34564))
* **New Resource:** `aws_s3control_access_grants_location` ([#34564](https://github.com/hashicorp/terraform-provider-aws/issues/34564))

ENHANCEMENTS:

* resource/aws_apprunner_auto_scaling_configuration_version: Add `has_associated_service` and `is_default` attributes ([#34292](https://github.com/hashicorp/terraform-provider-aws/issues/34292))
* resource/aws_apprunner_service: Add `network_configuration.ip_address_type` argument ([#34292](https://github.com/hashicorp/terraform-provider-aws/issues/34292))
* resource/aws_apprunner_service: Add `source_configuration.code_repository.source_directory` argument to support monorepos ([#34292](https://github.com/hashicorp/terraform-provider-aws/issues/34292))
* resource/aws_apprunner_service: Allow `health_check_configuration` to be updated in-place ([#34292](https://github.com/hashicorp/terraform-provider-aws/issues/34292))
* resource/aws_cloudwatch_event_rule: Add `state` parameter and deprecate `is_enabled` parameter ([#34510](https://github.com/hashicorp/terraform-provider-aws/issues/34510))
* resource/aws_elaticache_replication_group: Add `auth_token_update_strategy` argument ([#34460](https://github.com/hashicorp/terraform-provider-aws/issues/34460))
* resource/aws_lambda_function: Add support for `java21` `runtime` value ([#34476](https://github.com/hashicorp/terraform-provider-aws/issues/34476))
* resource/aws_lambda_function: Add support for `python3.12` `runtime` value ([#34533](https://github.com/hashicorp/terraform-provider-aws/issues/34533))
* resource/aws_lambda_layer_version: Add support for `java21` `compatible_runtimes` value ([#34476](https://github.com/hashicorp/terraform-provider-aws/issues/34476))
* resource/aws_lambda_layer_version: Add support for `python3.12` `compatible_runtimes` value ([#34533](https://github.com/hashicorp/terraform-provider-aws/issues/34533))
* resource/aws_s3_bucket_logging: Add `target_object_key_format` configuration block to support [automatic date-based partitioning](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerLogs.html#server-access-logging-overview) ([#34504](https://github.com/hashicorp/terraform-provider-aws/issues/34504))

BUG FIXES:

* resource/aws_appflow_flow: Fix `InvalidParameter: 2 validation error(s) found` error when `destination_flow_config` or `task` is updated ([#34456](https://github.com/hashicorp/terraform-provider-aws/issues/34456))
* resource/aws_appflow_flow: Fix `interface conversion: interface {} is nil, not map[string]interface {}` panic ([#34456](https://github.com/hashicorp/terraform-provider-aws/issues/34456))
* resource/aws_apprunner_service: Correctly set `service_url` for private services ([#34292](https://github.com/hashicorp/terraform-provider-aws/issues/34292))
* resource/aws_glue_trigger: Fix `ConcurrentModificationException: Workflow <workflowName> was modified while adding trigger <triggerName>` errors ([#34530](https://github.com/hashicorp/terraform-provider-aws/issues/34530))
* resource/aws_lb_target_group: Adds plan- and apply-time validation for invalid parameter combinations ([#34488](https://github.com/hashicorp/terraform-provider-aws/issues/34488))
* resource/aws_lexv2_bot_locale: Fix `voice_settings.engine` validation, value conversion errors ([#34532](https://github.com/hashicorp/terraform-provider-aws/issues/34532))
* resource/aws_lexv2models_bot: Properly send `type` argument on create and update when configured ([#34524](https://github.com/hashicorp/terraform-provider-aws/issues/34524))
* resource/aws_pipes_pipe: Fix error when zero value is sent to `source_parameters` on update ([#34487](https://github.com/hashicorp/terraform-provider-aws/issues/34487))

## 5.26.0 (November 16, 2023)

FEATURES:

* **New Data Source:** `aws_iot_registration_code` ([#15098](https://github.com/hashicorp/terraform-provider-aws/issues/15098))
* **New Resource:** `aws_iot_billing_group` ([#31237](https://github.com/hashicorp/terraform-provider-aws/issues/31237))
* **New Resource:** `aws_iot_ca_certificate` ([#15098](https://github.com/hashicorp/terraform-provider-aws/issues/15098))
* **New Resource:** `aws_iot_event_configurations` ([#31237](https://github.com/hashicorp/terraform-provider-aws/issues/31237))

ENHANCEMENTS:

* data-source/aws_autoscaling_group: Add `instance_maintenance_policy` attribute ([#34430](https://github.com/hashicorp/terraform-provider-aws/issues/34430))
* provider: Adds `https_proxy` and `no_proxy` parameters. ([#34243](https://github.com/hashicorp/terraform-provider-aws/issues/34243))
* resource/aws_autoscaling_group: Add `instance_maintenance_policy` configuration block ([#34430](https://github.com/hashicorp/terraform-provider-aws/issues/34430))
* resource/aws_finspace_kx_cluster: Increase default create and update timeouts to 4 hours to allow for increased startup times with large volumes of cached data ([#34398](https://github.com/hashicorp/terraform-provider-aws/issues/34398))
* resource/aws_finspace_kx_environment: Increase default delete timeout to 75 minutes ([#34398](https://github.com/hashicorp/terraform-provider-aws/issues/34398))
* resource/aws_iam_group_policy_attachment: Add plan-time validation of `policy_arn` ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iam_policy_attachment: Add plan-time validation of `policy_arn` ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iam_role_policy_attachment: Add plan-time validation of `policy_arn` ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iam_user_policy_attachment: Add plan-time validation of `policy_arn` ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iot_ca_certificate: Add `ca_certificate_id` attribute ([#15098](https://github.com/hashicorp/terraform-provider-aws/issues/15098))
* resource/aws_iot_policy: Add configurable timeouts ([#34329](https://github.com/hashicorp/terraform-provider-aws/issues/34329))
* resource/aws_iot_policy: When updating the resource, delete the oldest non-default version of the policy if creating a new version would exceed the maximum number of versions (5) ([#34329](https://github.com/hashicorp/terraform-provider-aws/issues/34329))
* resource/aws_lambda_function: Add support for `nodejs20.x` and `provided.al2023` `runtime` values ([#34401](https://github.com/hashicorp/terraform-provider-aws/issues/34401))
* resource/aws_lambda_layer_version: Add support for `nodejs20.x` and `provided.al2023` `compatible_runtimes` values ([#34401](https://github.com/hashicorp/terraform-provider-aws/issues/34401))
* resource/aws_quicksight_analysis: Add `definition.sheets.visuals.kpi_visual.chart_configuration.kpi_options.sparkline` attribute ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Add `definition.sheets.visuals.kpi_visual.chart_configuration.kpi_options.visual_layout_options` attribute ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Add `number_display_format_configuration` and `percentage_display_format_configuration` to nested `numeric_format_configuration` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Add `definition.sheets.visuals.kpi_visual.chart_configuration.kpi_options.sparkline` attribute ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Add `definition.sheets.visuals.kpi_visual.chart_configuration.kpi_options.visual_layout_options` attribute ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Add `number_display_format_configuration` and `percentage_display_format_configuration` to nested `numeric_format_configuration` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Add `definition.sheets.visuals.kpi_visual.chart_configuration.kpi_options.sparkline` attribute ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Add `definition.sheets.visuals.kpi_visual.chart_configuration.kpi_options.visual_layout_options` attribute ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Add `number_display_format_configuration` and `percentage_display_format_configuration` to nested `numeric_format_configuration` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_rds_cluster: Add `delete_automated_backups` argument ([#34309](https://github.com/hashicorp/terraform-provider-aws/issues/34309))

BUG FIXES:

* resource/aws_chime_voice_connector: Fix `read` error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chime_voice_connector_group: Fix `read` error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chime_voice_connector_logging: Fix `read` error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chime_voice_connector_origination: Fix `read` error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chime_voice_connector_termination: Fix `read` error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chime_voice_connector_termination_credentials: Fix `read` error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chimesdkmediapipelines_media_insights_pipeline_configuration: Fix eventual consistency error when resource is not created in `us-east-1` ([#34334](https://github.com/hashicorp/terraform-provider-aws/issues/34334))
* resource/aws_chimesdkvoice_sip_media_application: Fix eventual consistency errors when not using `us-east-1` ([#34426](https://github.com/hashicorp/terraform-provider-aws/issues/34426))
* resource/aws_chimesdkvoice_sip_rule: Fix eventual consistency errors when not using `us-east-1` ([#34426](https://github.com/hashicorp/terraform-provider-aws/issues/34426))
* resource/aws_elasticache_user: Fix `UserNotFound: ... is not available for tagging` errors on resource Read when there is a concurrent update to the user ([#34396](https://github.com/hashicorp/terraform-provider-aws/issues/34396))
* resource/aws_grafana_workspace_api_key: Change `key` to [`Sensitive`](https://developer.hashicorp.com/terraform/plugin/best-practices/sensitive-state#using-sensitive-flag-functionality) ([#34105](https://github.com/hashicorp/terraform-provider-aws/issues/34105))
* resource/aws_iam_group_policy_attachment: Retry `ConcurrentModificationException` errors on create and delete ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iam_policy_attachment: Retry `ConcurrentModificationException` errors on create and delete ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iam_role_policy_attachment: Retry `ConcurrentModificationException` errors on create and delete ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_iam_user_policy_attachment: Retry `ConcurrentModificationException` errors on create and delete ([#34378](https://github.com/hashicorp/terraform-provider-aws/issues/34378))
* resource/aws_inspector2_delegated_admin_account: Fix `errors: *target must be interface or implement error` panic ([#34424](https://github.com/hashicorp/terraform-provider-aws/issues/34424))
* resource/aws_inspector2_enabler: Fix `interface conversion: interface {} is nil, not map[string]inspector2.AccountResourceStatus` panic ([#34424](https://github.com/hashicorp/terraform-provider-aws/issues/34424))
* resource/aws_iot_ca_certificate: Change `ca_pem` and `certificate_pem` to [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) ([#15098](https://github.com/hashicorp/terraform-provider-aws/issues/15098))
* resource/aws_iot_policy: Retry `DeleteConflictException` errors on delete ([#34329](https://github.com/hashicorp/terraform-provider-aws/issues/34329))
* resource/aws_quicksight_analysis: Fix handling of the nested `number_scale`, `prefix`, and `suffix` integer arguments ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Fix handling of the nested `rolling_date` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Fix handling of the nested `select_all_options` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Fix handling of the nested `visual_ids` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Fixes to various optional blocks utilizing the shared column schema definition ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_analysis: Nested `column_index` and `row_index` arguments now properly handle zero values ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Fix handling of the nested `number_scale`, `prefix`, and `suffix` integer arguments ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Fix handling of the nested `rolling_date` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Fix handling of the nested `select_all_options` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Fix handling of the nested `visual_ids` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Fixes to various optional blocks utilizing the shared column schema definition ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_dashboard: Nested `column_index` and `row_index` arguments now properly handle zero values ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_data_set: Increase `permissions.actions` maximum item limit to 20, aligning with the AWS API limits ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_data_source: Set all parameters to update aws_quicksight_data_source ([#33061](https://github.com/hashicorp/terraform-provider-aws/issues/33061))
* resource/aws_quicksight_template: Fix handling of the nested `number_scale`, `prefix`, and `suffix` integer arguments ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Fix handling of the nested `rolling_date` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Fix handling of the nested `select_all_options` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Fix handling of the nested `visual_ids` argument ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Fixes to various optional blocks utilizing the shared column schema definition ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_quicksight_template: Nested `column_index` and `row_index` arguments now properly handle zero values ([#33931](https://github.com/hashicorp/terraform-provider-aws/issues/33931))
* resource/aws_sagemaker_user_profile: Change `default_user_settings.canvas_app_settings.identity_provider_oauth_settings` from TypeSet to TypeList, preventing `interface conversion: interface {} is *schema.Set, not []interface {}` panics ([#34418](https://github.com/hashicorp/terraform-provider-aws/issues/34418))
* resource/aws_synthetics_canary: Fix to properly suppress differences when `expression` is `rate(0 minutes)` ([#34084](https://github.com/hashicorp/terraform-provider-aws/issues/34084))
* resource/aws_vpn_connection: Fix `UnsupportedOperation: The tunnel inside ip version parameter is not currently supported in this region` error when creating connections in certain partitions and Regions ([#34420](https://github.com/hashicorp/terraform-provider-aws/issues/34420))

## 5.25.0 (November 10, 2023)

NOTES:

* resource/aws_cloudtrail: The resource's [import ID](https://developer.hashicorp.com/terraform/language/import#import-id) has changed from `name` to `arn` ([#30758](https://github.com/hashicorp/terraform-provider-aws/issues/30758))

FEATURES:

* **New Data Source:** `aws_apigatewayv2_vpc_link` ([#33974](https://github.com/hashicorp/terraform-provider-aws/issues/33974))
* **New Data Source:** `aws_athena_named_query` ([#24815](https://github.com/hashicorp/terraform-provider-aws/issues/24815))
* **New Data Source:** `aws_bedrock_foundation_model` ([#34148](https://github.com/hashicorp/terraform-provider-aws/issues/34148))
* **New Data Source:** `aws_bedrock_foundation_models` ([#34148](https://github.com/hashicorp/terraform-provider-aws/issues/34148))
* **New Resource:** `aws_athena_prepared_statement` ([#33417](https://github.com/hashicorp/terraform-provider-aws/issues/33417))
* **New Resource:** `aws_lexv2models_bot_locale` ([#33949](https://github.com/hashicorp/terraform-provider-aws/issues/33949))

ENHANCEMENTS:

* provider: Adds SSO API endpoint override parameter `endpoints.sso` ([#34302](https://github.com/hashicorp/terraform-provider-aws/issues/34302))
* resource/aws_appflow_connector_profile: Add `jwt_token` and `oauth2_grant_type` arguments to the `connector_profile_config.connector_profile_credentials.salesforce` block. ([#34248](https://github.com/hashicorp/terraform-provider-aws/issues/34248))
* resource/aws_autoscaling_group: Add plan-time validation of `initial_lifecycle_hook.default_result`, `initial_lifecycle_hook.heartbeat_timeout`, `initial_lifecycle_hook.lifecycle_transition`, `initial_lifecycle_hook.name`, `initial_lifecycle_hook.notification_target_arn` and `initial_lifecycle_hook.role_arn` ([#12145](https://github.com/hashicorp/terraform-provider-aws/issues/12145))
* resource/aws_autoscaling_lifecycle_hook: Add plan-time validation of `default_result`, `heartbeat_timeout`, `lifecycle_transition`, `name`, `notification_target_arn` and `role_arn` ([#12145](https://github.com/hashicorp/terraform-provider-aws/issues/12145))
* resource/aws_datasync_task: Add `task_report_config` argument ([#33861](https://github.com/hashicorp/terraform-provider-aws/issues/33861))
* resource/aws_db_instance: Add `postgres` as a valid `engine` value for blue/green deployments ([#34216](https://github.com/hashicorp/terraform-provider-aws/issues/34216))
* resource/aws_dms_endpoint: Add `pause_replication_tasks`, which when set to `true`, pauses associated running replication tasks, regardless if they are managed by Terraform, prior to modifying the endpoint (only tasks paused by the resource will be restarted after the modification completes) ([#34316](https://github.com/hashicorp/terraform-provider-aws/issues/34316))
* resource/aws_eks_cluster: Allow `vpc_config.security_group_ids` and `vpc_config.subnet_ids` to be updated in-place ([#32409](https://github.com/hashicorp/terraform-provider-aws/issues/32409))
* resource/aws_inspector2_organization_configuration: Add `lambda_code` argument to the `auto_enable` configuration block ([#34261](https://github.com/hashicorp/terraform-provider-aws/issues/34261))
* resource/aws_route53_record: Allow import of records with an empty record name. ([#34212](https://github.com/hashicorp/terraform-provider-aws/issues/34212))
* resource/aws_sagemaker_domain: Add `default_user_settings.canvas_app_settings.direct_deploy_settings`, `default_user_settings.canvas_app_settings.identity_provider_oauth_settings` and `default_user_settings.canvas_app_settings.kendra_settings` arguments ([#34265](https://github.com/hashicorp/terraform-provider-aws/issues/34265))
* resource/aws_sagemaker_domain: Change `default_space_settings.kernel_gateway_app_settings.custom_image`, `default_user_settings.kernel_gateway_app_settings.custom_image` and `default_user_settings.r_session_app_settings.custom_image` `MaxItems` from `30` to `200` ([#34265](https://github.com/hashicorp/terraform-provider-aws/issues/34265))
* resource/aws_sagemaker_feature_group: Add `offline_store_config.s3_storage_config.resolved_output_s3_uri`, `online_store_config.storage_type` and `online_store_config.ttl_duration` arguments ([#34283](https://github.com/hashicorp/terraform-provider-aws/issues/34283))
* resource/aws_sagemaker_feature_group: Allow `online_store_config.ttl_duration` to be updated in-place ([#34283](https://github.com/hashicorp/terraform-provider-aws/issues/34283))
* resource/aws_sagemaker_model: Add `container.model_data_source` and `primary_container.model_data_source` configuration blocks ([#34158](https://github.com/hashicorp/terraform-provider-aws/issues/34158))
* resource/aws_sagemaker_space: Change `space_settings.kernel_gateway_app_settings.custom_image` `MaxItems` from `30` to `200` ([#34265](https://github.com/hashicorp/terraform-provider-aws/issues/34265))
* resource/aws_sagemaker_user_profile: Add `default_user_settings.canvas_app_settings.direct_deploy_settings`, `default_user_settings.canvas_app_settings.identity_provider_oauth_settings` and `default_user_settings.canvas_app_settings.kendra_settings` arguments ([#34265](https://github.com/hashicorp/terraform-provider-aws/issues/34265))
* resource/aws_sns_topic: Add `archive_policy` argument and `beginning_archive_time` attribute to support [message archiving](https://docs.aws.amazon.com/sns/latest/dg/fifo-message-archiving-replay.html) ([#34252](https://github.com/hashicorp/terraform-provider-aws/issues/34252))
* resource/aws_sns_topic: Add `replay_policy` argument ([#34252](https://github.com/hashicorp/terraform-provider-aws/issues/34252))

BUG FIXES:

* provider: Fix `Value Conversion Error` panic for certain resources when `null` tag values are specified ([#34319](https://github.com/hashicorp/terraform-provider-aws/issues/34319))
* provider: Fixes parsing error in AWS shared config files with extra whitespace ([#34300](https://github.com/hashicorp/terraform-provider-aws/issues/34300))
* provider: Fixes poor performance when parsing AWS shared config files ([#34300](https://github.com/hashicorp/terraform-provider-aws/issues/34300))
* resource/aws_autoscaling_group: Change all `initial_lifecycle_hook` configuration block attributes to [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) ([#34260](https://github.com/hashicorp/terraform-provider-aws/issues/34260))
* resource/aws_cloudtrail: Change the `id` attribute from the trail's name to its ARN to support [organization trails](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/creating-trail-organization.html) ([#30758](https://github.com/hashicorp/terraform-provider-aws/issues/30758))
* resource/aws_cloudwatch_event_rule: Increase `event_pattern` max length for validation to 4096 ([#34270](https://github.com/hashicorp/terraform-provider-aws/issues/34270))
* resource/aws_sagemaker_domain: Fix updating `default_space_settings.r_studio_server_pro_app_settings.access_status` from `ENABLED` to `DISABLED` ([#34265](https://github.com/hashicorp/terraform-provider-aws/issues/34265))

## 5.24.0 (November  2, 2023)

NOTES:

* resource/aws_detective_organization_admin_account: Because we cannot easily test this functionality, it is best effort and we ask for community help in testing ([#25237](https://github.com/hashicorp/terraform-provider-aws/issues/25237))
* resource/aws_detective_organization_configuration: Because we cannot easily test this functionality, it is best effort and we ask for community help in testing ([#25237](https://github.com/hashicorp/terraform-provider-aws/issues/25237))

FEATURES:

* **New Data Source:** `aws_opensearchserverless_lifecycle_policy` ([#34144](https://github.com/hashicorp/terraform-provider-aws/issues/34144))
* **New Resource:** `aws_detective_organization_admin_account` ([#25237](https://github.com/hashicorp/terraform-provider-aws/issues/25237))
* **New Resource:** `aws_detective_organization_configuration` ([#25237](https://github.com/hashicorp/terraform-provider-aws/issues/25237))
* **New Resource:** `aws_opensearchserverless_lifecycle_policy` ([#34144](https://github.com/hashicorp/terraform-provider-aws/issues/34144))
* **New Resource:** `aws_redshift_resource_policy` ([#34149](https://github.com/hashicorp/terraform-provider-aws/issues/34149))
* **New Resource:** `aws_verifiedaccess_endpoint` ([#30763](https://github.com/hashicorp/terraform-provider-aws/issues/30763))

ENHANCEMENTS:

* resource/aws_amplify_app: Add `custom_headers` argument ([#31561](https://github.com/hashicorp/terraform-provider-aws/issues/31561))
* resource/aws_batch_job_definition: Add `node_properties` argument ([#34153](https://github.com/hashicorp/terraform-provider-aws/issues/34153))
* resource/aws_finspace_kx_cluster: In-place updates are now supported for the `code`, `database`, and `initialization_script` arguments. The update timeout has been increased to 30 minutes. ([#34220](https://github.com/hashicorp/terraform-provider-aws/issues/34220))
* resource/aws_iot_topic_rule: Add `kafka.header` and `error_action.kafka.header` arguments ([#34191](https://github.com/hashicorp/terraform-provider-aws/issues/34191))
* resource/aws_networkmanager_connect_attachment: Add `NO_ENCAP` as a valid `options.protocol` value ([#34109](https://github.com/hashicorp/terraform-provider-aws/issues/34109))
* resource/aws_networkmanager_connect_peer: Add `subnet_arn` argument to support [Tunnel-less Connect attachments](https://docs.aws.amazon.com/network-manager/latest/cloudwan/cloudwan-connect-attachment.html#cloudwan-connect-tlc) ([#34109](https://github.com/hashicorp/terraform-provider-aws/issues/34109))
* resource/aws_networkmanager_connect_peer: `inside_cidr_blocks` is Optional ([#34109](https://github.com/hashicorp/terraform-provider-aws/issues/34109))
* resource/aws_rds_cluster: Remove the provider default (previously, "1") and use the AWS default for `backup_retention_period` (also, "1") to allow integration with AWS Backup ([#34187](https://github.com/hashicorp/terraform-provider-aws/issues/34187))
* resource/aws_redshift_cluster: Add `snapshot_arn` argument ([#34181](https://github.com/hashicorp/terraform-provider-aws/issues/34181))
* resource/aws_redshift_cluster: Add the `manage_master_password` and `master_password_secret_kms_key_id` arguments to support managed admin credentials ([#34182](https://github.com/hashicorp/terraform-provider-aws/issues/34182))
* resource/aws_s3_object: Add `override_provider` configuration block, allowing tags inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) to be ignored ([#33262](https://github.com/hashicorp/terraform-provider-aws/issues/33262))
* resource/aws_secretsmanager_secret_rotation: The `rotation_lambda_arn` argument is now optional to support modifying the rotation schedule of AWS-managed secrets. ([#34180](https://github.com/hashicorp/terraform-provider-aws/issues/34180))

BUG FIXES:

* data-source/aws_vpc_ipam_pools: Add `id` attribute for individual IPAM pools ([#32133](https://github.com/hashicorp/terraform-provider-aws/issues/32133))
* resource/aws_alb_listener_rule: Fixed the `action.forward.target_group` argument minimum item requirement. Previously this was set to 2, but the AWS API allows specifying a single target group. ([#33727](https://github.com/hashicorp/terraform-provider-aws/issues/33727))
* resource/aws_amplify_branch: Remove [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) from `enable_performance_mode` ([#34141](https://github.com/hashicorp/terraform-provider-aws/issues/34141))
* resource/aws_lb_listener_rule: Fixed the `action.forward.target_group` argument minimum item requirement. Previously this was set to 2, but the AWS API allows specifying a single target group. ([#33727](https://github.com/hashicorp/terraform-provider-aws/issues/33727))
* resource/aws_quicksight_analysis: Fix "expected type to be integer" errors in `window_options.bounds.*` argument validatation functions ([#34230](https://github.com/hashicorp/terraform-provider-aws/issues/34230))
* resource/aws_quicksight_dashboard: Fix "expected type to be integer" errors in `window_options.bounds.*` argument validatation functions ([#34230](https://github.com/hashicorp/terraform-provider-aws/issues/34230))
* resource/aws_quicksight_template: Fix "expected type to be integer" errors in `window_options.bounds.*` argument validatation functions ([#34230](https://github.com/hashicorp/terraform-provider-aws/issues/34230))
* resource/aws_rds_cluster: Avoid an error on delete related to `unexpected state 'scaling-compute'` ([#34187](https://github.com/hashicorp/terraform-provider-aws/issues/34187))

## 5.23.1 (October 27, 2023)

BUG FIXES:

* data-source/aws_lambda_function: Add `vpc_config.ipv6_allowed_for_dual_stack` attribute, fixing `Invalid address to set: []string{"vpc_config", "0", "ipv6_allowed_for_dual_stack"}` errors ([#34134](https://github.com/hashicorp/terraform-provider-aws/issues/34134))

## 5.23.0 (October 26, 2023)

NOTES:

* provider: This release includes an update to the AWS SDK for Go v2 with breaking type changes to several services: `finspace`, `kafka`, `medialive`, `rds`, `s3control`, `timestreamwrite`, and `xray`. These changes primarily affect how arguments with default values are serialized for outbound requests, changing scalar types to pointers. See [this AWS SDK for Go V2 issue](https://github.com/aws/aws-sdk-go-v2/issues/2162) for additional context. The corresponding provider changes should make this breakfix transparent to users, but as with any breaking change there is the potential for missed edge cases. If errors are observed in the impacted resources, please link to this dependency update pull request in the bug report. ([#34096](https://github.com/hashicorp/terraform-provider-aws/issues/34096))

FEATURES:

* **New Resource:** `aws_iot_domain_configuration` ([#24765](https://github.com/hashicorp/terraform-provider-aws/issues/24765))

ENHANCEMENTS:

* data-source/aws_imagebuilder_image: Add `image_scanning_configuration` attribute ([#34049](https://github.com/hashicorp/terraform-provider-aws/issues/34049))
* resource/aws_config_config_rule: Add `evaluation_mode` attribute ([#34033](https://github.com/hashicorp/terraform-provider-aws/issues/34033))
* resource/aws_elasticache_replication_group: Add `ip_discovery` and `network_type` arguments ([#34019](https://github.com/hashicorp/terraform-provider-aws/issues/34019))
* resource/aws_imagebuilder_image: Add `image_scanning_configuration` configuration block ([#34049](https://github.com/hashicorp/terraform-provider-aws/issues/34049))
* resource/aws_kms_key: Add configurable timeouts ([#34112](https://github.com/hashicorp/terraform-provider-aws/issues/34112))
* resource/aws_lambda_function: Add `vpc_config.ipv6_allowed_for_dual_stack` argument ([#34045](https://github.com/hashicorp/terraform-provider-aws/issues/34045))
* resource/aws_lb: Add `dns_record_client_routing_policy` attribute to configure Availability Zonal DNS affinity on Network Load Balancer (NLB) ([#33992](https://github.com/hashicorp/terraform-provider-aws/issues/33992))
* resource/aws_lb_target_group: Add `target_health_state` configuration block ([#34070](https://github.com/hashicorp/terraform-provider-aws/issues/34070))
* resource/aws_lb_target_group: Remove default value (`false`) for `connection_termination` argument and mark as Computed, to support new default behavior for UDP/TCP_UDP target groups ([#34070](https://github.com/hashicorp/terraform-provider-aws/issues/34070))
* resource/aws_neptune_cluster: Add `slowquery` as a valid `enable_cloudwatch_logs_exports` value ([#34053](https://github.com/hashicorp/terraform-provider-aws/issues/34053))

BUG FIXES:

* provider/tags: Prevent crash when `tags_all` is null ([#34073](https://github.com/hashicorp/terraform-provider-aws/issues/34073))
* resource/aws_autoscaling_group: Fix error when `launch_template` name is updated. ([#34086](https://github.com/hashicorp/terraform-provider-aws/issues/34086))
* resource/aws_dms_s3_endpoint: Don't send the default value of `false` for `add_trailing_padding_character`, maintaining compatibility with older ([pre-3.4.7](https://docs.aws.amazon.com/dms/latest/userguide/CHAP_ReleaseNotes.html#CHAP_ReleaseNotes.DMS347)) DMS engine versions ([#34048](https://github.com/hashicorp/terraform-provider-aws/issues/34048))
* resource/aws_ecs_task_definition: Add `0` as a valid value for `volume.efs_volume_configuration.transit_encryption_port`, preventing unexpected drift ([#34020](https://github.com/hashicorp/terraform-provider-aws/issues/34020))
* resource/aws_identitystore_group: Fix updating `description` attribute when it is changed ([#34037](https://github.com/hashicorp/terraform-provider-aws/issues/34037))
* resource/aws_iot_indexing_configuration: Add `thing_indexing_configuration.filter` attribute, resolving `InvalidRequestException: NamedShadowNames Filter must not be empty for enabling NamedShadowIndexingMode` errors ([#26859](https://github.com/hashicorp/terraform-provider-aws/issues/26859))
* resource/aws_storagegateway_gateway: Support the value `0` (representing Sunday) for `maintenance_start_time.day_of_week` ([#34015](https://github.com/hashicorp/terraform-provider-aws/issues/34015))
* resource/aws_verifiedaccess_group: Fix `InvalidParameterValue: Policy Document cannot be provided when Policy Enabled is false or missing` errors when updating `policy_document` ([#34054](https://github.com/hashicorp/terraform-provider-aws/issues/34054))

## 5.22.0 (October 19, 2023)

FEATURES:

* **New Data Source:** `aws_media_convert_queue` ([#27075](https://github.com/hashicorp/terraform-provider-aws/issues/27075))
* **New Resource:** `aws_elasticsearch_vpc_endpoint` ([#33925](https://github.com/hashicorp/terraform-provider-aws/issues/33925))
* **New Resource:** `aws_msk_replicator` ([#33973](https://github.com/hashicorp/terraform-provider-aws/issues/33973))

ENHANCEMENTS:

* data-source/aws_ec2_client_vpn_endpoint: Add `self_service_portal_url` attribute ([#34007](https://github.com/hashicorp/terraform-provider-aws/issues/34007))
* resource/aws_alb: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_alb_target_group: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_cloudfront_public_key: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_db_option_group: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_docdb_cluster: Support import of `cluster_identifier_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_docdb_cluster_instance: Support import of `identifier_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_docdb_cluster_parameter_group: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_docdb_subnet_group: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_ec2_client_vpn_endpoint: Add `self_service_portal_url` attribute ([#34007](https://github.com/hashicorp/terraform-provider-aws/issues/34007))
* resource/aws_elb: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_emr_security_configuration: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_iam_group_policy: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_iam_role_policy: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_iam_user_policy: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_iot_provisioning_template: Add `type` attribute ([#33950](https://github.com/hashicorp/terraform-provider-aws/issues/33950))
* resource/aws_lb: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_lb_target_group: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_neptune_cluster: Support import of `cluster_identifier_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_neptune_cluster_instance: Support import of `identifier_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_neptune_cluster_parameter_group: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_neptune_event_subscription: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_pinpoint_app: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_rds_cluster: Support import of `cluster_identifier_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_rds_cluster_instance: Support import of `identifier_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_signer_signing_profile: Support import of `name_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_signer_signing_profile_permission: Add `signer:SignPayload` as a valid `action` value ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_signer_signing_profile_permission: Support import of `statement_id_prefix` argument ([#33852](https://github.com/hashicorp/terraform-provider-aws/issues/33852))
* resource/aws_transfer_server: Change `pre_authentication_login_banner` and `post_authentication_login_banner` length limits to 4096 ([#33937](https://github.com/hashicorp/terraform-provider-aws/issues/33937))
* resource/aws_wafv2_web_acl: Add `ja3_fingerprint` to `field_to_match` configuration blocks ([#33933](https://github.com/hashicorp/terraform-provider-aws/issues/33933))

BUG FIXES:

* data-source/aws_dms_certificate: Fix crash when certificate not found ([#34012](https://github.com/hashicorp/terraform-provider-aws/issues/34012))
* resource/aws_cloudformation_stack: Fix error when `computed` values are not set when there is no update ([#33969](https://github.com/hashicorp/terraform-provider-aws/issues/33969))
* resource/aws_codecommit_repository: Doesn't force replacement when renaming ([#32207](https://github.com/hashicorp/terraform-provider-aws/issues/32207))
* resource/aws_db_instance: Creating resource from snapshot or point-in-time recovery now handles `manage_master_user_password` and `master_user_secret_kms_key_id` attributes correctly ([#33699](https://github.com/hashicorp/terraform-provider-aws/issues/33699))
* resource/aws_elasticache_replication_group: Fix error when switching `engine_version` from `6.x` to a specific `6.<digit>` version number ([#33954](https://github.com/hashicorp/terraform-provider-aws/issues/33954))
* resource/aws_iam_role: Fix refreshing `permission_boundary` when deleted outside of Terraform ([#33963](https://github.com/hashicorp/terraform-provider-aws/issues/33963))
* resource/aws_iam_user: Fix refreshing `permission_boundary` when deleted outside of Terraform ([#33963](https://github.com/hashicorp/terraform-provider-aws/issues/33963))
* resource/aws_inspector2_enabler: Fix `Value at 'resourceTypes' failed to satisfy constraint` errors ([#33348](https://github.com/hashicorp/terraform-provider-aws/issues/33348))
* resource/aws_neptune_cluster_instance: Remove [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) from `engine_version` ([#33487](https://github.com/hashicorp/terraform-provider-aws/issues/33487))
* resource/aws_neptune_cluster_parameter_group: Fix condition where defined cluster parameters with system default values are seen as updates ([#33487](https://github.com/hashicorp/terraform-provider-aws/issues/33487))
* resource/aws_s3_bucket_object_lock_configuration: Fix `found resource` errors on Delete ([#33966](https://github.com/hashicorp/terraform-provider-aws/issues/33966))

## 5.21.0 (October 12, 2023)

FEATURES:

* **New Data Source:** `aws_servicequotas_templates` ([#33871](https://github.com/hashicorp/terraform-provider-aws/issues/33871))
* **New Resource:** `aws_ec2_image_block_public_access` ([#33810](https://github.com/hashicorp/terraform-provider-aws/issues/33810))
* **New Resource:** `aws_guardduty_organization_configuration_feature` ([#33913](https://github.com/hashicorp/terraform-provider-aws/issues/33913))
* **New Resource:** `aws_servicequotas_template_association` ([#33725](https://github.com/hashicorp/terraform-provider-aws/issues/33725))
* **New Resource:** `aws_verifiedaccess_group` ([#33297](https://github.com/hashicorp/terraform-provider-aws/issues/33297))
* **New Resource:** `aws_verifiedaccess_instance_logging_configuration` ([#33864](https://github.com/hashicorp/terraform-provider-aws/issues/33864))

ENHANCEMENTS:

* data-source/aws_dms_endpoint: Add `s3_settings.glue_catalog_generation` attribute ([#33778](https://github.com/hashicorp/terraform-provider-aws/issues/33778))
* data-source/aws_msk_cluster: Add `cluster_uuid` attribute ([#33805](https://github.com/hashicorp/terraform-provider-aws/issues/33805))
* resource/aws_codedeploy_deployment_group: Add `outdated_instances_strategy` argument ([#33844](https://github.com/hashicorp/terraform-provider-aws/issues/33844))
* resource/aws_dms_endpoint: Add `s3_settings.glue_catalog_generation` attribute ([#33778](https://github.com/hashicorp/terraform-provider-aws/issues/33778))
* resource/aws_dms_s3_endpoint: Add `glue_catalog_generation` attribute ([#33778](https://github.com/hashicorp/terraform-provider-aws/issues/33778))
* resource/aws_docdb_cluster: Add `allow_major_version_upgrade` argument ([#33790](https://github.com/hashicorp/terraform-provider-aws/issues/33790))
* resource/aws_docdb_cluster_instance: Add `copy_tags_to_snapshot` argument ([#31022](https://github.com/hashicorp/terraform-provider-aws/issues/31022))
* resource/aws_dynamodb_table: Add `import_table` configuration block ([#33802](https://github.com/hashicorp/terraform-provider-aws/issues/33802))
* resource/aws_msk_cluster: Add `cluster_uuid` attribute ([#33805](https://github.com/hashicorp/terraform-provider-aws/issues/33805))
* resource/aws_msk_serverless_cluster: Add `cluster_uuid` attribute ([#33805](https://github.com/hashicorp/terraform-provider-aws/issues/33805))
* resource/aws_networkmanager_core_network: Add `base_policy_document` argument ([#33712](https://github.com/hashicorp/terraform-provider-aws/issues/33712))
* resource/aws_redshiftserverless_workgroup: Allow `require_ssl` and `use_fips_ssl` `config_parameters` keys ([#33916](https://github.com/hashicorp/terraform-provider-aws/issues/33916))
* resource/aws_s3_bucket: Use configurable timeout for resource Delete ([#33845](https://github.com/hashicorp/terraform-provider-aws/issues/33845))
* resource/aws_verifiedaccess_instance: Add `fips_enabled` argument ([#33880](https://github.com/hashicorp/terraform-provider-aws/issues/33880))
* resource/aws_vpclattice_target_group: Add `config.lambda_event_structure_version` argument ([#33804](https://github.com/hashicorp/terraform-provider-aws/issues/33804))
* resource/aws_vpclattice_target_group: Make `config.port`, `config.protocol` and `config.vpc_identifier` optional ([#33804](https://github.com/hashicorp/terraform-provider-aws/issues/33804))
* resource/aws_wafv2_web_acl: Add `aws_managed_rules_acfp_rule_set` to `managed_rule_group_configs` configuration block ([#33915](https://github.com/hashicorp/terraform-provider-aws/issues/33915))

BUG FIXES:

* provider: Respect valid values for the `AWS_S3_US_EAST_1_REGIONAL_ENDPOINT` environment variable when configuring the S3 API client ([#33874](https://github.com/hashicorp/terraform-provider-aws/issues/33874))
* resource/aws_appflow_connector_profile: Fix various crashes ([#33856](https://github.com/hashicorp/terraform-provider-aws/issues/33856))
* resource/aws_db_parameter_group: Group names containing periods (`.`) no longer fail validation ([#33704](https://github.com/hashicorp/terraform-provider-aws/issues/33704))
* resource/aws_opensearchserverless_collection: Fix crash when error is returned ([#33918](https://github.com/hashicorp/terraform-provider-aws/issues/33918))
* resource/aws_rds_cluster_parameter_group: Group names containing periods (`.`) no longer fail validation ([#33704](https://github.com/hashicorp/terraform-provider-aws/issues/33704))

## 5.20.1 (October 10, 2023)

NOTES:

* provider: Build with [Terraform Plugin Framework v1.4.1](https://github.com/hashicorp/terraform-plugin-framework/blob/main/CHANGELOG.md#141-october-09-2023), fixing potential [initialization errors](https://github.com/hashicorp/terraform/issues/33990) when using v1.6 of the Terraform CLI.

## 5.20.0 (October  6, 2023)

FEATURES:

* **New Resource:** `aws_guardduty_detector_feature` ([#31463](https://github.com/hashicorp/terraform-provider-aws/issues/31463))
* **New Resource:** `aws_servicequotas_template` ([#33688](https://github.com/hashicorp/terraform-provider-aws/issues/33688))
* **New Resource:** `aws_sesv2_account_vdm_attributes` ([#33705](https://github.com/hashicorp/terraform-provider-aws/issues/33705))
* **New Resource:** `aws_verifiedaccess_instance_trust_provider_attachment` ([#33734](https://github.com/hashicorp/terraform-provider-aws/issues/33734))

ENHANCEMENTS:

* data-source/aws_guardduty_detector: Add `features` attribute ([#31463](https://github.com/hashicorp/terraform-provider-aws/issues/31463))
* resource/aws_finspace_kx_cluster: Increase default creation timeout to 45 minutes, default deletion timeout to 60 minutes ([#33745](https://github.com/hashicorp/terraform-provider-aws/issues/33745))
* resource/aws_finspace_kx_environment: Increase default deletion timeout to 45 minutes ([#33745](https://github.com/hashicorp/terraform-provider-aws/issues/33745))
* resource/aws_guardduty_filter: Add plan-time validation of `name` ([#21030](https://github.com/hashicorp/terraform-provider-aws/issues/21030))
* resource/aws_kinesis_firehose_delivery_stream: Add `opensearchserverless_configuration` and `msk_source_configuration` configuration blocks ([#33101](https://github.com/hashicorp/terraform-provider-aws/issues/33101))
* resource/aws_kinesis_firehose_delivery_stream: Add `opensearchserverless` as a valid `destination` value ([#33101](https://github.com/hashicorp/terraform-provider-aws/issues/33101))

BUG FIXES:

* data-source/aws_fsx_ontap_storage_virtual_machine: Fix crash when `active_directory_configuration.self_managed_active_directory_configuration.file_system_administrators_group` is not configured ([#33800](https://github.com/hashicorp/terraform-provider-aws/issues/33800))
* resource/aws_ec2_transit_gateway_route : Fix TGW route search filter to avoid routes being missed when more than 1,000 static routes are in a TGW route table ([#33765](https://github.com/hashicorp/terraform-provider-aws/issues/33765))
* resource/aws_fsx_ontap_storage_virtual_machine: Fix crash when `active_directory_configuration.self_managed_active_directory_configuration.file_system_administrators_group` is not configured ([#33800](https://github.com/hashicorp/terraform-provider-aws/issues/33800))
* resource/aws_medialive_channel: Fix VPC settings flatten/expand/docs. ([#33558](https://github.com/hashicorp/terraform-provider-aws/issues/33558))
* resource/aws_vpc_endpoint: Set `dns_options.dns_record_ip_type` to `Computed` to prevent diffs ([#33743](https://github.com/hashicorp/terraform-provider-aws/issues/33743))

## 5.19.0 (September 29, 2023)

BREAKING CHANGES:

* data-source/aws_s3_bucket_object: Following migration to [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/), the `metadata` attribute's [keys](https://developer.hashicorp.com/terraform/language/expressions/types#maps-objects) are always [returned in lowercase](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#HeadObjectOutput) ([#33660](https://github.com/hashicorp/terraform-provider-aws/issues/33660))
* data-source/aws_s3_object: Following migration to [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/), the `metadata` attribute's [keys](https://developer.hashicorp.com/terraform/language/expressions/types#maps-objects) are always [returned in lowercase](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/s3#HeadObjectOutput) ([#33660](https://github.com/hashicorp/terraform-provider-aws/issues/33660))

NOTES:

* data-source/aws_s3_bucket_object: The `metadata` attribute's keys are now always returned in lowercase. Please modify configurations as necessary ([#33660](https://github.com/hashicorp/terraform-provider-aws/issues/33660))
* data-source/aws_s3_object: The `metadata` attribute's keys are now always returned in lowercase. Please modify configurations as necessary ([#33660](https://github.com/hashicorp/terraform-provider-aws/issues/33660))
* resource/aws_iam_*: This release introduces additional validation of IAM policy JSON arguments to detect duplicate keys. Previously, arguments with duplicated keys resulted in all but one of the key values being overwritten. Since this results in unexpected IAM policies being submitted to AWS, we have updated the validation logic to error in these cases. This may cause existing IAM policy arguments to fail validation, however, those policies are likely not what was originally intended. ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))

FEATURES:

* **New Resource:** `aws_cleanrooms_configured_table` ([#33602](https://github.com/hashicorp/terraform-provider-aws/issues/33602))
* **New Resource:** `aws_dms_replication_config` ([#32908](https://github.com/hashicorp/terraform-provider-aws/issues/32908))
* **New Resource:** `aws_lexv2models_bot` ([#33475](https://github.com/hashicorp/terraform-provider-aws/issues/33475))
* **New Resource:** `aws_rds_custom_db_engine_version` ([#33285](https://github.com/hashicorp/terraform-provider-aws/issues/33285))

ENHANCEMENTS:

* resource/aws_cloud9_environment_ec2: Add `ubuntu-22.04-x86_64` and `resolve:ssm:/aws/service/cloud9/amis/ubuntu-22.04-x86_64` as valid values for `image_id` ([#33662](https://github.com/hashicorp/terraform-provider-aws/issues/33662))
* resource/aws_fsx_ontap_volume: Add `bypass_snaplock_enterprise_retention` argument and `snaplock_configuration` configuration block to support [SnapLock](https://docs.aws.amazon.com/fsx/latest/ONTAPGuide/snaplock.html) ([#32530](https://github.com/hashicorp/terraform-provider-aws/issues/32530))
* resource/aws_fsx_ontap_volume: Add `copy_tags_to_backups` and `snapshot_policy` arguments ([#32530](https://github.com/hashicorp/terraform-provider-aws/issues/32530))
* resource/aws_fsx_openzfs_volume: Add `delete_volume_options` argument ([#32530](https://github.com/hashicorp/terraform-provider-aws/issues/32530))
* resource/aws_lightsail_bucket: Add `force_delete` argument ([#33586](https://github.com/hashicorp/terraform-provider-aws/issues/33586))
* resource/aws_opensearch_outbound_connection: Add `connection_properties`, `connection_mode` and `accept_connection` arguments ([#32990](https://github.com/hashicorp/terraform-provider-aws/issues/32990))
* resource/aws_wafv2_rule_group: Add `rate_based_statement.custom_key` configuration block ([#33594](https://github.com/hashicorp/terraform-provider-aws/issues/33594))
* resource/aws_wafv2_web_acl: Add `rate_based_statement.custom_key` configuration block ([#33594](https://github.com/hashicorp/terraform-provider-aws/issues/33594))

BUG FIXES:

* resource/aws_batch_job_queue: Correctly validates elements of `compute_environments` as ARNs ([#33577](https://github.com/hashicorp/terraform-provider-aws/issues/33577))
* resource/aws_cloudfront_continuous_deployment_policy: Fix `IllegalUpdate` errors when updating a staging `aws_cloudfront_distribution` that is part of continuous deployment ([#33578](https://github.com/hashicorp/terraform-provider-aws/issues/33578))
* resource/aws_cloudfront_distribution: Fix `IllegalUpdate` errors when updating a staging distribution associated with an `aws_cloudfront_continuous_deployment_policy` ([#33578](https://github.com/hashicorp/terraform-provider-aws/issues/33578))
* resource/aws_cloudfront_distribution: Fix `PreconditionFailed` errors when destroying a distribution associated with an `aws_cloudfront_continuous_deployment_policy` ([#33578](https://github.com/hashicorp/terraform-provider-aws/issues/33578))
* resource/aws_cloudfront_distribution: Fix `StagingDistributionInUse` errors when destroying a distribution associated with an `aws_cloudfront_continuous_deployment_policy` ([#33578](https://github.com/hashicorp/terraform-provider-aws/issues/33578))
* resource/aws_datasync_location_fsx_ontap_file_system: Correct handling of `protocol.smb.domain`, `protocol.smb.user` and `protocol.smb.password` ([#33641](https://github.com/hashicorp/terraform-provider-aws/issues/33641))
* resource/aws_glacier_vault_lock: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_iam_group_policy: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_iam_policy: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_iam_role: Fail validation if duplicated keys are found in `assume_role_policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_iam_role_policy: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_iam_user_policy: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_mediastore_container_policy: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_s3_bucket_policy: Fix intermittent `couldn't find resource` errors on resource Create ([#33537](https://github.com/hashicorp/terraform-provider-aws/issues/33537))
* resource/aws_ssoadmin_permission_set_inline_policy: Fail validation if duplicated keys are found in `inline_policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_transfer_access: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))
* resource/aws_transfer_user: Fail validation if duplicated keys are found in `policy` ([#33570](https://github.com/hashicorp/terraform-provider-aws/issues/33570))

## 5.18.1 (September 26, 2023)

NOTES:

* documentation: Duplicate CDKTF guides with differing file extensions have been removed to resolve failures in the provider release workflow. ([#33630](https://github.com/hashicorp/terraform-provider-aws/issues/33630))

## 5.18.0 (September 21, 2023)

FEATURES:

* **New Data Source:** `aws_fsx_ontap_file_system` ([#32503](https://github.com/hashicorp/terraform-provider-aws/issues/32503))
* **New Data Source:** `aws_fsx_ontap_storage_virtual_machine` ([#32621](https://github.com/hashicorp/terraform-provider-aws/issues/32621))
* **New Data Source:** `aws_fsx_ontap_storage_virtual_machines` ([#32624](https://github.com/hashicorp/terraform-provider-aws/issues/32624))
* **New Data Source:** `aws_organizations_organizational_unit` ([#33408](https://github.com/hashicorp/terraform-provider-aws/issues/33408))
* **New Resource:** `aws_opensearch_package` ([#33227](https://github.com/hashicorp/terraform-provider-aws/issues/33227))
* **New Resource:** `aws_opensearch_package_association` ([#33227](https://github.com/hashicorp/terraform-provider-aws/issues/33227))

ENHANCEMENTS:

* resource/aws_fsx_ontap_storage_virtual_machine: Remove [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) from `active_directory_configuration.self_managed_active_directory_configuration.domain_name`, `active_directory_configuration.self_managed_active_directory_configuration.file_system_administrators_group` and `active_directory_configuration.self_managed_active_directory_configuration.organizational_unit_distinguished_name` allowing an SVM to join AD after creation ([#33466](https://github.com/hashicorp/terraform-provider-aws/issues/33466))

BUG FIXES:

* data-source/aws_sesv2_email_identity: Mark `dkim_signing_attributes.domain_signing_private_key` as sensitive ([#33477](https://github.com/hashicorp/terraform-provider-aws/issues/33477))
* resource/aws_db_instance: Fix so that `storage_throughput` can be changed when `iops` and `allocated_storage` are not changed ([#33529](https://github.com/hashicorp/terraform-provider-aws/issues/33529))
* resource/aws_db_option_group: Avoid erroneous differences being reported when an `option` `port` and/or `version` is not set ([#33511](https://github.com/hashicorp/terraform-provider-aws/issues/33511))
* resource/aws_fsx_ontap_storage_virtual_machine: Avoid recreating resource when `active_directory_configuration.self_managed_active_directory_configuration.file_system_administrators_group` is configured ([#33466](https://github.com/hashicorp/terraform-provider-aws/issues/33466))
* resource/aws_fsx_ontap_storage_virtual_machine: Change `file_system_id` to [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) ([#32621](https://github.com/hashicorp/terraform-provider-aws/issues/32621))
* resource/aws_s3_bucket_accelerate_configuration: Retry resource Delete on `OperationAborted: A conflicting conditional operation is currently in progress against this resource` errors ([#33531](https://github.com/hashicorp/terraform-provider-aws/issues/33531))
* resource/aws_s3_bucket_policy: Retry resource Delete on `OperationAborted: A conflicting conditional operation is currently in progress against this resource` errors ([#33531](https://github.com/hashicorp/terraform-provider-aws/issues/33531))
* resource/aws_s3_bucket_versioning: Retry resource Delete on `OperationAborted: A conflicting conditional operation is currently in progress against this resource` errors ([#33531](https://github.com/hashicorp/terraform-provider-aws/issues/33531))
* resource/aws_sesv2_email_identity: Mark `dkim_signing_attributes.domain_signing_private_key` as sensitive ([#33477](https://github.com/hashicorp/terraform-provider-aws/issues/33477))

## 5.17.0 (September 14, 2023)

NOTES:

* data-source/aws_s3_object: Migration to [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/) means that the edge case of specifying a single `/` as the value for `key` is no longer supported ([#33358](https://github.com/hashicorp/terraform-provider-aws/issues/33358))

FEATURES:

* **New Resource:** `aws_shield_application_layer_automatic_response` ([#33432](https://github.com/hashicorp/terraform-provider-aws/issues/33432))
* **New Resource:** `aws_verifiedaccess_instance` ([#33459](https://github.com/hashicorp/terraform-provider-aws/issues/33459))

ENHANCEMENTS:

* data-source/aws_s3_object: Add `checksum_mode` argument and `checksum_crc32`, `checksum_crc32c`, `checksum_sha1` and `checksum_sha256` attributes ([#33358](https://github.com/hashicorp/terraform-provider-aws/issues/33358))
* data-source/aws_s3control_multi_region_access_point: Add `details.region.bucket_account_id` attribute ([#33416](https://github.com/hashicorp/terraform-provider-aws/issues/33416))
* resource/aws_s3_object: Add `checksum_algorithm` argument and `checksum_crc32`, `checksum_crc32c`, `checksum_sha1` and `checksum_sha256` attributes ([#33358](https://github.com/hashicorp/terraform-provider-aws/issues/33358))
* resource/aws_s3_object_copy: Add `checksum_algorithm` argument and `checksum_crc32`, `checksum_crc32c`, `checksum_sha1` and `checksum_sha256` attributes ([#33358](https://github.com/hashicorp/terraform-provider-aws/issues/33358))
* resource/aws_s3control_multi_region_access_point: Add `details.region.bucket_account_id` argument to support [cross-account Multi-Region Access Points](https://docs.aws.amazon.com/AmazonS3/latest/userguide/multi-region-access-point-buckets.html) ([#33416](https://github.com/hashicorp/terraform-provider-aws/issues/33416))
* resource/aws_s3control_multi_region_access_point: Add `details.region.region` attribute ([#33416](https://github.com/hashicorp/terraform-provider-aws/issues/33416))
* resource/aws_schemas_schema: Add `JSONSchemaDraft4` schema type support ([#33442](https://github.com/hashicorp/terraform-provider-aws/issues/33442))
* resource/aws_transfer_connector: Add `sftp_config` argument and make `as2_config` optional ([#32741](https://github.com/hashicorp/terraform-provider-aws/issues/32741))
* resource/aws_wafv2_web_acl: Retry resource Update on `WAFOptimisticLockException` errors ([#33432](https://github.com/hashicorp/terraform-provider-aws/issues/33432))

BUG FIXES:

* resource/aws_dms_replication_task: Fix error when `replication_task_settings` is `nil` ([#33456](https://github.com/hashicorp/terraform-provider-aws/issues/33456))
* resource/aws_elasticache_cluster: Fix regression for `redis` engine types caused by the new `transit_encryption_enabled` argument ([#33451](https://github.com/hashicorp/terraform-provider-aws/issues/33451))
* resource/aws_neptune_cluster: Fix ignored `kms_key_arn` on restore from DB cluster snapshot ([#33413](https://github.com/hashicorp/terraform-provider-aws/issues/33413))
* resource/aws_servicecatalog_product: Allow import on `provisioning_artifact_parameters` attribute ([#33448](https://github.com/hashicorp/terraform-provider-aws/issues/33448))
* resource/aws_subnet: Fix destroy error when there is a lingering ENI for DMS ([#33375](https://github.com/hashicorp/terraform-provider-aws/issues/33375))

## 5.16.2 (September 11, 2023)

FEATURES:

* **New Data Source:** `aws_cognito_identity_pool` ([#33053](https://github.com/hashicorp/terraform-provider-aws/issues/33053))
* **New Resource:** `aws_verifiedaccess_trust_provider` ([#33195](https://github.com/hashicorp/terraform-provider-aws/issues/33195))

ENHANCEMENTS:

* resource/aws_autoscaling_group: Change the default values of `instance_refresh.preferences.scale_in_protected_instances` and `instance_refresh.preferences.standby_instances` from `Wait` to the [Amazon EC2 Auto Scaling console recommended value](https://docs.aws.amazon.com/autoscaling/ec2/userguide/understand-instance-refresh-default-values.html) of `Ignore` ([#33382](https://github.com/hashicorp/terraform-provider-aws/issues/33382))
* resource/aws_s3control_object_lambda_access_point: Add `alias` attribute ([#33388](https://github.com/hashicorp/terraform-provider-aws/issues/33388))

BUG FIXES:

* resource/aws_autoscaling_group: Fix `ValidationError` errors when starting Auto Scaling group instance refresh ([#33382](https://github.com/hashicorp/terraform-provider-aws/issues/33382))
* resource/aws_iot_topic_rule: Fix `InvalidParameter` errors on Update with Kafka destinations ([#33360](https://github.com/hashicorp/terraform-provider-aws/issues/33360))
* resource/aws_lightsail_certificate: Fix validation of `name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_database: Fix validation of `name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_disk: Fix validation of `name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_instance: Fix validation of `name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_lb: Fix validation of `lb_name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_lb_attachment: Fix validation of `lb_name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_lb_certificate: Fix validation of `lb_name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_lb_certificate_attachment: Fix validation of `lb_name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_lb_https_redirection_policy: Fix validation of `lb_name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))
* resource/aws_lightsail_lb_stickiness_policy: Fix validation of `lb_name` ([#33405](https://github.com/hashicorp/terraform-provider-aws/issues/33405))

## 5.16.1 (September  8, 2023)

BUG FIXES:

* data-source/aws_efs_file_system: Fix `Search returned 0 results` errors when there are more than 101 file systems in the configured Region ([#33336](https://github.com/hashicorp/terraform-provider-aws/issues/33336))
* resource/aws_db_instance_automated_backups_replication: Fix `unexpected state` errors on resource Create ([#33369](https://github.com/hashicorp/terraform-provider-aws/issues/33369))
* resource/aws_glue_catalog_table: Fix removal of `metadata_location` and `table_type` `parameters` when updating Iceberg tables ([#33374](https://github.com/hashicorp/terraform-provider-aws/issues/33374))
* resource/aws_service_discovery_instance: Fix validation error "expected to match regular expression" ([#33371](https://github.com/hashicorp/terraform-provider-aws/issues/33371))

## 5.16.0 (September  8, 2023)

NOTES:

* provider: Performance regression introduced in v5.14.0 should be largely mitigated ([#33317](https://github.com/hashicorp/terraform-provider-aws/issues/33317))

FEATURES:

* **New Resource:** `aws_shield_drt_access_log_bucket_association` ([#33328](https://github.com/hashicorp/terraform-provider-aws/issues/33328))
* **New Resource:** `aws_shield_drt_access_role_arn_association` ([#33328](https://github.com/hashicorp/terraform-provider-aws/issues/33328))

ENHANCEMENTS:

* data-source/aws_api_gateway_api_key: Add `customer_id` attribute ([#33281](https://github.com/hashicorp/terraform-provider-aws/issues/33281))
* data-source/aws_fsx_windows_file_system: Add `disk_iops_configuration` attribute ([#33303](https://github.com/hashicorp/terraform-provider-aws/issues/33303))
* data-source/aws_opensearch_domain: Add `software_update_options` attribute ([#32234](https://github.com/hashicorp/terraform-provider-aws/issues/32234))
* data-source/aws_s3_objects: Add `request_payer` argument and `request_charged` attribute ([#33304](https://github.com/hashicorp/terraform-provider-aws/issues/33304))
* data-source/aws_s3_objects: Add plan-time validation of `encoding_type` ([#33304](https://github.com/hashicorp/terraform-provider-aws/issues/33304))
* resource/aws_api_gateway_account: Add `api_key_version` and `features` attributes ([#33279](https://github.com/hashicorp/terraform-provider-aws/issues/33279))
* resource/aws_api_gateway_api_key: Add `customer_id` argument ([#33281](https://github.com/hashicorp/terraform-provider-aws/issues/33281))
* resource/aws_api_gateway_api_key: Allow updating `name` ([#33281](https://github.com/hashicorp/terraform-provider-aws/issues/33281))
* resource/aws_autoscaling_group: Add `scale_in_protected_instances` and `standby_instances` attributes to `instance_refresh.preferences` configuration block ([#33310](https://github.com/hashicorp/terraform-provider-aws/issues/33310))
* resource/aws_dms_endpoint: Add `redshift-serverless` as valid value for `engine_name` ([#33316](https://github.com/hashicorp/terraform-provider-aws/issues/33316))
* resource/aws_elasticache_cluster: Add `transit_encryption_enabled` argument, enabling in-transit encryption for Memcached clusters inside a VPC ([#26987](https://github.com/hashicorp/terraform-provider-aws/issues/26987))
* resource/aws_fsx_windows_file_system: Add `disk_iops_configuration` configuration block ([#33303](https://github.com/hashicorp/terraform-provider-aws/issues/33303))
* resource/aws_glue_catalog_table: Add `open_table_format_input` configuration block to support open table formats such as [Apache Iceberg](https://iceberg.apache.org/) ([#33274](https://github.com/hashicorp/terraform-provider-aws/issues/33274))
* resource/aws_medialive_channel: Implement expand/flatten functions for `automatic_input_failover_settings` in `input_attachments` ([#33129](https://github.com/hashicorp/terraform-provider-aws/issues/33129))
* resource/aws_opensearch_domain: Add `software_update_options` attribute ([#32234](https://github.com/hashicorp/terraform-provider-aws/issues/32234))
* resource/aws_ssm_association: Add `sync_compliance` attribute ([#23515](https://github.com/hashicorp/terraform-provider-aws/issues/23515))

BUG FIXES:

* data-source/aws_identitystore_group: Restore `filter` argument to prevent `UnknownOperationException` errors in certain Regions ([#33311](https://github.com/hashicorp/terraform-provider-aws/issues/33311))
* data-source/aws_identitystore_user: Restore `filter` argument to prevent `UnknownOperationException` errors in certain Regions ([#33311](https://github.com/hashicorp/terraform-provider-aws/issues/33311))
* data-source/aws_s3_objects: Respect configured `max_keys` value if it's greater than `1000` ([#33304](https://github.com/hashicorp/terraform-provider-aws/issues/33304))
* resource/aws_api_gateway_account: Allow setting `cloudwatch_role_arn` to an empty value and set it correctly on Read, allowing its value to be determined on import ([#33279](https://github.com/hashicorp/terraform-provider-aws/issues/33279))
* resource/aws_fsx_ontap_file_system: Increase maximum value of `disk_iops_configuration.iops` to `160000` ([#33263](https://github.com/hashicorp/terraform-provider-aws/issues/33263))
* resource/aws_servicecatalog_principal_portfolio_association: Fix `ResourceNotFoundException` errors on resource Delete when configured `principal_type` is `IAM_PATTERN` ([#32243](https://github.com/hashicorp/terraform-provider-aws/issues/32243))

## 5.15.0 (August 31, 2023)

ENHANCEMENTS:

* data-source/aws_efs_file_system: Add `name` attribute ([#33243](https://github.com/hashicorp/terraform-provider-aws/issues/33243))
* data-source/aws_lakeformation_data_lake_settings: Add `read_only_admins` attribute ([#33189](https://github.com/hashicorp/terraform-provider-aws/issues/33189))
* data-source/aws_opensearch_domain: Add `cluster_config.multi_az_with_standby_enabled` attribute ([#33031](https://github.com/hashicorp/terraform-provider-aws/issues/33031))
* resource/aws_cloudformation_stack_set: Support resource import with `call_as = "DELEGATED_ADMIN"` via _StackSetName_,_CallAs_ syntax for `import` block or `terraform import` command ([#19092](https://github.com/hashicorp/terraform-provider-aws/issues/19092))
* resource/aws_cloudformation_stack_set_instance: Support resource import with `call_as = "DELEGATED_ADMIN"` via _StackSetName_,_AccountID_,_Region_,_CallAs_ syntax for `import` block or `terraform import` command ([#19092](https://github.com/hashicorp/terraform-provider-aws/issues/19092))
* resource/aws_datasync_location_fsx_openzfs_file_system: Fix `setting protocol: Invalid address to set` errors ([#33225](https://github.com/hashicorp/terraform-provider-aws/issues/33225))
* resource/aws_efs_file_system: Add `name` attribute ([#33243](https://github.com/hashicorp/terraform-provider-aws/issues/33243))
* resource/aws_fsx_openzfs_file_system: Add `endpoint_ip_address_range`, `preferred_subnet_id` and `route_table_ids` arguments to support the [Multi-AZ deployment type](https://docs.aws.amazon.com/fsx/latest/OpenZFSGuide/availability-durability.html#choosing-single-or-multi) ([#33245](https://github.com/hashicorp/terraform-provider-aws/issues/33245))
* resource/aws_lakeformation_data_lake_settings: Add `read_only_admins` argument ([#33189](https://github.com/hashicorp/terraform-provider-aws/issues/33189))
* resource/aws_opensearch_domain: Add `cluster_config.multi_az_with_standby_enabled` argument ([#33031](https://github.com/hashicorp/terraform-provider-aws/issues/33031))
* resource/aws_wafv2_rule_group: Add `name_prefix` argument ([#33206](https://github.com/hashicorp/terraform-provider-aws/issues/33206))
* resource/aws_wafv2_web_acl: Add `statement.managed_rule_group_statement.managed_rule_group_configs.aws_managed_rules_atp_rule_set.enable_regex_in_path` argument ([#33217](https://github.com/hashicorp/terraform-provider-aws/issues/33217))

BUG FIXES:

* provider: Correctly use old and new tag values when updating `tags` that are `computed` ([#33226](https://github.com/hashicorp/terraform-provider-aws/issues/33226))
* resource/aws_appflow_connector_profile: Fix validation on `oauth2` in `custom_connector_profile` ([#33192](https://github.com/hashicorp/terraform-provider-aws/issues/33192))
* resource/aws_cloudformation_stack_set: Fix `Can only set RetainStacksOnAccountRemoval if AutoDeployment is enabled` errors ([#19092](https://github.com/hashicorp/terraform-provider-aws/issues/19092))
* resource/aws_cloudwatch_event_bus_policy: Fix error during plan when the associated aws_cloudwatch_event_bus resource is manually deleted ([#33203](https://github.com/hashicorp/terraform-provider-aws/issues/33203))
* resource/aws_codeartifact_domain: Change the type of asset_size_bytes to `TypeString` instead of `TypeInt` to prevent `value out of range` panic ([#33220](https://github.com/hashicorp/terraform-provider-aws/issues/33220))
* resource/aws_efs_file_system_policy: Retry IAM eventual consistency errors ([#21734](https://github.com/hashicorp/terraform-provider-aws/issues/21734))
* resource/aws_fsx_openzfs_file_system: Wait for administrative action completion when updating root volume ([#33245](https://github.com/hashicorp/terraform-provider-aws/issues/33245))
* resource/aws_iot_thing_type: Fix error during plan when resource is manually deleted ([#33203](https://github.com/hashicorp/terraform-provider-aws/issues/33203))
* resource/aws_kms_key: Fix `tag propagation: timeout while waiting for state to become 'TRUE'` errors when any tag value is empty (`""`) ([#33226](https://github.com/hashicorp/terraform-provider-aws/issues/33226))
* resource/aws_wafv2_web_acl: Prevent deletion of the AWS-managed `ShieldMitigationRuleGroup` rule on resource Update ([#33216](https://github.com/hashicorp/terraform-provider-aws/issues/33216))

## 5.14.0 (August 24, 2023)

NOTES:

* data-source/aws_iam_policy_document: In some cases, `statement.*.condition` blocks with the same `test` and `variable` arguments were incorrectly handled by the provider. Since this results in unexpected IAM Policies being submitted to AWS, we have updated the logic to merge `values` lists in this case. This may cause existing IAM Policy documents to report a difference. However, those policies are likely not what was originally intended. ([#33093](https://github.com/hashicorp/terraform-provider-aws/issues/33093))

FEATURES:

* **New Resource:** `aws_datasync_location_azure_blob` ([#32632](https://github.com/hashicorp/terraform-provider-aws/issues/32632))
* **New Resource:** `aws_datasync_location_fsx_ontap_file_system` ([#32632](https://github.com/hashicorp/terraform-provider-aws/issues/32632))

ENHANCEMENTS:

* data-source/aws_dms_endpoint: Fix crash when specified endpoint not found ([#33158](https://github.com/hashicorp/terraform-provider-aws/issues/33158))
* data-source/aws_dms_replication_instance: Add `network_type` attribute ([#33158](https://github.com/hashicorp/terraform-provider-aws/issues/33158))
* data-source/aws_ec2_network_insights_path: Add `destination_arn` and `source_arn` attributes ([#33168](https://github.com/hashicorp/terraform-provider-aws/issues/33168))
* resource/aws_dms_replication_instance: Add `network_type` argument ([#33158](https://github.com/hashicorp/terraform-provider-aws/issues/33158))
* resource/aws_ec2_network_insights_path: Add `destination_arn` and `source_arn` attributes ([#33168](https://github.com/hashicorp/terraform-provider-aws/issues/33168))
* resource/aws_finspace_kx_environment: Add `transit_gateway_configuration.*.attachment_network_acl_configuration` argument. ([#33123](https://github.com/hashicorp/terraform-provider-aws/issues/33123))
* resource/aws_medialive_channel: Updates schemas for `selector_settings` for `audio_selector` and `selector_settings` for `caption_selector` ([#32714](https://github.com/hashicorp/terraform-provider-aws/issues/32714))
* resource/aws_ssoadmin_account_assignment: Add configurable timeouts ([#33121](https://github.com/hashicorp/terraform-provider-aws/issues/33121))
* resource/aws_ssoadmin_customer_managed_policy_attachment: Add configurable timeouts ([#33121](https://github.com/hashicorp/terraform-provider-aws/issues/33121))
* resource/aws_ssoadmin_managed_policy_attachment: Add configurable timeouts ([#33121](https://github.com/hashicorp/terraform-provider-aws/issues/33121))
* resource/aws_ssoadmin_permission_set: Add configurable timeouts ([#33121](https://github.com/hashicorp/terraform-provider-aws/issues/33121))
* resource/aws_ssoadmin_permission_set_inline_policy: Add configurable timeouts ([#33121](https://github.com/hashicorp/terraform-provider-aws/issues/33121))
* resource/aws_ssoadmin_permissions_boundary_attachment: Add configurable timeouts ([#33121](https://github.com/hashicorp/terraform-provider-aws/issues/33121))

BUG FIXES:

* data-source/aws_iam_policy_document: Fix inconsistent handling of `condition` blocks with duplicated `test` and `variable` arguments ([#33093](https://github.com/hashicorp/terraform-provider-aws/issues/33093))
* resource/aws_ec2_host: Fixed a bug that caused resource recreation when specifying an `outpost_arn` without an `asset_id` ([#33142](https://github.com/hashicorp/terraform-provider-aws/issues/33142))
* resource/aws_ec2_network_insights_analysis: Fix `setting forward_path_components: Invalid address to set` errors ([#33168](https://github.com/hashicorp/terraform-provider-aws/issues/33168))
* resource/aws_ec2_network_insights_path: Avoid recreating resource when passing an ARN as `source` or `destination` ([#33168](https://github.com/hashicorp/terraform-provider-aws/issues/33168))
* resource/aws_ec2_network_insights_path: Retry `AnalysisExistsForNetworkInsightsPath` errors on resource Delete ([#33168](https://github.com/hashicorp/terraform-provider-aws/issues/33168))
* resource/aws_kms_key: Fix `tag propagation: timeout while waiting for state to become 'TRUE'` errors when [`ignore_tags`](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#ignore_tags) has been configured ([#33167](https://github.com/hashicorp/terraform-provider-aws/issues/33167))
* resource/aws_licensemanager_license_configuration: Surface `InvalidParameterValueException` errors during resource Delete ([#32845](https://github.com/hashicorp/terraform-provider-aws/issues/32845))
* resource/aws_msk_cluster_policy: Fix `Current cluster policy version needed for Update` errors ([#33118](https://github.com/hashicorp/terraform-provider-aws/issues/33118))
* resource/aws_quicksight_analysis: Change `definition.*.parameter_declarations` to a set type, preventing persistent differences ([#33120](https://github.com/hashicorp/terraform-provider-aws/issues/33120))
* resource/aws_quicksight_analysis: Fixed a bug that caused errors related to the `word_orientation` argument when using word cloud visuals. ([#33122](https://github.com/hashicorp/terraform-provider-aws/issues/33122))
* resource/aws_quicksight_analysis: Skip setting `definition.*.parameter_declarations.*.*_parameter_declaration.static_values` when empty, preventing persistent differences. ([#33161](https://github.com/hashicorp/terraform-provider-aws/issues/33161))
* resource/aws_quicksight_dashboard: Change `definition.*.parameter_declarations` to a set type, preventing persistent differences ([#33120](https://github.com/hashicorp/terraform-provider-aws/issues/33120))
* resource/aws_quicksight_dashboard: Fixed a bug that caused errors related to the `word_orientation` argument when using word cloud visuals. ([#33122](https://github.com/hashicorp/terraform-provider-aws/issues/33122))
* resource/aws_quicksight_dashboard: Skip setting `definition.*.parameter_declarations.*.*_parameter_declaration.static_values` when empty, preventing persistent differences. ([#33161](https://github.com/hashicorp/terraform-provider-aws/issues/33161))
* resource/aws_quicksight_template: Change `definition.*.parameter_declarations` to a set type, preventing persistent differences ([#33120](https://github.com/hashicorp/terraform-provider-aws/issues/33120))
* resource/aws_quicksight_template: Fixed a bug that caused errors related to the `word_orientation` argument when using word cloud visuals. ([#33122](https://github.com/hashicorp/terraform-provider-aws/issues/33122))
* resource/aws_quicksight_template: Skip setting `definition.*.parameter_declarations.*.*_parameter_declaration.static_values` when empty, preventing persistent differences. ([#33161](https://github.com/hashicorp/terraform-provider-aws/issues/33161))
* resource/aws_route53_zone: Skip disabling DNS SEC in unsupported partitions ([#33103](https://github.com/hashicorp/terraform-provider-aws/issues/33103))
* resource/aws_s3_object: Mark `acl` as Computed. This suppresses the diffs shown when migrating resources with no configured `acl` attribute value from v4.67.0 (or earlier) ([#33138](https://github.com/hashicorp/terraform-provider-aws/issues/33138))
* resource/aws_s3_object_copy: Mark `acl` as Computed. This suppresses the diffs shown when migrating resources with no configured `acl` attribute value from v4.67.0 (or earlier) ([#33138](https://github.com/hashicorp/terraform-provider-aws/issues/33138))
* resource/aws_securityhub_account: Remove default value (`SECURITY_CONTROL`) for `control_finding_generator` argument and mark as Computed ([#33095](https://github.com/hashicorp/terraform-provider-aws/issues/33095))

## 5.13.1 (August 18, 2023)

BUG FIXES:

* resource/aws_lambda_layer_version: Change `source_code_hash` back to [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew). This fixes `doesn't support update` errors ([#33097](https://github.com/hashicorp/terraform-provider-aws/issues/33097))
* resource/aws_organizations_organization: Fix `current Organization ID (o-xxxxxxxxxx) does not match` errors on resource Read ([#33091](https://github.com/hashicorp/terraform-provider-aws/issues/33091))

## 5.13.0 (August 18, 2023)

FEATURES:

* **New Resource:** `aws_msk_cluster_policy` ([#32848](https://github.com/hashicorp/terraform-provider-aws/issues/32848))
* **New Resource:** `aws_opensearch_vpc_endpoint` ([#32435](https://github.com/hashicorp/terraform-provider-aws/issues/32435))
* **New Resource:** `aws_ram_sharing_with_organization` ([#25433](https://github.com/hashicorp/terraform-provider-aws/issues/25433))

ENHANCEMENTS:

* data-source/aws_imagebuilder_image_pipeline: Add `image_scanning_configuration` attribute ([#33005](https://github.com/hashicorp/terraform-provider-aws/issues/33005))
* data-source/aws_ram_resource_share: Add `resource_arns` attribute ([#22591](https://github.com/hashicorp/terraform-provider-aws/issues/22591))
* provider: Adds the `s3_us_east_1_regional_endpoint` attribute to support using the regional S3 API endpoint in `us-east-1`. ([#33024](https://github.com/hashicorp/terraform-provider-aws/issues/33024))
* resource/aws_appstream_fleet: Retry ConcurrentModificationException errors during creation ([#32958](https://github.com/hashicorp/terraform-provider-aws/issues/32958))
* resource/aws_dms_endpoint: Add `babelfish` as an `engine_name` option ([#32975](https://github.com/hashicorp/terraform-provider-aws/issues/32975))
* resource/aws_imagebuilder_image_pipeline: Add `image_scanning_configuration` configuration block ([#33005](https://github.com/hashicorp/terraform-provider-aws/issues/33005))
* resource/aws_lb: Changes to `security_groups` for Network Load Balancers force a new resource if either the old or new set of security group IDs is empty ([#32987](https://github.com/hashicorp/terraform-provider-aws/issues/32987))
* resource/aws_rds_global_cluster: Add plan-time validation of `global_cluster_identifier` ([#30996](https://github.com/hashicorp/terraform-provider-aws/issues/30996))

BUG FIXES:

* data-source/aws_ecr_repository: Correctly set `most_recent_image_tags` when only a single image is found ([#31757](https://github.com/hashicorp/terraform-provider-aws/issues/31757))
* resource/aws_budgets_budget_action: No longer times out when creating a non-triggered action ([#33015](https://github.com/hashicorp/terraform-provider-aws/issues/33015))
* resource/aws_cloudformation_stack: Marks `outputs` as Computed when there are potential changes. ([#33059](https://github.com/hashicorp/terraform-provider-aws/issues/33059))
* resource/aws_cloudwatch_event_rule: Fix ARN-based partner event bus rule ID parsing error ([#30293](https://github.com/hashicorp/terraform-provider-aws/issues/30293))
* resource/aws_ecr_registry_scanning_configuration: Correctly delete rules on resource Update ([#31449](https://github.com/hashicorp/terraform-provider-aws/issues/31449))
* resource/aws_lambda_layer_version: Fix bug causing new version to be created on every apply when `source_code_hash` is used but not changed ([#32535](https://github.com/hashicorp/terraform-provider-aws/issues/32535))
* resource/aws_lb_listener_certificate: Remove from state when listener not found ([#32412](https://github.com/hashicorp/terraform-provider-aws/issues/32412))
* resource/aws_organizations_organization: Ensure that the Organization ID specified in `terraform import` is the current Organization ([#31796](https://github.com/hashicorp/terraform-provider-aws/issues/31796))
* resource/aws_quicksight_analysis: Adjust max length of `definition.*.calculated_fields.*.expression` to 32000 characters ([#33012](https://github.com/hashicorp/terraform-provider-aws/issues/33012))
* resource/aws_quicksight_analysis: Convert `definition.*.calculated_fields` to a set type, preventing persistent differences ([#33040](https://github.com/hashicorp/terraform-provider-aws/issues/33040))
* resource/aws_quicksight_analysis: Convert `permissions` argument to TypeSet, preventing persistent differences ([#33023](https://github.com/hashicorp/terraform-provider-aws/issues/33023))
* resource/aws_quicksight_analysis: Enable `font_configuration` to be set for table header styles ([#33018](https://github.com/hashicorp/terraform-provider-aws/issues/33018))
* resource/aws_quicksight_analysis: Enable `font_configuration` to be set for table header styles ([#33018](https://github.com/hashicorp/terraform-provider-aws/issues/33018))
* resource/aws_quicksight_analysis: Enable `font_configuration` to be set for table header styles ([#33018](https://github.com/hashicorp/terraform-provider-aws/issues/33018))
* resource/aws_quicksight_analysis: Raise limit for maximum allowed `visuals` blocks per sheet to 50 ([#32856](https://github.com/hashicorp/terraform-provider-aws/issues/32856))
* resource/aws_quicksight_dashboard: Adjust max length of `definition.*.calculated_fields.*.expression` to 32000 characters ([#33012](https://github.com/hashicorp/terraform-provider-aws/issues/33012))
* resource/aws_quicksight_dashboard: Convert `definition.*.calculated_fields` to a set type, preventing persistent differences ([#33040](https://github.com/hashicorp/terraform-provider-aws/issues/33040))
* resource/aws_quicksight_dashboard: Convert `permissions` argument to TypeSet, preventing persistent differences ([#33023](https://github.com/hashicorp/terraform-provider-aws/issues/33023))
* resource/aws_quicksight_data_set: Change permission attribute type from TypeList to TypeSet ([#32984](https://github.com/hashicorp/terraform-provider-aws/issues/32984))
* resource/aws_quicksight_template: Adjust max items of `definition.*.calculated_fields` to 500 ([#33012](https://github.com/hashicorp/terraform-provider-aws/issues/33012))
* resource/aws_quicksight_template: Adjust max length of `definition.*.calculated_fields.*.expression` to 32000 characters ([#33012](https://github.com/hashicorp/terraform-provider-aws/issues/33012))
* resource/aws_quicksight_template: Convert `definition.*.calculated_fields` to a set type, preventing persistent differences ([#33040](https://github.com/hashicorp/terraform-provider-aws/issues/33040))
* resource/aws_quicksight_template: Convert `permissions` argument to TypeSet, preventing persistent differences ([#33023](https://github.com/hashicorp/terraform-provider-aws/issues/33023))
* resource/aws_s3_bucket_logging: Fix perpetual drift when `expected_bucket_owner` is configured ([#32989](https://github.com/hashicorp/terraform-provider-aws/issues/32989))
* resource/aws_sagemaker_domain: Fix validation on `s3_kms_key_id` in `sharing_settings` and `kms_key_id` ([#32661](https://github.com/hashicorp/terraform-provider-aws/issues/32661))
* resource/aws_subnet: Fix allowing IPv6 to be enabled in an update after initial creation with IPv4 only ([#32896](https://github.com/hashicorp/terraform-provider-aws/issues/32896))
* resource/aws_wafv2_web_acl: Adds `rule_group_reference_statement.rule_action_override.action_to_use.challenge` argument ([#31127](https://github.com/hashicorp/terraform-provider-aws/issues/31127))

## 5.12.0 (August 10, 2023)

NOTES:

* data-source/aws_codecatalyst_dev_environment: Because we cannot easily test this functionality, it is best effort and we ask for community help in testing ([#32886](https://github.com/hashicorp/terraform-provider-aws/issues/32886))
* resource/aws_codecatalyst_dev_environment: Because we cannot easily test this functionality, it is best effort and we ask for community help in testing ([#32366](https://github.com/hashicorp/terraform-provider-aws/issues/32366))
* resource/aws_codecatalyst_project: Because we cannot easily test this functionality, it is best effort and we ask for community help in testing ([#32883](https://github.com/hashicorp/terraform-provider-aws/issues/32883))
* resource/aws_codecatalyst_source_repository: Because we cannot easily test this functionality, it is best effort and we ask for community help in testing ([#32899](https://github.com/hashicorp/terraform-provider-aws/issues/32899))

FEATURES:

* **New Data Source:** `aws_codecatalyst_dev_environment` ([#32886](https://github.com/hashicorp/terraform-provider-aws/issues/32886))
* **New Data Source:** `aws_ec2_transit_gateway_route_table_routes` ([#30771](https://github.com/hashicorp/terraform-provider-aws/issues/30771))
* **New Data Source:** `aws_msk_vpc_connection` ([#31062](https://github.com/hashicorp/terraform-provider-aws/issues/31062))
* **New Resource:** `aws_cloudfront_continuous_deployment_policy` ([#32936](https://github.com/hashicorp/terraform-provider-aws/issues/32936))
* **New Resource:** `aws_codecatalyst_dev_environment` ([#32366](https://github.com/hashicorp/terraform-provider-aws/issues/32366))
* **New Resource:** `aws_codecatalyst_project` ([#32883](https://github.com/hashicorp/terraform-provider-aws/issues/32883))
* **New Resource:** `aws_codecatalyst_source_repository` ([#32899](https://github.com/hashicorp/terraform-provider-aws/issues/32899))
* **New Resource:** `aws_msk_vpc_connection` ([#31062](https://github.com/hashicorp/terraform-provider-aws/issues/31062))

ENHANCEMENTS:

* data-source/aws_instance: Add `metadata_options.http_protocol_ipv6` attribute ([#32759](https://github.com/hashicorp/terraform-provider-aws/issues/32759))
* data-source/aws_rds_cluster: Add `db_system_id` attribute ([#32846](https://github.com/hashicorp/terraform-provider-aws/issues/32846))
* provider: Support `il-central-1` as a valid AWS Region ([#32878](https://github.com/hashicorp/terraform-provider-aws/issues/32878))
* resource/aws_autoscaling_group: Add `ignore_failed_scaling_activities` argument ([#32914](https://github.com/hashicorp/terraform-provider-aws/issues/32914))
* resource/aws_cloudfront_distribution: Add `continuous_deployment_policy_id` and `staging` arguments to support continuous deployments ([#32936](https://github.com/hashicorp/terraform-provider-aws/issues/32936))
* resource/aws_cloudwatch_composite_alarm: Add `actions_suppressor` configuration block ([#32751](https://github.com/hashicorp/terraform-provider-aws/issues/32751))
* resource/aws_cloudwatch_events_target: Add `sagemaker_pipeline_target` argument ([#32882](https://github.com/hashicorp/terraform-provider-aws/issues/32882))
* resource/aws_fms_admin_account: Add configurable timeouts ([#32860](https://github.com/hashicorp/terraform-provider-aws/issues/32860))
* resource/aws_glue_crawler: Add `hudi_target` argument ([#32898](https://github.com/hashicorp/terraform-provider-aws/issues/32898))
* resource/aws_instance: Add `http_protocol_ipv6` attribute to `metadata_options` configuration block ([#32759](https://github.com/hashicorp/terraform-provider-aws/issues/32759))
* resource/aws_lambda_event_source_mapping: Increased the maximum number of filters to 10 ([#32890](https://github.com/hashicorp/terraform-provider-aws/issues/32890))
* resource/aws_msk_broker: Add `bootstrap_brokers_vpc_connectivity_sasl_iam`, `bootstrap_brokers_vpc_connectivity_sasl_scram` and `bootstrap_brokers_vpc_connectivity_tls` attributes ([#31062](https://github.com/hashicorp/terraform-provider-aws/issues/31062))
* resource/aws_msk_broker: Add `vpc_connectivity` attribute to the `broker_node_group_info.connectivity_info` configuration block ([#31062](https://github.com/hashicorp/terraform-provider-aws/issues/31062))
* resource/aws_rds_cluster: Add `db_system_id` argument to support RDS Custom engine types ([#32846](https://github.com/hashicorp/terraform-provider-aws/issues/32846))
* resource/aws_rds_cluster_instance: Add `custom_iam_instance_profile` argument to allow RDS Custom users to specify an IAM Instance Profile for the RDS Cluster Instance ([#32846](https://github.com/hashicorp/terraform-provider-aws/issues/32846))
* resource/aws_rds_cluster_instance: Update `engine` plan-time validation to allow for RDS Custom engine types ([#32846](https://github.com/hashicorp/terraform-provider-aws/issues/32846))

BUG FIXES:

* data-source/aws_vpclattice_service: Avoid listing tags when the service has been shared to the current account via AWS Resource Access Manager (RAM) ([#32939](https://github.com/hashicorp/terraform-provider-aws/issues/32939))
* data-source/aws_vpclattice_service_network: Avoid listing tags when the service network has been shared to the current account via AWS Resource Access Manager (RAM) ([#32939](https://github.com/hashicorp/terraform-provider-aws/issues/32939))
* resource/aws_appstream_fleet: Increased upper limit of `max_user_duration_in_seconds` to 432000 ([#32933](https://github.com/hashicorp/terraform-provider-aws/issues/32933))
* resource/aws_cloudfront_distribution: Don't call `UpdateDistribution` API if only tags are updated ([#32865](https://github.com/hashicorp/terraform-provider-aws/issues/32865))
* resource/aws_db_instance: Fix crash creating resource with empty `restore_to_point_in_time` configuration block ([#32928](https://github.com/hashicorp/terraform-provider-aws/issues/32928))
* resource/aws_emr_cluster: Fix to allow empty `args` for `bootstrap_action` ([#32956](https://github.com/hashicorp/terraform-provider-aws/issues/32956))
* resource/aws_emr_instance_fleet: Fix fleet deletion failing for terminated clusters ([#32866](https://github.com/hashicorp/terraform-provider-aws/issues/32866))
* resource/aws_fms_policy: Prevent erroneous diffs on `security_service_policy_data.managed_service_data` ([#32860](https://github.com/hashicorp/terraform-provider-aws/issues/32860))
* resource/aws_instance: Fix `InvalidParameterCombination: Network interfaces and an instance-level security groups may not be specified on the same request` errors creating Instances with `subnet_id` configured and `launch_template` referencing an `aws_launch_template` with configured `vpc_security_group_ids` ([#32854](https://github.com/hashicorp/terraform-provider-aws/issues/32854))
* resource/aws_lb: Fix to avoid creating a load balancer with same name as an existing load balancer ([#32941](https://github.com/hashicorp/terraform-provider-aws/issues/32941))

## 5.11.0 (August  3, 2023)

FEATURES:

* **New Resource:** `aws_sagemaker_pipeline` ([#32527](https://github.com/hashicorp/terraform-provider-aws/issues/32527))

ENHANCEMENTS:

* data-source/aws_cloudtrail_service_account: Add service account ID for `il-central-1` AWS Region ([#32840](https://github.com/hashicorp/terraform-provider-aws/issues/32840))
* data-source/aws_db_cluster_snapshot: Add `tags` argument ([#31602](https://github.com/hashicorp/terraform-provider-aws/issues/31602))
* data-source/aws_db_instance: Add ability to filter by `tags` ([#32740](https://github.com/hashicorp/terraform-provider-aws/issues/32740))
* data-source/aws_db_instances: Add ability to filter by `tags` ([#32740](https://github.com/hashicorp/terraform-provider-aws/issues/32740))
* data-source/aws_db_snapshot: Add `tags` argument ([#31600](https://github.com/hashicorp/terraform-provider-aws/issues/31600))
* data-source/aws_elb_hosted_zone_id: Add hosted zone ID for `il-central-1` AWS Region ([#32840](https://github.com/hashicorp/terraform-provider-aws/issues/32840))
* data-source/aws_lb_hosted_zone_id: Add hosted zone IDs for `il-central-1` AWS Region ([#32840](https://github.com/hashicorp/terraform-provider-aws/issues/32840))
* data-source/aws_s3_bucket: Add hosted zone ID for `il-central-1` AWS Region ([#32840](https://github.com/hashicorp/terraform-provider-aws/issues/32840))
* data-source/aws_vpclattice_service: Add ability to find by `name` ([#32177](https://github.com/hashicorp/terraform-provider-aws/issues/32177))
* resource/aws_finspace_kx_cluster: Adjusted `savedown_storage_configuration.size` minimum value to `10` GB. ([#32800](https://github.com/hashicorp/terraform-provider-aws/issues/32800))
* resource/aws_lambda_function: Add support for `python3.11` `runtime` value ([#32729](https://github.com/hashicorp/terraform-provider-aws/issues/32729))
* resource/aws_lambda_layer_version: Add support for `python3.11` `compatible_runtimes` value ([#32729](https://github.com/hashicorp/terraform-provider-aws/issues/32729))
* resource/aws_networkfirewall_rule_group: Add support for `REJECT` action in stateful rule actions ([#32746](https://github.com/hashicorp/terraform-provider-aws/issues/32746))
* resource/aws_route_table: Allow an existing local route to be adopted or imported and the target to be updated ([#32794](https://github.com/hashicorp/terraform-provider-aws/issues/32794))
* resource/aws_sagemaker_endpoint: Add `deployment_config.rolling_update_policy` argument ([#32418](https://github.com/hashicorp/terraform-provider-aws/issues/32418))
* resource/aws_sagemaker_endpoint: Make `deployment_config.blue_green_update_policy` optional ([#32418](https://github.com/hashicorp/terraform-provider-aws/issues/32418))

BUG FIXES:

* data-source/aws_ecs_task_execution: Fixed bug that incorrectly mapped the value of `container_overrides.memory` to `container_overrides.memory_reservation` ([#32793](https://github.com/hashicorp/terraform-provider-aws/issues/32793))
* resource/aws_db_instance_automated_backups_replication: Fix `unexpected state 'Pending'` errors on resource Create ([#31600](https://github.com/hashicorp/terraform-provider-aws/issues/31600))
* resource/aws_ec2_transit_gateway_vpc_attachment: Change `transit_gateway_default_route_table_association` and `transit_gateway_default_route_table_propagation` to Computed ([#32821](https://github.com/hashicorp/terraform-provider-aws/issues/32821))
* resource/aws_emr_studio_session_mapping: Fix `InvalidRequestException: IdentityId is invalid` errors reading resources created with `identity_name` ([#32416](https://github.com/hashicorp/terraform-provider-aws/issues/32416))
* resource/aws_quicksight_analysis: Fix an error related to setting the value for `definition.sheets.visuals.insight_visual.insight_configuration.computation` ([#32791](https://github.com/hashicorp/terraform-provider-aws/issues/32791))
* resource/aws_quicksight_analysis: Fixed a bug that incorrectly determined the valid `select_all_options` values for `custom_filter_configuration`, `custom_filter_list_configuration`, `filter_list_configuration`, `numeric_equality_filter`, and `numeric_range_filter` ([#32822](https://github.com/hashicorp/terraform-provider-aws/issues/32822))
* resource/aws_quicksight_dashboard: Fix an error related to setting the value for `definition.sheets.visuals.insight_visual.insight_configuration.computation` ([#32791](https://github.com/hashicorp/terraform-provider-aws/issues/32791))
* resource/aws_quicksight_template: Fix an error related to setting the value for `definition.sheets.visuals.insight_visual.insight_configuration.computation` ([#32791](https://github.com/hashicorp/terraform-provider-aws/issues/32791))
* resource/aws_quicksight_template: Fixed a bug that incorrectly determined the valid `select_all_options` values for `custom_filter_configuration`, `custom_filter_list_configuration`, `filter_list_configuration`, `numeric_equality_filter`, and `numeric_range_filter` ([#32822](https://github.com/hashicorp/terraform-provider-aws/issues/32822))
* resource/aws_sfn_state_machine: Fix `Provider produced inconsistent final plan` errors for `publish` ([#32844](https://github.com/hashicorp/terraform-provider-aws/issues/32844))

## 5.10.0 (July 27, 2023)

FEATURES:

* **New Resource:** `aws_iam_security_token_service_preferences` ([#32091](https://github.com/hashicorp/terraform-provider-aws/issues/32091))

ENHANCEMENTS:

* data-source/aws_nat_gateway: Add `secondary_allocation_ids`, `secondary_private_ip_addresses` and `secondary_private_ip_address_count` attributes ([#31778](https://github.com/hashicorp/terraform-provider-aws/issues/31778))
* data-source/aws_transfer_server: Add `structured_log_destinations` attribute ([#32654](https://github.com/hashicorp/terraform-provider-aws/issues/32654))
* resource/aws_batch_compute_environment: `compute_resources.allocation_strategy`, `compute_resources.bid_percentage`, `compute_resources.ec2_configuration.image_id_override`, `compute_resources.ec2_configuration.image_type`, `compute_resources.ec2_key_pair`, `compute_resources.image_id`, `compute_resources.instance_role`, `compute_resources.launch_template.launch_template_id`
, `compute_resources.launch_template.launch_template_name`, `compute_resources.tags` and `compute_resources.type` can now be updated in-place ([#30438](https://github.com/hashicorp/terraform-provider-aws/issues/30438))
* resource/aws_glue_job: Add `command.runtime` attribute ([#32528](https://github.com/hashicorp/terraform-provider-aws/issues/32528))
* resource/aws_grafana_workspace: Allow `grafana_version` to be updated in-place ([#32679](https://github.com/hashicorp/terraform-provider-aws/issues/32679))
* resource/aws_kms_grant: Allow usage of service principal as grantee and revoker ([#32595](https://github.com/hashicorp/terraform-provider-aws/issues/32595))
* resource/aws_medialive_channel: Adds schemas for `caption_descriptions`, `global_configuration`, `motion_graphics_configuration`, and `nielsen_configuration` support to `encoder settings` ([#32233](https://github.com/hashicorp/terraform-provider-aws/issues/32233))
* resource/aws_nat_gateway: Add `secondary_allocation_ids`, `secondary_private_ip_addresses` and `secondary_private_ip_address_count` arguments ([#31778](https://github.com/hashicorp/terraform-provider-aws/issues/31778))
* resource/aws_nat_gateway: Add configurable timeouts ([#31778](https://github.com/hashicorp/terraform-provider-aws/issues/31778))
* resource/aws_networkfirewall_firewall_policy: Add `firewall_policy.policy_variables` configuration block to support Suricata HOME_NET variable override ([#32400](https://github.com/hashicorp/terraform-provider-aws/issues/32400))
* resource/aws_sagemaker_domain: Add `default_user_settings.canvas_app_settings.workspace_settings` attribute ([#32526](https://github.com/hashicorp/terraform-provider-aws/issues/32526))
* resource/aws_sagemaker_user_profile: Add `user_settings.canvas_app_settings.workspace_settings` attribute ([#32526](https://github.com/hashicorp/terraform-provider-aws/issues/32526))
* resource/aws_transfer_server: Add `structured_log_destinations` argument ([#32654](https://github.com/hashicorp/terraform-provider-aws/issues/32654))

BUG FIXES:

* resource/aws_account_primary_contact: Correct plan-time validation of `phone_number` ([#32715](https://github.com/hashicorp/terraform-provider-aws/issues/32715))
* resource/aws_apigatewayv2_authorizer: Skip setting authorizer TTL when there are no identity sources ([#32629](https://github.com/hashicorp/terraform-provider-aws/issues/32629))
* resource/aws_elasticache_parameter_group: Remove from state on resource Read if deleted outside of Terraform ([#32669](https://github.com/hashicorp/terraform-provider-aws/issues/32669))
* resource/aws_elasticsearch_domain: Omit `ebs_options.throughput` and `ebs_options.iops` for unsupported volume types ([#32659](https://github.com/hashicorp/terraform-provider-aws/issues/32659))
* resource/aws_finspace_kx_cluster: `database.cache_configurations.db_paths` argument is now optional ([#32579](https://github.com/hashicorp/terraform-provider-aws/issues/32579))
* resource/aws_finspace_kx_cluster: `database.cache_configurations` argument is now optional ([#32579](https://github.com/hashicorp/terraform-provider-aws/issues/32579))
* resource/aws_lambda_invocation: Fix plan failing with deferred input values ([#32706](https://github.com/hashicorp/terraform-provider-aws/issues/32706))
* resource/aws_lightsail_domain_entry: Add support for `AAAA` `type` value ([#32664](https://github.com/hashicorp/terraform-provider-aws/issues/32664))
* resource/aws_opensearch_domain: Correctly handle `off_peak_window_options.off_peak_window.window_start_time` value of `00:00` ([#32716](https://github.com/hashicorp/terraform-provider-aws/issues/32716))
* resource/aws_quicksight_analysis: Fix exception thrown when setting the value for `definition.sheets.visuals.pie_chart_visual.chart_configuration.data_labels.measure_label_visibility` ([#32668](https://github.com/hashicorp/terraform-provider-aws/issues/32668))
* resource/aws_quicksight_analysis: Grid layout `optimized_view_port_width` argument changed to Optional ([#32644](https://github.com/hashicorp/terraform-provider-aws/issues/32644))
* resource/aws_quicksight_dashboard: Fix exception thrown when setting the value for `definition.sheets.visuals.pie_chart_visual.chart_configuration.data_labels.measure_label_visibility` ([#32668](https://github.com/hashicorp/terraform-provider-aws/issues/32668))
* resource/aws_quicksight_dashboard: Grid layout `optimized_view_port_width` argument changed to Optional ([#32644](https://github.com/hashicorp/terraform-provider-aws/issues/32644))
* resource/aws_quicksight_template: Fix exception thrown when setting the value for `definition.sheets.visuals.pie_chart_visual.chart_configuration.data_labels.measure_label_visibility` ([#32668](https://github.com/hashicorp/terraform-provider-aws/issues/32668))
* resource/aws_quicksight_template: Grid layout `optimized_view_port_width` argument changed to Optional ([#32644](https://github.com/hashicorp/terraform-provider-aws/issues/32644))
* resource/aws_vpclattice_access_log_subscription: Avoid recreating resource when passing a non-wildcard CloudWatch Logs log group ARN as `destination_arn` ([#32186](https://github.com/hashicorp/terraform-provider-aws/issues/32186))
* resource/aws_vpclattice_access_log_subscription: Avoid recreating resource when passing an ARN as `resource_identifier` ([#32186](https://github.com/hashicorp/terraform-provider-aws/issues/32186))
* resource/aws_vpclattice_service_network_service_association: Avoid recreating resource when passing an ARN as `service_identifier` or `service_network_identifier` ([#32658](https://github.com/hashicorp/terraform-provider-aws/issues/32658))
* resource/aws_vpclattice_service_network_vpc_association: Avoid recreating resource when passing an ARN as `service_network_identifier` ([#32658](https://github.com/hashicorp/terraform-provider-aws/issues/32658))

## 5.9.0 (July 20, 2023)

FEATURES:

* **New Resource:** `aws_workspaces_connection_alias` ([#32482](https://github.com/hashicorp/terraform-provider-aws/issues/32482))

ENHANCEMENTS:

* data-source/aws_appmesh_gateway_route: Add `path` to the `spec.http_route.action.rewrite` and `spec.http2_route.action.rewrite` configuration blocks ([#32449](https://github.com/hashicorp/terraform-provider-aws/issues/32449))
* data-source/aws_db_instance: Add `max_allocated_storage` attribute ([#32477](https://github.com/hashicorp/terraform-provider-aws/issues/32477))
* data-source/aws_ec2_host: Add `asset_id` attribute ([#32388](https://github.com/hashicorp/terraform-provider-aws/issues/32388))
* resource/aws_appmesh_gateway_route: Add `path` to the `spec.http_route.action.rewrite` and `spec.http2_route.action.rewrite` configuration blocks ([#32449](https://github.com/hashicorp/terraform-provider-aws/issues/32449))
* resource/aws_cloudformation_stack_set_instance: Added the `stack_instance_summaries` attribute to track all account and stack IDs for deployments to organizational units. ([#24523](https://github.com/hashicorp/terraform-provider-aws/issues/24523))
* resource/aws_cloudformation_stack_set_instance: Changes to `deployment_targets` now force a new resource. ([#24523](https://github.com/hashicorp/terraform-provider-aws/issues/24523))
* resource/aws_connect_queue: add delete function ([#32538](https://github.com/hashicorp/terraform-provider-aws/issues/32538))
* resource/aws_connect_routing_profile: add delete function ([#32540](https://github.com/hashicorp/terraform-provider-aws/issues/32540))
* resource/aws_db_instance: Add `backup_target` attribute ([#32609](https://github.com/hashicorp/terraform-provider-aws/issues/32609))
* resource/aws_ec2_host: Add `asset_id` argument ([#32388](https://github.com/hashicorp/terraform-provider-aws/issues/32388))
* resource/aws_ec2_traffic_mirror_filter_rule: Fix crash when updating `rule_number` ([#32594](https://github.com/hashicorp/terraform-provider-aws/issues/32594))
* resource/aws_lightsail_key_pair: Add `tags` attribute ([#32606](https://github.com/hashicorp/terraform-provider-aws/issues/32606))
* resource/aws_signer_signing_profile: Add `signing_material` attribute. ([#32414](https://github.com/hashicorp/terraform-provider-aws/issues/32414))
* resource/aws_signer_signing_profile: Update `platform_id` validation. ([#32414](https://github.com/hashicorp/terraform-provider-aws/issues/32414))
* resource/aws_wafv2_web_acl: Add `association_config` argument ([#31668](https://github.com/hashicorp/terraform-provider-aws/issues/31668))

BUG FIXES:

* data-source/aws_dms_replication_instance: Fixed bug that caused `replication_instance_private_ips`, `replication_instance_public_ips`, and `vpc_security_group_ids` to always return `null` ([#32551](https://github.com/hashicorp/terraform-provider-aws/issues/32551))
* data-source/aws_mq_broker: Fix `setting user: Invalid address to set` errors ([#32593](https://github.com/hashicorp/terraform-provider-aws/issues/32593))
* data-source/aws_vpc_endpoint: Add `dns_options.private_dns_only_for_inbound_resolver_endpoint` ([#32517](https://github.com/hashicorp/terraform-provider-aws/issues/32517))
* resource/aws_appflow_flow: Fix tasks not updating properly due to empty task being processed ([#26614](https://github.com/hashicorp/terraform-provider-aws/issues/26614))
* resource/aws_cloudformation_stack_set_instance: Fix error when deploying to organizational units with no accounts. ([#24523](https://github.com/hashicorp/terraform-provider-aws/issues/24523))
* resource/aws_cognito_user_pool: Suppress diff when `schema.string_attribute_constraints` is omitted for `String` attribute types ([#32445](https://github.com/hashicorp/terraform-provider-aws/issues/32445))
* resource/aws_config_config_rule: Prevent crash from unhandled read error ([#32520](https://github.com/hashicorp/terraform-provider-aws/issues/32520))
* resource/aws_datasync_agent: Prevent persistent diffs when `private_link_endpoint` is not explicitly configured. ([#32546](https://github.com/hashicorp/terraform-provider-aws/issues/32546))
* resource/aws_globalaccelerator_custom_routing_endpoint_group: Respect configured `endpoint_group_region` value on resource Create ([#32393](https://github.com/hashicorp/terraform-provider-aws/issues/32393))
* resource/aws_pipes_pipe: Fix `Error: setting target_parameters: Invalid address to set` errors when creating pipes with ecs task targets ([#32432](https://github.com/hashicorp/terraform-provider-aws/issues/32432))
* resource/aws_pipes_pipe: Fix `ValidationException` errors when updating pipe ([#32622](https://github.com/hashicorp/terraform-provider-aws/issues/32622))
* resource/aws_quicksight_analysis: Correctly expand comparison method ([#32285](https://github.com/hashicorp/terraform-provider-aws/issues/32285))
* resource/aws_quicksight_folder: Fix misidentification of parent folder at grandchild level or deeper ([#32592](https://github.com/hashicorp/terraform-provider-aws/issues/32592))
* resource/aws_quicksight_group_membership: Allow non `default` value for namespace ([#32494](https://github.com/hashicorp/terraform-provider-aws/issues/32494))
* resource/aws_route53_cidr_location: Fix `Value Conversion Error` errors ([#32596](https://github.com/hashicorp/terraform-provider-aws/issues/32596))
* resource/aws_wafv2_web_acl: Fixed error handling `response_inspection` parameters ([#31111](https://github.com/hashicorp/terraform-provider-aws/issues/31111))

## 5.8.0 (July 13, 2023)

ENHANCEMENTS:

* data-source/aws_ssm_parameter: Add `insecure_value` attribute ([#30817](https://github.com/hashicorp/terraform-provider-aws/issues/30817))
* resource/aws_fms_policy: Add `policy_option` attribute for `security_service_policy_data` block ([#25362](https://github.com/hashicorp/terraform-provider-aws/issues/25362))
* resource/aws_iam_virtual_mfa_device: Add `enable_date` and `user_name` attributes ([#32462](https://github.com/hashicorp/terraform-provider-aws/issues/32462))

BUG FIXES:

* resource/aws_config_config_rule: Prevent crash on nil describe output ([#32439](https://github.com/hashicorp/terraform-provider-aws/issues/32439))
* resource/aws_mq_broker: default `replication_user` to `false` ([#32454](https://github.com/hashicorp/terraform-provider-aws/issues/32454))
* resource/aws_quicksight_analysis: Fix exception thrown when specifying `definition.sheets.visuals.bar_chart_visual.chart_configuration.category_axis.scrollbar_options.visible_range` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_analysis: Fix exception thrown when specifying `definition.sheets.visuals.pivot_table_visual.chart_configuration.field_options.selected_field_options.visibility` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_analysis: Fix exception thrown when specifying `definition.sheets.visuals.pivot_table_visual.chart_configuration.field_wells.pivot_table_aggregated_field_wells.rows` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_dashboard: Fix exception thrown when specifying `definition.sheets.visuals.bar_chart_visual.chart_configuration.category_axis.scrollbar_options.visible_range` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_dashboard: Fix exception thrown when specifying `definition.sheets.visuals.pivot_table_visual.chart_configuration.field_options.selected_field_options.visibility` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_dashboard: Fix exception thrown when specifying `definition.sheets.visuals.pivot_table_visual.chart_configuration.field_wells.pivot_table_aggregated_field_wells.rows` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_template: Fix exception thrown when specifying `definition.sheets.visuals.bar_chart_visual.chart_configuration.category_axis.scrollbar_options.visible_range` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_template: Fix exception thrown when specifying `definition.sheets.visuals.pivot_table_visual.chart_configuration.field_options.selected_field_options.visibility` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))
* resource/aws_quicksight_template: Fix exception thrown when specifying `definition.sheets.visuals.pivot_table_visual.chart_configuration.field_wells.pivot_table_aggregated_field_wells.rows` ([#32464](https://github.com/hashicorp/terraform-provider-aws/issues/32464))

## 5.7.0 (July  7, 2023)

FEATURES:

* **New Data Source:** `aws_opensearchserverless_security_config` ([#32321](https://github.com/hashicorp/terraform-provider-aws/issues/32321))
* **New Data Source:** `aws_opensearchserverless_security_policy` ([#32226](https://github.com/hashicorp/terraform-provider-aws/issues/32226))
* **New Data Source:** `aws_opensearchserverless_vpc_endpoint` ([#32276](https://github.com/hashicorp/terraform-provider-aws/issues/32276))
* **New Resource:** `aws_cleanrooms_collaboration` ([#31680](https://github.com/hashicorp/terraform-provider-aws/issues/31680))

ENHANCEMENTS:

* resource/aws_aws_keyspaces_table: Add `client_side_timestamps` configuration block ([#32339](https://github.com/hashicorp/terraform-provider-aws/issues/32339))
* resource/aws_glue_catalog_database: Add `target_database.region` argument ([#32283](https://github.com/hashicorp/terraform-provider-aws/issues/32283))
* resource/aws_glue_crawler: Add `iceberg_target` configuration block ([#32332](https://github.com/hashicorp/terraform-provider-aws/issues/32332))
* resource/aws_internetmonitor_monitor: Add `health_events_config` configuration block ([#32343](https://github.com/hashicorp/terraform-provider-aws/issues/32343))
* resource/aws_lambda_function: Support `code_signing_config_arn` in the `ap-east-1` AWS Region ([#32327](https://github.com/hashicorp/terraform-provider-aws/issues/32327))
* resource/aws_qldb_stream: Add configurable Create and Delete timeouts ([#32345](https://github.com/hashicorp/terraform-provider-aws/issues/32345))
* resource/aws_service_discovery_private_dns_namespace: Allow `description` to be updated in-place ([#32342](https://github.com/hashicorp/terraform-provider-aws/issues/32342))
* resource/aws_service_discovery_public_dns_namespace: Allow `description` to be updated in-place ([#32342](https://github.com/hashicorp/terraform-provider-aws/issues/32342))
* resource/aws_timestreamwrite_table: Add `schema` configuration block ([#32354](https://github.com/hashicorp/terraform-provider-aws/issues/32354))

BUG FIXES:

* provider: Correctly handle `forbidden_account_ids` ([#32352](https://github.com/hashicorp/terraform-provider-aws/issues/32352))
* resource/aws_kms_external_key: Correctly remove all tags ([#32371](https://github.com/hashicorp/terraform-provider-aws/issues/32371))
* resource/aws_kms_key: Correctly remove all tags ([#32371](https://github.com/hashicorp/terraform-provider-aws/issues/32371))
* resource/aws_kms_replica_external_key: Correctly remove all tags ([#32371](https://github.com/hashicorp/terraform-provider-aws/issues/32371))
* resource/aws_kms_replica_key: Correctly remove all tags ([#32371](https://github.com/hashicorp/terraform-provider-aws/issues/32371))
* resource/aws_secretsmanager_secret_rotation: Fix `InvalidParameterException: You cannot specify both rotation frequency and schedule expression together` errors on resource Update ([#31915](https://github.com/hashicorp/terraform-provider-aws/issues/31915))
* resource/aws_ssm_parameter: Skip Update if only `overwrite` parameter changes ([#32372](https://github.com/hashicorp/terraform-provider-aws/issues/32372))
* resource/aws_vpc_endpoint: Fix `InvalidParameter: PrivateDnsOnlyForInboundResolverEndpoint not supported for this service` errors creating S3 _Interface_ VPC endpoints ([#32355](https://github.com/hashicorp/terraform-provider-aws/issues/32355))

## 5.6.2 (June 30, 2023)

BUG FIXES:

* resource/aws_s3_bucket: Fix `InvalidArgument: Invalid attribute name specified` errors when listing S3 Bucket objects, caused by an [AWS SDK for Go regression](https://github.com/aws/aws-sdk-go/issues/4897) ([#32317](https://github.com/hashicorp/terraform-provider-aws/issues/32317))

## 5.6.1 (June 30, 2023)

BUG FIXES:

* provider: Prevent resource recreation if `tags` or `tags_all` are updated ([#32297](https://github.com/hashicorp/terraform-provider-aws/issues/32297))

## 5.6.0 (June 29, 2023)

FEATURES:

* **New Data Source:** `aws_opensearchserverless_access_policy` ([#32231](https://github.com/hashicorp/terraform-provider-aws/issues/32231))
* **New Data Source:** `aws_opensearchserverless_collection` ([#32247](https://github.com/hashicorp/terraform-provider-aws/issues/32247))
* **New Data Source:** `aws_sfn_alias` ([#32176](https://github.com/hashicorp/terraform-provider-aws/issues/32176))
* **New Data Source:** `aws_sfn_state_machine_versions` ([#32176](https://github.com/hashicorp/terraform-provider-aws/issues/32176))
* **New Resource:** `aws_ec2_instance_connect_endpoint` ([#31858](https://github.com/hashicorp/terraform-provider-aws/issues/31858))
* **New Resource:** `aws_sfn_alias` ([#32176](https://github.com/hashicorp/terraform-provider-aws/issues/32176))
* **New Resource:** `aws_transfer_agreement` ([#32203](https://github.com/hashicorp/terraform-provider-aws/issues/32203))
* **New Resource:** `aws_transfer_certificate` ([#32203](https://github.com/hashicorp/terraform-provider-aws/issues/32203))
* **New Resource:** `aws_transfer_connector` ([#32203](https://github.com/hashicorp/terraform-provider-aws/issues/32203))
* **New Resource:** `aws_transfer_profile` ([#32203](https://github.com/hashicorp/terraform-provider-aws/issues/32203))

ENHANCEMENTS:

* resource/aws_batch_compute_environment: Add `placement_group` attribute to the `compute_resources` configuration block ([#32200](https://github.com/hashicorp/terraform-provider-aws/issues/32200))
* resource/aws_emrserverless_application: Do not recreate the resource if `release_label` changes ([#32278](https://github.com/hashicorp/terraform-provider-aws/issues/32278))
* resource/aws_fis_experiment_template: Add `log_configuration` configuration block ([#32102](https://github.com/hashicorp/terraform-provider-aws/issues/32102))
* resource/aws_fis_experiment_template: Add `parameters` attribute to the `target` configuration block ([#32160](https://github.com/hashicorp/terraform-provider-aws/issues/32160))
* resource/aws_fis_experiment_template: Add support for `Pods` and `Tasks` to `action.*.target` ([#32152](https://github.com/hashicorp/terraform-provider-aws/issues/32152))
* resource/aws_lambda_event_source_mapping: The `queues` argument has changed from a set to a list with a maximum of one element. ([#31931](https://github.com/hashicorp/terraform-provider-aws/issues/31931))
* resource/aws_pipes_pipe: Add `activemq_broker_parameters`, `dynamodb_stream_parameters`, `kinesis_stream_parameters`, `managed_streaming_kafka_parameters`, `rabbitmq_broker_parameters`, `self_managed_kafka_parameters` and `sqs_queue_parameters` attributes to the `source_parameters` configuration block. NOTE: Because we cannot easily test all this functionality, it is best effort and we ask for community help in testing ([#31607](https://github.com/hashicorp/terraform-provider-aws/issues/31607))
* resource/aws_pipes_pipe: Add `batch_job_parameters`, `cloudwatch_logs_parameters`, `ecs_task_parameters`, `eventbridge_event_bus_parameters`, `http_parameters`, `kinesis_stream_parameters`, `lambda_function_parameters`, `redshift_data_parameters`, `sagemaker_pipeline_parameters`, `sqs_queue_parameters` and `step_function_state_machine_parameters` attributes to the `target_parameters` configuration block. NOTE: Because we cannot easily test all this functionality, it is best effort and we ask for community help in testing ([#31607](https://github.com/hashicorp/terraform-provider-aws/issues/31607))
* resource/aws_pipes_pipe: Add `enrichment_parameters` argument ([#31607](https://github.com/hashicorp/terraform-provider-aws/issues/31607))
* resource/aws_resourcegroups_group: `resource_query` no longer conflicts with `configuration` ([#30242](https://github.com/hashicorp/terraform-provider-aws/issues/30242))
* resource/aws_s3_bucket_logging: Retry on empty read of logging config ([#30916](https://github.com/hashicorp/terraform-provider-aws/issues/30916))
* resource/aws_sfn_state_machine: Add `description`, `publish`, `revision_id`, `state_machine_version_arn` and `version_description` attributes ([#32176](https://github.com/hashicorp/terraform-provider-aws/issues/32176))

BUG FIXES:

* resource/aws_db_instance: Fix resource Create returning instances not in the `available` state when `identifier_prefix` is specified ([#32287](https://github.com/hashicorp/terraform-provider-aws/issues/32287))
* resource/aws_resourcegroups_resource: Fix crash when resource Create fails ([#30242](https://github.com/hashicorp/terraform-provider-aws/issues/30242))
* resource/aws_route: Fix `reading Route in Route Table (rtb-1234abcd) with destination (1.2.3.4/5): couldn't find resource` errors when reading new resource ([#32196](https://github.com/hashicorp/terraform-provider-aws/issues/32196))
* resource/aws_vpc_security_group_egress_rule: `security_group_id` is Required ([#32148](https://github.com/hashicorp/terraform-provider-aws/issues/32148))
* resource/aws_vpc_security_group_ingress_rule: `security_group_id` is Required ([#32148](https://github.com/hashicorp/terraform-provider-aws/issues/32148))

## 5.5.0 (June 23, 2023)

NOTES:

* provider: Updates to Go 1.20, the last release that will run on any release of Windows 7, 8, Server 2008 and Server 2012. A future release will update to Go 1.21, and these platforms will no longer be supported. ([#32108](https://github.com/hashicorp/terraform-provider-aws/issues/32108))
* provider: Updates to Go 1.20, the last release that will run on macOS 10.13 High Sierra or 10.14 Mojave. A future release will update to Go 1.21, and these platforms will no longer be supported. ([#32108](https://github.com/hashicorp/terraform-provider-aws/issues/32108))
* provider: Updates to Go 1.20. The provider will now notice the `trust-ad` option in `/etc/resolv.conf` and, if set, will set the "authentic data" option in outgoing DNS requests in order to better match the behavior of the GNU libc resolver. ([#32108](https://github.com/hashicorp/terraform-provider-aws/issues/32108))

FEATURES:

* **New Data Source:** `aws_sesv2_email_identity` ([#32026](https://github.com/hashicorp/terraform-provider-aws/issues/32026))
* **New Data Source:** `aws_sesv2_email_identity_mail_from_attributes` ([#32026](https://github.com/hashicorp/terraform-provider-aws/issues/32026))
* **New Resource:** `aws_chimesdkvoice_sip_rule` ([#32070](https://github.com/hashicorp/terraform-provider-aws/issues/32070))
* **New Resource:** `aws_organizations_resource_policy` ([#32056](https://github.com/hashicorp/terraform-provider-aws/issues/32056))

ENHANCEMENTS:

* data-source/aws_organizations_organization: Return the full set of attributes when running as a [delegated administrator for AWS Organizations](https://docs.aws.amazon.com/organizations/latest/userguide/orgs_delegate_policies.html) ([#32056](https://github.com/hashicorp/terraform-provider-aws/issues/32056))
* provider: Mask all sensitive values that appear when `TF_LOG` level is `TRACE` ([#32174](https://github.com/hashicorp/terraform-provider-aws/issues/32174))
* resource/aws_config_configuration_recorder: Add `exclusion_by_resource_types` and `recording_strategy` attributes to the `recording_group` configuration block ([#32007](https://github.com/hashicorp/terraform-provider-aws/issues/32007))
* resource/aws_datasync_task: Add `object_tags` attribute to `options` configuration block ([#27811](https://github.com/hashicorp/terraform-provider-aws/issues/27811))
* resource/aws_networkmanager_attachment_accepter: Added support for Transit Gateway route table attachments ([#32023](https://github.com/hashicorp/terraform-provider-aws/issues/32023))
* resource/aws_ses_active_receipt_rule_set: Support import ([#27604](https://github.com/hashicorp/terraform-provider-aws/issues/27604))

BUG FIXES:

* resource/aws_api_gateway_rest_api: Fix crash when `binary_media_types` is `null` ([#32169](https://github.com/hashicorp/terraform-provider-aws/issues/32169))
* resource/aws_datasync_location_object_storage: Don't ignore `server_certificate` argument ([#27811](https://github.com/hashicorp/terraform-provider-aws/issues/27811))
* resource/aws_eip: Fix `reading EC2 EIP (eipalloc-abcd1234): couldn't find resource` errors when reading new resource ([#32016](https://github.com/hashicorp/terraform-provider-aws/issues/32016))
* resource/aws_quicksight_analysis: Fix schema mapping for string set elements ([#31903](https://github.com/hashicorp/terraform-provider-aws/issues/31903))
* resource/aws_redshiftserverless_workgroup: Fix `waiting for completion: unexpected state 'AVAILABLE'` errors when deleting resource ([#32067](https://github.com/hashicorp/terraform-provider-aws/issues/32067))
* resource/aws_route_table: Fix `reading Route Table (rtb-abcd1234): couldn't find resource` errors when reading new resource ([#30999](https://github.com/hashicorp/terraform-provider-aws/issues/30999))
* resource/aws_storagegateway_smb_file_share: Fix update error when `kms_encrypted` is `true` but `kms_key_arn` is not sent in the request ([#32171](https://github.com/hashicorp/terraform-provider-aws/issues/32171))

## 5.4.0 (June 15, 2023)

FEATURES:

* **New Data Source:** `aws_organizations_policies` ([#31545](https://github.com/hashicorp/terraform-provider-aws/issues/31545))
* **New Data Source:** `aws_organizations_policies_for_target` ([#31682](https://github.com/hashicorp/terraform-provider-aws/issues/31682))
* **New Resource:** `aws_chimesdkvoice_sip_media_application` ([#31937](https://github.com/hashicorp/terraform-provider-aws/issues/31937))
* **New Resource:** `aws_opensearchserverless_collection` ([#31091](https://github.com/hashicorp/terraform-provider-aws/issues/31091))
* **New Resource:** `aws_opensearchserverless_security_config` ([#28776](https://github.com/hashicorp/terraform-provider-aws/issues/28776))
* **New Resource:** `aws_opensearchserverless_vpc_endpoint` ([#28651](https://github.com/hashicorp/terraform-provider-aws/issues/28651))

ENHANCEMENTS:

* resource/aws_elb: Add configurable Create and Update timeouts ([#31976](https://github.com/hashicorp/terraform-provider-aws/issues/31976))
* resource/aws_glue_data_quality_ruleset: Add `catalog_id` argument to `target_table` block ([#31926](https://github.com/hashicorp/terraform-provider-aws/issues/31926))

BUG FIXES:

* provider: Fix `index out of range [0] with length 0` panic ([#32004](https://github.com/hashicorp/terraform-provider-aws/issues/32004))
* resource/aws_elb: Recreate the resource if `subnets` is updated to an empty list ([#31976](https://github.com/hashicorp/terraform-provider-aws/issues/31976))
* resource/aws_lambda_provisioned_concurrency_config: The `function_name` argument now properly handles ARN values ([#31933](https://github.com/hashicorp/terraform-provider-aws/issues/31933))
* resource/aws_quicksight_data_set: Allow physical table map to be optional ([#31863](https://github.com/hashicorp/terraform-provider-aws/issues/31863))
* resource/aws_ssm_default_patch_baseline: Fix `*conns.AWSClient is not ssm.ssmClient: missing method SSMClient` panic ([#31928](https://github.com/hashicorp/terraform-provider-aws/issues/31928))

## 5.3.0 (June 13, 2023)

NOTES:

* resource/aws_instance: The `metadata_options.http_endpoint` argument now correctly defaults to `enabled`. ([#24774](https://github.com/hashicorp/terraform-provider-aws/issues/24774))
* resource/aws_lambda_function: The `replace_security_groups_on_destroy` and `replacement_security_group_ids` attributes are being deprecated as AWS no longer supports this operation. These attributes now have no effect, and will be removed in a future major version. ([#31904](https://github.com/hashicorp/terraform-provider-aws/issues/31904))

FEATURES:

* **New Data Source:** `aws_quicksight_theme` ([#31900](https://github.com/hashicorp/terraform-provider-aws/issues/31900))
* **New Resource:** `aws_opensearchserverless_access_policy` ([#28518](https://github.com/hashicorp/terraform-provider-aws/issues/28518))
* **New Resource:** `aws_opensearchserverless_security_policy` ([#28470](https://github.com/hashicorp/terraform-provider-aws/issues/28470))
* **New Resource:** `aws_quicksight_theme` ([#31900](https://github.com/hashicorp/terraform-provider-aws/issues/31900))

ENHANCEMENTS:

* data-source/aws_redshift_cluster: Add `cluster_namespace_arn` attribute ([#31884](https://github.com/hashicorp/terraform-provider-aws/issues/31884))
* resource/aws_redshift_cluster: Add `cluster_namespace_arn` attribute ([#31884](https://github.com/hashicorp/terraform-provider-aws/issues/31884))
* resource/aws_vpc_endpoint: Add `private_dns_only_for_inbound_resolver_endpoint` attribute to the `dns_options` configuration block ([#31873](https://github.com/hashicorp/terraform-provider-aws/issues/31873))

BUG FIXES:

* resource/aws_ecs_task_definition: Fix to prevent persistent diff when `efs_volume_configuration` has both `root_volume` and `authorization_config` set. ([#26880](https://github.com/hashicorp/terraform-provider-aws/issues/26880))
* resource/aws_instance: Fix default for `metadata_options.http_endpoint` argument. ([#24774](https://github.com/hashicorp/terraform-provider-aws/issues/24774))
* resource/aws_keyspaces_keyspace: Correct plan time validation for `name` ([#31352](https://github.com/hashicorp/terraform-provider-aws/issues/31352))
* resource/aws_keyspaces_table: Correct plan time validation for `keyspace_name`, `table_name` and column names ([#31352](https://github.com/hashicorp/terraform-provider-aws/issues/31352))
* resource/aws_quicksight_analysis: Fix assignment of KPI visual field well target values ([#31901](https://github.com/hashicorp/terraform-provider-aws/issues/31901))
* resource/aws_redshift_cluster: Allow `availability_zone_relocation_enabled` to be `true` when `publicly_accessible` is `true` ([#31886](https://github.com/hashicorp/terraform-provider-aws/issues/31886))
* resource/aws_vpc: Fix `reading EC2 VPC (vpc-abcd1234) Attribute (enableDnsSupport): couldn't find resource` errors when reading new resource ([#31877](https://github.com/hashicorp/terraform-provider-aws/issues/31877))

## 5.2.0 (June  9, 2023)

NOTES:

* resource/aws_mwaa_environment: Upgrading your environment to a new major version of Apache Airflow forces replacement of the resource ([#31833](https://github.com/hashicorp/terraform-provider-aws/issues/31833))

FEATURES:

* **New Data Source:** `aws_budgets_budget` ([#31691](https://github.com/hashicorp/terraform-provider-aws/issues/31691))
* **New Data Source:** `aws_ecr_pull_through_cache_rule` ([#31696](https://github.com/hashicorp/terraform-provider-aws/issues/31696))
* **New Data Source:** `aws_guardduty_finding_ids` ([#31711](https://github.com/hashicorp/terraform-provider-aws/issues/31711))
* **New Data Source:** `aws_iam_principal_policy_simulation` ([#25569](https://github.com/hashicorp/terraform-provider-aws/issues/25569))
* **New Resource:** `aws_chimesdkvoice_global_settings` ([#31365](https://github.com/hashicorp/terraform-provider-aws/issues/31365))
* **New Resource:** `aws_finspace_kx_cluster` ([#31806](https://github.com/hashicorp/terraform-provider-aws/issues/31806))
* **New Resource:** `aws_finspace_kx_database` ([#31803](https://github.com/hashicorp/terraform-provider-aws/issues/31803))
* **New Resource:** `aws_finspace_kx_environment` ([#31802](https://github.com/hashicorp/terraform-provider-aws/issues/31802))
* **New Resource:** `aws_finspace_kx_user` ([#31804](https://github.com/hashicorp/terraform-provider-aws/issues/31804))

ENHANCEMENTS:

* data/aws_ec2_transit_gateway_connect_peer: Add `bgp_peer_address` and `bgp_transit_gateway_addresses` attributes ([#31752](https://github.com/hashicorp/terraform-provider-aws/issues/31752))
* provider: Adds `retry_mode` parameter ([#31745](https://github.com/hashicorp/terraform-provider-aws/issues/31745))
* resource/aws_chime_voice_connector: Add tagging support ([#31746](https://github.com/hashicorp/terraform-provider-aws/issues/31746))
* resource/aws_ec2_transit_gateway_connect_peer: Add `bgp_peer_address` and `bgp_transit_gateway_addresses` attributes ([#31752](https://github.com/hashicorp/terraform-provider-aws/issues/31752))
* resource/aws_ec2_transit_gateway_route_table_association: Add `replace_existing_association` argument ([#31452](https://github.com/hashicorp/terraform-provider-aws/issues/31452))
* resource/aws_fis_experiment_template: Add support for `Volumes` to `actions.*.target` ([#31499](https://github.com/hashicorp/terraform-provider-aws/issues/31499))
* resource/aws_instance: Add `instance_market_options` configuration block and `instance_lifecycle` and `spot_instance_request_id` attributes ([#31495](https://github.com/hashicorp/terraform-provider-aws/issues/31495))
* resource/aws_lambda_function: Add support for `ruby3.2` `runtime` value ([#31842](https://github.com/hashicorp/terraform-provider-aws/issues/31842))
* resource/aws_lambda_layer_version: Add support for `ruby3.2` `compatible_runtimes` value ([#31842](https://github.com/hashicorp/terraform-provider-aws/issues/31842))
* resource/aws_mwaa_environment: Consider `CREATING_SNAPSHOT` a valid pending state for resource update ([#31833](https://github.com/hashicorp/terraform-provider-aws/issues/31833))
* resource/aws_networkfirewall_firewall_policy: Add `stream_exception_policy` option to `firewall_policy.stateful_engine_options` ([#31541](https://github.com/hashicorp/terraform-provider-aws/issues/31541))
* resource/aws_redshiftserverless_workgroup: Additional supported values for `config_parameter.parameter_key` ([#31747](https://github.com/hashicorp/terraform-provider-aws/issues/31747))
* resource/aws_sagemaker_model: Add `container.model_package_name` and `primary_container.model_package_name` arguments ([#31755](https://github.com/hashicorp/terraform-provider-aws/issues/31755))

BUG FIXES:

* data-source/aws_redshift_cluster: Fix crash reading clusters in `modifying` state ([#31772](https://github.com/hashicorp/terraform-provider-aws/issues/31772))
* provider/default_tags: Fix perpetual diff when identical tags are moved from `default_tags` to resource `tags`, and vice versa ([#31826](https://github.com/hashicorp/terraform-provider-aws/issues/31826))
* resource/aws_autoscaling_group: Ignore any `Failed` scaling activities due to IAM eventual consistency ([#31282](https://github.com/hashicorp/terraform-provider-aws/issues/31282))
* resource/aws_dx_connection: Convert `vlan_id` from [`TypeString`](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types#typestring) to [`TypeInt`](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-types#typeint) in [Terraform state](https://developer.hashicorp.com/terraform/language/state) for existing resources. This fixes a regression introduced in [v5.1.0](https://github.com/hashicorp/terraform-provider-aws/blob/main/CHANGELOG.md#510-june--1-2023) causing `a number is required` errors ([#31735](https://github.com/hashicorp/terraform-provider-aws/issues/31735))
* resource/aws_globalaccelerator_endpoint_group: Fix bug updating `endpoint_configuration.weight` to `0` ([#31767](https://github.com/hashicorp/terraform-provider-aws/issues/31767))
* resource/aws_medialive_channel: Fix spelling in `hls_cdn_settings` expander. ([#31844](https://github.com/hashicorp/terraform-provider-aws/issues/31844))
* resource/aws_redshiftserverless_namespace: Fix perpetual `iam_roles` diffs when the namespace contains a workgroup ([#31749](https://github.com/hashicorp/terraform-provider-aws/issues/31749))
* resource/aws_redshiftserverless_workgroup: Change `config_parameter` from `TypeList` to `TypeSet` as order is not significant ([#31747](https://github.com/hashicorp/terraform-provider-aws/issues/31747))
* resource/aws_redshiftserverless_workgroup: Fix `ValidationException: Can't update multiple configurations at the same time` errors ([#31747](https://github.com/hashicorp/terraform-provider-aws/issues/31747))
* resource/aws_vpc_endpoint: Fix tagging error preventing use in ISO partitions ([#31801](https://github.com/hashicorp/terraform-provider-aws/issues/31801))

## 5.1.0 (June  1, 2023)

BREAKING CHANGES:

* resource/aws_iam_role: The `role_last_used` attribute has been removed. Use the `aws_iam_role` data source instead. ([#31656](https://github.com/hashicorp/terraform-provider-aws/issues/31656))

NOTES:

* resource/aws_autoscaling_group: The `load_balancers` and `target_group_arns` attributes have been changed to `Computed`. This means that omitting this argument is interpreted as ignoring any existing load balancer or target group attachments. To remove all load balancer or target group attachments an empty list should be specified. ([#31527](https://github.com/hashicorp/terraform-provider-aws/issues/31527))
* resource/aws_iam_role: The `role_last_used` attribute has been removed. Use the `aws_iam_role` data source instead. See the community feedback provided in the [linked issue](https://github.com/hashicorp/terraform-provider-aws/issues/30861) for additional justification on this change. As the attribute is read-only, unlikely to be used as an input to another resource, and available in the corresponding data source, a breaking change in a minor version was deemed preferable to a long deprecation/removal cycle in this circumstance. ([#31656](https://github.com/hashicorp/terraform-provider-aws/issues/31656))
* resource/aws_redshift_cluster: Ignores the parameter `aqua_configuration_status`, since the AWS API ignores it. Now always returns `auto`. ([#31612](https://github.com/hashicorp/terraform-provider-aws/issues/31612))

FEATURES:

* **New Data Source:** `aws_vpclattice_resource_policy` ([#31372](https://github.com/hashicorp/terraform-provider-aws/issues/31372))
* **New Resource:** `aws_autoscaling_traffic_source_attachment` ([#31527](https://github.com/hashicorp/terraform-provider-aws/issues/31527))
* **New Resource:** `aws_emrcontainers_job_template` ([#31399](https://github.com/hashicorp/terraform-provider-aws/issues/31399))
* **New Resource:** `aws_glue_data_quality_ruleset` ([#31604](https://github.com/hashicorp/terraform-provider-aws/issues/31604))
* **New Resource:** `aws_quicksight_analysis` ([#31542](https://github.com/hashicorp/terraform-provider-aws/issues/31542))
* **New Resource:** `aws_quicksight_dashboard` ([#31448](https://github.com/hashicorp/terraform-provider-aws/issues/31448))
* **New Resource:** `aws_resourcegroups_resource` ([#31430](https://github.com/hashicorp/terraform-provider-aws/issues/31430))

ENHANCEMENTS:

* data-source/aws_autoscaling_group: Add `traffic_source` attribute ([#31527](https://github.com/hashicorp/terraform-provider-aws/issues/31527))
* data-source/aws_opensearch_domain: Add `off_peak_window_options` attribute ([#30965](https://github.com/hashicorp/terraform-provider-aws/issues/30965))
* provider: Increases size of HTTP request bodies in logs to 1 KB ([#31718](https://github.com/hashicorp/terraform-provider-aws/issues/31718))
* resource/aws_appsync_graphql_api: Add `visibility` argument ([#31369](https://github.com/hashicorp/terraform-provider-aws/issues/31369))
* resource/aws_appsync_graphql_api: Add plan time validation for `log_config.cloudwatch_logs_role_arn` ([#31369](https://github.com/hashicorp/terraform-provider-aws/issues/31369))
* resource/aws_autoscaling_group: Add `traffic_source` configuration block ([#31527](https://github.com/hashicorp/terraform-provider-aws/issues/31527))
* resource/aws_cloudformation_stack_set: Add `managed_execution` argument ([#25210](https://github.com/hashicorp/terraform-provider-aws/issues/25210))
* resource/aws_fsx_ontap_volume: Add `skip_final_backup` argument ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_fsx_ontap_volume: Remove default value for `security_style` argument and mark as Computed ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_fsx_ontap_volume: Update `ontap_volume_type` attribute to be configurable ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_fsx_ontap_volume: `junction_path` is Optional ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_fsx_ontap_volume: `storage_efficiency_enabled` is Optional ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_grafana_workspace: Increase default Create and Update timeouts to 30 minutes ([#31422](https://github.com/hashicorp/terraform-provider-aws/issues/31422))
* resource/aws_lambda_invocation: Add lifecycle_scope CRUD to invoke on each resource state transition ([#29367](https://github.com/hashicorp/terraform-provider-aws/issues/29367))
* resource/aws_lambda_layer_version_permission: Add `skip_destroy` attribute ([#29571](https://github.com/hashicorp/terraform-provider-aws/issues/29571))
* resource/aws_lambda_provisioned_concurrency_configuration: Add `skip_destroy` argument ([#31646](https://github.com/hashicorp/terraform-provider-aws/issues/31646))
* resource/aws_opensearch_domain: Add `off_peak_window_options` configuration block ([#30965](https://github.com/hashicorp/terraform-provider-aws/issues/30965))
* resource/aws_sagemaker_endpoint_configuration: Add  and `shadow_production_variants.serverless_config.provisioned_concurrency` arguments ([#31398](https://github.com/hashicorp/terraform-provider-aws/issues/31398))
* resource/aws_transfer_server: Add support for `TransferSecurityPolicy-2023-05` `security_policy_name` value ([#31536](https://github.com/hashicorp/terraform-provider-aws/issues/31536))

BUG FIXES:

* data-source/aws_dx_connection: Fix the `vlan_id` being returned as null ([#31480](https://github.com/hashicorp/terraform-provider-aws/issues/31480))
* provider/tags: Fix crash when some `tags` are `null` and others are `computed` ([#31687](https://github.com/hashicorp/terraform-provider-aws/issues/31687))
* provider: Limits size of HTTP response bodies in logs to 4 KB ([#31718](https://github.com/hashicorp/terraform-provider-aws/issues/31718))
* resource/aws_autoscaling_group: Fix `The AutoRollback parameter cannot be set to true when the DesiredConfiguration parameter is empty` errors when refreshing instances ([#31715](https://github.com/hashicorp/terraform-provider-aws/issues/31715))
* resource/aws_autoscaling_group: Now ignores previous failed scaling activities ([#31551](https://github.com/hashicorp/terraform-provider-aws/issues/31551))
* resource/aws_cloudfront_distribution: Remove the upper limit on `origin_keepalive_timeout` ([#31608](https://github.com/hashicorp/terraform-provider-aws/issues/31608))
* resource/aws_connect_instance: Fix crash when reading instances with `CREATION_FAILED` status ([#31689](https://github.com/hashicorp/terraform-provider-aws/issues/31689))
* resource/aws_connect_security_profile: Set correct `tags` in state ([#31716](https://github.com/hashicorp/terraform-provider-aws/issues/31716))
* resource/aws_dx_connection: Fix the `vlan_id` being returned as null ([#31480](https://github.com/hashicorp/terraform-provider-aws/issues/31480))
* resource/aws_ecs_service: Fix crash when just `alarms` is updated ([#31683](https://github.com/hashicorp/terraform-provider-aws/issues/31683))
* resource/aws_fsx_ontap_volume: Change `storage_virtual_machine_id` to [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_fsx_ontap_volume: Change `volume_type` to [ForceNew](https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors#forcenew) ([#31544](https://github.com/hashicorp/terraform-provider-aws/issues/31544))
* resource/aws_kendra_index: Persist `user_group_resolution_mode` value to state after creation ([#31669](https://github.com/hashicorp/terraform-provider-aws/issues/31669))
* resource/aws_medialive_channel: Fix attribute spelling in `hls_cdn_settings` expand ([#31647](https://github.com/hashicorp/terraform-provider-aws/issues/31647))
* resource/aws_quicksight_data_set: Fix join_instruction not applied when creating dataset ([#31424](https://github.com/hashicorp/terraform-provider-aws/issues/31424))
* resource/aws_quicksight_data_set: Ignore failure to read refresh properties for non-SPICE datasets ([#31488](https://github.com/hashicorp/terraform-provider-aws/issues/31488))
* resource/aws_rbin_rule: Fix crash when multiple `resource_tags` blocks are configured ([#31393](https://github.com/hashicorp/terraform-provider-aws/issues/31393))
* resource/aws_rds_cluster: Correctly update `db_cluster_instance_class` ([#31709](https://github.com/hashicorp/terraform-provider-aws/issues/31709))
* resource/aws_redshift_cluster: No longer errors on deletion when status is `Maintenance` ([#31612](https://github.com/hashicorp/terraform-provider-aws/issues/31612))
* resource/aws_route53_vpc_association_authorization: Fix `ConcurrentModification` error ([#31588](https://github.com/hashicorp/terraform-provider-aws/issues/31588))
* resource/aws_s3_bucket_replication_configuration: Replication configs sometimes need more than a second or two. This resolves a race condition and adds retry logic when reading them. ([#30995](https://github.com/hashicorp/terraform-provider-aws/issues/30995))

## 5.0.1 (May 26, 2023)

BUG FIXES:

* provider/tags: Fix crash when tags are `null` ([#31587](https://github.com/hashicorp/terraform-provider-aws/issues/31587))

## 5.0.0 (May 25, 2023)

BREAKING CHANGES:

* data-source/aws_api_gateway_rest_api: `minimum_compression_size` is now a string type to allow values set via the `body` attribute to be properly computed. ([#30969](https://github.com/hashicorp/terraform-provider-aws/issues/30969))
* data-source/aws_connect_hours_of_operation: The `hours_of_operation_arn` attribute has been removed ([#31484](https://github.com/hashicorp/terraform-provider-aws/issues/31484))
* data-source/aws_db_instance: With the retirement of EC2-Classic the `db_security_groups` attribute has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* data-source/aws_elasticache_cluster: With the retirement of EC2-Classic the `security_group_names` attribute has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* data-source/aws_elasticache_replication_group: Remove `number_cache_clusters`, `replication_group_description` arguments -- use `num_cache_clusters`, and `description`, respectively, instead ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* data-source/aws_iam_policy_document: Don't add empty `statement.sid` values to `json` attribute value ([#28539](https://github.com/hashicorp/terraform-provider-aws/issues/28539))
* data-source/aws_iam_policy_document: `source_json` and `override_json` have been removed -- use `source_policy_documents` and `override_policy_documents`, respectively, instead ([#30829](https://github.com/hashicorp/terraform-provider-aws/issues/30829))
* data-source/aws_identitystore_group: The `filter` argument has been removed ([#31312](https://github.com/hashicorp/terraform-provider-aws/issues/31312))
* data-source/aws_identitystore_user: The `filter` argument has been removed ([#31312](https://github.com/hashicorp/terraform-provider-aws/issues/31312))
* data-source/aws_launch_configuration: With the retirement of EC2-Classic the `vpc_classic_link_id` and `vpc_classic_link_security_groups` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* data-source/aws_redshift_cluster: With the retirement of EC2-Classic the `cluster_security_groups` attribute has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* data-source/aws_secretsmanager_secret: The `rotation_enabled`, `rotation_lambda_arn` and `rotation_rules` attributes have been removed ([#31487](https://github.com/hashicorp/terraform-provider-aws/issues/31487))
* data-source/aws_vpc_peering_connection: With the retirement of EC2-Classic the `allow_classic_link_to_remote_vpc` and `allow_vpc_to_remote_classic_link` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* provider: The `assume_role.duration_seconds`, `assume_role_with_web_identity.duration_seconds`, `s3_force_path_style`, `shared_credentials_file` and `skip_get_ec2_platforms` attributes have been removed ([#31155](https://github.com/hashicorp/terraform-provider-aws/issues/31155))
* provider: The `aws_subnet_ids` data source has been removed ([#31140](https://github.com/hashicorp/terraform-provider-aws/issues/31140))
* provider: With the retirement of EC2-Classic the `aws_db_security_group` resource has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* provider: With the retirement of EC2-Classic the `aws_elasticache_security_group` resource has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* provider: With the retirement of EC2-Classic the `aws_redshift_security_group` resource has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* provider: With the retirement of Macie Classic the `aws_macie_member_account_association` resource has been removed ([#31058](https://github.com/hashicorp/terraform-provider-aws/issues/31058))
* provider: With the retirement of Macie Classic the `aws_macie_s3_bucket_association` resource has been removed ([#31058](https://github.com/hashicorp/terraform-provider-aws/issues/31058))
* resource/aws_acmpca_certificate_authority: The `status` attribute has been removed ([#31084](https://github.com/hashicorp/terraform-provider-aws/issues/31084))
* resource/aws_api_gateway_rest_api: `minimum_compression_size` is now a string type to allow values set via the `body` attribute to be properly computed. ([#30969](https://github.com/hashicorp/terraform-provider-aws/issues/30969))
* resource/aws_autoscaling_attachment: `alb_target_group_arn` has been removed -- use `lb_target_group_arn` instead ([#30828](https://github.com/hashicorp/terraform-provider-aws/issues/30828))
* resource/aws_autoscaling_group: Remove deprecated `tags` attribute ([#30842](https://github.com/hashicorp/terraform-provider-aws/issues/30842))
* resource/aws_budgets_budget: The `cost_filters` attribute has been removed ([#31395](https://github.com/hashicorp/terraform-provider-aws/issues/31395))
* resource/aws_ce_anomaly_subscription: The `threshold` attribute has been removed ([#30374](https://github.com/hashicorp/terraform-provider-aws/issues/30374))
* resource/aws_cloudwatch_event_target: The `ecs_target.propagate_tags` attribute now has no default value ([#25233](https://github.com/hashicorp/terraform-provider-aws/issues/25233))
* resource/aws_codebuild_project: The `secondary_sources.auth` and `source.auth` attributes have been removed ([#31483](https://github.com/hashicorp/terraform-provider-aws/issues/31483))
* resource/aws_connect_hours_of_operation: The `hours_of_operation_arn` attribute has been removed ([#31484](https://github.com/hashicorp/terraform-provider-aws/issues/31484))
* resource/aws_connect_queue: The `quick_connect_ids_associated` attribute has been removed ([#31376](https://github.com/hashicorp/terraform-provider-aws/issues/31376))
* resource/aws_connect_routing_profile: The `queue_configs_associated` attribute has been removed ([#31376](https://github.com/hashicorp/terraform-provider-aws/issues/31376))
* resource/aws_db_instance: Remove `name` - use `db_name` instead ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_db_instance: With the retirement of EC2-Classic the `security_group_names` attribute has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_db_instance: `id` is no longer the AWS database `identifier` - `id` is now the `dbi-resource-id`. Refer to `identifier` instead of `id` to use the database's identifier ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_default_vpc: With the retirement of EC2-Classic the `enable_classiclink` and `enable_classiclink_dns_support` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_dms_endpoint: `s3_settings.ignore_headers_row` has been removed ([#30452](https://github.com/hashicorp/terraform-provider-aws/issues/30452))
* resource/aws_docdb_cluster: `snapshot_identifier` change now properly forces replacement ([#29409](https://github.com/hashicorp/terraform-provider-aws/issues/29409))
* resource/aws_ec2_client_vpn_endpoint: The `status` attribute has been removed ([#31223](https://github.com/hashicorp/terraform-provider-aws/issues/31223))
* resource/aws_ec2_client_vpn_network_association: The `security_groups` attribute has been removed ([#31396](https://github.com/hashicorp/terraform-provider-aws/issues/31396))
* resource/aws_ec2_client_vpn_network_association: The `status` attribute has been removed ([#31223](https://github.com/hashicorp/terraform-provider-aws/issues/31223))
* resource/aws_ecs_cluster: The `capacity_providers` and `default_capacity_provider_strategy` attributes have been removed ([#31346](https://github.com/hashicorp/terraform-provider-aws/issues/31346))
* resource/aws_eip: With the retirement of EC2-Classic the `standard` domain is no longer supported ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_eip_association: With the retirement of EC2-Classic the `standard` domain is no longer supported ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_elasticache_cluster: With the retirement of EC2-Classic the `security_group_names` attribute has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_elasticache_replication_group: Remove `availability_zones`, `number_cache_clusters`, `replication_group_description` arguments -- use `preferred_cache_cluster_azs`, `num_cache_clusters`, and `description`, respectively, instead ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* resource/aws_elasticache_replication_group: Remove `cluster_mode` configuration block -- use top-level `num_node_groups` and `replicas_per_node_group` instead ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* resource/aws_kinesis_firehose_delivery_stream: Remove `s3_configuration` attribute from the root of the resource. `s3_configuration` is now a part of the following blocks: `elasticsearch_configuration`, `opensearch_configuration`, `redshift_configuration`, `splunk_configuration`, and `http_endpoint_configuration` ([#31138](https://github.com/hashicorp/terraform-provider-aws/issues/31138))
* resource/aws_kinesis_firehose_delivery_stream: Remove `s3` as an option for `destination`. Use `extended_s3` instead ([#31138](https://github.com/hashicorp/terraform-provider-aws/issues/31138))
* resource/aws_kinesis_firehose_delivery_stream: Rename `extended_s3_configuration.0.s3_backup_configuration.0.buffer_size` and `extended_s3_configuration.0.s3_backup_configuration.0.buffer_interval` to `extended_s3_configuration.0.s3_backup_configuration.0.buffering_size` and `extended_s3_configuration.0.s3_backup_configuration.0.buffering_interval`, respectively ([#31141](https://github.com/hashicorp/terraform-provider-aws/issues/31141))
* resource/aws_kinesis_firehose_delivery_stream: Rename `redshift_configuration.0.s3_backup_configuration.0.buffer_size` and `redshift_configuration.0.s3_backup_configuration.0.buffer_interval` to `redshift_configuration.0.s3_backup_configuration.0.buffering_size` and `redshift_configuration.0.s3_backup_configuration.0.buffering_interval`, respectively ([#31141](https://github.com/hashicorp/terraform-provider-aws/issues/31141))
* resource/aws_kinesis_firehose_delivery_stream: Rename `s3_configuration.0.buffer_size` and `s3_configuration.0.buffer_internval` to `s3_configuration.0.buffering_size` and `s3_configuration.0.buffering_internval`, respectively ([#31141](https://github.com/hashicorp/terraform-provider-aws/issues/31141))
* resource/aws_launch_configuration: With the retirement of EC2-Classic the `vpc_classic_link_id` and `vpc_classic_link_security_groups` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_lightsail_instance: The `ipv6_address` attribute has been removed ([#31489](https://github.com/hashicorp/terraform-provider-aws/issues/31489))
* resource/aws_medialive_multiplex_program: The `statemux_settings` attribute has been removed. Use `statmux_settings` argument instead ([#31034](https://github.com/hashicorp/terraform-provider-aws/issues/31034))
* resource/aws_msk_cluster: The `broker_node_group_info.ebs_volume_size` attribute has been removed ([#31324](https://github.com/hashicorp/terraform-provider-aws/issues/31324))
* resource/aws_neptune_cluster: `snapshot_identifier` change now properly forces replacement ([#29409](https://github.com/hashicorp/terraform-provider-aws/issues/29409))
* resource/aws_networkmanager_core_network: Removed `policy_document` argument -- use `aws_networkmanager_core_network_policy_attachment` resource instead ([#30875](https://github.com/hashicorp/terraform-provider-aws/issues/30875))
* resource/aws_rds_cluster: The `engine` argument is now required and has no default ([#31112](https://github.com/hashicorp/terraform-provider-aws/issues/31112))
* resource/aws_rds_cluster: `snapshot_identifier` change now properly forces replacement ([#29409](https://github.com/hashicorp/terraform-provider-aws/issues/29409))
* resource/aws_rds_cluster_instance: The `engine` argument is now required and has no default ([#31112](https://github.com/hashicorp/terraform-provider-aws/issues/31112))
* resource/aws_redshift_cluster: With the retirement of EC2-Classic the `cluster_security_groups` attribute has been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_route: `instance_id` can no longer be set in configurations. Use `network_interface_id` instead, for example, setting `network_interface_id` to `aws_instance.test.primary_network_interface_id`. ([#30804](https://github.com/hashicorp/terraform-provider-aws/issues/30804))
* resource/aws_route_table: `route.*.instance_id` can no longer be set in configurations. Use `route.*.network_interface_id` instead, for example, setting `network_interface_id` to `aws_instance.test.primary_network_interface_id`. ([#30804](https://github.com/hashicorp/terraform-provider-aws/issues/30804))
* resource/aws_secretsmanager_secret: The `rotation_enabled`, `rotation_lambda_arn` and `rotation_rules` attributes have been removed ([#31487](https://github.com/hashicorp/terraform-provider-aws/issues/31487))
* resource/aws_security_group: With the retirement of EC2-Classic non-VPC security groups are no longer supported ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_security_group_rule: With the retirement of EC2-Classic non-VPC security groups are no longer supported ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_servicecatalog_product: Changes to any `provisioning_artifact_parameters` arguments now properly trigger a replacement. This fixes incorrect behavior, but may technically be breaking for configurations expecting non-functional in-place updates. ([#31061](https://github.com/hashicorp/terraform-provider-aws/issues/31061))
* resource/aws_vpc: With the retirement of EC2-Classic the `enable_classiclink` and `enable_classiclink_dns_support` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_vpc_peering_connection: With the retirement of EC2-Classic the `allow_classic_link_to_remote_vpc` and `allow_vpc_to_remote_classic_link` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_vpc_peering_connection_accepter: With the retirement of EC2-Classic the `allow_classic_link_to_remote_vpc` and `allow_vpc_to_remote_classic_link` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_vpc_peering_connection_options: With the retirement of EC2-Classic the `allow_classic_link_to_remote_vpc` and `allow_vpc_to_remote_classic_link` attributes have been removed ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))
* resource/aws_wafv2_web_acl: The `statement.managed_rule_group_statement.excluded_rule` and `statement.rule_group_reference_statement.excluded_rule` attributes have been removed ([#31374](https://github.com/hashicorp/terraform-provider-aws/issues/31374))
* resource/aws_wafv2_web_acl_logging_configuration: The `redacted_fields.all_query_arguments`, `redacted_fields.body` and `redacted_fields.single_query_argument` attributes have been removed ([#31486](https://github.com/hashicorp/terraform-provider-aws/issues/31486))

NOTES:

* data-source/aws_elasticache_replication_group: Update configurations to use `description` instead of the `replication_group_description` argument ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* data-source/aws_elasticache_replication_group: Update configurations to use `num_cache_clusters` instead of the `number_cache_clusters` argument ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* data-source/aws_opensearch_domain: The `kibana_endpoint` attribute has been deprecated. All configurations using `kibana_endpoint` should be updated to use the `dashboard_endpoint` attribute instead ([#31490](https://github.com/hashicorp/terraform-provider-aws/issues/31490))
* data-source/aws_quicksight_data_set: The `tags_all` attribute has been deprecated and will be removed in a future version ([#31162](https://github.com/hashicorp/terraform-provider-aws/issues/31162))
* data-source/aws_redshift_service_account: The `aws_redshift_service_account` data source has been deprecated and will be removed in a future version. AWS documentation [states that](https://docs.aws.amazon.com/redshift/latest/mgmt/db-auditing.html#db-auditing-bucket-permissions) a [service principal name](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_principal.html#principal-services) should be used instead of an AWS account ID in any relevant IAM policy ([#31006](https://github.com/hashicorp/terraform-provider-aws/issues/31006))
* data-source/aws_service_discovery_service: The `tags_all` attribute has been deprecated and will be removed in a future version ([#31162](https://github.com/hashicorp/terraform-provider-aws/issues/31162))
* resource/aws_api_gateway_rest_api: Update configurations with `minimum_compression_size` set to pass the value as a string. Valid values remain the same. ([#30969](https://github.com/hashicorp/terraform-provider-aws/issues/30969))
* resource/aws_autoscaling_attachment: Update configurations to use `lb_target_group_arn` instead of `alb_target_group_arn` which has been removed ([#30828](https://github.com/hashicorp/terraform-provider-aws/issues/30828))
* resource/aws_db_event_subscription: Configurations that define `source_ids` using the `id` attribute of `aws_db_instance` must be updated to use `identifier` instead - for example, `source_ids = [aws_db_instance.example.id]` must be updated to `source_ids = [aws_db_instance.example.identifier]` ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_db_instance: Configurations that define `replicate_source_db` using the `id` attribute of `aws_db_instance` must be updated to use `identifier` instead - for example, `replicate_source_db = aws_db_instance.example.id` must be updated to `replicate_source_db = aws_db_instance.example.identifier` ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_db_instance: The change of what `id` is, namely, a DBI Resource ID now versus DB Identifier previously, has far-reaching consequences. Configurations that refer to, for example, `aws_db_instance.example.id` will now have errors and must be changed to use `identifier` instead, for example, `aws_db_instance.example.identifier` ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_db_instance_role_association: Configurations that define `db_instance_identifier` using the `id` attribute of `aws_db_instance` must be updated to use `identifier` instead - for example, `db_instance_identifier = aws_db_instance.example.id` must be updated to `db_instance_identifier = aws_db_instance.example.identifier` ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_db_proxy_target: Configurations that define `db_instance_identifier` using the `id` attribute of `aws_db_instance` must be updated to use `identifier` instead - for example, `db_instance_identifier = aws_db_instance.example.id` must be updated to `db_instance_identifier = aws_db_instance.example.identifier` ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_db_snapshot: Configurations that define `db_instance_identifier` using the `id` attribute of `aws_db_instance` must be updated to use `identifier` instead - for example, `db_instance_identifier = aws_db_instance.example.id` must be updated to `db_instance_identifier = aws_db_instance.example.identifier` ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_docdb_cluster: Changes to the `snapshot_identifier` attribute will now trigger a replacement, rather than an in-place update. This corrects the previous behavior which resulted in a successful apply, but did not actually restore the cluster from the designated snapshot. ([#29409](https://github.com/hashicorp/terraform-provider-aws/issues/29409))
* resource/aws_dx_gateway_association: The `vpn_gateway_id` attribute has been deprecated. All configurations using `vpn_gateway_id` should be updated to use the `associated_gateway_id` attribute instead ([#31384](https://github.com/hashicorp/terraform-provider-aws/issues/31384))
* resource/aws_elasticache_replication_group: Update configurations to use `description` instead of the `replication_group_description` argument ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* resource/aws_elasticache_replication_group: Update configurations to use `num_cache_clusters` instead of the `number_cache_clusters` argument ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* resource/aws_elasticache_replication_group: Update configurations to use `preferred_cache_cluster_azs` instead of the `availability_zones` argument ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* resource/aws_elasticache_replication_group: Update configurations to use top-level `num_node_groups` and `replicas_per_node_group` instead of `cluster_mode.0.num_node_groups` and `cluster_mode.0.replicas_per_node_group`, respectively ([#31008](https://github.com/hashicorp/terraform-provider-aws/issues/31008))
* resource/aws_flow_log: The `log_group_name` attribute has been deprecated. All configurations using `log_group_name` should be updated to use the `log_destination` attribute instead ([#31382](https://github.com/hashicorp/terraform-provider-aws/issues/31382))
* resource/aws_guardduty_organization_configuration: The `auto_enable` argument has been deprecated. Use the `auto_enable_organization_members` argument instead. ([#30736](https://github.com/hashicorp/terraform-provider-aws/issues/30736))
* resource/aws_neptune_cluster: Changes to the `snapshot_identifier` attribute will now trigger a replacement, rather than an in-place update. This corrects the previous behavior which resulted in a successful apply, but did not actually restore the cluster from the designated snapshot. ([#29409](https://github.com/hashicorp/terraform-provider-aws/issues/29409))
* resource/aws_networkmanager_core_network: Update configurations to use the `aws_networkmanager_core_network_policy_attachment` resource instead of the `policy_document` argument ([#30875](https://github.com/hashicorp/terraform-provider-aws/issues/30875))
* resource/aws_opensearch_domain: The `engine_version` attribute no longer has a default value. When omitted, the underlying AWS API will use the latest OpenSearch engine version. ([#31568](https://github.com/hashicorp/terraform-provider-aws/issues/31568))
* resource/aws_opensearch_domain: The `kibana_endpoint` attribute has been deprecated. All configurations using `kibana_endpoint` should be updated to use the `dashboard_endpoint` attribute instead ([#31490](https://github.com/hashicorp/terraform-provider-aws/issues/31490))
* resource/aws_rds_cluster: Changes to the `snapshot_identifier` attribute will now trigger a replacement, rather than an in-place update. This corrects the previous behavior which resulted in a successful apply, but did not actually restore the cluster from the designated snapshot. ([#29409](https://github.com/hashicorp/terraform-provider-aws/issues/29409))
* resource/aws_rds_cluster: Configurations not including the `engine` argument must be updated to include `engine` as it is now required. Previously, not including `engine` was equivalent to `engine = "aurora"` and created a MySQL-5.6-compatible cluster ([#31112](https://github.com/hashicorp/terraform-provider-aws/issues/31112))
* resource/aws_rds_cluster_instance: Configurations not including the `engine` argument must be updated to include `engine` as it is now required. Previously, not including `engine` was equivalent to `engine = "aurora"` and created a MySQL-5.6-compatible cluster instance ([#31112](https://github.com/hashicorp/terraform-provider-aws/issues/31112))
* resource/aws_route: Since `instance_id` can no longer be set in configurations, use `network_interface_id` instead. For example, set `network_interface_id` to `aws_instance.test.primary_network_interface_id`. ([#30804](https://github.com/hashicorp/terraform-provider-aws/issues/30804))
* resource/aws_route_table: Since `route.*.instance_id` can no longer be set in configurations, use `route.*.network_interface_id` instead. For example, set `network_interface_id` to `aws_instance.test.primary_network_interface_id`. ([#30804](https://github.com/hashicorp/terraform-provider-aws/issues/30804))
* resource/aws_ssm_association: The `instance_id` attribute has been deprecated. All configurations using `instance_id` should be updated to use the `targets` attribute instead ([#31380](https://github.com/hashicorp/terraform-provider-aws/issues/31380))

ENHANCEMENTS:

* provider: Allow `computed` `tags` on resources ([#30793](https://github.com/hashicorp/terraform-provider-aws/issues/30793))
* provider: Allow `default_tags` and resource `tags` to include zero values `""` ([#30793](https://github.com/hashicorp/terraform-provider-aws/issues/30793))
* provider: Duplicate `default_tags` can now be included and will be overwritten by resource `tags` ([#30793](https://github.com/hashicorp/terraform-provider-aws/issues/30793))
* resource/aws_db_instance: Updates to `identifier` and `identifier_prefix` will no longer cause the database instance to be destroyed and recreated ([#31232](https://github.com/hashicorp/terraform-provider-aws/issues/31232))
* resource/aws_eip: Deprecate `vpc` attribute. Use `domain` instead ([#31567](https://github.com/hashicorp/terraform-provider-aws/issues/31567))
* resource/aws_guardduty_organization_configuration: Add `auto_enable_organization_members` attribute ([#30736](https://github.com/hashicorp/terraform-provider-aws/issues/30736))
* resource/aws_kinesis_firehose_delivery_stream: Add `s3_configuration` to `elasticsearch_configuration`, `opensearch_configuration`, `redshift_configuration`, `splunk_configuration`, and `http_endpoint_configuration` ([#31138](https://github.com/hashicorp/terraform-provider-aws/issues/31138))
* resource/aws_opensearch_domain: Removed `engine_version` default value ([#31568](https://github.com/hashicorp/terraform-provider-aws/issues/31568))
* resource/aws_wafv2_web_acl: Support `rule_action_override` on `rule_group_reference_statement` ([#31374](https://github.com/hashicorp/terraform-provider-aws/issues/31374))

BUG FIXES:

* resource/aws_ecs_capacity_provider: Allow an `instance_warmup_period` of `0` in the `auto_scaling_group_provider.managed_scaling` configuration block ([#24005](https://github.com/hashicorp/terraform-provider-aws/issues/24005))
* resource/aws_launch_template: Remove default values in `metadata_options` to allow default condition ([#30545](https://github.com/hashicorp/terraform-provider-aws/issues/30545))
* resource/aws_s3_bucket: Fix bucket_regional_domain_name not including region for buckets in us-east-1 ([#25724](https://github.com/hashicorp/terraform-provider-aws/issues/25724))
* resource/aws_s3_object: Remove `acl` default in order to work with S3 buckets that have ACL disabled ([#27197](https://github.com/hashicorp/terraform-provider-aws/issues/27197))
* resource/aws_s3_object_copy: Remove `acl` default in order to work with S3 buckets that have ACL disabled ([#27197](https://github.com/hashicorp/terraform-provider-aws/issues/27197))
* resource/aws_servicecatalog_product: Changes to `provisioning_artifact_parameters` arguments now properly trigger a replacement ([#31061](https://github.com/hashicorp/terraform-provider-aws/issues/31061))
* resource/aws_vpc_peering_connection: Fix crash in `vpcPeeringConnectionOptionsEqual` ([#30966](https://github.com/hashicorp/terraform-provider-aws/issues/30966))

## Previous Releases

For information on prior major releases, see their changelogs:

* [4.x](https://github.com/hashicorp/terraform-provider-aws/blob/release/4.x/CHANGELOG.md)
* [3.x](https://github.com/hashicorp/terraform-provider-aws/blob/release/3.x/CHANGELOG.md)
* [2.x and earlier](https://github.com/hashicorp/terraform-provider-aws/blob/release/2.x/CHANGELOG.md)
