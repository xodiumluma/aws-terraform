// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package fsx

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/fsx"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
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

// @SDKResource("aws_fsx_windows_file_system", name="Windows File System")
// @Tags(identifierAttribute="arn")
func ResourceWindowsFileSystem() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceWindowsFileSystemCreate,
		ReadWithoutTimeout:   resourceWindowsFileSystemRead,
		UpdateWithoutTimeout: resourceWindowsFileSystemUpdate,
		DeleteWithoutTimeout: resourceWindowsFileSystemDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("skip_final_backup", false)

				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"active_directory_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"self_managed_active_directory"},
			},
			"aliases": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 50,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.All(
						validation.StringLenBetween(4, 253),
						// validation.StringMatch(regexache.MustCompile(`^[0-9A-Za-z]([.][0-9A-Za-z][0-9A-Za-z-]*[0-9A-Za-z])+$`), "must be in the fqdn format hostname.domain"),
					),
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"audit_log_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audit_log_destination": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidARN,
							StateFunc:    windowsAuditLogStateFunc,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return strings.HasPrefix(old, fmt.Sprintf("%s:", new))
							},
						},
						"file_access_audit_log_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      fsx.WindowsAccessAuditLogLevelDisabled,
							ValidateFunc: validation.StringInSlice(fsx.WindowsAccessAuditLogLevel_Values(), false),
						},
						"file_share_access_audit_log_level": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      fsx.WindowsAccessAuditLogLevelDisabled,
							ValidateFunc: validation.StringInSlice(fsx.WindowsAccessAuditLogLevel_Values(), false),
						},
					},
				},
			},
			"automatic_backup_retention_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(0, 90),
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"copy_tags_to_backups": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"daily_automatic_backup_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(5, 5),
					validation.StringMatch(regexache.MustCompile(`^([01]\d|2[0-3]):?([0-5]\d)$`), "must be in the format HH:MM"),
				),
			},
			"deployment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      fsx.WindowsDeploymentTypeSingleAz1,
				ValidateFunc: validation.StringInSlice(fsx.WindowsDeploymentType_Values(), false),
			},
			"disk_iops_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iops": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 350000),
						},
						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      fsx.DiskIopsConfigurationModeAutomatic,
							ValidateFunc: validation.StringInSlice(fsx.DiskIopsConfigurationMode_Values(), false),
						},
					},
				},
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidARN,
			},
			"network_interface_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_file_server_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"preferred_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"remote_administration_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 50,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"self_managed_active_directory": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"active_directory_id"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_ips": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 2,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.IsIPAddress,
							},
						},
						"domain_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"file_system_administrators_group": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "Domain Admins",
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"organizational_unit_distinguished_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 2000),
						},
						"password": {
							Type:         schema.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
						"username": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 256),
						},
					},
				},
			},
			"skip_final_backup": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"storage_capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(32, 65536),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      fsx.StorageTypeSsd,
				ValidateFunc: validation.StringInSlice(fsx.StorageType_Values(), false),
			},
			"subnet_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
			"throughput_capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(8, 2048),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"weekly_maintenance_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(7, 7),
					validation.StringMatch(regexache.MustCompile(`^[1-7]:([01]\d|2[0-3]):?([0-5]\d)$`), "must be in the format d:HH:MM"),
				),
			},
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceWindowsFileSystemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	inputC := &fsx.CreateFileSystemInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		FileSystemType:     aws.String(fsx.FileSystemTypeWindows),
		StorageCapacity:    aws.Int64(int64(d.Get("storage_capacity").(int))),
		SubnetIds:          flex.ExpandStringList(d.Get("subnet_ids").([]interface{})),
		Tags:               getTagsIn(ctx),
		WindowsConfiguration: &fsx.CreateFileSystemWindowsConfiguration{
			AutomaticBackupRetentionDays: aws.Int64(int64(d.Get("automatic_backup_retention_days").(int))),
			CopyTagsToBackups:            aws.Bool(d.Get("copy_tags_to_backups").(bool)),
			ThroughputCapacity:           aws.Int64(int64(d.Get("throughput_capacity").(int))),
		},
	}
	inputB := &fsx.CreateFileSystemFromBackupInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		SubnetIds:          flex.ExpandStringList(d.Get("subnet_ids").([]interface{})),
		Tags:               getTagsIn(ctx),
		WindowsConfiguration: &fsx.CreateFileSystemWindowsConfiguration{
			AutomaticBackupRetentionDays: aws.Int64(int64(d.Get("automatic_backup_retention_days").(int))),
			CopyTagsToBackups:            aws.Bool(d.Get("copy_tags_to_backups").(bool)),
			ThroughputCapacity:           aws.Int64(int64(d.Get("throughput_capacity").(int))),
		},
	}

	if v, ok := d.GetOk("active_directory_id"); ok {
		inputC.WindowsConfiguration.ActiveDirectoryId = aws.String(v.(string))
		inputB.WindowsConfiguration.ActiveDirectoryId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("aliases"); ok {
		inputC.WindowsConfiguration.Aliases = flex.ExpandStringSet(v.(*schema.Set))
		inputB.WindowsConfiguration.Aliases = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("audit_log_configuration"); ok && len(v.([]interface{})) > 0 {
		inputC.WindowsConfiguration.AuditLogConfiguration = expandWindowsAuditLogCreateConfiguration(v.([]interface{}))
		inputB.WindowsConfiguration.AuditLogConfiguration = expandWindowsAuditLogCreateConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("daily_automatic_backup_start_time"); ok {
		inputC.WindowsConfiguration.DailyAutomaticBackupStartTime = aws.String(v.(string))
		inputB.WindowsConfiguration.DailyAutomaticBackupStartTime = aws.String(v.(string))
	}

	if v, ok := d.GetOk("disk_iops_configuration"); ok && len(v.([]interface{})) > 0 {
		inputC.WindowsConfiguration.DiskIopsConfiguration = expandWindowsDiskIopsConfiguration(v.([]interface{}))
		inputB.WindowsConfiguration.DiskIopsConfiguration = expandWindowsDiskIopsConfiguration(v.([]interface{}))
	}

	if v, ok := d.GetOk("deployment_type"); ok {
		inputC.WindowsConfiguration.DeploymentType = aws.String(v.(string))
		inputB.WindowsConfiguration.DeploymentType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("kms_key_id"); ok {
		inputC.KmsKeyId = aws.String(v.(string))
		inputB.KmsKeyId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("preferred_subnet_id"); ok {
		inputC.WindowsConfiguration.PreferredSubnetId = aws.String(v.(string))
		inputB.WindowsConfiguration.PreferredSubnetId = aws.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_ids"); ok {
		inputC.SecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
		inputB.SecurityGroupIds = flex.ExpandStringSet(v.(*schema.Set))
	}

	if v, ok := d.GetOk("self_managed_active_directory"); ok {
		inputC.WindowsConfiguration.SelfManagedActiveDirectoryConfiguration = expandSelfManagedActiveDirectoryConfigurationCreate(v.([]interface{}))
		inputB.WindowsConfiguration.SelfManagedActiveDirectoryConfiguration = expandSelfManagedActiveDirectoryConfigurationCreate(v.([]interface{}))
	}

	if v, ok := d.GetOk("storage_type"); ok {
		inputC.StorageType = aws.String(v.(string))
		inputB.StorageType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("weekly_maintenance_start_time"); ok {
		inputC.WindowsConfiguration.WeeklyMaintenanceStartTime = aws.String(v.(string))
		inputB.WindowsConfiguration.WeeklyMaintenanceStartTime = aws.String(v.(string))
	}

	if v, ok := d.GetOk("backup_id"); ok {
		backupID := v.(string)
		inputB.BackupId = aws.String(backupID)

		output, err := conn.CreateFileSystemFromBackupWithContext(ctx, inputB)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating FSx for Windows File Server File System from backup (%s): %s", backupID, err)
		}

		d.SetId(aws.StringValue(output.FileSystem.FileSystemId))
	} else {
		output, err := conn.CreateFileSystemWithContext(ctx, inputC)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "creating FSx for Windows File Server File System: %s", err)
		}

		d.SetId(aws.StringValue(output.FileSystem.FileSystemId))
	}

	if _, err := waitFileSystemCreated(ctx, conn, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for FSx for Windows File Server File System (%s) create: %s", d.Id(), err)
	}

	return append(diags, resourceWindowsFileSystemRead(ctx, d, meta)...)
}

func resourceWindowsFileSystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	filesystem, err := FindWindowsFileSystemByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] FSx for Windows File Server File System (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading FSx for Windows File Server File System (%s): %s", d.Id(), err)
	}

	windowsConfig := filesystem.WindowsConfiguration

	d.Set("active_directory_id", windowsConfig.ActiveDirectoryId)
	d.Set("aliases", aws.StringValueSlice(expandAliasValues(windowsConfig.Aliases)))
	d.Set("arn", filesystem.ResourceARN)
	if err := d.Set("audit_log_configuration", flattenWindowsAuditLogConfiguration(windowsConfig.AuditLogConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting audit_log_configuration: %s", err)
	}
	d.Set("automatic_backup_retention_days", windowsConfig.AutomaticBackupRetentionDays)
	d.Set("copy_tags_to_backups", windowsConfig.CopyTagsToBackups)
	d.Set("daily_automatic_backup_start_time", windowsConfig.DailyAutomaticBackupStartTime)
	d.Set("deployment_type", windowsConfig.DeploymentType)
	if err := d.Set("disk_iops_configuration", flattenWindowsDiskIopsConfiguration(windowsConfig.DiskIopsConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting disk_iops_configuration: %s", err)
	}
	d.Set("dns_name", filesystem.DNSName)
	d.Set("kms_key_id", filesystem.KmsKeyId)
	d.Set("network_interface_ids", aws.StringValueSlice(filesystem.NetworkInterfaceIds))
	d.Set("owner_id", filesystem.OwnerId)
	d.Set("preferred_file_server_ip", windowsConfig.PreferredFileServerIp)
	d.Set("preferred_subnet_id", windowsConfig.PreferredSubnetId)
	d.Set("remote_administration_endpoint", windowsConfig.RemoteAdministrationEndpoint)
	if err := d.Set("self_managed_active_directory", flattenSelfManagedActiveDirectoryConfiguration(d, windowsConfig.SelfManagedActiveDirectoryConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting self_managed_active_directory: %s", err)
	}
	d.Set("storage_capacity", filesystem.StorageCapacity)
	d.Set("storage_type", filesystem.StorageType)
	d.Set("subnet_ids", aws.StringValueSlice(filesystem.SubnetIds))
	d.Set("throughput_capacity", windowsConfig.ThroughputCapacity)
	d.Set("vpc_id", filesystem.VpcId)
	d.Set("weekly_maintenance_start_time", windowsConfig.WeeklyMaintenanceStartTime)

	setTagsOut(ctx, filesystem.Tags)

	return diags
}

func resourceWindowsFileSystemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	if d.HasChange("aliases") {
		o, n := d.GetChange("aliases")
		os, ns := o.(*schema.Set), n.(*schema.Set)
		add, del := flex.ExpandStringValueSet(ns.Difference(os)), flex.ExpandStringValueSet(os.Difference(ns))

		if len(add) > 0 {
			input := &fsx.AssociateFileSystemAliasesInput{
				Aliases:      aws.StringSlice(add),
				FileSystemId: aws.String(d.Id()),
			}

			_, err := conn.AssociateFileSystemAliasesWithContext(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "associating FSx for Windows File Server File System (%s) aliases: %s", d.Id(), err)
			}

			if _, err := waitFileSystemAdministrativeActionCompleted(ctx, conn, d.Id(), fsx.AdministrativeActionTypeFileSystemAliasAssociation, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for FSx for Windows File Server File System (%s) administrative action (%s) complete: %s", d.Id(), fsx.AdministrativeActionTypeFileSystemAliasAssociation, err)
			}
		}

		if len(del) > 0 {
			input := &fsx.DisassociateFileSystemAliasesInput{
				Aliases:      aws.StringSlice(del),
				FileSystemId: aws.String(d.Id()),
			}

			_, err := conn.DisassociateFileSystemAliasesWithContext(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "disassociating FSx for Windows File Server File System (%s) aliases: %s", d.Id(), err)
			}

			if _, err := waitFileSystemAdministrativeActionCompleted(ctx, conn, d.Id(), fsx.AdministrativeActionTypeFileSystemAliasDisassociation, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for FSx for Windows File Server File System (%s) administrative action (%s) complete: %s", d.Id(), fsx.AdministrativeActionTypeFileSystemAliasDisassociation, err)
			}
		}
	}

	// Increase ThroughputCapacity first to avoid errors like
	// "BadRequest: Unable to perform the storage capacity update. Updating storage capacity requires your file system to have at least 16 MB/s of throughput capacity."
	if d.HasChange("throughput_capacity") {
		o, n := d.GetChange("throughput_capacity")
		if o, n := o.(int), n.(int); n > o {
			input := &fsx.UpdateFileSystemInput{
				ClientRequestToken: aws.String(id.UniqueId()),
				FileSystemId:       aws.String(d.Id()),
				WindowsConfiguration: &fsx.UpdateFileSystemWindowsConfiguration{
					ThroughputCapacity: aws.Int64(int64(n)),
				},
			}

			startTime := time.Now()
			_, err := conn.UpdateFileSystemWithContext(ctx, input)

			if err != nil {
				return sdkdiag.AppendErrorf(diags, "updating FSx for Windows File Server File System (%s) ThroughputCapacity: %s", d.Id(), err)
			}

			if _, err := waitFileSystemUpdated(ctx, conn, d.Id(), startTime, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for FSx Windows File Server File System (%s) update: %s", d.Id(), err)
			}

			if _, err := waitFileSystemAdministrativeActionCompleted(ctx, conn, d.Id(), fsx.AdministrativeActionTypeFileSystemUpdate, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return sdkdiag.AppendErrorf(diags, "waiting for FSx Windows File Server File System (%s) administrative action (%s) complete: %s", d.Id(), fsx.AdministrativeActionTypeFileSystemUpdate, err)
			}
		}
	}

	if d.HasChangesExcept("aliases", "tags", "tags_all") {
		input := &fsx.UpdateFileSystemInput{
			ClientRequestToken:   aws.String(id.UniqueId()),
			FileSystemId:         aws.String(d.Id()),
			WindowsConfiguration: &fsx.UpdateFileSystemWindowsConfiguration{},
		}

		if d.HasChange("audit_log_configuration") {
			input.WindowsConfiguration.AuditLogConfiguration = expandWindowsAuditLogCreateConfiguration(d.Get("audit_log_configuration").([]interface{}))
		}

		if d.HasChange("automatic_backup_retention_days") {
			input.WindowsConfiguration.AutomaticBackupRetentionDays = aws.Int64(int64(d.Get("automatic_backup_retention_days").(int)))
		}

		if d.HasChange("daily_automatic_backup_start_time") {
			input.WindowsConfiguration.DailyAutomaticBackupStartTime = aws.String(d.Get("daily_automatic_backup_start_time").(string))
		}

		if d.HasChange("disk_iops_configuration") {
			input.WindowsConfiguration.DiskIopsConfiguration = expandWindowsDiskIopsConfiguration(d.Get("disk_iops_configuration").([]interface{}))
		}

		if d.HasChange("self_managed_active_directory") {
			input.WindowsConfiguration.SelfManagedActiveDirectoryConfiguration = expandSelfManagedActiveDirectoryConfigurationUpdate(d.Get("self_managed_active_directory").([]interface{}))
		}

		if d.HasChange("storage_capacity") {
			input.StorageCapacity = aws.Int64(int64(d.Get("storage_capacity").(int)))
		}

		if d.HasChange("throughput_capacity") {
			input.WindowsConfiguration.ThroughputCapacity = aws.Int64(int64(d.Get("throughput_capacity").(int)))
		}

		if d.HasChange("weekly_maintenance_start_time") {
			input.WindowsConfiguration.WeeklyMaintenanceStartTime = aws.String(d.Get("weekly_maintenance_start_time").(string))
		}

		startTime := time.Now()
		_, err := conn.UpdateFileSystemWithContext(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating FSx for Windows File Server File System (%s): %s", d.Id(), err)
		}

		if _, err := waitFileSystemUpdated(ctx, conn, d.Id(), startTime, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for FSx Windows File Server File System (%s) update: %s", d.Id(), err)
		}

		if _, err := waitFileSystemAdministrativeActionCompleted(ctx, conn, d.Id(), fsx.AdministrativeActionTypeFileSystemUpdate, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return sdkdiag.AppendErrorf(diags, "waiting for FSx Windows File Server File System (%s) administrative action (%s) complete: %s", d.Id(), fsx.AdministrativeActionTypeFileSystemUpdate, err)
		}
	}

	return append(diags, resourceWindowsFileSystemRead(ctx, d, meta)...)
}

func resourceWindowsFileSystemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).FSxConn(ctx)

	input := &fsx.DeleteFileSystemInput{
		ClientRequestToken: aws.String(id.UniqueId()),
		FileSystemId:       aws.String(d.Id()),
		WindowsConfiguration: &fsx.DeleteFileSystemWindowsConfiguration{
			SkipFinalBackup: aws.Bool(d.Get("skip_final_backup").(bool)),
		},
	}

	log.Printf("[DEBUG] Deleting FSx for Windows File Server File System: %s", d.Id())
	_, err := conn.DeleteFileSystemWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, fsx.ErrCodeFileSystemNotFound) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting FSx for Windows File Server File System (%s): %s", d.Id(), err)
	}

	if _, err := waitFileSystemDeleted(ctx, conn, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return sdkdiag.AppendErrorf(diags, "waiting for FSx for Windows File Server File System (%s) delete: %s", d.Id(), err)
	}

	return diags
}

func expandAliasValues(aliases []*fsx.Alias) []*string {
	var alternateDNSNames []*string

	for _, alias := range aliases {
		aName := alias.Name
		alternateDNSNames = append(alternateDNSNames, aName)
	}

	return alternateDNSNames
}

func expandSelfManagedActiveDirectoryConfigurationCreate(l []interface{}) *fsx.SelfManagedActiveDirectoryConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	data := l[0].(map[string]interface{})
	req := &fsx.SelfManagedActiveDirectoryConfiguration{
		DomainName: aws.String(data["domain_name"].(string)),
		DnsIps:     flex.ExpandStringSet(data["dns_ips"].(*schema.Set)),
		Password:   aws.String(data["password"].(string)),
		UserName:   aws.String(data["username"].(string)),
	}

	if v, ok := data["file_system_administrators_group"]; ok && v.(string) != "" {
		req.FileSystemAdministratorsGroup = aws.String(v.(string))
	}

	if v, ok := data["organizational_unit_distinguished_name"]; ok && v.(string) != "" {
		req.OrganizationalUnitDistinguishedName = aws.String(v.(string))
	}

	return req
}

func expandSelfManagedActiveDirectoryConfigurationUpdate(l []interface{}) *fsx.SelfManagedActiveDirectoryConfigurationUpdates {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	data := l[0].(map[string]interface{})
	req := &fsx.SelfManagedActiveDirectoryConfigurationUpdates{}

	if v, ok := data["dns_ips"].(*schema.Set); ok && v.Len() > 0 {
		req.DnsIps = flex.ExpandStringSet(v)
	}

	if v, ok := data["password"].(string); ok && v != "" {
		req.Password = aws.String(v)
	}

	if v, ok := data["username"].(string); ok && v != "" {
		req.UserName = aws.String(v)
	}

	return req
}

func flattenSelfManagedActiveDirectoryConfiguration(d *schema.ResourceData, adopts *fsx.SelfManagedActiveDirectoryAttributes) []map[string]interface{} {
	if adopts == nil {
		return []map[string]interface{}{}
	}

	// Since we are in a configuration block and the FSx API does not return
	// the password, we need to set the value if we can or Terraform will
	// show a difference for the argument from empty string to the value.
	// This is not a pattern that should be used normally.
	// See also: flattenEmrKerberosAttributes

	m := map[string]interface{}{
		"dns_ips":                                aws.StringValueSlice(adopts.DnsIps),
		"domain_name":                            aws.StringValue(adopts.DomainName),
		"file_system_administrators_group":       aws.StringValue(adopts.FileSystemAdministratorsGroup),
		"organizational_unit_distinguished_name": aws.StringValue(adopts.OrganizationalUnitDistinguishedName),
		"password":                               d.Get("self_managed_active_directory.0.password").(string),
		"username":                               aws.StringValue(adopts.UserName),
	}

	return []map[string]interface{}{m}
}

func expandWindowsAuditLogCreateConfiguration(l []interface{}) *fsx.WindowsAuditLogCreateConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	data := l[0].(map[string]interface{})
	req := &fsx.WindowsAuditLogCreateConfiguration{
		FileAccessAuditLogLevel:      aws.String(data["file_access_audit_log_level"].(string)),
		FileShareAccessAuditLogLevel: aws.String(data["file_share_access_audit_log_level"].(string)),
	}

	if v, ok := data["audit_log_destination"].(string); ok && v != "" {
		req.AuditLogDestination = aws.String(windowsAuditLogStateFunc(v))
	}

	return req
}

func flattenWindowsAuditLogConfiguration(adopts *fsx.WindowsAuditLogConfiguration) []map[string]interface{} {
	if adopts == nil {
		return []map[string]interface{}{}
	}

	m := map[string]interface{}{
		"file_access_audit_log_level":       aws.StringValue(adopts.FileAccessAuditLogLevel),
		"file_share_access_audit_log_level": aws.StringValue(adopts.FileShareAccessAuditLogLevel),
	}

	if adopts.AuditLogDestination != nil {
		m["audit_log_destination"] = aws.StringValue(adopts.AuditLogDestination)
	}

	return []map[string]interface{}{m}
}

func expandWindowsDiskIopsConfiguration(l []interface{}) *fsx.DiskIopsConfiguration {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	data := l[0].(map[string]interface{})
	req := &fsx.DiskIopsConfiguration{}

	if v, ok := data["iops"].(int); ok {
		req.Iops = aws.Int64(int64(v))
	}

	if v, ok := data["mode"].(string); ok && v != "" {
		req.Mode = aws.String(v)
	}

	return req
}

func flattenWindowsDiskIopsConfiguration(rs *fsx.DiskIopsConfiguration) []interface{} {
	if rs == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{}

	if rs.Iops != nil {
		m["iops"] = aws.Int64Value(rs.Iops)
	}
	if rs.Mode != nil {
		m["mode"] = aws.StringValue(rs.Mode)
	}

	return []interface{}{m}
}

func windowsAuditLogStateFunc(v interface{}) string {
	value := v.(string)
	// API returns the specific log stream arn instead of provided log group
	logArn, _ := arn.Parse(value)
	if logArn.Service == "logs" {
		parts := strings.SplitN(logArn.Resource, ":", 3)
		if len(parts) == 3 {
			return strings.TrimSuffix(value, fmt.Sprintf(":%s", parts[2]))
		} else {
			return value
		}
	}
	return value
}

func FindWindowsFileSystemByID(ctx context.Context, conn *fsx.FSx, id string) (*fsx.FileSystem, error) {
	output, err := findFileSystemByIDAndType(ctx, conn, id, fsx.FileSystemTypeWindows)

	if err != nil {
		return nil, err
	}

	if output.WindowsConfiguration == nil {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	return output, nil
}
