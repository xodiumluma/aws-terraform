---
subcategory: "API Gateway"
layout: "aws"
page_title: "AWS: aws_api_gateway_account"
description: |-
  Provides a settings of an API Gateway Account.
---

# Resource: aws_api_gateway_account

Provides a settings of an API Gateway Account. Settings is applied region-wide per `provider` block.

-> **Note:** As there is no API method for deleting account settings or resetting it to defaults, destroying this resource will keep your account settings intact

## Example Usage

```terraform
resource "aws_api_gateway_account" "demo" {
  cloudwatch_role_arn = aws_iam_role.cloudwatch.arn
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["apigateway.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "cloudwatch" {
  name               = "api_gateway_cloudwatch_global"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "cloudwatch" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:DescribeLogGroups",
      "logs:DescribeLogStreams",
      "logs:PutLogEvents",
      "logs:GetLogEvents",
      "logs:FilterLogEvents",
    ]

    resources = ["*"]
  }
}
resource "aws_iam_role_policy" "cloudwatch" {
  name   = "default"
  role   = aws_iam_role.cloudwatch.id
  policy = data.aws_iam_policy_document.cloudwatch.json
}
```

## Argument Reference

The following argument is supported:

* `cloudwatch_role_arn` - (Optional) ARN of an IAM role for CloudWatch (to allow logging & monitoring). See more [in AWS Docs](https://docs.aws.amazon.com/apigateway/latest/developerguide/how-to-stage-settings.html#how-to-stage-settings-console). Logging & monitoring can be enabled/disabled and otherwise tuned on the API Gateway Stage level.

## Attributes Reference

The following attribute is exported:

* `throttle_settings` - Account-Level throttle settings. See exported fields below.

`throttle_settings` block exports the following:

* `burst_limit` - Absolute maximum number of times API Gateway allows the API to be called per second (RPS).
* `rate_limit` - Number of times API Gateway allows the API to be called per second on average (RPS).

## Import

API Gateway Accounts can be imported using the word `api-gateway-account`, e.g.,

```
$ terraform import aws_api_gateway_account.demo api-gateway-account
```
