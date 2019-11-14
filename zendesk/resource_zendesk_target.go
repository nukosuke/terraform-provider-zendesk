package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/targets
func resourceZendeskTarget() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return createTarget(d, zd)
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return readTarget(d, zd)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return updateTarget(d, zd)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return deleteTarget(d, zd)
		},
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
					"http_target",
					"url_target_v2", // synonym of http_target
					//"yammer_target",
				}, false),
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			// email_target
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// http_target
			"target_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"method": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"get",
					"patch",
					"put",
					"post",
					"delete",
				}, false),
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
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

func createTarget(d identifiableGetterSetter, zd client.TargetAPI) error {
	target, err := unmarshalTarget(d)
	if err != nil {
		return err
	}

	// Actual API request
	ctx := context.Background()
	target, err = zd.CreateTarget(ctx, target)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", target.ID))
	return marshalTarget(target, d)
}

func readTarget(d identifiableGetterSetter, zd client.TargetAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	target, err := zd.GetTarget(ctx, id)
	if err != nil {
		return err
	}

	return marshalTarget(target, d)
}

func updateTarget(d identifiableGetterSetter, zd client.TargetAPI) error {
	target, err := unmarshalTarget(d)
	if err != nil {
		return err
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	// ActualAPI request
	ctx := context.Background()
	target, err = zd.UpdateTarget(ctx, id, target)
	if err != nil {
		return err
	}

	return marshalTarget(target, d)
}

func deleteTarget(d identifiable, zd client.TargetAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteTarget(ctx, id)
}
