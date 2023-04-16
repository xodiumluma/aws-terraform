package elb

import ( // nosemgrep:ci.aws-sdk-go-multiple-service-imports
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_elb", name="Classic Load Balancer")
// @Tags(identifierAttribute="id")
func ResourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLoadBalancerCreate,
		ReadWithoutTimeout:   resourceLoadBalancerRead,
		UpdateWithoutTimeout: resourceLoadBalancerUpdate,
		DeleteWithoutTimeout: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: verify.SetTagsDiff,

		Schema: map[string]*schema.Schema{
			"access_logs": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bucket_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      60,
							ValidateFunc: ValidAccessLogsInterval,
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"connection_draining": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"connection_draining_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"cross_zone_load_balancing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"desync_mitigation_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "defensive",
				ValidateFunc: validation.StringInSlice([]string{
					"monitor",
					"defensive",
					"strictest",
				}, false),
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_check": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"healthy_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(2, 10),
						},
						"interval": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(5, 300),
						},
						"target": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidHeathCheckTarget,
						},
						"timeout": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(2, 60),
						},
						"unhealthy_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(2, 10),
						},
					},
				},
			},
			"idle_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(1, 4000),
			},
			"instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"internal": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"listener": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"instance_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateListenerProtocol(),
						},
						"lb_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},
						"lb_protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateListenerProtocol(),
						},
						"ssl_certificate_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: verify.ValidARN,
						},
					},
				},
				Set: ListenerHash,
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  ValidName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validNamePrefix,
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_security_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ELBConn()

	var elbName string
	if v, ok := d.GetOk("name"); ok {
		elbName = v.(string)
	} else {
		if v, ok := d.GetOk("name_prefix"); ok {
			elbName = id.PrefixedUniqueId(v.(string))
		} else {
			elbName = id.PrefixedUniqueId("tf-lb-")
		}
		d.Set("name", elbName)
	}

	// Expand the "listener" set to aws-sdk-go compat []*elb.Listener
	listeners, err := ExpandListeners(d.Get("listener").(*schema.Set).List())
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating ELB Classic Load Balancer (%s): %s", elbName, err)
	}
	// Provision the elb
	input := &elb.CreateLoadBalancerInput{
		LoadBalancerName: aws.String(elbName),
		Listeners:        listeners,
		Tags:             GetTagsIn(ctx),
	}

	if _, ok := d.GetOk("internal"); ok {
		input.Scheme = aws.String("internal")
	}

	if v, ok := d.GetOk("availability_zones"); ok {
		input.AvailabilityZones = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("security_groups"); ok {
		input.SecurityGroups = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("subnets"); ok {
		input.Subnets = flex.ExpandStringSet(v.(*schema.Set))
	}

	_, err = tfresource.RetryWhenAWSErrCodeEquals(ctx, 5*time.Minute, func() (interface{}, error) {
		return conn.CreateLoadBalancerWithContext(ctx, input)
	}, elb.ErrCodeCertificateNotFoundException)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating ELB Classic Load Balancer (%s): %s", elbName, err)
	}

	// Assign the elb's unique identifier for use later
	d.SetId(elbName)
	log.Printf("[INFO] ELB ID: %s", d.Id())

	return append(diags, resourceLoadBalancerUpdate(ctx, d, meta)...)
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ELBConn()

	lb, err := FindLoadBalancerByName(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] ELB Classic Load Balancer (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading ELB Classic Load Balancer (%s): %s", d.Id(), err)
	}

	arn := arn.ARN{
		Partition: meta.(*conns.AWSClient).Partition,
		Region:    meta.(*conns.AWSClient).Region,
		Service:   "elasticloadbalancing",
		AccountID: meta.(*conns.AWSClient).AccountID,
		Resource:  fmt.Sprintf("loadbalancer/%s", d.Id()),
	}
	d.Set("arn", arn.String())

	if err := flattenLoadBalancerResource(ctx, d, meta.(*conns.AWSClient).EC2Conn(), conn, lb); err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}
	return diags
}

// flattenLoadBalancerResource takes a *elb.LoadBalancerDescription and populates all respective resource fields.
func flattenLoadBalancerResource(ctx context.Context, d *schema.ResourceData, ec2conn *ec2.EC2, elbconn *elb.ELB, lb *elb.LoadBalancerDescription) error {
	describeAttrsOpts := &elb.DescribeLoadBalancerAttributesInput{
		LoadBalancerName: aws.String(d.Id()),
	}
	describeAttrsResp, err := elbconn.DescribeLoadBalancerAttributesWithContext(ctx, describeAttrsOpts)
	if err != nil {
		return fmt.Errorf("Error retrieving ELB: %s", err)
	}

	lbAttrs := describeAttrsResp.LoadBalancerAttributes

	d.Set("name", lb.LoadBalancerName)
	d.Set("dns_name", lb.DNSName)
	d.Set("zone_id", lb.CanonicalHostedZoneNameID)

	var scheme bool
	if lb.Scheme != nil {
		scheme = aws.StringValue(lb.Scheme) == "internal"
	}
	d.Set("internal", scheme)
	d.Set("availability_zones", flex.FlattenStringList(lb.AvailabilityZones))
	d.Set("instances", flattenInstances(lb.Instances))
	d.Set("listener", flattenListeners(lb.ListenerDescriptions))
	d.Set("security_groups", flex.FlattenStringList(lb.SecurityGroups))

	if lb.SourceSecurityGroup != nil {
		group := lb.SourceSecurityGroup.GroupName
		if v := aws.StringValue(lb.SourceSecurityGroup.OwnerAlias); v != "" {
			group = aws.String(v + "/" + aws.StringValue(lb.SourceSecurityGroup.GroupName))
		}
		d.Set("source_security_group", group)

		// Manually look up the ELB Security Group ID, since it's not provided
		if lb.VPCId != nil {
			sg, err := tfec2.FindSecurityGroupByNameAndVPCIDAndOwnerID(ctx, ec2conn, aws.StringValue(lb.SourceSecurityGroup.GroupName), aws.StringValue(lb.VPCId), aws.StringValue(lb.SourceSecurityGroup.OwnerAlias))
			if err != nil {
				return fmt.Errorf("Error looking up ELB Security Group ID: %w", err)
			} else {
				d.Set("source_security_group_id", sg.GroupId)
			}
		}
	}
	d.Set("subnets", flex.FlattenStringList(lb.Subnets))
	if lbAttrs.ConnectionSettings != nil {
		d.Set("idle_timeout", lbAttrs.ConnectionSettings.IdleTimeout)
	}
	d.Set("connection_draining", lbAttrs.ConnectionDraining.Enabled)
	d.Set("connection_draining_timeout", lbAttrs.ConnectionDraining.Timeout)
	d.Set("cross_zone_load_balancing", lbAttrs.CrossZoneLoadBalancing.Enabled)
	if lbAttrs.AccessLog != nil {
		// The AWS API does not allow users to remove access_logs, only disable them.
		// During creation of the ELB, Terraform sets the access_logs to disabled,
		// so there should not be a case where lbAttrs.AccessLog above is nil.

		// Here we do not record the remove value of access_log if:
		// - there is no access_log block in the configuration
		// - the remote access_logs are disabled
		//
		// This indicates there is no access_log in the configuration.
		// - externally added access_logs will be enabled, so we'll detect the drift
		// - locally added access_logs will be in the config, so we'll add to the
		// API/state
		// See https://github.com/hashicorp/terraform/issues/10138
		_, n := d.GetChange("access_logs")
		elbal := lbAttrs.AccessLog
		nl := n.([]interface{})
		if len(nl) == 0 && !*elbal.Enabled {
			elbal = nil
		}
		if err := d.Set("access_logs", flattenAccessLog(elbal)); err != nil {
			return fmt.Errorf("reading ELB Classic Load Balancer (%s): setting access_logs: %w", d.Id(), err)
		}
	}

	for _, attr := range lbAttrs.AdditionalAttributes {
		switch aws.StringValue(attr.Key) {
		case "elb.http.desyncmitigationmode":
			d.Set("desync_mitigation_mode", attr.Value)
		}
	}

	// There's only one health check, so save that to state as we
	// currently can
	if aws.StringValue(lb.HealthCheck.Target) != "" {
		d.Set("health_check", FlattenHealthCheck(lb.HealthCheck))
	}

	return nil
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ELBConn()

	if d.HasChange("listener") {
		o, n := d.GetChange("listener")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		remove, _ := ExpandListeners(os.Difference(ns).List())
		add, err := ExpandListeners(ns.Difference(os).List())
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating ELB Classic Load Balancer (%s): %s", d.Id(), err)
		}

		if len(remove) > 0 {
			ports := make([]*int64, 0, len(remove))
			for _, listener := range remove {
				ports = append(ports, listener.LoadBalancerPort)
			}

			deleteListenersOpts := &elb.DeleteLoadBalancerListenersInput{
				LoadBalancerName:  aws.String(d.Id()),
				LoadBalancerPorts: ports,
			}

			log.Printf("[DEBUG] ELB Delete Listeners opts: %s", deleteListenersOpts)
			_, err := conn.DeleteLoadBalancerListenersWithContext(ctx, deleteListenersOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure removing outdated ELB listeners: %s", err)
			}
		}

		if len(add) > 0 {
			input := &elb.CreateLoadBalancerListenersInput{
				LoadBalancerName: aws.String(d.Id()),
				Listeners:        add,
			}

			// Occasionally AWS will error with a 'duplicate listener', without any
			// other listeners on the ELB. Retry here to eliminate that.
			err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
				_, err := conn.CreateLoadBalancerListenersWithContext(ctx, input)
				if err != nil {
					if tfawserr.ErrCodeEquals(err, elb.ErrCodeDuplicateListenerException) {
						log.Printf("[DEBUG] Duplicate listener found for ELB (%s), retrying", d.Id())
						return retry.RetryableError(err)
					}
					if tfawserr.ErrMessageContains(err, elb.ErrCodeCertificateNotFoundException, "Server Certificate not found for the key: arn") {
						log.Printf("[DEBUG] SSL Cert not found for given ARN, retrying")
						return retry.RetryableError(err)
					}

					// Didn't recognize the error, so shouldn't retry.
					return retry.NonRetryableError(err)
				}
				// Successful creation
				return nil
			})
			if tfresource.TimedOut(err) {
				_, err = conn.CreateLoadBalancerListenersWithContext(ctx, input)
			}
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure adding new or updated ELB listeners: %s", err)
			}
		}
	}

	// If we currently have instances, or did have instances,
	// we want to figure out what to add and remove from the load
	// balancer
	if d.HasChange("instances") {
		o, n := d.GetChange("instances")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := ExpandInstanceString(os.Difference(ns).List())
		add := ExpandInstanceString(ns.Difference(os).List())

		if len(add) > 0 {
			registerInstancesOpts := elb.RegisterInstancesWithLoadBalancerInput{
				LoadBalancerName: aws.String(d.Id()),
				Instances:        add,
			}

			_, err := conn.RegisterInstancesWithLoadBalancerWithContext(ctx, &registerInstancesOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure registering instances with ELB: %s", err)
			}
		}
		if len(remove) > 0 {
			deRegisterInstancesOpts := elb.DeregisterInstancesFromLoadBalancerInput{
				LoadBalancerName: aws.String(d.Id()),
				Instances:        remove,
			}

			_, err := conn.DeregisterInstancesFromLoadBalancerWithContext(ctx, &deRegisterInstancesOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure deregistering instances from ELB: %s", err)
			}
		}
	}

	if d.HasChanges("cross_zone_load_balancing", "idle_timeout", "access_logs", "desync_mitigation_mode") {
		attrs := elb.ModifyLoadBalancerAttributesInput{
			LoadBalancerName: aws.String(d.Get("name").(string)),
			LoadBalancerAttributes: &elb.LoadBalancerAttributes{
				AdditionalAttributes: []*elb.AdditionalAttribute{
					{
						Key:   aws.String("elb.http.desyncmitigationmode"),
						Value: aws.String(d.Get("desync_mitigation_mode").(string)),
					},
				},
				CrossZoneLoadBalancing: &elb.CrossZoneLoadBalancing{
					Enabled: aws.Bool(d.Get("cross_zone_load_balancing").(bool)),
				},
				ConnectionSettings: &elb.ConnectionSettings{
					IdleTimeout: aws.Int64(int64(d.Get("idle_timeout").(int))),
				},
			},
		}

		logs := d.Get("access_logs").([]interface{})
		if len(logs) == 1 {
			l := logs[0].(map[string]interface{})
			attrs.LoadBalancerAttributes.AccessLog = &elb.AccessLog{
				Enabled:        aws.Bool(l["enabled"].(bool)),
				EmitInterval:   aws.Int64(int64(l["interval"].(int))),
				S3BucketName:   aws.String(l["bucket"].(string)),
				S3BucketPrefix: aws.String(l["bucket_prefix"].(string)),
			}
		} else if len(logs) == 0 {
			// disable access logs
			attrs.LoadBalancerAttributes.AccessLog = &elb.AccessLog{
				Enabled: aws.Bool(false),
			}
		}

		log.Printf("[DEBUG] ELB Modify Load Balancer Attributes Request: %#v", attrs)
		_, err := conn.ModifyLoadBalancerAttributesWithContext(ctx, &attrs)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "Failure configuring ELB attributes: %s", err)
		}
	}

	// We have to do these changes separately from everything else since
	// they have some weird undocumented rules. You can't set the timeout
	// without having connection draining to true, so we set that to true,
	// set the timeout, then reset it to false if requested.
	if d.HasChanges("connection_draining", "connection_draining_timeout") {
		// We do timeout changes first since they require us to set draining
		// to true for a hot second.
		if d.HasChange("connection_draining_timeout") {
			attrs := elb.ModifyLoadBalancerAttributesInput{
				LoadBalancerName: aws.String(d.Get("name").(string)),
				LoadBalancerAttributes: &elb.LoadBalancerAttributes{
					ConnectionDraining: &elb.ConnectionDraining{
						Enabled: aws.Bool(true),
						Timeout: aws.Int64(int64(d.Get("connection_draining_timeout").(int))),
					},
				},
			}

			_, err := conn.ModifyLoadBalancerAttributesWithContext(ctx, &attrs)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure configuring ELB attributes: %s", err)
			}
		}

		// Then we always set connection draining even if there is no change.
		// This lets us reset to "false" if requested even with a timeout
		// change.
		attrs := elb.ModifyLoadBalancerAttributesInput{
			LoadBalancerName: aws.String(d.Get("name").(string)),
			LoadBalancerAttributes: &elb.LoadBalancerAttributes{
				ConnectionDraining: &elb.ConnectionDraining{
					Enabled: aws.Bool(d.Get("connection_draining").(bool)),
				},
			},
		}

		_, err := conn.ModifyLoadBalancerAttributesWithContext(ctx, &attrs)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "Failure configuring ELB attributes: %s", err)
		}
	}

	if d.HasChange("health_check") {
		hc := d.Get("health_check").([]interface{})
		if len(hc) > 0 {
			check := hc[0].(map[string]interface{})
			configureHealthCheckOpts := elb.ConfigureHealthCheckInput{
				LoadBalancerName: aws.String(d.Id()),
				HealthCheck: &elb.HealthCheck{
					HealthyThreshold:   aws.Int64(int64(check["healthy_threshold"].(int))),
					UnhealthyThreshold: aws.Int64(int64(check["unhealthy_threshold"].(int))),
					Interval:           aws.Int64(int64(check["interval"].(int))),
					Target:             aws.String(check["target"].(string)),
					Timeout:            aws.Int64(int64(check["timeout"].(int))),
				},
			}
			_, err := conn.ConfigureHealthCheckWithContext(ctx, &configureHealthCheckOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure configuring health check for ELB: %s", err)
			}
		}
	}

	if d.HasChange("security_groups") {
		applySecurityGroupsOpts := elb.ApplySecurityGroupsToLoadBalancerInput{
			LoadBalancerName: aws.String(d.Id()),
			SecurityGroups:   flex.ExpandStringSet(d.Get("security_groups").(*schema.Set)),
		}

		_, err := conn.ApplySecurityGroupsToLoadBalancerWithContext(ctx, &applySecurityGroupsOpts)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "Failure applying security groups to ELB: %s", err)
		}
	}

	if d.HasChange("availability_zones") {
		o, n := d.GetChange("availability_zones")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removed := flex.ExpandStringSet(os.Difference(ns))
		added := flex.ExpandStringSet(ns.Difference(os))

		if len(added) > 0 {
			enableOpts := &elb.EnableAvailabilityZonesForLoadBalancerInput{
				LoadBalancerName:  aws.String(d.Id()),
				AvailabilityZones: added,
			}

			log.Printf("[DEBUG] ELB enable availability zones opts: %s", enableOpts)
			_, err := conn.EnableAvailabilityZonesForLoadBalancerWithContext(ctx, enableOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure enabling ELB availability zones: %s", err)
			}
		}

		if len(removed) > 0 {
			disableOpts := &elb.DisableAvailabilityZonesForLoadBalancerInput{
				LoadBalancerName:  aws.String(d.Id()),
				AvailabilityZones: removed,
			}

			log.Printf("[DEBUG] ELB disable availability zones opts: %s", disableOpts)
			_, err := conn.DisableAvailabilityZonesForLoadBalancerWithContext(ctx, disableOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure disabling ELB availability zones: %s", err)
			}
		}
	}

	if d.HasChange("subnets") {
		o, n := d.GetChange("subnets")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)

		removed := flex.ExpandStringSet(os.Difference(ns))
		added := flex.ExpandStringSet(ns.Difference(os))

		if len(removed) > 0 {
			detachOpts := &elb.DetachLoadBalancerFromSubnetsInput{
				LoadBalancerName: aws.String(d.Id()),
				Subnets:          removed,
			}

			log.Printf("[DEBUG] ELB detach subnets opts: %s", detachOpts)
			_, err := conn.DetachLoadBalancerFromSubnetsWithContext(ctx, detachOpts)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure removing ELB subnets: %s", err)
			}
		}

		if len(added) > 0 {
			attachOpts := &elb.AttachLoadBalancerToSubnetsInput{
				LoadBalancerName: aws.String(d.Id()),
				Subnets:          added,
			}

			log.Printf("[DEBUG] ELB attach subnets opts: %s", attachOpts)
			err := retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
				_, err := conn.AttachLoadBalancerToSubnetsWithContext(ctx, attachOpts)
				if err != nil {
					if tfawserr.ErrMessageContains(err, elb.ErrCodeInvalidConfigurationRequestException, "cannot be attached to multiple subnets in the same AZ") {
						// eventually consistent issue with removing a subnet in AZ1 and
						// immediately adding a new one in the same AZ
						log.Printf("[DEBUG] retrying az association")
						return retry.RetryableError(err)
					}
					return retry.NonRetryableError(err)
				}
				return nil
			})
			if tfresource.TimedOut(err) {
				_, err = conn.AttachLoadBalancerToSubnetsWithContext(ctx, attachOpts)
			}
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "Failure adding ELB subnets: %s", err)
			}
		}
	}

	return append(diags, resourceLoadBalancerRead(ctx, d, meta)...)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ELBConn()

	log.Printf("[INFO] Deleting ELB Classic Load Balancer: %s", d.Id())
	_, err := conn.DeleteLoadBalancerWithContext(ctx, &elb.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(d.Id()),
	})

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting ELB Classic Load Balancer (%s): %s", d.Id(), err)
	}

	err = cleanupNetworkInterfaces(ctx, meta.(*conns.AWSClient).EC2Conn(), d.Id())

	if err != nil {
		diags = sdkdiag.AppendWarningf(diags, "cleaning up ELB Classic Load Balancer (%s) ENIs: %s", d.Id(), err)
	}

	return diags
}

func FindLoadBalancerByName(ctx context.Context, conn *elb.ELB, name string) (*elb.LoadBalancerDescription, error) {
	input := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: aws.StringSlice([]string{name}),
	}

	output, err := conn.DescribeLoadBalancersWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, elb.ErrCodeAccessPointNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.LoadBalancerDescriptions) == 0 || output.LoadBalancerDescriptions[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.LoadBalancerDescriptions); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	// Eventual consistency check.
	if aws.StringValue(output.LoadBalancerDescriptions[0].LoadBalancerName) != name {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return output.LoadBalancerDescriptions[0], nil
}

func ListenerHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["instance_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["instance_protocol"].(string))))
	buf.WriteString(fmt.Sprintf("%d-", m["lb_port"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["lb_protocol"].(string))))

	if v, ok := m["ssl_certificate_id"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return create.StringHashcode(buf.String())
}

func ValidAccessLogsInterval(v interface{}, k string) (ws []string, errors []error) {
	value := v.(int)

	// Check if the value is either 5 or 60 (minutes).
	if value != 5 && value != 60 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Access Logs interval \"%d\". "+
				"Valid intervals are either 5 or 60 (minutes).",
			k, value))
	}
	return
}

func ValidHeathCheckTarget(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// Parse the Health Check target value.
	matches := regexp.MustCompile(`\A(\w+):(\d+)(.+)?\z`).FindStringSubmatch(value)

	// Check if the value contains a valid target.
	if matches == nil || len(matches) < 1 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check: %s",
			k, value))

		// Invalid target? Return immediately,
		// there is no need to collect other
		// errors.
		return ws, errors
	}

	// Check if the value contains a valid protocol.
	if !isValidProtocol(matches[1]) {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check protocol %q. "+
				"Valid protocols are either %q, %q, %q, or %q.",
			k, matches[1], "TCP", "SSL", "HTTP", "HTTPS"))
	}

	// Check if the value contains a valid port range.
	port, _ := strconv.Atoi(matches[2])
	if port < 1 || port > 65535 {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid Health Check target port \"%d\". "+
				"Valid port is in the range from 1 to 65535 inclusive.",
			k, port))
	}

	switch strings.ToLower(matches[1]) {
	case "tcp", "ssl":
		// Check if value is in the form <PROTOCOL>:<PORT> for TCP and/or SSL.
		if matches[3] != "" {
			errors = append(errors, fmt.Errorf(
				"%q cannot contain a path in the Health Check target: %s",
				k, value))
		}

	case "http", "https":
		// Check if value is in the form <PROTOCOL>:<PORT>/<PATH> for HTTP and/or HTTPS.
		if matches[3] == "" {
			errors = append(errors, fmt.Errorf(
				"%q must contain a path in the Health Check target: %s",
				k, value))
		}

		// Cannot be longer than 1024 multibyte characters.
		if len([]rune(matches[3])) > 1024 {
			errors = append(errors, fmt.Errorf("%q cannot contain a path longer "+
				"than 1024 characters in the Health Check target: %s",
				k, value))
		}
	}

	return ws, errors
}

func isValidProtocol(s string) bool {
	if s == "" {
		return false
	}
	s = strings.ToLower(s)

	validProtocols := map[string]bool{
		"http":  true,
		"https": true,
		"ssl":   true,
		"tcp":   true,
	}

	if _, ok := validProtocols[s]; !ok {
		return false
	}

	return true
}

func validateListenerProtocol() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"HTTP",
		"HTTPS",
		"SSL",
		"TCP",
	}, true)
}

// ELB automatically creates ENI(s) on creation
// but the cleanup is asynchronous and may take time
// which then blocks IGW, SG or VPC on deletion
// So we make the cleanup "synchronous" here
func cleanupNetworkInterfaces(ctx context.Context, conn *ec2.EC2, name string) error {
	// https://aws.amazon.com/premiumsupport/knowledge-center/elb-find-load-balancer-IP/.
	networkInterfaces, err := tfec2.FindNetworkInterfacesByAttachmentInstanceOwnerIDAndDescription(ctx, conn, "amazon-elb", "ELB "+name)

	if err != nil {
		return err
	}

	var errs *multierror.Error

	for _, networkInterface := range networkInterfaces {
		if networkInterface.Attachment == nil {
			continue
		}

		attachmentID := aws.StringValue(networkInterface.Attachment.AttachmentId)
		networkInterfaceID := aws.StringValue(networkInterface.NetworkInterfaceId)

		err = tfec2.DetachNetworkInterface(ctx, conn, networkInterfaceID, attachmentID, tfec2.NetworkInterfaceDetachedTimeout)

		if err != nil {
			errs = multierror.Append(errs, err)

			continue
		}

		err = tfec2.DeleteNetworkInterface(ctx, conn, networkInterfaceID)

		if err != nil {
			errs = multierror.Append(errs, err)

			continue
		}
	}

	return errs.ErrorOrNil()
}
