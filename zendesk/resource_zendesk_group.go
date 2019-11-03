package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/groups
func resourceZendeskGroup() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return createGroup(d, zd)
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return readGroup(d, zd)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return updateGroup(d, zd)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return deleteGroup(d, zd)
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func createGroup(d identifiableGetterSetter, zd client.GroupAPI) error {
	group, err := unmarshalGroup(d)
	if err != nil {
		return err
	}

	// Actual API request
	ctx := context.Background()
	group, err = zd.CreateGroup(ctx, group)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", group.ID))
	return marshalGroup(group, d)
}

func readGroup(d identifiableGetterSetter, zd client.GroupAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	group, err := zd.GetGroup(ctx, id)
	if err != nil {
		return err
	}

	return marshalGroup(group, d)
}

func updateGroup(d identifiableGetterSetter, zd client.GroupAPI) error {
	group, err := unmarshalGroup(d)
	if err != nil {
		return err
	}

	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	// ActualAPI request
	ctx := context.Background()
	group, err = zd.UpdateGroup(ctx, id, group)
	if err != nil {
		return err
	}

	return marshalGroup(group, d)
}

func deleteGroup(d identifiable, zd client.GroupAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteGroup(ctx, id)
}
