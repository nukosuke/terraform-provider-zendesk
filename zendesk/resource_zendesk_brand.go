package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/api-reference/ticketing/account-configuration/brands/
func resourceZendeskBrand() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a brand resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return createBrand(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return readBrand(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return updateBrand(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return deleteBrand(ctx, d, zd)
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The API url of this brand.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "The name of the brand.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"brand_url": {
				Description: "The url of the brand.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"has_help_center": {
				Description: "If the brand has a Help Center.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"help_center_state": {
				Description: `The state of the Help Center. Allowed values are "enabled", "disabled", or "restricted".`,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active": {
				Description: "If the brand is set as active.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"default": {
				Description: "Is the brand the default brand for this account.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"logo_attachment_id": {
				Description: "Logo attachment id for the brand.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ticket_form_ids": {
				Description: "The ids of ticket forms that are available for use by a brand.",
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
			"subdomain": {
				Description: "The subdomain of the brand.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"host_mapping": {
				Description: "The hostmapping to this brand, if any. Only admins view this property.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"signature_template": {
				Description: "The signature template for a brand.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func marshalBrand(brand client.Brand, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":                brand.URL,
		"name":               brand.Name,
		"brand_url":          brand.BrandURL,
		"has_help_center":    brand.HasHelpCenter,
		"help_center_state":  brand.HelpCenterState,
		"active":             brand.Active,
		"default":            brand.Default,
		"logo_attachment_id": brand.Logo.ID,
		"ticket_form_ids":    brand.TicketFormIDs,
		"subdomain":          brand.Subdomain,
		"host_mapping":       brand.HostMapping,
		"signature_template": brand.SignatureTemplate,
	}

	return setSchemaFields(d, fields)
}

func unmarshalBrand(d identifiableGetterSetter) (client.Brand, error) {
	brand := client.Brand{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return brand, fmt.Errorf("could not parse brand id %s: %v", v, err)
		}
		brand.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		brand.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		brand.Name = v.(string)
	}

	if v, ok := d.GetOk("brand_url"); ok {
		brand.BrandURL = v.(string)
	}

	if v, ok := d.GetOk("has_help_center"); ok {
		brand.HasHelpCenter = v.(bool)
	}

	if v, ok := d.GetOk("help_center_state"); ok {
		brand.HelpCenterState = v.(string)
	}

	if v, ok := d.GetOk("active"); ok {
		brand.Active = v.(bool)
	}

	if v, ok := d.GetOk("default"); ok {
		brand.Default = v.(bool)
	}

	if v, ok := d.GetOk("logo_attachment_id"); ok {
		brand.Logo.ID = v.(int64)
	}

	if v, ok := d.GetOk("ticket_form_ids"); ok {
		ticketFormIDs := v.(*schema.Set).List()
		for _, ticketFormID := range ticketFormIDs {
			brand.TicketFormIDs = append(brand.TicketFormIDs, int64(ticketFormID.(int)))
		}
	}

	if v, ok := d.GetOk("subdomain"); ok {
		brand.Subdomain = v.(string)
	}

	if v, ok := d.GetOk("host_mapping"); ok {
		brand.HostMapping = v.(string)
	}

	if v, ok := d.GetOk("signature_template"); ok {
		brand.SignatureTemplate = v.(string)
	}

	return brand, nil
}

func createBrand(ctx context.Context, d identifiableGetterSetter, zd client.BrandAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	brand, err := unmarshalBrand(d)
	if err != nil {
		return diag.FromErr(err)
	}

	brand, err = zd.CreateBrand(ctx, brand)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", brand.ID))

	err = marshalBrand(brand, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readBrand(ctx context.Context, d identifiableGetterSetter, zd client.BrandAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	brand, err := zd.GetBrand(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalBrand(brand, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateBrand(ctx context.Context, d identifiableGetterSetter, zd client.BrandAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	brand, err := unmarshalBrand(d)
	if err != nil {
		return diag.FromErr(err)
	}

	brand, err = zd.UpdateBrand(ctx, id, brand)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalBrand(brand, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteBrand(ctx context.Context, d identifiableGetterSetter, zd client.BrandAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteBrand(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
