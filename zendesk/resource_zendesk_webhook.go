package zendesk

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/api-reference/webhooks/webhooks-api/webhooks
func resourceZendeskWebhook() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a webhook resource.",
		CreateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*client.Client)
			return createWebhook(ctx, data, zd)
		},
		ReadContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*client.Client)
			return readWebhook(ctx, data, zd)
		},
		UpdateContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*client.Client)
			return updateWebhook(ctx, data, zd)
		},
		DeleteContext: func(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(*client.Client)
			return deleteWebhook(ctx, data, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: webhookStateContext,
		},
		Schema: map[string]*schema.Schema{
			"authentication": {
				Description: "Adds authentication to the webhook's HTTP requests.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Type of authentication.",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: validation.StringInSlice(
								[]string{"api_key", "basic_auth", "bearer_token"},
								false,
							),
						},
						"data": {
							Description: "Authentication data.",
							Type:        schema.TypeMap,
							Required:    true,
							Sensitive:   true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"add_position": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "header",
							ValidateFunc: validation.StringInSlice(
								[]string{"header"},
								false,
							),
						},
					},
				},
			},
			"description": {
				Description: "Webhook description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"endpoint": {
				Description: "The destination URL that the webhook notifies when Zendesk events occur.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"http_method": {
				Description: `HTTP method used for the webhook's requests. To subscribe the webhook to Zendesk events, this must be "POST". Allowed values are "GET", "POST", "PUT", "PATCH", or "DELETE".`,
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						http.MethodGet,
						http.MethodPost,
						http.MethodPut,
						http.MethodPatch,
						http.MethodDelete,
					},
					false,
				),
			},
			"name": {
				Description: "Webhook name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"request_format": {
				Description: `The format of the data that the webhook will send. To subscribe the webhook to Zendesk events, this must be "json". Allowed values are "json", "xml", or "form_encoded".`,
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{"json", "xml", "form_encoded"},
					false,
				),
			},
			"status": {
				Description: `Current status of the webhook. Allowed values are "active", or "inactive".`,
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validation.StringInSlice(
					[]string{"active", "inactive"},
					false,
				),
			},
			"subscriptions": {
				Description: `Event subscriptions for the webhook. To subscribe the webhook to Zendesk events, specify one or more event types. For supported event type values, see Webhook event types. To connect the webhook to a trigger or automation, specify only "conditional_ticket_events" in the array.`,
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"signing_secret": {
				Description: "Signing secret used to verify webhook requests.",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// unmarshalTicketField parses the provided ResourceData and returns a ticket field
func unmarshalWebhook(d identifiableGetterSetter) (*client.Webhook, error) {
	var wh client.Webhook

	if v := d.Id(); v != "" {
		wh.ID = v
	}

	if v, ok := d.GetOk("name"); ok {
		wh.Name = v.(string)
	}

	if v, ok := d.GetOk("http_method"); ok {
		wh.HTTPMethod = v.(string)
	}

	if v, ok := d.GetOk("endpoint"); ok {
		wh.Endpoint = v.(string)
	}

	if v, ok := d.GetOk("request_format"); ok {
		wh.RequestFormat = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		wh.Status = v.(string)
	}

	if v, ok := d.GetOk("authentication"); ok {
		authentication := v.([]any)[0].(map[string]any)
		authType := authentication["type"].(string)

		wh.Authentication = &client.WebhookAuthentication{
			Type:        authType,
			Data:        authentication["data"].(map[string]any),
			AddPosition: authentication["add_position"].(string),
		}
	}

	if v, ok := d.GetOk("description"); ok {
		wh.Description = v.(string)
	}

	if v, ok := d.GetOk("subscriptions"); ok {
		subscriptions := v.(*schema.Set).List()
		wh.Subscriptions = make([]string, len(subscriptions))
		for i, sub := range subscriptions {
			wh.Subscriptions[i] = sub.(string)
		}
	}

	return &wh, nil
}

// marshalTicketField encodes the provided form into the provided resource data
func marshalWebhook(wh *client.Webhook, d identifiableGetterSetter) error {
	fields := map[string]any{
		"description":    wh.Description,
		"endpoint":       wh.Endpoint,
		"http_method":    wh.HTTPMethod,
		"name":           wh.Name,
		"request_format": wh.RequestFormat,
		"status":         wh.Status,
		"subscriptions":  wh.Subscriptions,
	}

	if wh.Authentication != nil {
		auth := map[string]any{
			"type":         wh.Authentication.Type,
			"data":         wh.Authentication.Data.(map[string]any),
			"add_position": wh.Authentication.AddPosition,
		}
		fields["authentication"] = []map[string]any{auth}
	}

	if !wh.CreatedAt.IsZero() {
		fields["created_at"] = wh.CreatedAt.String()
		fields["created_by"] = wh.CreatedBy
	}

	if !wh.UpdatedAt.IsZero() {
		fields["updated_at"] = wh.UpdatedAt.String()
		fields["updated_by"] = wh.UpdatedBy
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func createWebhook(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	wh, err := unmarshalWebhook(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	wh, err = zd.CreateWebhook(ctx, wh)
	if err != nil {
		return diag.FromErr(err)
	}

	// Patch from created resource
	d.SetId(wh.ID)

	err = marshalWebhook(wh, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readWebhook(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	wh, err := zd.GetWebhook(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalWebhook(wh, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateWebhook(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	wh, err := unmarshalWebhook(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.UpdateWebhook(ctx, wh.ID, wh)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalWebhook(wh, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteWebhook(ctx context.Context, d identifiable, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	err := zd.DeleteWebhook(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func webhookStateContext(ctx context.Context, d *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	zd := i.(client.WebhookAPI)
	wh, err := zd.GetWebhook(ctx, d.Id())
	if err != nil {
		return nil, err
	}

	if err := marshalWebhook(wh, d); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
