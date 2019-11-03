package zendesk

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

const (
	accountVar = "ZENDESK_ACCEPTANCE_TEST_ACCOUNT"
	emailVar   = "ZENDESK_ACCEPTANCE_TEST_EMAIL"
	tokenVar   = "ZENDESK_ACCEPTANCE_TEST_TOKEN"
)

// Provider returns provider instance for Zendesk
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		// https://developer.zendesk.com/rest_api/docs/support/introduction#security-and-authentication
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(accountVar, ""),
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(emailVar, ""),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(tokenVar, ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zendesk_brand":        resourceZendeskBrand(),
			"zendesk_group":        resourceZendeskGroup(),
			"zendesk_ticket_field": resourceZendeskTicketField(),
			"zendesk_ticket_form":  resourceZendeskTicketForm(),
			"zendesk_trigger":      resourceZendeskTrigger(),
			"zendesk_attachment":   resourceZendeskAttachment(),
			"zendesk_organization": resourceZendeskOrganization(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"zendesk_ticket_field": dataSourceZendeskTicketField(),
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
