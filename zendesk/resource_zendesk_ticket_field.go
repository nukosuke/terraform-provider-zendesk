package zendesk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// https://developer.zendesk.com/rest_api/docs/core/ticket_fields
func resourceZendeskTicketField() *schema.Resource {
	return &schema.Resource{
		Create: resourceZendeskTicketFieldCreate,
		Read:   resourceZendeskTicketFieldRead,
		Update: resourceZendeskTicketFieldUpdate,
		Delete: resourceZendeskTicketFieldDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				//TODO: not empty and included in
				// "checkbox", "date", "decimal", "integer", "regexp", "tagger", "text", or "textarea"
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"raw_title": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"raw_description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				//Computed?
				// positions 0 to 7 are reserved for system fields
				//TODO: Validation
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				//Computed?
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
				//Computed?
			},
			"collapsed_for_agents": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"regexp_for_validation": {
				Type:     schema.TypeString,
				Optional: true,
				// Regular expression field only
				//TODO: validation
			},
			"title_in_portal": {
				Type:     schema.TypeString,
				Optional: true,
				// The title of the ticket field is mandatory when it's visible to end users
				//TODO: validation
			},
			"raw_title_in_portal": {
				Type:     schema.TypeString,
				Optional: true,
				//TODO: same to title_in_portal
			},
			"visible_in_portal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"editable_in_portal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"required_in_portal": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_field_options": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			// Required only for "tagger" type
			"custom_field_options": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				//TODO: empty is invalid form
			},
			// "priority" and "status" fields only
			"sub_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				//TODO: validation
			},
			"removable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"agent_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceZendeskTicketFieldCreate(d *schema.ResourceData, meta interface{}) error {
	//TODO: fetch client from meta
	//client, _ := zd.NewClient(nil)
	d.SetId("1")
	return resourceZendeskTicketFieldRead(d, meta)
}

func resourceZendeskTicketFieldRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTicketFieldUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTicketFieldDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
