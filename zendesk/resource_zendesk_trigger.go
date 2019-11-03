package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/triggers
func resourceZendeskTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceZendeskTriggerCreate,
		Read:   resourceZendeskTriggerRead,
		Update: resourceZendeskTriggerUpdate,
		Delete: func(data *schema.ResourceData, i interface{}) error {
			zd := i.(client.TriggerAPI)
			return resourceZendeskTriggerDelete(data, zd)
		},
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
				Default:  true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			// Both the "all" and "any" parameter are optional, but at least one of them must be supplied
			"all": triggerConditionSchema(),
			"any": triggerConditionSchema(),
			"action": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
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
				Default:  "",
			},
		},
	}
}

func resourceZendeskTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	zd := meta.(*client.Client)
	trg := client.Trigger{
		Title:       d.Get("title").(string),
		Active:      d.Get("active").(bool),
		Description: d.Get("description").(string),
	}

	// Conditions
	alls := d.Get("all").(*schema.Set).List()
	for _, all := range alls {
		trg.Conditions.All = append(trg.Conditions.All, client.TriggerCondition{
			Field:    all.(map[string]interface{})["field"].(string),
			Operator: all.(map[string]interface{})["operator"].(string),
			Value:    all.(map[string]interface{})["value"].(string),
		})
	}
	anys := d.Get("any").(*schema.Set).List()
	for _, any := range anys {
		trg.Conditions.Any = append(trg.Conditions.Any, client.TriggerCondition{
			Field:    any.(map[string]interface{})["field"].(string),
			Operator: any.(map[string]interface{})["operator"].(string),
			Value:    any.(map[string]interface{})["value"].(string),
		})
	}

	// Actions
	actions := d.Get("action").(*schema.Set).List()
	for _, action := range actions {
		trg.Actions = append(trg.Actions, client.TriggerAction{
			Field: action.(map[string]interface{})["field"].(string),
			Value: action.(map[string]interface{})["value"].(string),
		})
		// TODO: notification_user specific handling
	}

	// Actual API request
	ctx := context.Background()
	trg, err := zd.CreateTrigger(ctx, trg)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", trg.ID))
	return nil
}

func resourceZendeskTriggerRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTriggerDelete(d identifiable, zd client.TriggerAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteTrigger(ctx, id)
}

func triggerConditionSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			},
		},
		Optional: true,
	}
}
