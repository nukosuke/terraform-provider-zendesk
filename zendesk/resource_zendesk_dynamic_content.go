package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/api-reference/ticketing/ticket-management/dynamic_content/
func resourceZendeskDynamicContentItem() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(client.DynamicContentAPI)
			return createDynamicContentItem(d, zd)
		},
		Read: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.DynamicContentAPI)
			return readDynamicContentItem(d, zd)
		},
		Update: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.DynamicContentAPI)
			return updateDynamicContentItem(d, zd)
		},
		Delete: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.DynamicContentAPI)
			return deleteDynamicContentItem(d, zd)
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"default_locale": {
				// TODO: convert locale_id (int64) to locale string using zendesk.LocaleType
				Type:     schema.TypeInt,
				Required: true,
			},
			"variant": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
						// TODO: convert locale_id (int64) to locale string using zendesk.LocaleType
						"locale": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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
		// FIXME: see TODO above
		item.DefaultLocaleID = v.(int64)
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
				LocaleID: variant["locale"].(int64),
			})
		}
		item.Variants = variants
	}

	return item, nil
}

func createDynamicContentItem(d identifiableGetterSetter, zd client.DynamicContentAPI) error {
	item, err := unmarshalDynamicContentItem(d)
	if err != nil {
		return err
	}

	ctx := context.Background()
	item, err = zd.CreateDynamicContentItem(ctx, item)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", item.ID))
	return marshalDynamicContentItem(item, d)
}

func readDynamicContentItem(d identifiableGetterSetter, zd client.DynamicContentAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	item, err := zd.GetDynamicContentItem(ctx, id)
	if err != nil {
		return err
	}

	return marshalDynamicContentItem(item, d)
}

func updateDynamicContentItem(d identifiableGetterSetter, zd client.DynamicContentAPI) error {
	item, err := unmarshalDynamicContentItem(d)
	if err != nil {
		return err
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	item, err = zd.UpdateDynamicContentItem(ctx, id, item)
	if err != nil {
		return err
	}

	return marshalDynamicContentItem(item, d)
}

func deleteDynamicContentItem(d identifiableGetterSetter, zd client.DynamicContentAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteDynamicContentItem(ctx, id)
}
