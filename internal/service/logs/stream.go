// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_cloudwatch_log_stream")
func resourceStream() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceStreamCreate,
		ReadWithoutTimeout:   resourceStreamRead,
		DeleteWithoutTimeout: resourceStreamDelete,

		Importer: &schema.ResourceImporter{
			State: resourceStreamImport,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validStreamName,
			},
		},
	}
}

func resourceStreamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	name := d.Get("name").(string)
	input := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(d.Get("log_group_name").(string)),
		LogStreamName: aws.String(name),
	}

	_, err := conn.CreateLogStream(ctx, input)

	if err != nil {
		return diag.Errorf("creating CloudWatch Logs Log Stream (%s): %s", name, err)
	}

	d.SetId(name)

	_, err = tfresource.RetryWhenNotFound(ctx, propagationTimeout, func() (interface{}, error) {
		return findLogStreamByTwoPartKey(ctx, conn, d.Get("log_group_name").(string), d.Id())
	})

	if err != nil {
		return diag.Errorf("waiting for CloudWatch Logs Log Stream (%s) create: %s", d.Id(), err)
	}

	return resourceStreamRead(ctx, d, meta)
}

func resourceStreamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	ls, err := findLogStreamByTwoPartKey(ctx, conn, d.Get("log_group_name").(string), d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudWatch Logs Log Stream (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading CloudWatch Logs Log Stream (%s): %s", d.Id(), err)
	}

	d.Set("arn", ls.Arn)
	d.Set("name", ls.LogStreamName)

	return nil
}

func resourceStreamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	log.Printf("[INFO] Deleting CloudWatch Logs Log Stream: %s", d.Id())
	_, err := conn.DeleteLogStream(ctx, &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  aws.String(d.Get("log_group_name").(string)),
		LogStreamName: aws.String(d.Id()),
	})

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting CloudWatch Logs Log Stream (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceStreamImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("wrong format of import ID (%s), use: 'log-group-name:log-stream-name'", d.Id())
	}

	logGroupName := parts[0]
	logStreamName := parts[1]

	d.SetId(logStreamName)
	d.Set("log_group_name", logGroupName)

	return []*schema.ResourceData{d}, nil
}

func findLogStreamByTwoPartKey(ctx context.Context, conn *cloudwatchlogs.Client, logGroupName, name string) (*types.LogStream, error) { // nosemgrep:ci.logs-in-func-name
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(logGroupName),
		LogStreamNamePrefix: aws.String(name),
	}

	pages := cloudwatchlogs.NewDescribeLogStreamsPaginator(conn, input)
	for pages.HasMorePages() {
		page, err := pages.NextPage(ctx)

		if errs.IsA[*types.ResourceNotFoundException](err) {
			return nil, &retry.NotFoundError{
				LastError:   err,
				LastRequest: input,
			}
		}

		if err != nil {
			return nil, err
		}

		for _, v := range page.LogStreams {
			if aws.ToString(v.LogStreamName) == name {
				return &v, nil
			}
		}
	}

	return nil, tfresource.NewEmptyResultError(input)
}

func validStreamName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if regexache.MustCompile(`:`).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"colons not allowed in %q:", k))
	}
	if len(value) < 1 || len(value) > 512 {
		errors = append(errors, fmt.Errorf(
			"%q must be between 1 and 512 characters: %q", k, value))
	}

	return
}
