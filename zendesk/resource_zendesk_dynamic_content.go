package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

var reversedLocaleTypes map[string]int64

func init() {
	reversedLocaleTypes = reverseLocaleTypes()
}

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/dynamic_content/
func resourceZendeskDynamicContentItem() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a dynamic content item resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(client.DynamicContentAPI)
			return createDynamicContentItem(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.DynamicContentAPI)
			return readDynamicContentItem(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.DynamicContentAPI)
			return updateDynamicContentItem(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.DynamicContentAPI)
			return deleteDynamicContentItem(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The unique name of the item.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"default_locale": {
				Description: "The default locale for the item. Must be one of the [locales the account has active](https://developer.zendesk.com/api-reference/ticketing/account-configuration/locales/#list-locales).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"variant": {
				Description: "Variant within this item.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Description: "If the variant is active and useable.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
						},
						"content": {
							Description: "The content of the variant.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"locale": {
							Description: "The locale of the variant.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Required: true,
			},
		},
	}
}

func marshalDynamicContentItem(item client.DynamicContentItem, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"name":           item.Name,
		"default_locale": item.DefaultLocaleID,
	}

	var variants []map[string]interface{}
	for _, v := range item.Variants {
		m := map[string]interface{}{
			"active":  v.Active,
			"content": v.Content,
			"locale":  v.LocaleID,
		}
		variants = append(variants, m)
	}
	fields["variant"] = variants

	return setSchemaFields(d, fields)
}

func unmarshalDynamicContentItem(d identifiableGetterSetter) (client.DynamicContentItem, error) {
	item := client.DynamicContentItem{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return item, fmt.Errorf("could not parse dynamic content item id %s: %v", v, err)
		}
		item.ID = id
	}

	if v, ok := d.GetOk("name"); ok {
		item.Name = v.(string)
	}

	if v, ok := d.GetOk("default_locale"); ok {
		item.DefaultLocaleID = reversedLocaleTypes[v.(string)]
	}

	if v, ok := d.GetOk("variant"); ok {
		tfstateVariants := v.(*schema.Set).List()
		variants := []client.DynamicContentVariant{}

		for _, vari := range tfstateVariants {
			variant, ok := vari.(map[string]interface{})
			if !ok {
				return item, fmt.Errorf("could not parse 'variant' for dynamic content item %v", item)
			}
			variants = append(variants, client.DynamicContentVariant{
				Active:   variant["active"].(bool),
				Content:  variant["content"].(string),
				LocaleID: reversedLocaleTypes[variant["locale"].(string)],
			})
		}
		item.Variants = variants
	}

	return item, nil
}

func createDynamicContentItem(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	item, err := unmarshalDynamicContentItem(d)
	if err != nil {
		return diag.FromErr(err)
	}

	item, err = zd.CreateDynamicContentItem(ctx, item)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", item.ID))

	err = marshalDynamicContentItem(item, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readDynamicContentItem(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	item, err := zd.GetDynamicContentItem(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalDynamicContentItem(item, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateDynamicContentItem(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	item, err := unmarshalDynamicContentItem(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	item, err = zd.UpdateDynamicContentItem(ctx, id, item)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalDynamicContentItem(item, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteDynamicContentItem(ctx context.Context, d identifiableGetterSetter, zd client.DynamicContentAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteDynamicContentItem(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func reverseLocaleTypes() map[string]int64 {
	rlocs := map[string]int64{}
	for i := client.LocaleENUS; i <= client.LocaleENPH; i++ {
		loctxt := client.LocaleTypeText(i)
		if loctxt != "" {
			rlocs[loctxt] = int64(i)
		}
	}
	return rlocs
}
