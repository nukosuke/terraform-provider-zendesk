package zendesk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

const (
	accountVar = "ZENDESK_ACCEPTANCE_TEST_ACCOUNT"
	emailVar   = "ZENDESK_ACCEPTANCE_TEST_EMAIL"
	tokenVar   = "ZENDESK_ACCEPTANCE_TEST_TOKEN"
)

// Provider returns provider instance for Zendesk
func Provider() *schema.Provider {
	return &schema.Provider{
		// https://developer.zendesk.com/rest_api/docs/support/introduction#security-and-authentication
		Schema: map[string]*schema.Schema{
			"account": {
				Description: "Account name of your Zendesk instance.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(accountVar, ""),
			},
			"email": {
				Description: "Email address of agent user who have permission to access the API.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(emailVar, ""),
			},
			"token": {
				Description: "[API token](https://developer.zendesk.com/rest_api/docs/support/introduction#api-token) for your Zendesk instance.",
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(tokenVar, ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"zendesk_automation":   resourceZendeskAutomation(),
			"zendesk_brand":        resourceZendeskBrand(),
			"zendesk_group":        resourceZendeskGroup(),
			"zendesk_ticket_field": resourceZendeskTicketField(),
			"zendesk_ticket_form":  resourceZendeskTicketForm(),
			"zendesk_trigger":      resourceZendeskTrigger(),
			"zendesk_target":       resourceZendeskTarget(),
			"zendesk_attachment":   resourceZendeskAttachment(),
			"zendesk_organization": resourceZendeskOrganization(),
			"zendesk_sla_policy":   resourceZendeskSLAPolicy(),
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
