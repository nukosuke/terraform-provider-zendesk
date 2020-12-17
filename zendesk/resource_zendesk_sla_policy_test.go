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

func TestMarshalSLAPolicy(t *testing.T) {
	expected := zendesk.SLAPolicy{
		Title:       "title",
		Description: "blabla",
		Active:      true,
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalSLAPolicy(expected, m)
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

func TestUnmarshalSLAPolicy(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "100",
		mapGetterSetter: mapGetterSetter{
			"title":       "Auto reply",
			"description": "reply automatically",
			"active":      true,
		},
	}

	sla, err := unmarshalSLAPolicy(m)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if v := m.Get("title"); sla.Title != v {
		t.Fatalf("sla policy had title value %v. should have been %v", sla.Title, v)
	}
}

func TestCreateSLAPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	out := zendesk.SLAPolicy{
		ID:    12345,
		Title: "sla policy",
	}

	m.EXPECT().CreateSLAPolicy(gomock.Any(), gomock.Any()).Return(out, nil)
	if err := createSLAPolicy(i, m); err != nil {
		t.Fatal("CreateSLAPolicy return an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("CreateSLAPolicy did not set resource id. Id was %s", v)
	}

	if v := i.Get("title"); v != "sla policy" {
		t.Fatalf("CreateSLAPolicy did not set resource title. title was %s", v)
	}
}

func TestReadSLAPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := zendesk.SLAPolicy{
		Title:  "sla policy",
		Active: true,
	}
	m.EXPECT().GetSLAPolicy(gomock.Any(), gomock.Eq(int64(12345))).Return(expected, nil)
	if err := readSLAPolicy(i, m); err != nil {
		t.Fatalf("GetSLAPolicy received an error when calling: %v", err)
	}

	active := i.Get("active").(bool)
	if !active {
		t.Fatal("Did not set active field")
	}
}

func TestUpdateSLAPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: mapGetterSetter{},
	}

	m.EXPECT().UpdateSLAPolicy(gomock.Any(), gomock.Eq(int64(12345)), gomock.Any()).Return(zendesk.SLAPolicy{}, nil)
	if err := updateSLAPolicy(i, m); err != nil {
		t.Fatalf("updateSLAPolicy returned an error %v", err)
	}
}

func TestDeleteSLAPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteSLAPolicy(gomock.Any(), gomock.Eq(int64(1234))).Return(nil)
	err := deleteSLAPolicy(d, c)
	if err != nil {
		t.Fatalf("Got error from resource delete: %v", err)
	}
}

func testSLAPolicyDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.SLAPolicyAPI)

	for k, r := range s.RootModule().Resources {
		if r.Type != "zendesk_sla_policy" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		ctx := context.Background()
		_, err = client.GetSLAPolicy(ctx, id)
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

func TestAccSLAPolicyExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testSLAPolicyDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "sla_policies.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_sla_policy.incidents_sla_policy", "title", "Incidents"),
					resource.TestCheckResourceAttr("zendesk_sla_policy.incidents_sla_policy", "active", "true"),
					resource.TestCheckResourceAttrSet("zendesk_sla_policy.incidents_sla_policy", "all.#"),
				),
			},
		},
	})
}
