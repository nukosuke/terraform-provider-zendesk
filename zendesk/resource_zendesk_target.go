package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/targets
func resourceZendeskTarget() *schema.Resource {
	return &schema.Resource{
		Description: `[Beta] Provides a target resource. (HTTP target is deprecated. See https://support.zendesk.com/hc/en-us/articles/4408826284698 for details.)`,
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return createTarget(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return readTarget(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return updateTarget(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return deleteTarget(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
					//"basecamp_target",
					//"campfire_target",
					//"clickatell_target",
					"email_target",
					//"flowdock_target",
					//"get_satisfaction_target",
					//"jira_target",
					//"pivotal_target",
					//"twitter_target",
					//"url_target",
					"http_target",   // DEPRECATED. will be removed in future.
					"url_target_v2", // DEPRECATED. synonym of http_target
					//"yammer_target",
				}, false),
			},
			"title": {
				Description: "A name for the target.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"active": {
				Description: "Whether or not the target is activated.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},

			// email_target
			"email": {
				Description: `Email address for "email_target"`,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"subject": {
				Description: `Email subject for "email_target"`,
				Type:        schema.TypeString,
				Optional:    true,
			},

			// http_target
			"target_url": {
				Description: "The URL for the target. Some target types commonly use this field.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"method": {
				Description: "HTTP method.",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"get",
					"patch",
					"put",
					"post",
					"delete",
				}, false),
			},
			"username": {
				Description: "Username of the account which the target recognize.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "Password of the account which the target authenticate.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"content_type": {
				Description: "Content-Type for http_target",
				Deprecated:  "http_target is deprecated",
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"application/json",
					"application/xml",
					"application/x-www-form-urlencoded",
				}, false),
			},
		},
	}
}

func marshalTarget(target client.Target, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":    target.URL,
		"type":   target.Type,
		"title":  target.Title,
		"active": target.Active,
		// email_target
		"email":   target.Email,
		"subject": target.Subject,
		// http_target
		"target_url":   target.TargetURL,
		"method":       target.Method,
		"username":     target.Username,
		"password":     target.Password,
		"content_type": target.ContentType,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalTarget(d identifiableGetterSetter) (client.Target, error) {
	target := client.Target{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return target, fmt.Errorf("could not parse target id %s: %v", v, err)
		}
		target.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		target.URL = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		target.Type = v.(string)
	}

	if v, ok := d.GetOk("title"); ok {
		target.Title = v.(string)
	}

	if v, ok := d.GetOk("active"); ok {
		target.Active = v.(bool)
	}

	// email_target

	if v, ok := d.GetOk("email"); ok {
		target.Email = v.(string)
	}

	if v, ok := d.GetOk("subject"); ok {
		target.Subject = v.(string)
	}

	// http_target

	if v, ok := d.GetOk("target_url"); ok {
		target.TargetURL = v.(string)
	}

	if v, ok := d.GetOk("method"); ok {
		target.Method = v.(string)
	}

	if v, ok := d.GetOk("username"); ok {
		target.Username = v.(string)
	}

	if v, ok := d.GetOk("password"); ok {
		target.Password = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok {
		target.ContentType = v.(string)
	}

	return target, nil
}

func createTarget(ctx context.Context, d identifiableGetterSetter, zd client.TargetAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	target, err := unmarshalTarget(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	target, err = zd.CreateTarget(ctx, target)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", target.ID))

	err = marshalTarget(target, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readTarget(ctx context.Context, d identifiableGetterSetter, zd client.TargetAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	target, err := zd.GetTarget(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTarget(target, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateTarget(ctx context.Context, d identifiableGetterSetter, zd client.TargetAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	target, err := unmarshalTarget(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// ActualAPI request
	target, err = zd.UpdateTarget(ctx, id, target)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalTarget(target, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteTarget(ctx context.Context, d identifiable, zd client.TargetAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteTarget(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
