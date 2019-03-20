package zendesk

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
)

func testTicketFormDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.TicketFieldAPI)

	for _, r := range s.RootModule().Resources {
		if r.Type != "zendesk_ticket_field" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		_, err = client.GetTicketField(id)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed ticket field")
		}

		zd, ok := err.(zendesk.Error)
		if !ok {
			return fmt.Errorf("error %v cannot be asserted as a zendesk error", err)
		}

		if zd.Status() != http.StatusNotFound {
			return fmt.Errorf(`did not get a "not found error" after destroy. error was %v`, zd)
		}

	}

	return nil
}

func TestAccTicketFormExample(t *testing.T) {
	// TODO: remove this skip on upgrade
	t.Skip("the test zendesk account is currently a trial account and forms cannot be created")
	configs := []string{
		readExampleConfig(t, "ticket_fields.tf"),
		readExampleConfig(t, "ticket_forms.tf"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testTicketFieldDestroyed,
			testTicketFormDestroyed,
		),
		Steps: []resource.TestStep{
			{
				Config: concatExampleConfig(t, configs...),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_ticket_form.form-1", "name", "Form 1"),
					resource.TestCheckResourceAttr("zendesk_ticket_form.form-2", "name", "Form 2"),
				),
			},
		},
	})
}
