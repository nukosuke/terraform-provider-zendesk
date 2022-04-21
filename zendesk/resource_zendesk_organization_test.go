package zendesk

import (
	"context"
	"strconv"
	"testing"

	. "github.com/golang/mock/gomock"
	//"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	//"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestMarshalOrganization(t *testing.T) {
	expectedURL := "https://example.com"
	expectedName := "Rebel Alliance"

	g := zendesk.Organization{
		URL:  expectedURL,
		Name: expectedName,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalOrganization(g, m)
	if err != nil {
		t.Fatalf("Could marshal map %v", err)
	}

	v, ok := m.GetOk("url")
	if !ok {
		t.Fatalf("Failed to get url value")
	}
	if v != expectedURL {
		t.Fatalf("organization had incorrect url value %v. should have been %v", v, expectedURL)
	}

	v, ok = m.GetOk("name")
	if !ok {
		t.Fatalf("Failed to get name value")
	}
	if v != expectedName {
		t.Fatalf("organization had incorrect name value %v. should have been %v", v, expectedName)
	}
}

func TestUnmarshalOrganization(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "1234",
		mapGetterSetter: mapGetterSetter{
			"url":  "https://example.zendesk.com/api/v2/organizations/361898904439.json",
			"name": "name",
		},
	}

	g, err := unmarshalOrganization(m)
	if err != nil {
		t.Fatalf("Could not marshal map %v", err)
	}

	if v := m.Get("url"); g.URL != v {
		t.Fatalf("organization had url value %v. should have been %v", g.URL, v)
	}

	if v := m.Get("name"); g.Name != v {
		t.Fatalf("organization had name value %v. should have been %v", g.Name, v)
	}
}

func TestReadOrganization(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	id := 1234
	gs := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
		id:              strconv.Itoa(id),
	}

	field := zendesk.Organization{
		ID:   int64(id),
		URL:  "foo",
		Name: "bar",
	}

	m.EXPECT().GetOrganization(Any(), Any()).Return(field, nil)
	if diags := readOrganization(context.Background(), gs, m); len(diags) != 0 {
		t.Fatalf("readOrganization returned an error: %v", diags)
	}

	if v := gs.mapGetterSetter["url"]; v != field.URL {
		t.Fatalf("url field %v does not have expected value %v", v, field.URL)
	}

	if v := gs.mapGetterSetter["name"]; v != field.Name {
		t.Fatalf("name field %v does not have expected value %v", v, field.Name)
	}
}

func TestCreateOrganization(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		mapGetterSetter: make(mapGetterSetter),
	}

	out := zendesk.Organization{
		ID: 12345,
	}

	m.EXPECT().CreateOrganization(Any(), Any()).Return(out, nil)
	if diags := createOrganization(context.Background(), i, m); len(diags) != 0 {
		t.Fatalf("create organization returned an error: %v", diags)
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}
}

func TestUpdateOrganization(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: make(mapGetterSetter),
	}

	m.EXPECT().UpdateOrganization(Any(), Eq(int64(12345)), Any()).Return(zendesk.Organization{}, nil)
	if diags := updateOrganization(context.Background(), i, m); len(diags) != 0 {
		t.Fatalf("updateOrganization returned an error: %v", diags)
	}
}

func TestDeleteOrganization(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id: "12345",
	}

	m.EXPECT().DeleteOrganization(Any(), Eq(int64(12345))).Return(nil)
	if diags := deleteOrganization(context.Background(), i, m); len(diags) != 0 {
		t.Fatalf("deleteOrganization returned an error: %v", diags)
	}
}
