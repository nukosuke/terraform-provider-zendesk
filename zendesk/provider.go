package zendesk

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// Provider returns provider instance for Zendesk
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// https://developer.zendesk.com/rest_api/docs/support/introduction#security-and-authentication
		Schema: map[string]*schema.Schema{
			"account": {
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
			"zendesk_attachment":   resourceZendeskAttachment(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Account: d.Get("account").(string),
		Email:   d.Get("email").(string),
		Token:   d.Get("token").(string),
	}

	// Create & configure Zendesk API client
	zd, err := client.NewClient(nil) // TODO: set UserAgent to terraform/version
	if err != nil {
		return nil, err
	}

	if err = zd.SetSubdomain(config.Account); err != nil {
		return nil, err
	}
	zd.SetCredential(client.NewAPITokenCredential(config.Email, config.Token))

	return zd, nil
}
