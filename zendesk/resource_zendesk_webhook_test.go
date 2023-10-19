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

func TestMarshalWebhook(t *testing.T) {
	expected := zendesk.Webhook{
		Name:          "title",
		Description:   "blabla",
		Status:        "active",
		Subscriptions: []string{"conditional_ticket_events"},
		RequestFormat: "json",
		Endpoint:      "example.org",
		HTTPMethod:    http.MethodPost,
		Authentication: &zendesk.WebhookAuthentication{
			Type: "basic_auth",
			Data: map[string]any{
				"username": "tester",
				"password": "foobar",
			},
		},
	}
	m := &identifiableMapGetterSetter{
		mapGetterSetter: mapGetterSetter{},
	}

	err := marshalWebhook(&expected, m)
	if err != nil {
		t.Fatalf("Failed to marshal map %v", err)
	}

	v, ok := m.GetOk("name")
	if !ok {
		t.Fatal("Failed to get title value")
	}
	if v != expected.Name {
		t.Fatalf("webhook had incorrect name value %v. should have been %v", v, expected.Name)
	}

	v, ok = m.GetOk("description")
	if !ok {
		t.Fatal("Failed to get description value")
	}
	if v != expected.Description {
		t.Fatalf("webhook had incorrect description value %v. should have been %v", v, expected.Description)
	}

	v, ok = m.GetOk("status")
	if !ok {
		t.Fatal("Failed to get status value")
	}
	if v != expected.Status {
		t.Fatalf("webhook had incorrect active value %v. should have been %v", v, expected.Status)
	}

	v, ok = m.GetOk("subscriptions")
	if !ok {
		t.Fatal("Failed to get subscriptions value")
	}
}

func TestUnmarshalWebhook(t *testing.T) {
	m := &identifiableMapGetterSetter{
		id: "100",
		mapGetterSetter: mapGetterSetter{
			"name":           "No auth",
			"description":    "Webhook without auth",
			"status":         "active",
			"http_method":    http.MethodPost,
			"endpoint":       "http://example.org",
			"request_format": "json",
		},
	}

	whk, err := unmarshalWebhook(m)
	if err != nil {
		t.Fatalf("unmarshal returned an error: %v", err)
	}

	if v := m.Get("name"); whk.Name != v {
		t.Fatalf("webhook had name value %v. should have been %v", whk.Name, v)
	}
}

func TestCreateWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	out := &zendesk.Webhook{
		ID:   "12345",
		Name: "webhook",
	}

	m.EXPECT().CreateWebhook(gomock.Any(), gomock.Any()).Return(out, nil)
	if diags := createWebhook(context.Background(), i, m); len(diags) != 0 {
		t.Fatal("CreateWebhook return an error")
	}

	if v := i.Id(); v != "12345" {
		t.Fatalf("CreateWebhook did not set resource id. Id was %q", v)
	}

	if v := i.Get("name"); v != "webhook" {
		t.Fatalf("CreateWebhook did not set resource name. name was %q", v)
	}
}

func TestReadWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := mock.NewClient(ctrl)
	i := newIdentifiableGetterSetter()
	i.SetId("12345")

	expected := &zendesk.Webhook{
		ID:     "12345",
		Name:   "webhook",
		Status: "active",
	}
	m.EXPECT().GetWebhook(gomock.Any(), gomock.Eq("12345")).Return(expected, nil)
	if diags := readWebhook(context.Background(), i, m); len(diags) != 0 {
		t.Fatalf("GetWebhook received an error when calling: %v", diags)
	}
}

func TestUpdateWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := mock.NewClient(ctrl)
	i := &identifiableMapGetterSetter{
		id:              "12345",
		mapGetterSetter: mapGetterSetter{},
	}

	m.EXPECT().UpdateWebhook(gomock.Any(), gomock.Eq("12345"), gomock.Any()).Return(nil)
	if diags := updateWebhook(context.Background(), i, m); len(diags) != 0 {
		t.Fatalf("updateWebhook returned an error %v", diags)
	}
}

func TestDeleteWebhook(t *testing.T) {
	ctrl := gomock.NewController(t)

	c := mock.NewClient(ctrl)
	d := newIdentifiableGetterSetter()

	d.SetId("1234")

	c.EXPECT().DeleteWebhook(gomock.Any(), gomock.Eq("1234")).Return(nil)
	diags := deleteWebhook(context.Background(), d, c)
	if len(diags) != 0 {
		t.Fatalf("Got error from resource delete: %v", diags)
	}
}

func testWebhookDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.WebhookAPI)

	for k, r := range s.RootModule().Resources {
		if r.Type != "zendesk_webhook" {
			continue
		}

		ctx := context.Background()
		_, err := client.GetWebhook(ctx, r.Primary.ID)
		if err == nil {
			return fmt.Errorf("did not get error from zendesk when trying to fetch the destroyed webhook. resource name %s", k)
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

func TestAccWebhookExample(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testWebhookDestroyed,
		Steps: []resource.TestStep{
			{
				Config: readExampleConfig(t, "resources/zendesk_webhook/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_webhook.example-webhook", "name", "Example Webhook without authentication"),
					resource.TestCheckResourceAttr("zendesk_webhook.example-webhook", "endpoint", "https://example.com/status/200"),
					resource.TestCheckResourceAttr("zendesk_webhook.example-webhook", "http_method", "true"),
					resource.TestCheckResourceAttr("zendesk_webhook.example-webhook", "request_format", "json"),
					resource.TestCheckResourceAttr("zendesk_webhook.example-webhook", "status", "active"),
					resource.TestCheckResourceAttrSet("zendesk_webhook.example-webhook", "subscriptions.#"),
				),
			},
		},
	})
}
