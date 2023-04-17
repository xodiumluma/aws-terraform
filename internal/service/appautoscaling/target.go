package appautoscaling

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/applicationautoscaling"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

// @SDKResource("aws_appautoscaling_target")
func ResourceTarget() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceTargetPut,
		ReadWithoutTimeout:   resourceTargetRead,
		UpdateWithoutTimeout: resourceTargetPut,
		DeleteWithoutTimeout: resourceTargetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTargetImport,
		},

		Schema: map[string]*schema.Schema{
			"max_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"min_capacity": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scalable_dimension": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTargetPut(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AppAutoScalingConn()

	var targetOpts applicationautoscaling.RegisterScalableTargetInput

	targetOpts.MaxCapacity = aws.Int64(int64(d.Get("max_capacity").(int)))
	targetOpts.MinCapacity = aws.Int64(int64(d.Get("min_capacity").(int)))
	targetOpts.ResourceId = aws.String(d.Get("resource_id").(string))
	targetOpts.ScalableDimension = aws.String(d.Get("scalable_dimension").(string))
	targetOpts.ServiceNamespace = aws.String(d.Get("service_namespace").(string))

	if roleArn, exists := d.GetOk("role_arn"); exists {
		targetOpts.RoleARN = aws.String(roleArn.(string))
	}

	log.Printf("[DEBUG] Application autoscaling target create configuration %s", targetOpts)
	var err error
	err = retry.RetryContext(ctx, propagationTimeout, func() *retry.RetryError {
		_, err = conn.RegisterScalableTargetWithContext(ctx, &targetOpts)

		if err != nil {
			if tfawserr.ErrMessageContains(err, applicationautoscaling.ErrCodeValidationException, "Unable to assume IAM role") {
				return retry.RetryableError(err)
			}
			if tfawserr.ErrMessageContains(err, applicationautoscaling.ErrCodeValidationException, "ECS service doesn't exist") {
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}

		return nil
	})
	if tfresource.TimedOut(err) {
		_, err = conn.RegisterScalableTargetWithContext(ctx, &targetOpts)
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Application AutoScaling Target: %s", err)
	}

	d.SetId(d.Get("resource_id").(string))

	return append(diags, resourceTargetRead(ctx, d, meta)...)
}

func resourceTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AppAutoScalingConn()

	outputRaw, err := tfresource.RetryWhenNewResourceNotFound(ctx, 2*time.Minute,
		func() (interface{}, error) {
			return FindTargetByThreePartKey(ctx, conn, d.Id(), d.Get("service_namespace").(string), d.Get("scalable_dimension").(string))
		},
		d.IsNewResource(),
	)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Application AutoScaling Target (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Application AutoScaling Target (%s): %s", d.Id(), err)
	}

	t := outputRaw.(*applicationautoscaling.ScalableTarget)

	d.Set("max_capacity", t.MaxCapacity)
	d.Set("min_capacity", t.MinCapacity)
	d.Set("resource_id", t.ResourceId)
	d.Set("role_arn", t.RoleARN)
	d.Set("scalable_dimension", t.ScalableDimension)
	d.Set("service_namespace", t.ServiceNamespace)

	return diags
}

func resourceTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).AppAutoScalingConn()

	input := &applicationautoscaling.DeregisterScalableTargetInput{
		ResourceId:        aws.String(d.Id()),
		ScalableDimension: aws.String(d.Get("scalable_dimension").(string)),
		ServiceNamespace:  aws.String(d.Get("service_namespace").(string)),
	}

	log.Printf("[INFO] Deleting Application AutoScaling Target: %s", d.Id())
	_, err := conn.DeregisterScalableTargetWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, applicationautoscaling.ErrCodeObjectNotFoundException) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Application AutoScaling Target (%s): %s", d.Id(), err)
	}

	_, err = tfresource.RetryUntilNotFound(ctx, 5*time.Minute, func() (interface{}, error) {
		return FindTargetByThreePartKey(ctx, conn, d.Id(), d.Get("service_namespace").(string), d.Get("scalable_dimension").(string))
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for Application AutoScaling Target (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func FindTargetByThreePartKey(ctx context.Context, conn *applicationautoscaling.ApplicationAutoScaling, resourceID, namespace, dimension string) (*applicationautoscaling.ScalableTarget, error) {
	input := &applicationautoscaling.DescribeScalableTargetsInput{
		ResourceIds:       aws.StringSlice([]string{resourceID}),
		ScalableDimension: aws.String(dimension),
		ServiceNamespace:  aws.String(namespace),
	}
	var output []*applicationautoscaling.ScalableTarget

	err := conn.DescribeScalableTargetsPagesWithContext(ctx, input, func(page *applicationautoscaling.DescribeScalableTargetsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.ScalableTargets {
			if v != nil {
				output = append(output, v)
			}
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	if len(output) == 0 || output[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	target := output[0]

	if aws.StringValue(target.ResourceId) != resourceID || aws.StringValue(target.ScalableDimension) != dimension || aws.StringValue(target.ServiceNamespace) != namespace {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return target, nil
}

func resourceTargetImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	idParts := strings.Split(d.Id(), "/")

	if len(idParts) < 3 {
		return nil, fmt.Errorf("unexpected format (%q), expected <service-namespace>/<resource-id>/<scalable-dimension>", d.Id())
	}

	serviceNamespace := idParts[0]
	resourceId := strings.Join(idParts[1:len(idParts)-1], "/")
	scalableDimension := idParts[len(idParts)-1]

	if serviceNamespace == "" || resourceId == "" || scalableDimension == "" {
		return nil, fmt.Errorf("unexpected format (%q), expected <service-namespace>/<resource-id>/<scalable-dimension>", d.Id())
	}

	d.Set("service_namespace", serviceNamespace)
	d.Set("resource_id", resourceId)
	d.Set("scalable_dimension", scalableDimension)
	d.SetId(resourceId)

	return []*schema.ResourceData{d}, nil
}
