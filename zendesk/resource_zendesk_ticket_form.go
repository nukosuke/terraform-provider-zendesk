package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/ticket_forms
func resourceZendeskTicketForm() *schema.Resource {
	return &schema.Resource{
		Create: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(*client.Client)
			return createTicketForm(data, zd)
		},
		Read: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(*client.Client)
			return readTicketForm(data, zd)
		},
		Update: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(*client.Client)
			return updateTicketForm(data, zd)
		},
		Delete: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(*client.Client)
			return deleteTicketForm(data, zd)
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"end_user_visible": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ticket_field_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"in_all_brands": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: true,
			},
			"restricted_brand_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
		},
	}
}

// unmarshalTicketField parses the provided ResourceData and returns a ticket field
func unmarshalTicketForm(d identifiableGetterSetter) (client.TicketForm, error) {
	tf := client.TicketForm{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return tf, fmt.Errorf("could not parse ticket field id %s: %v", v, err)
		}
		tf.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		tf.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		tf.Name = v.(string)
		tf.RawName = v.(string)
	}

	if v, ok := d.GetOk("display_name"); ok {
		tf.DisplayName = v.(string)
		tf.RawDisplayName = v.(string)
	}

	if v, ok := d.GetOk("position"); ok {
		tf.Position = int64(v.(int))
	}

	if v, ok := d.GetOk("active"); ok {
		tf.Active = v.(bool)
	}

	if v, ok := d.GetOk("end_user_visible"); ok {
		tf.EndUserVisible = v.(bool)
	}

	if v, ok := d.GetOk("default"); ok {
		tf.Default = v.(bool)
	}

	if v, ok := d.GetOk("in_all_brands"); ok {
		tf.InAllBrands = v.(bool)
	}

	if v, ok := d.GetOk("ticket_field_ids"); ok {
		ticketFieldIDs := v.(*schema.Set).List()
		for _, ticketFieldID := range ticketFieldIDs {
			tf.TicketFieldIDs = append(tf.TicketFieldIDs, int64(ticketFieldID.(int)))
		}
	}

	if v, ok := d.GetOk("restricted_brand_ids"); ok {
		brandIDs := v.(*schema.Set).List()
		for _, id := range brandIDs {
			tf.TicketFieldIDs = append(tf.RestrictedBrandIDs, int64(id.(int)))
		}
	}

	return tf, nil
}

// marshalTicketField encodes the provided form into the provided resource data
func marshalTicketForm(f client.TicketForm, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":                  f.URL,
		"name":                 f.Name,
		"display_name":         f.DisplayName,
		"position":             f.Position,
		"active":               f.Active,
		"end_user_visible":     f.EndUserVisible,
		"default":              f.Default,
		"ticket_field_ids":     f.TicketFieldIDs,
		"in_all_brands":        f.InAllBrands,
		"restricted_brand_ids": f.RestrictedBrandIDs,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func createTicketForm(d identifiableGetterSetter, zd client.TicketFormAPI) error {
	tf, err := unmarshalTicketForm(d)
	if err != nil {
		return err
	}

	// Actual API request
	ctx := context.Background()
	tf, err = zd.CreateTicketForm(ctx, tf)
	if err != nil {
		return err
	}

	// Patch from created resource
	d.SetId(fmt.Sprintf("%d", tf.ID))
	return marshalTicketForm(tf, d)
}

func readTicketForm(d identifiableGetterSetter, zd client.TicketFormAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	tf, err := zd.GetTicketForm(ctx, id)
	if err != nil {
		return err
	}

	return marshalTicketForm(tf, d)
}

func updateTicketForm(d identifiableGetterSetter, zd client.TicketFormAPI) error {
	tf, err := unmarshalTicketForm(d)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tf, err = zd.UpdateTicketForm(ctx, tf.ID, tf)
	if err != nil {
		return err
	}

	return marshalTicketForm(tf, d)
}

func deleteTicketForm(d identifiable, zd client.TicketFormAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteTicketForm(ctx, id)
}
