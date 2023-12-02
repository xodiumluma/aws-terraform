// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package logs

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

// @SDKResource("aws_cloudwatch_log_subscription_filter")
func resourceSubscriptionFilter() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceSubscriptionFilterPut,
		ReadWithoutTimeout:   resourceSubscriptionFilterRead,
		UpdateWithoutTimeout: resourceSubscriptionFilterPut,
		DeleteWithoutTimeout: resourceSubscriptionFilterDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSubscriptionFilterImport,
		},

		Schema: map[string]*schema.Schema{
			"destination_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"distribution": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          types.DistributionByLogStream,
				ValidateDiagFunc: enum.Validate[types.Distribution](),
			},
			"filter_pattern": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
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
				ValidateFunc: validation.StringLenBetween(1, 512),
			},
			"role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidARN,
			},
		},
	}
}

func resourceSubscriptionFilterPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	logGroupName := d.Get("log_group_name").(string)
	name := d.Get("name").(string)
	input := &cloudwatchlogs.PutSubscriptionFilterInput{
		DestinationArn: aws.String(d.Get("destination_arn").(string)),
		FilterName:     aws.String(name),
		FilterPattern:  aws.String(d.Get("filter_pattern").(string)),
		LogGroupName:   aws.String(logGroupName),
	}

	if v, ok := d.GetOk("distribution"); ok {
		input.Distribution = types.Distribution(v.(string))
	}

	if v, ok := d.GetOk("role_arn"); ok {
		input.RoleArn = aws.String(v.(string))
	}

	_, err := tfresource.RetryWhen(ctx, 5*time.Minute,
		func() (interface{}, error) {
			return conn.PutSubscriptionFilter(ctx, input)
		},
		func(err error) (bool, error) {
			if errs.IsAErrorMessageContains[*types.InvalidParameterException](err, "Could not deliver test message to specified") {
				return true, err
			}

			if errs.IsAErrorMessageContains[*types.InvalidParameterException](err, "Could not execute the lambda function") {
				return true, err
			}

			if errs.IsAErrorMessageContains[*types.OperationAbortedException](err, "Please try again") {
				return true, err
			}

			return false, err
		})

	if err != nil {
		return diag.Errorf("putting CloudWatch Logs Subscription Filter (%s): %s", name, err)
	}

	d.SetId(subscriptionFilterID(logGroupName))

	return nil
}

func resourceSubscriptionFilterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	subscriptionFilter, err := findSubscriptionFilterByTwoPartKey(ctx, conn, d.Get("log_group_name").(string), d.Get("name").(string))

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] CloudWatch Logs Subscription Filter (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("reading CloudWatch Logs Subscription Filter (%s): %s", d.Id(), err)
	}

	d.Set("destination_arn", subscriptionFilter.DestinationArn)
	d.Set("distribution", subscriptionFilter.Distribution)
	d.Set("filter_pattern", subscriptionFilter.FilterPattern)
	d.Set("log_group_name", subscriptionFilter.LogGroupName)
	d.Set("name", subscriptionFilter.FilterName)
	d.Set("role_arn", subscriptionFilter.RoleArn)

	return nil
}

func resourceSubscriptionFilterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LogsClient(ctx)

	log.Printf("[INFO] Deleting CloudWatch Logs Subscription Filter: %s", d.Id())
	_, err := conn.DeleteSubscriptionFilter(ctx, &cloudwatchlogs.DeleteSubscriptionFilterInput{
		FilterName:   aws.String(d.Get("name").(string)),
		LogGroupName: aws.String(d.Get("log_group_name").(string)),
	})

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil
	}

	if err != nil {
		return diag.Errorf("deleting CloudWatch Logs Subscription Filter (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceSubscriptionFilterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), "|")
	if len(idParts) < 2 {
		return nil, fmt.Errorf("unexpected format of ID (%q), expected <log-group-name>|<filter-name>", d.Id())
	}

	logGroupName := idParts[0]
	filterNamePrefix := idParts[1]

	d.Set("log_group_name", logGroupName)
	d.Set("name", filterNamePrefix)
	d.SetId(subscriptionFilterID(filterNamePrefix))

	return []*schema.ResourceData{d}, nil
}

func subscriptionFilterID(log_group_name string) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s-", log_group_name)) // only one filter allowed per log_group_name at the moment

	return fmt.Sprintf("cwlsf-%d", create.StringHashcode(buf.String()))
}

func findSubscriptionFilterByTwoPartKey(ctx context.Context, conn *cloudwatchlogs.Client, logGroupName, name string) (*types.SubscriptionFilter, error) {
	input := &cloudwatchlogs.DescribeSubscriptionFiltersInput{
		FilterNamePrefix: aws.String(name),
		LogGroupName:     aws.String(logGroupName),
	}

	pages := cloudwatchlogs.NewDescribeSubscriptionFiltersPaginator(conn, input)
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

		for _, v := range page.SubscriptionFilters {
			if aws.ToString(v.FilterName) == name {
				return &v, nil
			}
		}
	}

	return nil, tfresource.NewEmptyResultError(input)
}
