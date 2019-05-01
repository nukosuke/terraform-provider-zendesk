package zendesk

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestCreateTicketForm(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	out := zendesk.TicketForm{
		ID:   12345,
		Name: "foo",
	}

	m.EXPECT().CreateTicketForm(Any()).Return(out, nil)
	if err := createTicketForm(i, m); err != nil {
		t.Fatal("create ticket field returned an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}

	if v := i.Get("name"); v != "foo" {
		t.Fatalf("Create did not set resource name. name was %v", v)
	}
}

func TestDeleteTicketForm(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expectedID := int64(12345)
	m.EXPECT().DeleteTicketForm(Eq(expectedID)).Return(nil)
	if err := deleteTicketForm(i, m); err != nil {
		t.Fatal("create ticket field returned an error")
	}
}

func TestReadTicketForm(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := zendesk.TicketForm{
		Name:     "foobar",
		Position: int64(1),
	}
	m.EXPECT().GetTicketForm(Eq(int64(12345))).Return(expected, nil)
	if err := readTicketForm(i, m); err != nil {
		t.Fatalf("recieved an error when calling read ticket form: %v", err)
	}
}

func TestUnmarshalTicketForm(t *testing.T) {

	d := &identifiableMapGetterSetter{
		id: "47",
		mapGetterSetter: mapGetterSetter{
			"url":              "https://company.zendesk.com/api/v2/ticket_forms/47.json",
			"name":             "Snowboard Problem",
			"display_name":     "Snowboard Damage",
			"end_user_visible": true,
			"position":         9999,
			"active":           true,
			"default":          false,
			"in_all_brands":    false,
		},
	}

	tf, err := unmarshalTicketForm(d)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if tf.Name != d.Get("name") {
		t.Fatalf("ticket did not have the correct name")
	}
}

func testTicketFormDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.TicketFieldAPI)

	for _, r := range s.RootModule().Resources {
		if r.Type != "zendesk_ticket_field" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
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
	configs := []string{
		readExampleConfig(t, "ticket_fields.tf"),
		readExampleConfig(t, "ticket_forms.tf"),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testSystemFieldVariablePreCheck(t)
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
