package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/rest_api/docs/support/organizations
func resourceZendeskOrganization() *schema.Resource {
	return &schema.Resource{
		Create: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return createOrganization(d, zd)
		},
		Read: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return readOrganization(d, zd)
		},
		Update: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return updateOrganization(d, zd)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			zd := meta.(*client.Client)
			return deleteOrganization(d, zd)
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
			"domain_names": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"shared_tickets": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"shared_comments": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tags": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
		},
	}
}

func marshalOrganization(org client.Organization, d identifiableGetterSetter) error {
	fields := map[string]interface{}{
		"url":             org.URL,
		"name":            org.Name,
		"domain_names":    org.DomainNames,
		"group_id":        org.GroupID,
		"shared_tickets":  org.SharedTickets,
		"shared_comments": org.SharedComments,
		"tags":            org.Tags,
	}

	return setSchemaFields(d, fields)
}

func unmarshalOrganization(d identifiableGetterSetter) (client.Organization, error) {
	org := client.Organization{}

	if v := d.Id(); v != "" {
		id, err := atoi64(v)
		if err != nil {
			return org, fmt.Errorf("could not parse organization id %s: %v", v, err)
		}
		org.ID = id
	}

	if v, ok := d.GetOk("url"); ok {
		org.URL = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		org.Name = v.(string)
	}

	if v, ok := d.GetOk("domain_names"); ok {
		domainNames := v.(*schema.Set).List()
		for _, domainName := range domainNames {
			org.DomainNames = append(org.DomainNames, domainName.(string))
		}
	}

	if v, ok := d.GetOk("group_id"); ok {
		org.GroupID = int64(v.(int))
	}

	if v, ok := d.GetOk("shared_tickets"); ok {
		org.SharedTickets = v.(bool)
	}

	if v, ok := d.GetOk("shared_comments"); ok {
		org.SharedComments = v.(bool)
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		for _, tag := range tags {
			org.Tags = append(org.Tags, tag.(string))
		}
	}

	return org, nil
}

func createOrganization(d identifiableGetterSetter, zd client.OrganizationAPI) error {
	org, err := unmarshalOrganization(d)
	if err != nil {
		return err
	}

	ctx := context.Background()
	org, err = zd.CreateOrganization(ctx, org)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", org.ID))
	return marshalOrganization(org, d)
}

func readOrganization(d identifiableGetterSetter, zd client.OrganizationAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	org, err := zd.GetOrganization(ctx, id)
	if err != nil {
		return err
	}

	return marshalOrganization(org, d)
}

func updateOrganization(d identifiableGetterSetter, zd client.OrganizationAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	org, err := unmarshalOrganization(d)
	if err != nil {
		return err
	}

	ctx := context.Background()
	org, err = zd.UpdateOrganization(ctx, id, org)
	if err != nil {
		return err
	}

	return marshalOrganization(org, d)
}

func deleteOrganization(d identifiableGetterSetter, zd client.OrganizationAPI) error {
	id, err := atoi64(d.Id())
	if err != nil {
		return err
	}

	ctx := context.Background()
	return zd.DeleteOrganization(ctx, id)
}
