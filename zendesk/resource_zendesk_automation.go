package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/api-reference/ticketing/business-rules/automations/
func resourceZendeskAutomation() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an automation resource.",
		Create: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.AutomationAPI)
			return createAutomation(d, zd)
		},
		Read: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.AutomationAPI)
			return readAutomation(d, zd)
		},
		Update: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.AutomationAPI)
			return updateAutomation(d, zd)
		},
		Delete: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.AutomationAPI)
			return deleteAutomation(d, zd)
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"title": {
				Description: "The title of the automation.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"active": {
				Description: "Whether the automation is active.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"position": {
				Description: "The position of the automation which specifies the order it will be executed.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			// Both the "all" and "any" parameter are optional, but at least one of them must be supplied
			"all": automationConditionSchema("Logical AND. All the conditions must be met."),
			"any": automationConditionSchema("Logical OR. Any condition can be met."),
			"action": {
				Description: "What the automation will do.",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Description: "The name of a ticket field to modify.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "The new value of the field.",
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

// Marshal the zendesk client object to the terraform schema
func marshalAutomation(automation client.Automation, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"title":    automation.Title,
		"active":   automation.Active,
		"position": automation.Position,
	}

	var alls []map[string]interface{}
	for _, v := range automation.Conditions.All {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		alls = append(alls, m)
	}
	fields["all"] = alls

	var anys []map[string]interface{}
	for _, v := range automation.Conditions.Any {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		anys = append(anys, m)
	}
	fields["any"] = anys

	var actions []map[string]interface{}
	for _, action := range automation.Actions {

		// If the automation value is a string, leave it be
		// If it's a list, marshal it to a string
		var stringVal string
		switch action.Value.(type) {
		case []interface{}:
			tmp, err := json.Marshal(action.Value)
			if err != nil {
				return fmt.Errorf("error decoding automation action value: %s", err)
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
func unmarshalAutomation(d identifiableGetterSetter) (client.Automation, error) {
	automation := client.Automation{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return automation, fmt.Errorf("could not parse automation id %s: %v", v, err)
		}
		automation.ID = id
	}

	if v, ok := d.GetOk("title"); ok {
		automation.Title = v.(string)
	}

	if v, ok := d.GetOk("active"); ok {
		automation.Active = v.(bool)
	}

	if v, ok := d.GetOk("all"); ok {
		allConditions := v.(*schema.Set).List()
		conditions := []client.AutomationCondition{}
		for _, c := range allConditions {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return automation, fmt.Errorf("could not parse 'all' conditions for automation %v", automation)
			}
			conditions = append(conditions, client.AutomationCondition{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		automation.Conditions.All = conditions
	}

	if v, ok := d.GetOk("any"); ok {
		anyConditions := v.(*schema.Set).List()
		conditions := []client.AutomationCondition{}
		for _, c := range anyConditions {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return automation, fmt.Errorf("could not parse 'any' conditions for automation %v", automation)
			}
			conditions = append(conditions, client.AutomationCondition{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		automation.Conditions.Any = conditions
	}

	if v, ok := d.GetOk("action"); ok {
		automationActions := v.(*schema.Set).List()
		actions := []client.AutomationAction{}
		for _, a := range automationActions {
			action, ok := a.(map[string]interface{})
			if !ok {
				return automation, fmt.Errorf("could not parse actions for automation %v", automation)
			}

			// If the action value is a list, unmarshal it
			var actionValue interface{}
			if strings.HasPrefix(action["value"].(string), "[") {
				err := json.Unmarshal([]byte(action["value"].(string)), &actionValue)
				if err != nil {
					return automation, fmt.Errorf("error unmarshalling automation action value: %s", err)
				}
			} else {
				actionValue = action["value"]
			}

			actions = append(actions, client.AutomationAction{
				Field: action["field"].(string),
				Value: actionValue,
			})
		}
		automation.Actions = actions
	}

	return automation, nil
}

func createAutomation(d identifiableGetterSetter, zd client.AutomationAPI) error {
	automation, err := unmarshalAutomation(d)
	if err != nil {
		return err
	}

	ctx := context.Background()
	automation, err = zd.CreateAutomation(ctx, automation)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", automation.ID))
	return marshalAutomation(automation, d)
}

func readAutomation(d identifiableGetterSetter, zd client.AutomationAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	automation, err := zd.GetAutomation(ctx, id)
	if err != nil {
		return err
	}

	return marshalAutomation(automation, d)
}

func updateAutomation(d identifiableGetterSetter, zd client.AutomationAPI) error {
	automation, err := unmarshalAutomation(d)
	if err != nil {
		return err
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	automation, err = zd.UpdateAutomation(ctx, id, automation)
	if err != nil {
		return err
	}

	return marshalAutomation(automation, d)
}

func deleteAutomation(d identifiable, zd client.AutomationAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteAutomation(ctx, id)
}

func automationConditionSchema(desc string) *schema.Schema {
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
