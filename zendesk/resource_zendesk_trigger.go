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
			"conditions": {
				Type: schema.TypeSet,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"all": {
							Type:     schema.TypeSet,
							Elem:     schema.Resource{Schema: triggerConditionSchema()},
							Optional: true,
						},
						"any": {
							Type:     schema.TypeSet,
							Elem:     schema.Resource{Schema: triggerConditionSchema()},
							Optional: true,
						},
					},
				},
				Required: true,
			},
			"actions": {
				Type: schema.TypeSet,
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Required: true,
			},
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

func triggerConditionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"field": {
			Type:     schema.TypeString,
			Required: true,
		},
		"operator": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}
