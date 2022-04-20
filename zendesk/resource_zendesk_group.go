package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/groups
func resourceZendeskGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provides a group resource.",
		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return createGroup(ctx, d, zd)
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return readGroup(ctx, d, zd)
		},
		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return updateGroup(ctx, d, zd)
		},
		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			zd := meta.(*client.Client)
			return deleteGroup(ctx, d, zd)
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "Group name.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func marshalGroup(group client.Group, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":  group.URL,
		"name": group.Name,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}

	return nil
}

func unmarshalGroup(d identifiableGetterSetter) (client.Group, error) {
	group := client.Group{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return group, fmt.Errorf("could not parse group id %s: %v", v, err)
		}
		group.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		group.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		group.Name = v.(string)
	}

	return group, nil
}

func createGroup(ctx context.Context, d identifiableGetterSetter, zd client.GroupAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	group, err := unmarshalGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Actual API request
	group, err = zd.CreateGroup(ctx, group)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", group.ID))

	err = marshalGroup(group, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func readGroup(ctx context.Context, d identifiableGetterSetter, zd client.GroupAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := zd.GetGroup(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalGroup(group, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func updateGroup(ctx context.Context, d identifiableGetterSetter, zd client.GroupAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	group, err := unmarshalGroup(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// ActualAPI request
	group, err = zd.UpdateGroup(ctx, id, group)
	if err != nil {
		return diag.FromErr(err)
	}

	err = marshalGroup(group, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func deleteGroup(ctx context.Context, d identifiable, zd client.GroupAPI) diag.Diagnostics {
	var diags diag.Diagnostics

	id, err := atoi64(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = zd.DeleteGroup(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
