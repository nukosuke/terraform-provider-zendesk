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

func TestMarshalAutomation(t *testing.T) {
	expected := zendesk.Automation{
		Title:    "title",
		Active:   true,
		Position: 1,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalAutomation(expected, m)
	if err != nil {
		t.Fatalf("Failed to marshal map %v", err)
	}

	v, ok := m.GetOk("title")
	if !ok {
		t.Fatal("Failed to get title value")
	}
	if v != expected.Title {
		t.Fatalf("automation had incorrect title value %v. should have been %v", v, expected.Title)
	}
	v, ok = m.GetOk("active")
	if !ok {
		t.Fatal("Failed to get active value")
	}
	if v != expected.Active {
		t.Fatalf("automation had incorrect active value %v. should have been %v", v, expected.Active)
	}
}

func TestUnmarshalAutomation(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "100",
		mapGetterSetter: mapGetterSetter{
			"title":       "Auto reply",
			"description": "reply automatically",
			"active":      true,
		},
	}

	automation, err := unmarshalAutomation(m)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if v := m.Get("title"); automation.Title != v {
		t.Fatalf("automation had title value %v. should have been %v", automation.Title, v)
	}
}

func TestCreateAutomation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	out := zendesk.Automation{
		ID:    12345,
		Title: "automation",
	}

	m.EXPECT().CreateAutomation(gomock.Any(), gomock.Any()).Return(out, nil)
	if err := createAutomation(i, m); err != nil {
		t.Fatal("CreateAutomation return an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("CreateAutomation did not set resource id. Id was %s", v)
	}

	if v := i.Get("title"); v != "automation" {
		t.Fatalf("CreateAutomation did not set resource title. title was %s", v)
	}
}

func TestReadAutomation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := zendesk.Automation{
		Title:  "automation",
		Active: true,
	}
	m.EXPECT().GetAutomation(gomock.Any(), gomock.Eq(int64(12345))).Return(expected, nil)
	if err := readAutomation(i, m); err != nil {
		t.Fatalf("GetAutomation received an error when calling: %v", err)
	}
}

func TestUpdateAutomation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: mapGetterSetter{},
	}

	m.EXPECT().UpdateAutomation(gomock.Any(), gomock.Eq(int64(12345)), gomock.Any()).Return(zendesk.Automation{}, nil)
	if err := updateAutomation(i, m); err != nil {
		t.Fatalf("updateAutomation returned an error %v", err)
	}
}

func TestDeleteAutomation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteAutomation(gomock.Any(), gomock.Eq(int64(1234))).Return(nil)
	err := deleteAutomation(d, c)
	if err != nil {
		t.Fatalf("Got error from resource delete: %v", err)
	}
}

func testAutomationDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.AutomationAPI)

	for k, r := range s.RootModule().Resources {
		if r.Type != "zendesk_automation" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		ctx := context.Background()
		_, err = client.GetAutomation(ctx, id)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed automation. resource name %s", k)
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

func TestAccAutomationExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAutomationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "resources/zendesk_automation/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zendesk_automation.auto-close-automation",
						"title",
						"Close ticket 4 days after status is set to solved",
					),
					resource.TestCheckResourceAttr("zendesk_automation.auto-close-automation",
						"active",
						"true",
					),
					resource.TestCheckResourceAttrSet("zendesk_automation.auto-close-automation", "all.#"),
					resource.TestCheckResourceAttrSet("zendesk_automation.auto-close-automation", "action.#"),
				),
			},
		},
	})
}
