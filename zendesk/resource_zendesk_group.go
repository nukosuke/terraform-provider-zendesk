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
	return nil
}

func resourceZendeskGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceZendeskGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
