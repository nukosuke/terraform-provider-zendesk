package zendesk

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	client "github.com/nukosuke/go-zendesk/zendesk"
)

// https://developer.zendesk.com/api-reference/ticketing/users/users/
func resourceZendeskUser() *schema.Resource {
	return &schema.Resource{
		Description: "[Experimental] Provides a user resouce.",
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The user's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"phone": {
				Description: "The user's primary phone number. The phone number should comply with the E.164 international telephone numbering plan. Example +15551234567.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"alias": {
				Description: "An alias displayed to end users.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"details": {
				Description: "Any details you want to store about the user, such as an address.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"notes": {
				Description: "Any notes you want to store about the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"role": {
				Description: `The user's role. Possible values are "end-user", "agent", or "admin".`,
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"end-user",
					"agent",
					"admin",
				}, false),
			},
			"custom_role_id": {
				Description: "A custom role if the user is an agent on the Enterprise plan or above.",
				Type:        schema.TypeSet,
				Optional:    true,
			},
			"default_group_id": {
				Description: "The id of the user's default group.",
				Type:        schema.TypeSet,
				Optional:    true,
			},
			"ticket_restriction": {
				Description: `Specifies which tickets the user has access to. Possible values are: "organization", "groups", "assigned", "requested", null.`,
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"organization",
					"groups",
					"assigned",
					"requested",
				}, false),
			},
			"time_zone": {
				Description: "The user's time zone.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tags": {
				Description: "The user's tags. Only present if your account has user tagging enabled.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}
