---
subcategory: "Lightsail"
layout: "aws"
page_title: "AWS: aws_lightsail_lb_certificate"
description: |-
  Provides a Lightsail Load Balancer
---

# Resource: aws_lightsail_lb_certificate

Creates a Lightsail load balancer Certificate resource.

## Example Usage

```terraform
resource "aws_lightsail_lb" "test" {
  name              = "test-load-balancer"
  health_check_path = "/"
  instance_port     = "80"
  tags = {
    foo = "bar"
  }
}

resource "aws_lightsail_lb_certificate" "test" {
  name        = "test-load-balancer-certificate"
  lb_name     = aws_lightsail_lb.test.id
  domain_name = "test.com"
}
```

## Argument Reference

This resource supports the following arguments:

* `domain_name` - (Required) The domain name (e.g., example.com) for your SSL/TLS certificate.
* `lb_name` - (Required) The load balancer name where you want to create the SSL/TLS certificate.
* `name` - (Required) The SSL/TLS certificate name.
* `name` - (Required) The SSL/TLS certificate name.
* `subject_alternative_names` - (Optional) Set of domains that should be SANs in the issued certificate. `domain_name` attribute is automatically added as a Subject Alternative Name.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `id` - A combination of attributes to create a unique id: `lb_name`,`name`
* `arn` - The ARN of the lightsail certificate.
* `created_at` - The timestamp when the instance was created.
* `domain_validation_options` - Set of domain validation objects which can be used to complete certificate validation. Can have more than one element, e.g., if SANs are defined.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import `aws_lightsail_lb_certificate` using the id attribute. For example:

```terraform
import {
  to = aws_lightsail_lb_certificate.test
  id = "example-load-balancer,example-load-balancer-certificate"
}
```

Using `terraform import`, import `aws_lightsail_lb_certificate` using the id attribute. For example:

```console
% terraform import aws_lightsail_lb_certificate.test example-load-balancer,example-load-balancer-certificate
```
