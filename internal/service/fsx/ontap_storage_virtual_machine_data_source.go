// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fsx

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// @SDKDataSource("aws_fsx_ontap_storage_virtual_machine", name="ONTAP Storage Virtual Machine")
func DataSourceONTAPStorageVirtualMachine() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceONTAPStorageVirtualMachineRead,

		Schema: map[string]*schema.Schema{
			"active_directory_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"netbios_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_managed_active_directory_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_ips": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"domain_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"file_system_administrators_group": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"organizational_unit_distinguished_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iscsi": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"management": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"nfs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"smb": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": DataSourceStorageVirtualMachineFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lifecycle_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_transition_reason": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subtype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tftags.TagsSchemaComputed(),
			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceONTAPStorageVirtualMachineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)
	defaultTagsConfig := meta.(*conns.AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*conns.AWSClient).IgnoreTagsConfig

	input := &fsx.DescribeStorageVirtualMachinesInput{}

	if id, ok := d.GetOk("id"); ok {
		input.StorageVirtualMachineIds = []*string{aws.String(id.(string))}
	}

	input.Filters = BuildStorageVirtualMachineFiltersDataSource(
		d.Get("filter").(*schema.Set),
	)

	if len(input.Filters) == 0 {
		input.Filters = nil
	}

	svm, err := findStorageVirtualMachine(ctx, conn, input, tfslices.PredicateTrue[*fsx.StorageVirtualMachine]())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading FSx ONTAP Storage Virtual Machine: %s", err)
	}

	d.SetId(aws.StringValue(svm.StorageVirtualMachineId))
	if err := d.Set("active_directory_configuration", flattenSvmActiveDirectoryConfiguration(d, svm.ActiveDirectoryConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting svm_active_directory: %s", err)
	}
	d.Set("arn", svm.ResourceARN)
	d.Set("creation_time", svm.CreationTime.Format(time.RFC3339))
	if err := d.Set("endpoints", flattenSvmEndpoints(svm.Endpoints)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting endpoints: %s", err)
	}
	d.Set("file_system_id", svm.FileSystemId)
	d.Set("lifecycle_status", svm.Lifecycle)
	if err := d.Set("lifecycle_transition_reason", flattenLifecycleTransitionReason(svm.LifecycleTransitionReason)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting lifecycle_transition_reason: %s", err)
	}
	d.Set("name", svm.Name)
	d.Set("subtype", svm.Subtype)
	d.Set("uuid", svm.UUID)

	tags := KeyValueTags(ctx, svm.Tags).IgnoreAWS().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting tags: %s", err)
	}

	return diags
}

func flattenLifecycleTransitionReason(rs *fsx.LifecycleTransitionReason) []interface{} {
	if rs == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if rs.Message != nil {
		m["message"] = aws.StringValue(rs.Message)
	}

	return []interface{}{m}
}
