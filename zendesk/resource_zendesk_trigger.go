package zendesk

import "github.com/hashicorp/terraform/helper/schema"

// https://developer.zendesk.com/rest_api/docs/support/triggers
func resourceZendeskTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceZendeskTriggerCreate,
		Read:   resourceZendeskTriggerRead,
		Update: resourceZendeskTriggerUpdate,
		Delete: resourceZendeskTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			//"conditions": {
			// Type:
			//},
			//"actions": {
			// Type:
			//},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceZendeskTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTriggerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
