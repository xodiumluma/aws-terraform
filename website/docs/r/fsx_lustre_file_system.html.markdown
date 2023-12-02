---
subcategory: "FSx"
layout: "aws"
page_title: "AWS: aws_fsx_lustre_file_system"
description: |-
  Manages a FSx Lustre File System.
---

# Resource: aws_fsx_lustre_file_system

Manages a FSx Lustre File System. See the [FSx Lustre Guide](https://docs.aws.amazon.com/fsx/latest/LustreGuide/what-is.html) for more information.

~> **NOTE:** `auto_import_policy`, `export_path`, `import_path` and `imported_file_chunk_size` are not supported with the `PERSISTENT_2` deployment type. Use `aws_fsx_data_repository_association` instead.

## Example Usage

```terraform
resource "aws_fsx_lustre_file_system" "example" {
  import_path      = "s3://${aws_s3_bucket.example.bucket}"
  storage_capacity = 1200
  subnet_ids       = [aws_subnet.example.id]
}
```

## Argument Reference

This resource supports the following arguments:

* `storage_capacity` - (Optional) The storage capacity (GiB) of the file system. Minimum of `1200`. See more details at [Allowed values for Fsx storage capacity](https://docs.aws.amazon.com/fsx/latest/APIReference/API_CreateFileSystem.html#FSx-CreateFileSystem-request-StorageCapacity). Update is allowed only for `SCRATCH_2`, `PERSISTENT_1` and `PERSISTENT_2` deployment types, See more details at [Fsx Storage Capacity Update](https://docs.aws.amazon.com/fsx/latest/APIReference/API_UpdateFileSystem.html#FSx-UpdateFileSystem-request-StorageCapacity). Required when not creating filesystem for a backup.
* `subnet_ids` - (Required) A list of IDs for the subnets that the file system will be accessible from. File systems currently support only one subnet. The file server is also launched in that subnet's Availability Zone.
* `backup_id` - (Optional) The ID of the source backup to create the filesystem from.
* `export_path` - (Optional) S3 URI (with optional prefix) where the root of your Amazon FSx file system is exported. Can only be specified with `import_path` argument and the path must use the same Amazon S3 bucket as specified in `import_path`. Set equal to `import_path` to overwrite files on export. Defaults to `s3://{IMPORT BUCKET}/FSxLustre{CREATION TIMESTAMP}`. Only supported on `PERSISTENT_1` deployment types.
* `import_path` - (Optional) S3 URI (with optional prefix) that you're using as the data repository for your FSx for Lustre file system. For example, `s3://example-bucket/optional-prefix/`. Only supported on `PERSISTENT_1` deployment types.
* `imported_file_chunk_size` - (Optional) For files imported from a data repository, this value determines the stripe count and maximum amount of data per file (in MiB) stored on a single physical disk. Can only be specified with `import_path` argument. Defaults to `1024`. Minimum of `1` and maximum of `512000`. Only supported on `PERSISTENT_1` deployment types.
* `security_group_ids` - (Optional) A list of IDs for the security groups that apply to the specified network interfaces created for file system access. These security groups will apply to all network interfaces.
* `tags` - (Optional) A map of tags to assign to the file system. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `weekly_maintenance_start_time` - (Optional) The preferred start time (in `d:HH:MM` format) to perform weekly maintenance, in the UTC time zone.
* `deployment_type` - (Optional) - The filesystem deployment type. One of: `SCRATCH_1`, `SCRATCH_2`, `PERSISTENT_1`, `PERSISTENT_2`.
* `kms_key_id` - (Optional) ARN for the KMS Key to encrypt the file system at rest, applicable for `PERSISTENT_1` and `PERSISTENT_2` deployment_type. Defaults to an AWS managed KMS Key.
* `per_unit_storage_throughput` - (Optional) - Describes the amount of read and write throughput for each 1 tebibyte of storage, in MB/s/TiB, required for the `PERSISTENT_1` and `PERSISTENT_2` deployment_type. Valid values for `PERSISTENT_1` deployment_type and `SSD` storage_type are 50, 100, 200. Valid values for `PERSISTENT_1` deployment_type and `HDD` storage_type are 12, 40. Valid values for `PERSISTENT_2` deployment_type and ` SSD` storage_type are 125, 250, 500, 1000.
* `automatic_backup_retention_days` - (Optional) The number of days to retain automatic backups. Setting this to 0 disables automatic backups. You can retain automatic backups for a maximum of 90 days. only valid for `PERSISTENT_1` and `PERSISTENT_2` deployment_type.
* `storage_type` - (Optional) - The filesystem storage type. Either `SSD` or `HDD`, defaults to `SSD`. `HDD` is only supported on `PERSISTENT_1` deployment types.
* `drive_cache_type` - (Optional) - The type of drive cache used by `PERSISTENT_1` filesystems that are provisioned with `HDD` storage_type. Required for `HDD` storage_type, set to either `READ` or `NONE`.
* `daily_automatic_backup_start_time` - (Optional) A recurring daily time, in the format HH:MM. HH is the zero-padded hour of the day (0-23), and MM is the zero-padded minute of the hour. For example, 05:00 specifies 5 AM daily. only valid for `PERSISTENT_1` and `PERSISTENT_2` deployment_type. Requires `automatic_backup_retention_days` to be set.
* `auto_import_policy` - (Optional) How Amazon FSx keeps your file and directory listings up to date as you add or modify objects in your linked S3 bucket. see [Auto Import Data Repo](https://docs.aws.amazon.com/fsx/latest/LustreGuide/autoimport-data-repo.html) for more details. Only supported on `PERSISTENT_1` deployment types.
* `copy_tags_to_backups` - (Optional) A boolean flag indicating whether tags for the file system should be copied to backups. Applicable for `PERSISTENT_1` and `PERSISTENT_2` deployment_type. The default value is false.
* `data_compression_type` - (Optional) Sets the data compression configuration for the file system. Valid values are `LZ4` and `NONE`. Default value is `NONE`. Unsetting this value reverts the compression type back to `NONE`.
* `file_system_type_version` - (Optional) Sets the Lustre version for the file system that you're creating. Valid values are 2.10 for `SCRATCH_1`, `SCRATCH_2` and `PERSISTENT_1` deployment types. Valid values for 2.12 include all deployment types.
* `log_configuration` - (Optional) The Lustre logging configuration used when creating an Amazon FSx for Lustre file system. When logging is enabled, Lustre logs error and warning events for data repositories associated with your file system to Amazon CloudWatch Logs.
* `root_squash_configuration` - (Optional) The Lustre root squash configuration used when creating an Amazon FSx for Lustre file system. When enabled, root squash restricts root-level access from clients that try to access your file system as a root user.

### log_configuration

* `destination` - (Optional) The Amazon Resource Name (ARN) that specifies the destination of the logs. The name of the Amazon CloudWatch Logs log group must begin with the `/aws/fsx` prefix. If you do not provide a destination, Amazon FSx will create and use a log stream in the CloudWatch Logs `/aws/fsx/lustre` log group.
* `level` - (Optional) Sets which data repository events are logged by Amazon FSx. Valid values are `WARN_ONLY`, `FAILURE_ONLY`, `ERROR_ONLY`, `WARN_ERROR` and `DISABLED`. Default value is `DISABLED`.

### root_squash_configuration

* `no_squash_nids` - (Optional) When root squash is enabled, you can optionally specify an array of NIDs of clients for which root squash does not apply. A client NID is a Lustre Network Identifier used to uniquely identify a client. You can specify the NID as either a single address or a range of addresses: 1. A single address is described in standard Lustre NID format by specifying the client’s IP address followed by the Lustre network ID (for example, 10.0.1.6@tcp). 2. An address range is described using a dash to separate the range (for example, 10.0.[2-10].[1-255]@tcp).
* `root_squash` - (Optional) You enable root squash by setting a user ID (UID) and group ID (GID) for the file system in the format UID:GID (for example, 365534:65534). The UID and GID values can range from 0 to 4294967294.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - Amazon Resource Name of the file system.
* `dns_name` - DNS name for the file system, e.g., `fs-12345678.fsx.us-west-2.amazonaws.com`
* `id` - Identifier of the file system, e.g., `fs-12345678`
* `network_interface_ids` - Set of Elastic Network Interface identifiers from which the file system is accessible. As explained in the [documentation](https://docs.aws.amazon.com/fsx/latest/LustreGuide/mounting-on-premises.html), the first network interface returned is the primary network interface.
* `mount_name` - The value to be used when mounting the filesystem.
* `owner_id` - AWS account identifier that created the file system.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).
* `vpc_id` - Identifier of the Virtual Private Cloud for the file system.

## Timeouts

[Configuration options](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts):

* `create` - (Default `30m`)
* `update` - (Default `30m`)
* `delete` - (Default `30m`)

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import FSx File Systems using the `id`. For example:

```terraform
import {
  to = aws_fsx_lustre_file_system.example
  id = "fs-543ab12b1ca672f33"
}
```

Using `terraform import`, import FSx File Systems using the `id`. For example:

```console
% terraform import aws_fsx_lustre_file_system.example fs-543ab12b1ca672f33
```

Certain resource arguments, like `security_group_ids`, do not have a FSx API method for reading the information after creation. If the argument is set in the Terraform configuration on an imported resource, Terraform will always show a difference. To workaround this behavior, either omit the argument from the Terraform configuration or use [`ignore_changes`](https://www.terraform.io/docs/configuration/meta-arguments/lifecycle.html#ignore_changes) to hide the difference. For example:

```terraform
resource "aws_fsx_lustre_file_system" "example" {
  # ... other configuration ...
  security_group_ids = [aws_security_group.example.id]

  # There is no FSx API for reading security_group_ids
  lifecycle {
    ignore_changes = [security_group_ids]
  }
}
```
