---
subcategory: "SSM (Systems Manager)"
layout: "aws"
page_title: "AWS: aws_ssm_patch_baseline"
description: |-
  Provides an SSM Patch Baseline data source
---

# Data Source: aws_ssm_patch_baseline

Provides an SSM Patch Baseline data source. Useful if you wish to reuse the default baselines provided.

## Example Usage

To retrieve a baseline provided by AWS:

```terraform
data "aws_ssm_patch_baseline" "centos" {
  owner            = "AWS"
  name_prefix      = "AWS-"
  operating_system = "CENTOS"
}
```

To retrieve a baseline on your account:

```terraform
data "aws_ssm_patch_baseline" "default_custom" {
  owner            = "Self"
  name_prefix      = "MyCustomBaseline"
  default_baseline = true
  operating_system = "WINDOWS"
}
```

## Argument Reference

This data source supports the following arguments:

* `owner` - (Required) Owner of the baseline. Valid values: `All`, `AWS`, `Self` (the current account).
* `name_prefix` - (Optional) Filter results by the baseline name prefix.
* `default_baseline` - (Optional) Filters the results against the baselines default_baseline field.
* `operating_system` - (Optional) Specified OS for the baseline. Valid values: `AMAZON_LINUX`, `AMAZON_LINUX_2`, `UBUNTU`, `REDHAT_ENTERPRISE_LINUX`, `SUSE`, `CENTOS`, `ORACLE_LINUX`, `DEBIAN`, `MACOS`, `RASPBIAN` and `ROCKY_LINUX`.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `approved_patches` - List of explicitly approved patches for the baseline.
* `approved_patches_compliance_level` - The compliance level for approved patches.
* `approved_patches_enable_non_security` - Indicates whether the list of approved patches includes non-security updates that should be applied to the instances.
* `approval_rule` - List of rules used to include patches in the baseline.
    * `approve_after_days` - The number of days after the release date of each patch matched by the rule the patch is marked as approved in the patch baseline.
    * `approve_until_date` - The cutoff date for auto approval of released patches. Any patches released on or before this date are installed automatically. Date is formatted as `YYYY-MM-DD`. Conflicts with `approve_after_days`
    * `compliance_level` - The compliance level for patches approved by this rule.
    * `enable_non_security` - Boolean enabling the application of non-security updates.
    * `patch_filter` - The patch filter group that defines the criteria for the rule.
        * `key` - The key for the filter.
        * `values` - The value for the filter.
* `global_filter` - Set of global filters used to exclude patches from the baseline.
    * `key` - The key for the filter.
    * `values` - The value for the filter.
* `id` - ID of the baseline.
* `name` - Name of the baseline.
* `description` - Description of the baseline.
* `rejected_patches` - List of rejected patches.
* `rejected_patches_action` - The action specified to take on patches included in the `rejected_patches` list.
* `source` - Information about the patches to use to update the managed nodes, including target operating systems and source repositories.
    * `configuration` - The value of the yum repo configuration.
    * `name` - The name specified to identify the patch source.
    * `products` - The specific operating system versions a patch repository applies to.
