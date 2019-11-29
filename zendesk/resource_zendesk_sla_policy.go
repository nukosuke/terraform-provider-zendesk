package zendesk

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/slaPolicies
func resourceZendeskSLAPolicy() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.SLAPolicyAPI)
			return createSLAPolicy(d, zd)
		},
		Read: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.SLAPolicyAPI)
			return readSLAPolicy(d, zd)
		},
		Update: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.SLAPolicyAPI)
			return updateSLAPolicy(d, zd)
		},
		Delete: func(d *schema.ResourceData, i interface{}) error {
			zd := i.(client.SLAPolicyAPI)
			return deleteSLAPolicy(d, zd)
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
				Computed: true,
			},
			// Both the "all" and "any" parameter are optional, but at least one of them must be supplied
			"all": slaPolicyFilterSchema(),
			"any": slaPolicyFilterSchema(),
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

// Marshal the zendesk client object to the terraform schema
func marshalSLAPolicy(slaPolicy client.SLAPolicy, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"title":          slaPolicy.Title,
		"active":         slaPolicy.Active,
		"position":       slaPolicy.Position,
		"description":    slaPolicy.Description,
		"policy_metrics": slaPolicy.PolicyMetrics,
	}

	var alls []map[string]interface{}
	for _, v := range slaPolicy.Filter.All {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		alls = append(alls, m)
	}
	fields["all"] = alls

	var anys []map[string]interface{}
	for _, v := range slaPolicy.Filter.Any {
		m := map[string]interface{}{
			"field":    v.Field,
			"operator": v.Operator,
			"value":    v.Value,
		}
		anys = append(anys, m)
	}
	fields["any"] = anys

	var metrics []client.SLAPolicyMetric

	fields["policy_metrics"] = metrics
	return setSchemaFields(d, fields)
}

// Unmarshal the terraform schema to the Zendesk client object
func unmarshalSLAPolicy(d identifiableGetterSetter) (client.SLAPolicy, error) {
	sla := client.SLAPolicy{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return sla, fmt.Errorf("could not parse slaPolicy id %s: %v", v, err)
		}
		sla.ID = id
	}

	if v, ok := d.GetOk("title"); ok {
		sla.Title = v.(string)
	}

	if v, ok := d.GetOk("active"); ok {
		sla.Active = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		sla.Description = v.(string)
	}

	if v, ok := d.GetOk("all"); ok {
		allFilters := v.(*schema.Set).List()
		filters := []client.SLAPolicyFilter{}
		for _, c := range allFilters {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return sla, fmt.Errorf("could not parse 'all' filters for slaPolicy %v", sla)
			}
			filters = append(filters, client.SLAPolicyFilter{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		sla.Filter.All = filters
	}

	if v, ok := d.GetOk("any"); ok {
		anyFilters := v.(*schema.Set).List()
		filters := []client.SLAPolicyFilter{}
		for _, c := range anyFilters {
			condition, ok := c.(map[string]interface{})
			if !ok {
				return sla, fmt.Errorf("could not parse 'any' filters for slaPolicy %v", sla)
			}
			filters = append(filters, client.SLAPolicyFilter{
				Field:    condition["field"].(string),
				Operator: condition["operator"].(string),
				Value:    condition["value"].(string),
			})
		}
		sla.Filter.Any = filters
	}

	if v, ok := d.GetOk("policy_metrics"); ok {
		slaPolicyMetrics := v.(*schema.Set).List()
		metrics := []client.SLAPolicyMetric{}
		for _, a := range slaPolicyMetrics {
			metric, ok := a.(map[string]interface{})
			if !ok {
				return sla, fmt.Errorf("could not parse metrics for slaPolicy %v", sla)
			}

			metrics = append(metrics, client.SLAPolicyMetric{
				Priority:      metric["priority"].(string),
				Metric:        metric["metric"].(string),
				Target:        metric["target"].(int),
				BusinessHours: metric["business_hours"].(bool),
			})
		}
		sla.PolicyMetrics = metrics
	}

	return sla, nil
}

func createSLAPolicy(d identifiableGetterSetter, zd client.SLAPolicyAPI) error {
	sla, err := unmarshalSLAPolicy(d)
	if err != nil {
		return err
	}

	ctx := context.Background()
	sla, err = zd.CreateSLAPolicy(ctx, sla)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", sla.ID))
	return marshalSLAPolicy(sla, d)
}

func readSLAPolicy(d identifiableGetterSetter, zd client.SLAPolicyAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	slaPolicy, err := zd.GetSLAPolicy(ctx, id)
	if err != nil {
		return err
	}

	return marshalSLAPolicy(slaPolicy, d)
}

func updateSLAPolicy(d identifiableGetterSetter, zd client.SLAPolicyAPI) error {
	slaPolicy, err := unmarshalSLAPolicy(d)
	if err != nil {
		return err
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	slaPolicy, err = zd.UpdateSLAPolicy(ctx, id, slaPolicy)
	if err != nil {
		return err
	}

	return marshalSLAPolicy(slaPolicy, d)
}

func deleteSLAPolicy(d identifiable, zd client.SLAPolicyAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteSLAPolicy(ctx, id)
}

func slaPolicyFilterSchema() *schema.Schema {
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
