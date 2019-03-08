package zendesk

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
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
				ValidateFunc: validation.StringInSlice([]string{
					"checkbox",
					"date",
					"decimal",
					"integer",
					"regexp",
					"tagger",
					"text",
					"textarea",
				}, false),
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"raw_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"raw_description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				// positions 0 to 7 are reserved for system fields
				ValidateFunc: validation.IntAtLeast(8),
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"required": {
				Type:     schema.TypeBool,
				Optional: true,
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
			// https://developer.zendesk.com/rest_api/docs/support/ticket_fields#updating-drop-down-field-options
			"custom_field_option": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
				//TODO: empty is invalid form
			},
			// "priority" and "status" fields only
			"sub_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				//TODO: validation
			},
			//TODO: this is not necessary because it's only for system field
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
	zd := meta.(*client.Client)
	tf := client.TicketField{
		Type:  d.Get("type").(string),
		Title: d.Get("title").(string),
	}

	// Handle type specific value
	switch d.Get("type") {
	case "regexp":
		tf.RegexpForValidation = d.Get("regexp_for_validation").(string)
	case "tagger":
		options := d.Get("custom_field_option").(*schema.Set).List()

		for _, option := range options {
			tf.CustomFieldOptions = append(tf.CustomFieldOptions, client.CustomFieldOption{
				Name:  option.(map[string]interface{})["name"].(string),
				Value: option.(map[string]interface{})["value"].(string),
			})
		}
	default:
		// nop
	}

	// Actual API request
	tf, err := zd.CreateTicketField(tf)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", tf.ID))
	d.Set("url", tf.URL)
	return resourceZendeskTicketFieldRead(d, meta)
}

func setSchemaFields(d *schema.ResourceData, m map[string]interface{}) error {
	for k, v := range m {
		err := d.Set(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceZendeskTicketFieldRead(d *schema.ResourceData, meta interface{}) error {
	zd := meta.(*client.Client)
	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return err
	}

	field, err := zd.GetTicketField(id)
	if err != nil {
		return err
	}

	fields := map[string]interface{}{
		"url":                   field.URL,
		"type":                  field.Type,
		"title":                 field.Title,
		"raw_title":             field.RawTitle,
		"description":           field.Description,
		"raw_description":       field.RawDescription,
		"position":              field.Position,
		"active":                field.Active,
		"required":              field.Required,
		"collapsed_for_agents":  field.CollapsedForAgents,
		"regexp_for_validation": field.RegexpForValidation,
		"title_in_portal":       field.TitleInPortal,
		"raw_title_in_portal":   field.RawTitleInPortal,
		"visible_in_portal":     field.VisibleInPortal,
		"editable_in_portal":    field.EditableInPortal,
		"required_in_portal":    field.Required,
		"tag":                   field.Tag,
		"system_field_options":  field.SystemFieldOptions,
		"custom_field_option":   field.CustomFieldOptions,
		"sub_type_id":           field.SubTypeID,
		"removable":             field.Removable,
		"agent_description":     field.AgentDescription,
	}

	err = setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func resourceZendeskTicketFieldUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskTicketFieldDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
