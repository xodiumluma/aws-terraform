---
subcategory: "Chime SDK Voice"
layout: "aws"
page_title: "AWS: aws_chimesdkvoice_global_settings"
description: |-
  Terraform resource for managing Amazon Chime SDK Voice Global Settings.
---

# Resource: aws_chimesdkvoice_global_settings

Terraform resource for managing Amazon Chime SDK Voice Global Settings.

## Example Usage

### Basic Usage

```terraform
resource "aws_chimesdkvoice_global_settings" "example" {
  voice_connector {
    cdr_bucket = "example-bucket-name"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `voice_connector` - (Required) The Voice Connector settings. See [voice_connector](#voice_connector).

### `voice_connector`

The Amazon Chime SDK Voice Connector settings. Includes any Amazon S3 buckets designated for storing call detail records.

* `cdr_bucket` - (Optional) The S3 bucket that stores the Voice Connector's call detail records.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - AWS account ID for which the settings are applied.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import AWS Chime SDK Voice Global Settings using the `id` (AWS account ID). For example:

```terraform
import {
  to = aws_chimesdkvoice_global_settings.example
  id = "123456789012"
}
```

Using `terraform import`, import AWS Chime SDK Voice Global Settings using the `id` (AWS account ID). For example:

```console
% terraform import aws_chimesdkvoice_global_settings.example 123456789012
```
