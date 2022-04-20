package zendesk

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	. "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestMarshalGroup(t *testing.T) {
	expectedURL := "https://example.com"
	expectedName := "Support"

	g := zendesk.Group{
		URL:  expectedURL,
		Name: expectedName,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalGroup(g, m)
	if err != nil {
		t.Fatalf("Could marshal map %v", err)
	}

	v, ok := m.GetOk("url")
	if !ok {
		t.Fatalf("Failed to get url value")
	}
	if v != expectedURL {
		t.Fatalf("group had incorrect url value %v. should have been %v", v, expectedURL)
	}

	v, ok = m.GetOk("name")
	if !ok {
		t.Fatalf("Failed to get name value")
	}
	if v != expectedName {
		t.Fatalf("group had incorrect name value %v. should have been %v", v, expectedName)
	}
}

func TestUnmarshalGroup(t *testing.T) {
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

	m.EXPECT().GetGroup(Any(), Any()).Return(field, nil)
	if err := readGroup(context.Background(), gs, m); err != nil {
		t.Fatalf("readGroup returned an error: %v", err)
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

	m.EXPECT().CreateGroup(Any(), Any()).Return(out, nil)
	if err := createGroup(context.Background(), i, m); err != nil {
		t.Fatalf("create group returned an error: %v", err)
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("Create did not set resource id. Id was %s", v)
	}
}

func TestUpdateGroup(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: make(mapGetterSetter),
	}

	m.EXPECT().UpdateGroup(Any(), Eq(int64(12345)), Any()).Return(zendesk.Group{}, nil)
	if err := updateGroup(context.Background(), i, m); err != nil {
		t.Fatalf("updateGroup returned an error: %v", err)
	}
}

func TestDeleteGroup(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id: "12345",
	}

	m.EXPECT().DeleteGroup(Any(), Eq(int64(12345))).Return(nil)
	if err := deleteGroup(context.Background(), i, m); err != nil {
		t.Fatalf("deleteGroup returned an error: %v", err)
	}
}

func testGroupDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.GroupAPI)

	for _, r := range s.RootModule().Resources {
		if r.Type != "zendesk_group" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		group, err := client.GetGroup(context.Background(), id)
		if err != nil {
			return err
		}

		if !group.Deleted {
			return fmt.Errorf("group %d is not marked as deleted", group.ID)
		}
	}

	return nil
}

func TestAccGroupExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "resources/zendesk_group/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_group.moderator-group", "name", "Moderator"),
					resource.TestCheckResourceAttrSet("zendesk_group.moderator-group", "url"),
					resource.TestCheckResourceAttr("zendesk_group.developer-group", "name", "Developer"),
					resource.TestCheckResourceAttrSet("zendesk_group.developer-group", "url"),
				),
			},
		},
	})
}
