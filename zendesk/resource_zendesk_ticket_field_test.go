package zendesk

import (
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestReadTicketField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	id := 1234
	gs := identifiableMapGetterSetter{
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

	m.EXPECT().GetTicketField(gomock.Any()).Return(field, nil)
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
}

func TestDeleteTicketField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := identifiableMapGetterSetter{
		id: "12345",
	}

	m.EXPECT().DeleteTicketField(gomock.Eq(int64(12345))).Return(nil)
	if err := deleteTicketField(i, m); err != nil {
		t.Fatal("readTicketField returned an error")
	}
}

func TestMarshalTicketField(t *testing.T) {
	m := identifiableMapGetterSetter{
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

}
