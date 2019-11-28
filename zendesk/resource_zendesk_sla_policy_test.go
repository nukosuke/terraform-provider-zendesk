package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

func TestMarshalSlaPolicy(t *testing.T) {
	expected := zendesk.SlaPolicy{
		Title:       "title",
		Description: "blabla",
		Active:      true,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalSlaPolicy(expected, m)
	if err != nil {
		t.Fatalf("Failed to marshal map %v", err)
	}

	v, ok := m.GetOk("title")
	if !ok {
		t.Fatal("Failed to get title value")
	}
	if v != expected.Title {
		t.Fatalf("sla policy had incorrect title value %v. should have been %v", v, expected.Title)
	}

	v, ok = m.GetOk("description")
	if !ok {
		t.Fatal("Failed to get description value")
	}
	if v != expected.Description {
		t.Fatalf("sla policy had incorrect description value %v. should have been %v", v, expected.Description)
	}

	v, ok = m.GetOk("active")
	if !ok {
		t.Fatal("Failed to get active value")
	}
	if v != expected.Active {
		t.Fatalf("sla policy had incorrect active value %v. should have been %v", v, expected.Active)
	}
}

func TestUnmarshalSlaPolicy(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "100",
		mapGetterSetter: mapGetterSetter{
			"title":       "Auto reply",
			"description": "reply automatically",
			"active":      true,
		},
	}

	sla, err := unmarshalSlaPolicy(m)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if v := m.Get("title"); sla.Title != v {
		t.Fatalf("sla policy had title value %v. should have been %v", sla.Title, v)
	}
}

func TestCreateSlaPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	out := zendesk.SlaPolicy{
		ID:    12345,
		Title: "sla policy",
	}

	m.EXPECT().CreateSlaPolicy(gomock.Any(), gomock.Any()).Return(out, nil)
	if err := createSlaPolicy(i, m); err != nil {
		t.Fatal("CreateSlaPolicy return an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("CreateSlaPolicy did not set resource id. Id was %s", v)
	}

	if v := i.Get("title"); v != "sla policy" {
		t.Fatalf("CreateSlaPolicy did not set resource title. title was %s", v)
	}
}

func TestReadSlaPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := zendesk.SlaPolicy{
		Title:  "sla policy",
		Active: true,
	}
	m.EXPECT().GetSlaPolicy(gomock.Any(), gomock.Eq(int64(12345))).Return(expected, nil)
	if err := readSlaPolicy(i, m); err != nil {
		t.Fatalf("GetSlaPolicy received an error when calling: %v", err)
	}
}

func TestUpdateSlaPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: mapGetterSetter{},
	}

	m.EXPECT().UpdateSlaPolicy(gomock.Any(), gomock.Eq(int64(12345)), gomock.Any()).Return(zendesk.SlaPolicy{}, nil)
	if err := updateSlaPolicy(i, m); err != nil {
		t.Fatalf("updateSlaPolicy returned an error %v", err)
	}
}

func TestDeleteSlaPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteSlaPolicy(gomock.Any(), gomock.Eq(int64(1234))).Return(nil)
	err := deleteSlaPolicy(d, c)
	if err != nil {
		t.Fatalf("Got error from resource delete: %v", err)
	}
}

func testSlaPolicyDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.SlaPolicyAPI)

	for k, r := range s.RootModule().Resources {
		if r.Type != "zendesk_sla_policy" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		ctx := context.Background()
		_, err = client.GetSlaPolicy(ctx, id)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed sla policy. resource name %s", k)
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

func TestAccSlaPolicyExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testSlaPolicyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "sla_policies.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_sla_policy.auto-reply-sla_policy", "title", "Auto Reply SlaPolicy"),
					resource.TestCheckResourceAttr("zendesk_sla_policy.auto-reply-sla_policy", "active", "true"),
					resource.TestCheckResourceAttrSet("zendesk_sla_policy.auto-reply-sla_policy", "all.#"),
					resource.TestCheckResourceAttrSet("zendesk_sla_policy.auto-reply-sla_policy", "action.#"),
				),
			},
		},
	})
}
