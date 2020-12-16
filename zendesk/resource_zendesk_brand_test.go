package zendesk

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	. "github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nukosuke/go-zendesk/zendesk"
	"github.com/nukosuke/go-zendesk/zendesk/mock"
)

var testBrand = zendesk.Brand{
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

const brandConfig = `
resource "zendesk_brand" "acc_brand" {
  name            = "T-%d"
  active          = true
  subdomain       = "d3v-terraform-provider-t%d"
}
`

func TestCreateBrand(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)

	m.EXPECT().CreateBrand(Any(), Any()).Return(testBrand, nil)

	i := newIdentifiableGetterSetter()
	err := createBrand(i, m)
	if err != nil {
		t.Fatalf("Create brand returned an error %v", err)
	}

	if i.Id() != fmt.Sprintf("%d", testBrand.ID) {
		t.Fatalf("Created object does not have the correct brand id. Was: %s. Expected %d", i.Id(), testBrand.ID)
	}

	if i.Get("logo_attachment_id") != testBrand.Logo.ID {
		t.Fatalf("Created object does not have the correct logo id. Was: %d. Expected %d", i.Get("logo_attachment_id"), testBrand.Logo.ID)
	}
}

func TestReadBrand(t *testing.T) {
	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	m.EXPECT().GetBrand(Any(), testBrand.ID).Return(testBrand, nil)
	i := newIdentifiableGetterSetter()
	i.SetId(fmt.Sprintf("%d", testBrand.ID))

	err := readBrand(i, m)
	if err != nil {
		t.Fatalf("readBrand returned an error: %v", err)
	}

	if v := i.Get("subdomain"); v != testBrand.Subdomain {
		t.Fatalf("Subdomain was not set to the expected value. Was: %s Expected %s", v, testBrand.Subdomain)
	}
}

func TestUpdateBrand(t *testing.T) {
	updatedBrand := testBrand
	updatedBrand.Name = "1234"

	i := newIdentifiableGetterSetter()
	i.SetId(fmt.Sprintf("%d", testBrand.ID))

	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	m.EXPECT().UpdateBrand(Any(), testBrand.ID, Any()).Return(updatedBrand, nil)

	err := updateBrand(i, m)
	if err != nil {
		t.Fatalf("update brand returned an error: %v", err)
	}

	if v := i.Get("name"); v != updatedBrand.Name {
		t.Fatalf("Update did not set name to the expected value. Was %s expected %s", v, updatedBrand.Name)
	}
}

func TestDeleteBrand(t *testing.T) {
	id := int64(1234)
	i := newIdentifiableGetterSetter()
	i.SetId(fmt.Sprintf("%d", id))

	ctrl := NewController(t)
	defer ctrl.Finish()

	m := mock.NewClient(ctrl)
	m.EXPECT().DeleteBrand(Any(), id).Return(nil)

	err := deleteBrand(i, m)
	if err != nil {
		t.Fatalf("delete brand returned an error: %v", err)
	}
}

func testBrandDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(zendesk.BrandAPI)

	for _, r := range s.RootModule().Resources {
		if r.Type != "zendesk_brand" {
			continue
		}

		id, err := atoi64(r.Primary.ID)
		if err != nil {
			return err
		}

		brand, err := client.GetBrand(context.Background(), id)
		if err != nil {
			zd, ok := err.(zendesk.Error)
			if !ok {
				return fmt.Errorf("error %v cannot be asserted as a zendesk error", err)
			}

			if zd.Status() != http.StatusNotFound {
				return fmt.Errorf("did not get a not found error after destroy. error was %v", zd)
			}
		} else {
			if brand.Active {
				return fmt.Errorf("brand named %s is still active", brand.Name)
			}
		}
	}

	return nil
}

func TestAccZendeskBrand(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Int31n(199) + 801
	config := fmt.Sprintf(brandConfig, n, n)

	expectedSubdomain := fmt.Sprintf("d3v-terraform-provider-t%d", n)
	expectedName := fmt.Sprintf("T-%d", n)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testBrandDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zendesk_brand.acc_brand", "name", expectedName),
					resource.TestCheckResourceAttr("zendesk_brand.acc_brand", "subdomain", expectedSubdomain),
				),
			},
		},
	})
}
