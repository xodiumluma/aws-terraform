package lightsail

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_lightsail_domain_entry")
func ResourceDomainEntry() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceDomainEntryCreate,
		ReadWithoutTimeout:   resourceDomainEntryRead,
		DeleteWithoutTimeout: resourceDomainEntryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"is_alias": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"A",
					"CNAME",
					"MX",
					"NS",
					"SOA",
					"SRV",
					"TXT",
				}, false),
				ForceNew: true,
			},
		},
	}
}

func resourceDomainEntryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailConn()
	name := d.Get("name").(string)
	req := &lightsail.CreateDomainEntryInput{
		DomainName: aws.String(d.Get("domain_name").(string)),

		DomainEntry: &lightsail.DomainEntry{
			IsAlias: aws.Bool(d.Get("is_alias").(bool)),
			Name:    aws.String(expandDomainEntryName(name, d.Get("domain_name").(string))),
			Target:  aws.String(d.Get("target").(string)),
			Type:    aws.String(d.Get("type").(string)),
		},
	}

	resp, err := conn.CreateDomainEntryWithContext(ctx, req)

	if err != nil {
		return create.DiagError(names.Lightsail, lightsail.OperationTypeCreateDomain, ResDomainEntry, name, err)
	}

	diag := expandOperations(ctx, conn, []*lightsail.Operation{resp.Operation}, lightsail.OperationTypeCreateDomain, ResDomainEntry, name)

	if diag != nil {
		return diag
	}

	// Generate an ID
	vars := []string{
		name,
		d.Get("domain_name").(string),
		d.Get("type").(string),
		d.Get("target").(string),
	}

	d.SetId(strings.Join(vars, "_"))

	return resourceDomainEntryRead(ctx, d, meta)
}

func resourceDomainEntryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailConn()

	entry, err := FindDomainEntryById(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		create.LogNotFoundRemoveState(names.Lightsail, create.ErrActionReading, ResDomainEntry, d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.Lightsail, create.ErrActionReading, ResDomainEntry, d.Id(), err)
	}

	domainName := expandDomainNameFromId(d.Id())

	d.Set("name", flattenDomainEntryName(aws.StringValue(entry.Name), domainName))
	d.Set("domain_name", domainName)
	d.Set("type", entry.Type)
	d.Set("is_alias", entry.IsAlias)
	d.Set("target", entry.Target)

	return nil
}

func resourceDomainEntryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).LightsailConn()

	resp, err := conn.DeleteDomainEntryWithContext(ctx, &lightsail.DeleteDomainEntryInput{
		DomainName:  aws.String(expandDomainNameFromId(d.Id())),
		DomainEntry: expandDomainEntry(d.Id()),
	})

	if err != nil && tfawserr.ErrCodeEquals(err, lightsail.ErrCodeNotFoundException) {
		return nil
	}

	if err != nil {
		return create.DiagError(names.Lightsail, create.ErrActionDeleting, ResDomainEntry, d.Id(), err)
	}

	diag := expandOperations(ctx, conn, []*lightsail.Operation{resp.Operation}, lightsail.OperationTypeDeleteDomain, ResDomainEntry, d.Id())

	if diag != nil {
		return diag
	}

	return nil
}

func expandDomainEntry(id string) *lightsail.DomainEntry {
	id_parts := strings.Split(id, "_")
	idLength := len(id_parts)
	var index int
	var name string

	if idLength == 5 {
		index = 1
		name = "_" + id_parts[index+0]
	} else {
		index = 0
		name = id_parts[index+0]
	}

	domainName := id_parts[index+1]
	recordType := id_parts[index+2]
	recordTarget := id_parts[index+3]

	entry := &lightsail.DomainEntry{
		Name:   aws.String(expandDomainEntryName(name, domainName)),
		Type:   aws.String(recordType),
		Target: aws.String(recordTarget),
	}

	return entry
}

func expandDomainNameFromId(id string) string {
	id_parts := strings.Split(id, "_")
	idLength := len(id_parts)
	var index int

	if idLength == 5 {
		index = 1
	} else {
		index = 0
	}

	domainName := id_parts[index+1]

	return domainName
}

func expandDomainEntryName(name, domainName string) string {
	rn := strings.ToLower(strings.TrimSuffix(name, "."))
	domainName = strings.TrimSuffix(domainName, ".")
	if !strings.HasSuffix(rn, domainName) {
		if len(name) == 0 {
			rn = domainName
		} else {
			rn = strings.Join([]string{rn, domainName}, ".")
		}
	}
	return rn
}

func flattenDomainEntryName(name, domainName string) string {
	rn := strings.ToLower(strings.TrimSuffix(name, "."))
	domainName = strings.TrimSuffix(domainName, ".")
	if strings.HasSuffix(rn, domainName) {
		rn = strings.TrimSuffix(rn, fmt.Sprintf(".%s", domainName))
	}
	return rn
}
