package zendesk

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	. "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestReadTicketField(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	id := 1234
	gs := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
		id:              strconv.Itoa(id),
	}

	now := time.Now()
	field := zendesk.TicketField{
		ID:                  int64(id),
		URL:                 "foobar",
		Type:                "decimal",
		Title:               "foobar",
		RawTitle:            "foobar",
		Description:         "foobar",
		RawDescription:      "foobar",
		Position:            int64(50),
		Active:              true,
		Required:            true,
		CollapsedForAgents:  true,
		RegexpForValidation: "regex",
		TitleInPortal:       "title",
		RawTitleInPortal:    "title",
		VisibleInPortal:     true,
		EditableInPortal:    true,
		RequiredInPortal:    true,
		Tag:                 "foobar",
		CreatedAt:           &now,
		UpdatedAt:           &now,
		SubTypeID:           int64(12345),
		Removable:           true,
		AgentDescription:    "foo",
		SystemFieldOptions: []zendesk.TicketFieldSystemFieldOption{{
			Name:  "Open",
			Value: "open",
		}},
		CustomFieldOptions: []zendesk.CustomFieldOption{{
			ID:    360013088874,
			Name:  "Option 1",
			Value: "opt1",
		}},
	}

	m.EXPECT().GetTicketField(Any()).Return(field, nil)
	if err := readTicketField(gs, m); err != nil {
		t.Fatal("readTicketField returned an error")
	}

	if v := gs.mapGetterSetter["url"]; v != field.URL {
		t.Fatalf("url field %v does not have expected value %v", v, field.URL)
	}

	if v := gs.mapGetterSetter["type"]; v != field.Type {
		t.Fatalf("type field %v does not have expected value %v", v, field.Type)
	}

	if v := gs.mapGetterSetter["title"]; v != field.Title {
		t.Fatalf("type field %v does not have expected value %v", v, field.Title)
	}

	if v := gs.mapGetterSetter["tag"]; v != field.Tag {
		t.Fatalf("type field %v does not have expected value %v", v, field.Title)
	}
}

func TestDeleteTicketField(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id: "12345",
	}

	m.EXPECT().DeleteTicketField(Eq(int64(12345))).Return(nil)
	if err := deleteTicketField(i, m); err != nil {
		t.Fatal("readTicketField returned an error")
	}
}

func TestUpdateTicketField(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: make(mapGetterSetter),
	}

	m.EXPECT().UpdateTicketField(Eq(int64(12345)), Any()).Return(zendesk.TicketField{}, nil)
	if err := updateTicketField(i, m); err != nil {
		t.Fatal("readTicketField returned an error")
	}
}

func TestCreateTicketField(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
	}

	out := zendesk.TicketField{
		ID: 12345,
	}

	m.EXPECT().CreateTicketField(Any()).Return(out, nil)
	if err := createTicketField(i, m); err != nil {
		t.Fatal("create ticket field returned an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}
}

func TestMarshalTicketField(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "1234",
		mapGetterSetter: mapGetterSetter{
			"url":                   "https://example.zendesk.com/api/v2/ticket_fields/360011737434.json",
			"type":                  "subject",
			"title":                 "title",
			"description":           "description",
			"position":              1,
			"active":                true,
			"required":              false,
			"collapsed_for_agents":  false,
			"regexp_for_validation": "+w{2}",
			"title_in_portal":       "portal",
			"visible_in_portal":     true,
			"editable_in_portal":    true,
			"required_in_portal":    true,
			"tag":                   "tag",
			"removable":             false,
			"agent_description":     "hey agents",
			"sub_type_id":           0,
		},
	}

	tf, err := unmarshalTicketField(m)
	if err != nil {
		t.Fatalf("Could marshal map %v", err)
	}

	if v := m.Get("url"); tf.URL != v {
		t.Fatalf("ticket had url value %v. shouldhave been %v", tf.URL, v)
	}

	if v := m.Get("title"); tf.Title != v && tf.RawTitle != v {
		t.Fatalf("ticket had incorrect title value %v or raw title %v. should have been %v", tf.Title, tf.RawTitle, v)
	}

	if v := m.Get("description"); tf.Description != v && tf.RawDescription != v {
		t.Fatalf("ticket had incorrect description value %v or raw description %v. should have been %v", tf.Description, tf.RawDescription, v)
	}

	if v := m.Get("title_in_portal"); tf.TitleInPortal != v && tf.RawTitleInPortal != v {
		t.Fatalf("ticket had incorrect title in portal value %v or raw title in portal %v. should have been %v", tf.Description, tf.RawDescription, v)
	}

	if v := m.Get("tag"); tf.Tag != v || tf.Tag == "" {
		t.Fatalf("ticket had incorrect tag %v. should have been %v", tf.Tag, v)
	}

}

func testTicketFieldDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.TicketFieldAPI)

	for k, r := range s.RootModule().Resources {
		if strings.HasPrefix(k, "data") {
			continue
		}

		if r.Type != "zendesk_ticket_field" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 64)
		if err != nil {
			return err
		}

		_, err = client.GetTicketField(id)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed ticket field named %s", k)
		}

		zd, ok := err.(zendesk.Error)
		if !ok {
			return fmt.Errorf("error %v cannot be asserted as a zendesk error", err)
		}

		if zd.Status() != http.StatusNotFound {
			return fmt.Errorf("did not get a not found error after destroy. error was %v", zd)
		}

	}

	return nil
}

func TestAccTicketFieldExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testSystemFieldVariablePreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testTicketFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "ticket_fields.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_ticket_field.checkbox-field", "title", "Checkbox Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.date-field", "title", "Date Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.decimal-field", "title", "Decimal Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.integer-field", "title", "Integer Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.regexp-field", "title", "Regexp Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.tagger-field", "title", "Tagger Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.text-field", "title", "Text Field"),
					resource.TestCheckResourceAttr("zendesk_ticket_field.textarea-field", "title", "Textarea Field"),
				),
			},
		},
	})
}
