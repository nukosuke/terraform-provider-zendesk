package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/triggers
func resourceZendeskTrigger() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a trigger resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.TriggerAPI)
			return createTrigger(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.TriggerAPI)
			return readTrigger(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.TriggerAPI)
			return updateTrigger(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.TriggerAPI)
			return deleteTrigger(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Description: "The title of the trigger.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"active": {
				Description: "Whether the trigger is active.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"position": {
				Description: "Position of the trigger, determines the order they will execute in.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			// Both the "all" and "any" parameter are optional, but at least one of them must be supplied
			"all": triggerConditionSchema("Logical AND. All the conditions must be met."),
			"any": triggerConditionSchema("Logical OR. Any condition can be met."),
			"action": {
				Description: "What the trigger will do.",
				Type:        schema.TypeSet,
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
				Description: "The description of the trigger.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
		},
	}
}

// Marshal the zendesk client object to the terraform schema
func marshalTrigger(trigger client.Trigger, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"title":       trigger.Title,
		"active":      trigger.Active,
		"position":    trigger.Position,
		"description": trigger.Description,
	}

	var alls []map[string]interface{}
	for _, v := range trigger.Conditions.All {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		alls = append(alls, m)
	}
	fields["all"] = alls

	var anys []map[string]interface{}
	for _, v := range trigger.Conditions.Any {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		anys = append(anys, m)
	}
	fields["any"] = anys

	var actions []map[string]interface{}
	for _, action := range trigger.Actions {

		// If the trigger value is a string, leave it be
		// If it's a list, marshal it to a string
		var stringVal string
		switch action.Value.(type) {
		case []interface{}:
			tmp, err := json.Marshal(action.Value)
			if err != nil {
				return fmt.Errorf("error decoding trigger action value: %s", err)
			}
			stringVal = string(tmp)
		case string:
			stringVal = action.Value.(string)
		}

		m := map[string]interface{}{
			"field": action.Field,
			"value": stringVal,
		}
		actions = append(actions, m)
	}
	fields["action"] = actions
	return setSchemaFields(d, fields)
}

// Unmarshal the terraform schema to the Zendesk client object
func unmarshalTrigger(d identifiableGetterSetter) (client.Trigger, error) {
	trg := client.Trigger{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return trg, fmt.Errorf("could not parse trigger id %s: %v", v, err)
		}
		trg.ID = id
	}

	if v, ok := d.GetOk("title"); ok {
		trg.Title = v.(string)
	}

	if v, ok := d.GetOk("active"); ok {
		trg.Active = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		trg.Description = v.(string)
	}

	if v, ok := d.GetOk("all"); ok {
		allConditions := v.(*schema.Set).List()
		conditions := []client.TriggerCondition{}
		for _, c := range allConditions {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return trg, fmt.Errorf("could not parse 'all' conditions for trigger %v", trg)
			}
			conditions = append(conditions, client.TriggerCondition{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		trg.Conditions.All = conditions
	}

	if v, ok := d.GetOk("any"); ok {
		anyConditions := v.(*schema.Set).List()
		conditions := []client.TriggerCondition{}
		for _, c := range anyConditions {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return trg, fmt.Errorf("could not parse 'any' conditions for trigger %v", trg)
			}
			conditions = append(conditions, client.TriggerCondition{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		trg.Conditions.Any = conditions
	}

	if v, ok := d.GetOk("action"); ok {
		triggerActions := v.(*schema.Set).List()
		actions := []client.TriggerAction{}
		for _, a := range triggerActions {
			action, ok := a.(map[string]interface{})
			if !ok {
				return trg, fmt.Errorf("could not parse actions for trigger %v", trg)
			}

			// If the action value is a list, unmarshal it
			var actionValue interface{}
			if strings.HasPrefix(action["value"].(string), "[") {
				err := json.Unmarshal([]byte(action["value"].(string)), &actionValue)
				if err != nil {
					return trg, fmt.Errorf("error unmarshalling trigger action value: %s", err)
				}
			} else {
				actionValue = action["value"]
			}

			actions = append(actions, client.TriggerAction{
				Field: action["field"].(string),
				Value: actionValue,
			})
		}
		trg.Actions = actions
	}

	return trg, nil
}

func createTrigger(ctx context.Context, d identifiableGetterSetter, zd client.TriggerAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	trg, err := unmarshalTrigger(d)
	if err != nil {
		return diag.FromErr(err)
	}

	trg, err = zd.CreateTrigger(ctx, trg)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", trg.ID))

	err = marshalTrigger(trg, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readTrigger(ctx context.Context, d identifiableGetterSetter, zd client.TriggerAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	trigger, err := zd.GetTrigger(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTrigger(trigger, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateTrigger(ctx context.Context, d identifiableGetterSetter, zd client.TriggerAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	trigger, err := unmarshalTrigger(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	trigger, err = zd.UpdateTrigger(ctx, id, trigger)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTrigger(trigger, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteTrigger(ctx context.Context, d identifiable, zd client.TriggerAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteTrigger(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func triggerConditionSchema(desc string) *schema.Schema {
	return &schema.Schema{
		Description: desc,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field": {
					Description: "The name of a ticket field.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"operator": {
					Description: "A comparison operator.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"value": {
					Description: "The value of a ticket field.",
					Type:        schema.TypeString,
					Required:    true,
				},
			},
		},
		Optional: true,
	}
}
