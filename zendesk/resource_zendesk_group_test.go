package zendesk

import (
	"strconv"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestReadGroup(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	id := 1234
	gs := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
		id:              strconv.Itoa(id),
	}

	field := zendesk.Group{
		ID:   int64(id),
		URL:  "foo",
		Name: "bar",
	}

	m.EXPECT().GetGroup(Any()).Return(field, nil)
	if err := readGroup(gs, m); err != nil {
		t.Fatal("readGroup returned an error")
	}

	if v := gs.mapGetterSetter["url"]; v != field.URL {
		t.Fatalf("url field %v does not have expected value %v", v, field.URL)
	}

	if v := gs.mapGetterSetter["name"]; v != field.Name {
		t.Fatalf("name field %v does not have expected value %v", v, field.Name)
	}
}

func TestCreateGroup(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
	}

	out := zendesk.Group{
		ID: 12345,
	}

	m.EXPECT().CreateGroup(Any()).Return(out, nil)
	if err := createGroup(i, m); err != nil {
		t.Fatal("create group returned an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}
}

func TestMarshalGroup(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "1234",
		mapGetterSetter: mapGetterSetter{
			"url":  "https://example.zendesk.com/api/v2/ticket_fields/360011737434.json",
			"name": "name",
		},
	}

	g, err := unmarshalGroup(m)
	if err != nil {
		t.Fatalf("Could marshal map %v", err)
	}

	if v := m.Get("url"); g.URL != v {
		t.Fatalf("group had url value %v. should have been %v", g.URL, v)
	}

	if v := m.Get("name"); g.Name != v {
		t.Fatalf("group had name value %v. should have been %v", g.Name, v)
	}
}
