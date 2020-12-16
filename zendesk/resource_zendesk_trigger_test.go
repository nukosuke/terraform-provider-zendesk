package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestMarshalTrigger(t *testing.T) {
	expected := zendesk.Trigger{
		Title:       "title",
		Description: "blabla",
		Active:      true,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalTrigger(expected, m)
	if err != nil {
		t.Fatalf("Failed to marshal map %v", err)
	}

	v, ok := m.GetOk("title")
	if !ok {
		t.Fatal("Failed to get title value")
	}
	if v != expected.Title {
		t.Fatalf("trigger had incorrect title value %v. should have been %v", v, expected.Title)
	}

	v, ok = m.GetOk("description")
	if !ok {
		t.Fatal("Failed to get description value")
	}
	if v != expected.Description {
		t.Fatalf("trigger had incorrect description value %v. should have been %v", v, expected.Description)
	}

	v, ok = m.GetOk("active")
	if !ok {
		t.Fatal("Failed to get active value")
	}
	if v != expected.Active {
		t.Fatalf("trigger had incorrect active value %v. should have been %v", v, expected.Active)
	}
}

func TestUnmarshalTrigger(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "100",
		mapGetterSetter: mapGetterSetter{
			"title":       "Auto reply",
			"description": "reply automatically",
			"active":      true,
		},
	}

	trg, err := unmarshalTrigger(m)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if v := m.Get("title"); trg.Title != v {
		t.Fatalf("trigger had title value %v. should have been %v", trg.Title, v)
	}
}

func TestCreateTrigger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	out := zendesk.Trigger{
		ID:    12345,
		Title: "trigger",
	}

	m.EXPECT().CreateTrigger(gomock.Any(), gomock.Any()).Return(out, nil)
	if err := createTrigger(i, m); err != nil {
		t.Fatal("CreateTrigger return an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("CreateTrigger did not set resource id. Id was %s", v)
	}

	if v := i.Get("title"); v != "trigger" {
		t.Fatalf("CreateTrigger did not set resource title. title was %s", v)
	}
}

func TestReadTrigger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := zendesk.Trigger{
		Title:  "trigger",
		Active: true,
	}
	m.EXPECT().GetTrigger(gomock.Any(), gomock.Eq(int64(12345))).Return(expected, nil)
	if err := readTrigger(i, m); err != nil {
		t.Fatalf("GetTrigger received an error when calling: %v", err)
	}
}

func TestUpdateTrigger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: mapGetterSetter{},
	}

	m.EXPECT().UpdateTrigger(gomock.Any(), gomock.Eq(int64(12345)), gomock.Any()).Return(zendesk.Trigger{}, nil)
	if err := updateTrigger(i, m); err != nil {
		t.Fatalf("updateTrigger returned an error %v", err)
	}
}

func TestDeleteTrigger(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteTrigger(gomock.Any(), gomock.Eq(int64(1234))).Return(nil)
	err := deleteTrigger(d, c)
	if err != nil {
		t.Fatalf("Got error from resource delete: %v", err)
	}
}

func testTriggerDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.TriggerAPI)

	for k, r := range s.RootModule().Resources {
		if r.Type != "zendesk_trigger" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		ctx := context.Background()
		_, err = client.GetTrigger(ctx, id)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed trigger. resource name %s", k)
		}

		zdresp, ok := err.(zendesk.Error)
		if !ok {
			return fmt.Errorf("error %v cannot be asserted as a zendesk error", err)
		}

		if zdresp.Status() != http.StatusNotFound {
			return fmt.Errorf("did not get a not found error after destroy. error was %v", zdresp)
		}
	}
	return nil
}

func TestAccTriggerExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testTriggerDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "triggers.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_trigger.auto-reply-trigger", "title", "Auto Reply Trigger"),
					resource.TestCheckResourceAttr("zendesk_trigger.auto-reply-trigger", "active", "true"),
					resource.TestCheckResourceAttrSet("zendesk_trigger.auto-reply-trigger", "all.#"),
					resource.TestCheckResourceAttrSet("zendesk_trigger.auto-reply-trigger", "action.#"),
				),
			},
		},
	})
}
