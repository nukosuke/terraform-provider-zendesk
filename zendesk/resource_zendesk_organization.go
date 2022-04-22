package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/organizations
func resourceZendeskOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "[Beta] Provides an organization resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return createOrganization(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return readOrganization(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return updateOrganization(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return deleteOrganization(ctx, d, zd)
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Organization name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"domain_names": {
				Description: "A list of domain names associated with this organization.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"group_id": {
				Description: "New tickets from users in this organization are automatically put in this group.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"shared_tickets": {
				Description: "Whether end users in this organization are able to see each other's tickets.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"shared_comments": {
				Description: "End users in this organization are able to see each other's comments on tickets.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"tags": {
				Description: "The tags of the organization.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func marshalOrganization(org client.Organization, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":             org.URL,
		"name":            org.Name,
		"domain_names":    org.DomainNames,
		"group_id":        org.GroupID,
		"shared_tickets":  org.SharedTickets,
		"shared_comments": org.SharedComments,
		"tags":            org.Tags,
	}

	return setSchemaFields(d, fields)
}

func unmarshalOrganization(d identifiableGetterSetter) (client.Organization, error) {
	org := client.Organization{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return org, fmt.Errorf("could not parse organization id %s: %v", v, err)
		}
		org.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		org.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		org.Name = v.(string)
	}

	if v, ok := d.GetOk("domain_names"); ok {
		domainNames := v.(*schema.Set).List()
		for _, domainName := range domainNames {
			org.DomainNames = append(org.DomainNames, domainName.(string))
		}
	}

	if v, ok := d.GetOk("group_id"); ok {
		org.GroupID = int64(v.(int))
	}

	if v, ok := d.GetOk("shared_tickets"); ok {
		org.SharedTickets = v.(bool)
	}

	if v, ok := d.GetOk("shared_comments"); ok {
		org.SharedComments = v.(bool)
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		for _, tag := range tags {
			org.Tags = append(org.Tags, tag.(string))
		}
	}

	return org, nil
}

func createOrganization(ctx context.Context, d identifiableGetterSetter, zd client.OrganizationAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	org, err := unmarshalOrganization(d)
	if err != nil {
		return diag.FromErr(err)
	}

	org, err = zd.CreateOrganization(ctx, org)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", org.ID))

	err = marshalOrganization(org, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readOrganization(ctx context.Context, d identifiableGetterSetter, zd client.OrganizationAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	org, err := zd.GetOrganization(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOrganization(org, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateOrganization(ctx context.Context, d identifiableGetterSetter, zd client.OrganizationAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	org, err := unmarshalOrganization(d)
	if err != nil {
		return diag.FromErr(err)
	}

	org, err = zd.UpdateOrganization(ctx, id, org)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalOrganization(org, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteOrganization(ctx context.Context, d identifiableGetterSetter, zd client.OrganizationAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteOrganization(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
