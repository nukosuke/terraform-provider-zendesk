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
		t.Fatalf("type field %v does not have expected value %v", v, field.URL)
	}

	if v := gs.mapGetterSetter["title"]; v != field.Title {
		t.Fatalf("type field %v does not have expected value %v", v, field.URL)
	}
}
