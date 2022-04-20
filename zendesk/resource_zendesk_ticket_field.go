package zendesk

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/core/ticket_fields
func resourceZendeskTicketField() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a ticket field resource.",
		CreateContext: resourceZendeskTicketFieldCreate,
		ReadContext:   resourceZendeskTicketFieldRead,
		UpdateContext: resourceZendeskTicketFieldUpdate,
		DeleteContext: resourceZendeskTicketFieldDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Description: "The URL for this ticket field.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "System or custom field type. Editable for custom field types and only on creation.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"checkbox",
					"date",
					"decimal",
					"integer",
					"multiselect",
					"partialcreditcard",
					"regexp",
					"tagger",
					"text",
					"textarea",
				}, false),
			},
			"title": {
				Description: "The title of the ticket field.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Describes the purpose of the ticket field to users.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"position": {
				Description: "The relative position of the ticket field on a ticket. Note that for accounts with ticket forms, positions are controlled by the different forms.",
				Type:        schema.TypeInt,
				Optional:    true,
				// positions 0 to 7 are reserved for system fields
				ValidateFunc: validation.IntAtLeast(8),
				Computed:     true,
			},
			"active": {
				Description: "Whether this field is available.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"required": {
				Description: "If true, agents must enter a value in the field to change the ticket status to solved.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"collapsed_for_agents": {
				Description: "If true, the field is shown to agents by default. If false, the field is hidden alongside infrequently used fields. Classic interface only.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"regexp_for_validation": {
				Description: `For "regexp" fields only. The validation pattern for a field value to be deemed valid.`,
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				// Regular expression field only
				//TODO: validation
			},
			"title_in_portal": {
				Description: "The title of the ticket field for end users in Help Center.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"visible_in_portal": {
				Description: "Whether this field is visible to end users in Help Center.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"editable_in_portal": {
				Description: "Whether this field is editable by end users in Help Center.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"required_in_portal": {
				Description: "If true, end users must enter a value in the field to create the request.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"tag": {
				Description: `For "checkbox" fields only. A tag added to tickets when the checkbox field is selected.`,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"system_field_options": {
				Description: `Presented for a system ticket field of type "tickettype", "priority" or "status".`,
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "System field option name.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"value": {
							Description: "System field option value.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
				Computed: true,
			},
			// https://developer.zendesk.com/api-reference/ticketing/tickets/ticket_fields/#updating-drop-down-field-options
			"custom_field_option": {
				Description: `Required and presented for a custom ticket field of type "multiselect" or "tagger".`,
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Custom field option name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"value": {
							Description: "Custom field option value.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "Custom field option id.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Optional: true,
				//TODO: empty is invalid form
			},
			// "priority" and "status" fields only
			"sub_type_id": {
				Description: `For system ticket fields of type "priority" and "status". Defaults to 0. A "priority" sub type of 1 removes the "Low" and "Urgent" options. A "status" sub type of 1 adds the "On-Hold" option.`,
				Type:        schema.TypeInt,
				Optional:    true,
				//TODO: validation
			},
			// NOTE: Maybe this is not necessary because it's only for system field
			"removable": {
				Description: "If false, this field is a system field that must be present on all tickets.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"agent_description": {
				Description: "A description of the ticket field that only agents can see.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

// marshalTicketField encodes the provided ticket field into the provided resource data
func marshalTicketField(field client.TicketField, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":                   field.URL,
		"type":                  field.Type,
		"title":                 field.Title,
		"description":           field.Description,
		"position":              field.Position,
		"active":                field.Active,
		"required":              field.Required,
		"collapsed_for_agents":  field.CollapsedForAgents,
		"regexp_for_validation": field.RegexpForValidation,
		"title_in_portal":       field.TitleInPortal,
		"visible_in_portal":     field.VisibleInPortal,
		"editable_in_portal":    field.EditableInPortal,
		"required_in_portal":    field.RequiredInPortal,
		"tag":                   field.Tag,
		"sub_type_id":           field.SubTypeID,
		"removable":             field.Removable,
		"agent_description":     field.AgentDescription,
	}

	// set system field options
	systemFieldOptions := make([]map[string]interface{}, 0)
	for _, v := range field.SystemFieldOptions {
		m := map[string]interface{}{
			"name":  v.Name,
			"value": v.Value,
		}
		systemFieldOptions = append(systemFieldOptions, m)
	}

	fields["system_field_options"] = systemFieldOptions

	// Set custom field options
	customFieldOptions := make([]map[string]interface{}, 0)
	for _, v := range field.CustomFieldOptions {
		m := map[string]interface{}{
			"name":  v.Name,
			"value": v.Value,
			"id":    v.ID,
		}
		customFieldOptions = append(customFieldOptions, m)
	}

	fields["custom_field_option"] = customFieldOptions

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalTicketField parses the provided ResourceData and returns a ticket field
func unmarshalTicketField(d identifiableGetterSetter) (client.TicketField, error) {
	tf := client.TicketField{}

	if v := d.Id(); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return tf, fmt.Errorf("could not parse ticket field id %s: %v", v, err)
		}
		tf.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		tf.URL = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		tf.Type = v.(string)
	}

	if v, ok := d.GetOk("title"); ok {
		tf.Title = v.(string)
		tf.RawTitle = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		tf.Description = v.(string)
		tf.RawDescription = v.(string)
	}

	if v, ok := d.GetOk("position"); ok {
		tf.Position = int64(v.(int))
	}

	if v, ok := d.GetOk("active"); ok {
		tf.Active = v.(bool)
	}

	if v, ok := d.GetOk("required"); ok {
		tf.Required = v.(bool)
	}

	if v, ok := d.GetOk("regexp_for_validation"); ok {
		tf.RegexpForValidation = v.(string)
	}

	if v, ok := d.GetOk("title_in_portal"); ok {
		tf.TitleInPortal = v.(string)
		tf.RawTitleInPortal = v.(string)
	}

	if v, ok := d.GetOk("visible_in_portal"); ok {
		tf.VisibleInPortal = v.(bool)
	}

	if v, ok := d.GetOk("editable_in_portal"); ok {
		tf.EditableInPortal = v.(bool)
	}

	if v, ok := d.GetOk("required_in_portal"); ok {
		tf.RequiredInPortal = v.(bool)
	}

	if v, ok := d.GetOk("tag"); ok {
		tf.Tag = v.(string)
	}

	if v, ok := d.GetOk("sub_type_id"); ok {
		tf.SubTypeID = int64(v.(int))
	}

	if v, ok := d.GetOk("removable"); ok {
		tf.Removable = v.(bool)
	}

	if v, ok := d.GetOk("agent_description"); ok {
		tf.AgentDescription = v.(string)
	}

	if v, ok := d.GetOk("custom_field_option"); ok {
		options := v.(*schema.Set).List()
		customFieldOptions := make([]client.CustomFieldOption, 0)
		for _, o := range options {
			option, ok := o.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse custom options for field %v", tf)
			}

			customFieldOptions = append(customFieldOptions, client.CustomFieldOption{
				Name:  option["name"].(string),
				Value: option["value"].(string),
				ID:    int64(option["id"].(int)),
			})
		}

		tf.CustomFieldOptions = customFieldOptions
	}

	if v, ok := d.GetOk("system_field_options"); ok {
		options := v.(*schema.Set).List()
		systemFieldOptions := make([]client.TicketFieldSystemFieldOption, 0)
		for _, o := range options {
			option, ok := o.(map[string]interface{})
			if !ok {
				return tf, fmt.Errorf("could not parse system options for field %v", tf)
			}

			systemFieldOptions = append(systemFieldOptions, client.TicketFieldSystemFieldOption{
				Name:  option["name"].(string),
				Value: option["value"].(string),
			})
		}

		tf.SystemFieldOptions = systemFieldOptions
	}

	return tf, nil
}

func resourceZendeskTicketFieldCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*client.Client)
	return createTicketField(ctx, d, zd)
}

func createTicketField(ctx context.Context, d identifiableGetterSetter, zd client.TicketFieldAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalTicketField(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.CreateTicketField(ctx, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))

	err = marshalTicketField(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskTicketFieldRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*client.Client)
	return readTicketField(ctx, d, zd)
}

func readTicketField(ctx context.Context, d identifiableGetterSetter, zd client.TicketFieldAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	field, err := zd.GetTicketField(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTicketField(field, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskTicketFieldUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*client.Client)
	return updateTicketField(ctx, d, zd)
}

func updateTicketField(ctx context.Context, d identifiableGetterSetter, zd client.TicketFieldAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	tf, err := unmarshalTicketField(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	tf, err = zd.UpdateTicketField(ctx, id, tf)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTicketField(tf, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceZendeskTicketFieldDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	zd := meta.(*client.Client)
	return deleteTicketField(ctx, d, zd)
}

func deleteTicketField(ctx context.Context, d identifiable, zd client.TicketFieldAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteTicketField(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
