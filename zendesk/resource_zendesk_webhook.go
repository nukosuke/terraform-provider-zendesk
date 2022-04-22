package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/api-reference/event-connectors/webhooks/webhooks/
func resourceZendeskWebhook() *schema.Resource {
	return &schema.Resource{
		Description: "[Experimental] Provides a webhook resource. This feature is still experimental and has not been fully tested. Do not use in production environments.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.WebhookAPI)
			return createWebhook(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.WebhookAPI)
			return readWebhook(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.WebhookAPI)
			return updateWebhook(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
			zd := i.(client.WebhookAPI)
			return deleteWebhook(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Webhook name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Webhook description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"status": {
				Description: `Current status of the webhook. Allowed values are "active" or "inactive". Default is "active".`,
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "active",
				ValidateFunc: validation.StringInSlice([]string{
					"active",
					"inactive",
				}, false),
			},
			"endpoint": {
				Description:  "The destination URL that the webhook notifies when Zendesk events occur.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"http_method": {
				Description: `The HTTP method used by the webhook. Allowed values are "GET", "POST", "PUT", "PATCH", or "DELETE". Default is "POST"`,
				Type:        schema.TypeString,
				Optional:    true,
				Default:     http.MethodPost,
				ValidateFunc: validation.StringInSlice([]string{
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodPatch,
					http.MethodDelete,
				}, false),
			},
			"request_format": {
				Description: `The format of the data that the webhook will send. Allowed values are "json", "xml", or "form_encoded". Default is "json"`,
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "json",
				ValidateFunc: validation.StringInSlice([]string{
					"json",
					"xml",
					"form_encoded",
				}, false),
			},
			"authentication": {
				Description: "Authentication data that enables the integration with the destination system. Supports basic authentication and bearer token authentication.",
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: `Authentication type. Allowed values are "basic_auth" or "bearer_token".`,
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: validation.StringInSlice([]string{
								"basic_auth",
								"bearer_token",
							}, false),
						},
						"add_position": {
							Description: `Where to add credentials. Allowed value is only "header" currently.`,
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: validation.StringInSlice([]string{
								"header",
							}, false),
						},
						"data": {
							Description:  `Authentication data as JSON string. This field generally includes credentials username and password for "basic_auth", token for "bearer_token".`,
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsJSON,
							Sensitive:    true,
						},
					},
				},
			},
			"subscriptions": {
				Description: `Zendesk event subscriptions. Allowed value is "conditional_ticket_events".`,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"conditional_ticket_events",
					}, false),
				},
				Optional: true,
			},
		},
	}
}

func marshalWebhook(webhook *client.Webhook, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"name":           webhook.Name,
		"description":    webhook.Description,
		"status":         webhook.Status,
		"endpoint":       webhook.Endpoint,
		"http_method":    webhook.HTTPMethod,
		"request_format": webhook.RequestFormat,
	}

	if webhook.Authentication != nil {
		data, err := json.Marshal(webhook.Authentication.Data)
		if err != nil {
			return err
		}

		auth := map[string]interface{}{
			"type":         webhook.Authentication.Type,
			"add_position": webhook.Authentication.AddPosition,
			"data":         string(data),
		}
		fields["authentication"] = []map[string]interface{}{auth}
	}

	if len(webhook.Subscriptions) > 0 {
		fields["subscriptions"] = webhook.Subscriptions
	}

	return setSchemaFields(d, fields)
}

func unmarshalWebhook(d identifiableGetterSetter) (*client.Webhook, error) {
	webhook := &client.Webhook{}

	if v := d.Id(); v != "" {
		webhook.ID = v
	}

	if v, ok := d.GetOk("name"); ok {
		webhook.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		webhook.Description = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		webhook.Status = v.(string)
	}

	if v, ok := d.GetOk("endpoint"); ok {
		webhook.Endpoint = v.(string)
	}

	if v, ok := d.GetOk("http_method"); ok {
		webhook.HTTPMethod = v.(string)
	}

	if v, ok := d.GetOk("request_format"); ok {
		webhook.RequestFormat = v.(string)
	}

	if v, ok := d.GetOk("authentication"); ok {
		authset, ok := v.(*schema.Set)
		if !ok {
			return nil, fmt.Errorf("Failed to cast authentication set")
		}

		auth, ok := authset.List()[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Failed to cast authentication map")
		}

		var data interface{}
		err := json.Unmarshal([]byte(auth["data"].(string)), &data)
		if err != nil {
			return nil, err
		}

		webhook.Authentication = &client.WebhookAuthentication{
			Type:        auth["type"].(string),
			AddPosition: auth["add_position"].(string),
			Data:        data,
		}
	}

	if v, ok := d.GetOk("subscriptions"); ok {
		subs := v.(*schema.Set).List()
		for _, sub := range subs {
			webhook.Subscriptions = append(webhook.Subscriptions, sub.(string))
		}
	}

	return webhook, nil
}

func createWebhook(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	webhook, err := unmarshalWebhook(d)
	if err != nil {
		return diag.FromErr(err)
	}

	webhook, err = zd.CreateWebhook(ctx, webhook)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s", webhook.ID))

	err = marshalWebhook(webhook, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readWebhook(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	webhook, err := zd.GetWebhook(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalWebhook(webhook, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateWebhook(ctx context.Context, d identifiableGetterSetter, zd client.WebhookAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	webhook, err := unmarshalWebhook(d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.UpdateWebhook(ctx, d.Id(), webhook)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalWebhook(webhook, d)
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
