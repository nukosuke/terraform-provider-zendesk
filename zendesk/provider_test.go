package zendesk

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"zendesk": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(emailVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", emailVar)
	}
	if v := os.Getenv(tokenVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", tokenVar)
	}
	if v := os.Getenv(accountVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", accountVar)
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
