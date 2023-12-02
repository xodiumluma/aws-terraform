---
subcategory: "Kinesis"
layout: "aws"
page_title: "AWS: aws_kinesis_stream"
description: |-
  Provides a Kinesis Stream data source.
---

# Data Source: aws_kinesis_stream

Use this data source to get information about a Kinesis Stream for use in other
resources.

For more details, see the [Amazon Kinesis Documentation][1].

## Example Usage

```terraform
data "aws_kinesis_stream" "stream" {
  name = "stream-name"
}
```

## Argument Reference

* `name` - (Required) Name of the Kinesis Stream.

## Attribute Reference

`id` is set to the ARN of the Kinesis Stream. In addition, the following attributes
are exported:

* `arn` - ARN of the Kinesis Stream (same as id).
* `name` - Name of the Kinesis Stream.
* `creation_timestamp` - Approximate UNIX timestamp that the stream was created.
* `status` - Current status of the stream. The stream status is one of CREATING, DELETING, ACTIVE, or UPDATING.
* `retention_period` - Length of time (in hours) data records are accessible after they are added to the stream.
* `open_shards` - List of shard ids in the OPEN state. See [Shard State][2] for more.
* `closed_shards` - List of shard ids in the CLOSED state. See [Shard State][2] for more.
* `shard_level_metrics` - List of shard-level CloudWatch metrics which are enabled for the stream. See [Monitoring with CloudWatch][3] for more.
* `stream_mode_details` - [Capacity mode][4] of the data stream. Detailed below.
* `tags` - Map of tags to assigned to the stream.

### stream_mode_details Configuration Block

* `stream_mode` - Capacity mode of the stream. Either `ON_DEMAND` or `PROVISIONED`.

[1]: https://aws.amazon.com/documentation/kinesis/
[2]: https://docs.aws.amazon.com/streams/latest/dev/kinesis-using-sdk-java-after-resharding.html#kinesis-using-sdk-java-resharding-data-routing
[3]: https://docs.aws.amazon.com/streams/latest/dev/monitoring-with-cloudwatch.html
[4]: https://docs.aws.amazon.com/streams/latest/dev/how-do-i-size-a-stream.html
