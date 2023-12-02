---
subcategory: "ECS (Elastic Container)"
layout: "aws"
page_title: "AWS: aws_ecs_cluster"
description: |-
  Provides an ECS cluster.
---

# Resource: aws_ecs_cluster

Provides an ECS cluster.

## Example Usage

```terraform
resource "aws_ecs_cluster" "foo" {
  name = "white-hart"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }
}
```

### Example with Log Configuration

```terraform
resource "aws_kms_key" "example" {
  description             = "example"
  deletion_window_in_days = 7
}

resource "aws_cloudwatch_log_group" "example" {
  name = "example"
}

resource "aws_ecs_cluster" "test" {
  name = "example"

  configuration {
    execute_command_configuration {
      kms_key_id = aws_kms_key.example.arn
      logging    = "OVERRIDE"

      log_configuration {
        cloud_watch_encryption_enabled = true
        cloud_watch_log_group_name     = aws_cloudwatch_log_group.example.name
      }
    }
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `configuration` - (Optional) The execute command configuration for the cluster. Detailed below.
* `name` - (Required) Name of the cluster (up to 255 letters, numbers, hyphens, and underscores)
* `service_connect_defaults` - (Optional) Configures a default Service Connect namespace. Detailed below.
* `setting` - (Optional) Configuration block(s) with cluster settings. For example, this can be used to enable CloudWatch Container Insights for a cluster. Detailed below.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

### `configuration`

* `execute_command_configuration` - (Optional) The details of the execute command configuration. Detailed below.

#### `execute_command_configuration`

* `kms_key_id` - (Optional) The AWS Key Management Service key ID to encrypt the data between the local client and the container.
* `log_configuration` - (Optional) The log configuration for the results of the execute command actions Required when `logging` is `OVERRIDE`. Detailed below.
* `logging` - (Optional) The log setting to use for redirecting logs for your execute command results. Valid values are `NONE`, `DEFAULT`, and `OVERRIDE`.

##### `log_configuration`

* `cloud_watch_encryption_enabled` - (Optional) Whether or not to enable encryption on the CloudWatch logs. If not specified, encryption will be disabled.
* `cloud_watch_log_group_name` - (Optional) The name of the CloudWatch log group to send logs to.
* `s3_bucket_name` - (Optional) The name of the S3 bucket to send logs to.
* `s3_bucket_encryption_enabled` - (Optional) Whether or not to enable encryption on the logs sent to S3. If not specified, encryption will be disabled.
* `s3_key_prefix` - (Optional) An optional folder in the S3 bucket to place logs in.

### `setting`

* `name` - (Required) Name of the setting to manage. Valid values: `containerInsights`.
* `value` -  (Required) The value to assign to the setting. Valid values are `enabled` and `disabled`.

### `service_connect_defaults`

* `namespace` - (Required) The ARN of the [`aws_service_discovery_http_namespace`](/docs/providers/aws/r/service_discovery_http_namespace.html) that's used when you create a service and don't specify a Service Connect configuration.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - ARN that identifies the cluster.
* `id` - ARN that identifies the cluster.
* `tags_all` - Map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import ECS clusters using the `name`. For example:

```terraform
import {
  to = aws_ecs_cluster.stateless
  id = "stateless-app"
}
```

Using `terraform import`, import ECS clusters using the `name`. For example:

```console
% terraform import aws_ecs_cluster.stateless stateless-app
```
