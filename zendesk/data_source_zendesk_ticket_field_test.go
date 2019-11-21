package zendesk

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

const systemFieldConfig = `
data "zendesk_ticket_field" "assignee" {
	type = "%s"
}
`

func TestTicketFieldDataSourceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mock.NewClient(ctrl)

	m := newIdentifiableGetterSetter()
	fieldtype := "subject"

	err := m.Set("type", fieldtype)
	if err != nil {
		t.Fatalf("Read system field returned an error. %v", err)
	}

	out := zendesk.TicketField{
		ID:    1234,
		Type:  "subject",
		Title: "Subject",
		URL:   "foobar",
	}

	c.EXPECT().GetTicketFields(gomock.Any()).Return([]zendesk.TicketField{out}, zendesk.Page{}, nil)
	c.EXPECT().GetTicketField(gomock.Any(), gomock.Any()).Return(out, nil)

	err = readTicketFieldDataSource(m, c)
	if err != nil {
		t.Fatalf("Read system field returned an error. %v", err)
	}

	if v, ok := strconv.Atoi(m.Id()); ok != nil || int64(v) != out.ID {
		t.Fatalf("Read system field did not set ID field. Expected %v, Got %v", out.ID, v)
	}

	if v, ok := m.GetOk("url"); !ok || v.(string) != out.URL {
		t.Fatalf("Read system field did not set URL field. Expected %v, Got %v", out.URL, v)
	}
}

func TestAccTicketFieldDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(systemFieldConfig, "assignee"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.zendesk_ticket_field.assignee", "id"),
					resource.TestCheckResourceAttrSet("data.zendesk_ticket_field.assignee", "url"),
					resource.TestCheckResourceAttr("data.zendesk_ticket_field.assignee", "type", "assignee"),
				),
			},
		},
	})
}
