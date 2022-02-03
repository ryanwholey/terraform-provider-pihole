package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PIHOLE_URL"); v == "" {
		t.Fatal("PIHOLE_URL must be set for acceptance tests")
	}

	if v := os.Getenv("PIHOLE_PASSWORD"); v == "" {
		t.Fatal("PIHOLE_PASSWORD must be set for acceptance tests")
	}
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"pihole": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = Provider()
}
