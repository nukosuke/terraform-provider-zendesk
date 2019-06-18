package zendesk

import (
	"fmt"
	. "github.com/golang/mock/gomock"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
	"testing"
)

func TestCreateBrand(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)

	b := zendesk.Brand{
		ID:              47,
		URL:             "https://company.zendesk.com/api/v2/brands/47.json",
		Name:            "Brand 1",
		BrandURL:        "https://brand1.com",
		HasHelpCenter:   true,
		HelpCenterState: "enabled",
		Active:          true,
		Default:         true,
		Logo: zendesk.Attachment{
			ID:          928374,
			FileName:    "brand1_logo.png",
			ContentURL:  "https://company.zendesk.com/logos/brand1_logo.png",
			ContentType: "image/png",
			Size:        166144,
		},
		Subdomain:         "brand1",
		HostMapping:       "brand1.com",
		SignatureTemplate: "{{agent.signature}}",
	}

	m.EXPECT().CreateBrand(Any()).Return(b, nil)

	i := newIdentifiableGetterSetter()
	err := createBrand(i, m)
	if err != nil {
		t.Fatalf("Create brand returned an error %v", err)
	}

	if i.Id() != fmt.Sprintf("%d", b.ID) {
		t.Fatalf("Created object does not have the correct brand id. Was: %s. Expected %d", i.Id(), b.ID)
	}

	if i.Get("logo_attachment_id") != b.Logo.ID {
		t.Fatalf("Created object does not have the correct logo id. Was: %d. Expected %d", i.Get("logo_attachment_id"), b.Logo.ID)
	}

}
