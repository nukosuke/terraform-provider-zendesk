package zendesk

import "github.com/hashicorp/terraform/helper/schema"

// https://developer.zendesk.com/rest_api/docs/support/ticket_forms
func resourceZendeskTicketForm() *schema.Resource {
	return &schema.Resource{
		Create: resourceZendeskTicketFormCreate,
		Read:   resourceZendeskTicketFormRead,
		Update: resourceZendeskTicketFormUpdate,
		Delete: resourceZendeskTicketFormDelete,
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
			"raw_name": {
				Type:     schema.TypeString,
				Optional: true,
				//computed?
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"raw_display_name": {
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

func resourceZendeskTicketFormCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTicketFormRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTicketFormUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTicketFormDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
