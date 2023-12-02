// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/hashicorp/aws-sdk-go-base/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKDataSource("aws_sqs_queue")
func dataSourceQueue() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceQueueRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceQueueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SQSClient(ctx)
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	name := d.Get("name").(string)
	urlOutput, err := findQueueURLByName(ctx, conn, name)

	if err != nil {
		return diag.Errorf("reading SQS Queue (%s) URL: %s", name, err)
	}

	queueURL := aws.ToString(urlOutput)
	attributesOutput, err := findQueueAttributeByTwoPartKey(ctx, conn, queueURL, types.QueueAttributeNameQueueArn)

	if err != nil {
		return diag.Errorf("reading SQS Queue (%s) ARN attribute: %s", queueURL, err)
	}

	d.SetId(queueURL)
	d.Set("arn", attributesOutput)
	d.Set("url", queueURL)

	tags, err := listTags(ctx, conn, queueURL)

	if errs.IsUnsupportedOperationInPartitionError(meta.(*conns.AWSClient).Partition, err) {
		// Some partitions may not support tagging, giving error
		log.Printf("[WARN] failed listing tags for SQS Queue (%s): %s", d.Id(), err)
		return nil
	}

	if err != nil {
		return diag.Errorf("listing tags for SQS Queue (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAWS().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return diag.Errorf("setting tags: %s", err)
	}

	return nil
}

func findQueueURLByName(ctx context.Context, conn *sqs.Client, name string) (*string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	}

	output, err := conn.GetQueueUrl(ctx, input)

	if tfawserr.ErrCodeEquals(err, errCodeQueueDoesNotExist) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.QueueUrl == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.QueueUrl, nil
}
