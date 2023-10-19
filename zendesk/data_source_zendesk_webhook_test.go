package zendesk

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

const webhookConfig = `
data "zendesk_webhook" "test_webhook" {
	id = "%s"
}
`

func TestWebhookDataSourceRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	c := mock.NewClient(ctrl)

	m := newIdentifiableGetterSetter()
	id := "1234"

	err := m.Set("id", id)
	if err != nil {
		t.Fatalf("Read webhook field returned an error. %v", err)
	}

	createdAt := time.Now().Add(-1 * time.Hour)
	updatedAt := createdAt.Add(30 * time.Minute)

	out := &zendesk.Webhook{
		Authentication: &zendesk.WebhookAuthentication{
			Type:        "basic_auth",
			AddPosition: "header",
			Data: map[string]any{
				"username": "user",
				"password": "xxxx",
			},
		},
		CreatedAt:     createdAt,
		CreatedBy:     "111",
		Description:   "test webhook",
		Endpoint:      "http://example.com/status/200",
		HTTPMethod:    http.MethodPost,
		ID:            "1234",
		Name:          "my test webhook",
		RequestFormat: "json",
		Status:        "active",
		Subscriptions: []string{"conditional_ticket_events"},
		UpdatedAt:     updatedAt,
		UpdatedBy:     "999",
	}

	c.EXPECT().GetWebhook(gomock.Any(), id).Return(out, nil)

	diags := readWebhookDataSource(context.Background(), m, c)
	if len(diags) != 0 {
		t.Fatalf("Read system field returned an error. %v", diags)
	}

	if v := m.Id(); v != out.ID {
		t.Fatalf("Read test_webhook did not set ID field. Expected %v, Got %v", out.ID, v)
	}

	if v, ok := m.GetOk("created_at"); !ok || v.(string) != createdAt.String() {
		t.Fatalf("Read test_webhook did not set CreatedAt field. Expected %v, Got %v", out.CreatedAt, v)
	}
}

func TestAccWebhookDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(webhookConfig, "1234"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.zendesk_webhook.test_webhook", "id", "1234"),
					resource.TestCheckResourceAttrSet("data.zendesk_webhook.test_webhook", "created_at"),
				),
			},
		},
	})
}
