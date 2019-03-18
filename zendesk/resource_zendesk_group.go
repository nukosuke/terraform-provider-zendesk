package zendesk

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/groups
func resourceZendeskGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceZendeskGroupCreate,
		Read:   resourceZendeskGroupRead,
		Update: resourceZendeskGroupUpdate,
		Delete: resourceZendeskGroupDelete,
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

func marshalGroup(field client.Group, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url": field.URL,
		"name": field.Name,
	}

	err := setSchemaFields(d, fields)
	if err != nil {
		return err
	}
	
	return nil
}

func unmarshalGroup(d identifiableGetterSetter) (client.Group, error) {
	group := client.Group{
		URL: d.Get("url").(string),
		Name: d.Get("name").(string),
	}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return group, fmt.Errorf("could not parse group id %s: %v", v, err)
		}
		group.ID = id
	}

	return group, nil
}

func resourceZendeskGroupCreate(d *schema.ResourceData, meta interface{}) error {
	zd := meta.(*client.Client)
	group := client.Group{
		Name: d.Get("name").(string),
	}

	// Actual API request
	group, err := zd.CreateGroup(group)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", group.ID))
	return nil
}

func resourceZendeskGroupRead(d *schema.ResourceData, meta interface{}) error {
	zd := meta.(*client.Client)
	return readGroup(d, zd)
}

func readGroup(d identifiableGetterSetter, zd client.GroupAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	field, err := zd.GetGroup(id)
	if err != nil {
		return err
	}

	return marshalGroup(field, d)
}

func resourceZendeskGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
