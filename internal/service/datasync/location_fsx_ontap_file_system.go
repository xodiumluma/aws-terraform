// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasync

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/datasync"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_datasync_location_fsx_ontap_file_system", name="Location FSx for NetApp ONTAP File System")
// @Tags(identifierAttribute="id")
func ResourceLocationFSxONTAPFileSystem() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceLocationFSxONTAPFileSystemCreate,
		ReadWithoutTimeout:   resourceLocationFSxONTAPFileSystemRead,
		UpdateWithoutTimeout: resourceLocationFSxONTAPFileSystemUpdate,
		DeleteWithoutTimeout: resourceLocationFSxONTAPFileSystemDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), "#")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected DataSyncLocationArn#FsxSVMArn", d.Id())
				}

				DSArn := idParts[0]
				FSxArn := idParts[1]

				d.Set("fsx_filesystem_arn", FSxArn)
				d.SetId(DSArn)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fsx_filesystem_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nfs": {
							Type:         schema.TypeList,
							Optional:     true,
							ForceNew:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"protocol.0.nfs", "protocol.0.smb"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_options": {
										Type:     schema.TypeList,
										Required: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:     schema.TypeString,
													Default:  datasync.NfsVersionNfs3,
													Optional: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														datasync.NfsVersionNfs3,
													}, false),
												},
											},
										},
									},
								},
							},
						},
						"smb": {
							Type:         schema.TypeList,
							Optional:     true,
							ForceNew:     true,
							MaxItems:     1,
							ExactlyOneOf: []string{"protocol.0.nfs", "protocol.0.smb"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain": {
										Type:         schema.TypeString,
										Optional:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringLenBetween(1, 253),
									},
									"mount_options": {
										Type:     schema.TypeList,
										Required: true,
										ForceNew: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"version": {
													Type:     schema.TypeString,
													Default:  datasync.SmbVersionAutomatic,
													Optional: true,
													ForceNew: true,
													ValidateFunc: validation.StringInSlice([]string{
														datasync.SmbVersionAutomatic,
														datasync.SmbVersionSmb2,
														datasync.SmbVersionSmb3,
														datasync.SmbVersionSmb20,
													}, false),
												},
											},
										},
									},
									"password": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										Sensitive:    true,
										ValidateFunc: validation.StringLenBetween(1, 104),
									},
									"user": {
										Type:         schema.TypeString,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.StringLenBetween(1, 104),
									},
								},
							},
						},
					},
				},
			},
			"security_group_arns": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidARN,
				},
			},
			"storage_virtual_machine_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"subdirectory": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceLocationFSxONTAPFileSystemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.CreateLocationFsxOntapInput{
		Protocol:                 expandProtocol(d.Get("protocol").([]interface{})),
		SecurityGroupArns:        flex.ExpandStringSet(d.Get("security_group_arns").(*schema.Set)),
		StorageVirtualMachineArn: aws.String(d.Get("storage_virtual_machine_arn").(string)),
		Tags:                     getTagsIn(ctx),
	}

	if v, ok := d.GetOk("subdirectory"); ok {
		input.Subdirectory = aws.String(v.(string))
	}

	output, err := conn.CreateLocationFsxOntapWithContext(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating DataSync Location FSx for NetApp ONTAP File System: %s", err)
	}

	d.SetId(aws.StringValue(output.LocationArn))

	return append(diags, resourceLocationFSxONTAPFileSystemRead(ctx, d, meta)...)
}

func resourceLocationFSxONTAPFileSystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	output, err := FindLocationFSxONTAPByARN(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] DataSync Location FSx for NetApp ONTAP File System (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading DataSync Location FSx for NetApp ONTAP File System (%s): %s", d.Id(), err)
	}

	uri := aws.StringValue(output.LocationUri)
	subdirectory, err := subdirectoryFromLocationURI(uri)
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	d.Set("arn", output.LocationArn)
	d.Set("creation_time", output.CreationTime.Format(time.RFC3339))
	d.Set("fsx_filesystem_arn", output.FsxFilesystemArn)
	// SMB Password is not returned from the API.
	if output.Protocol != nil && output.Protocol.SMB != nil && aws.StringValue(output.Protocol.SMB.Password) == "" {
		if smbPassword := d.Get("protocol.0.smb.0.password").(string); smbPassword != "" {
			output.Protocol.SMB.Password = aws.String(smbPassword)
		}
	}
	if err := d.Set("protocol", flattenProtocol(output.Protocol)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting protocol: %s", err)
	}
	d.Set("security_group_arns", aws.StringValueSlice(output.SecurityGroupArns))
	d.Set("storage_virtual_machine_arn", output.StorageVirtualMachineArn)
	d.Set("subdirectory", subdirectory)
	d.Set("uri", uri)

	return diags
}

func resourceLocationFSxONTAPFileSystemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Tags only.

	return append(diags, resourceLocationFSxONTAPFileSystemRead(ctx, d, meta)...)
}

func resourceLocationFSxONTAPFileSystemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).DataSyncConn(ctx)

	input := &datasync.DeleteLocationInput{
		LocationArn: aws.String(d.Id()),
	}

	log.Printf("[DEBUG] Deleting DataSync Location FSx for NetApp ONTAP File System: %s", d.Id())
	_, err := conn.DeleteLocationWithContext(ctx, input)

	if tfawserr.ErrMessageContains(err, datasync.ErrCodeInvalidRequestException, "not found") {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting DataSync Location FSx for NetApp ONTAP File System (%s): %s", d.Id(), err)
	}

	return diags
}

func FindLocationFSxONTAPByARN(ctx context.Context, conn *datasync.DataSync, arn string) (*datasync.DescribeLocationFsxOntapOutput, error) {
	input := &datasync.DescribeLocationFsxOntapInput{
		LocationArn: aws.String(arn),
	}

	output, err := conn.DescribeLocationFsxOntapWithContext(ctx, input)

	if tfawserr.ErrMessageContains(err, datasync.ErrCodeInvalidRequestException, "not found") {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}
