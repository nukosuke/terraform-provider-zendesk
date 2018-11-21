package zendesk

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider returns provider instance for Zendesk
func Provider() *schema.Provider {
	return &schema.Provider{
		// https://developer.zendesk.com/rest_api/docs/support/introduction#security-and-authentication
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"token": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zendesk_ticket_field": resourceZendeskTicketField(),
			"zendesk_ticket_form":  resourceZendeskTicketForm(),
			"zendesk_trigger":      resourceZendeskTrigger(),
		},
	}
}
