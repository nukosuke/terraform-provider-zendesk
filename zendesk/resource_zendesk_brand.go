package zendesk

import (
	//"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	//client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/brands
func resourceZendeskBrand() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			return nil
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return nil
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			return nil
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			return nil
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
			"brand_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"has_help_center": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"help_center_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"logo": {
				// TODO
			},
			"ticket_form_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
			},
			"subdomain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host_mapping": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signature_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
