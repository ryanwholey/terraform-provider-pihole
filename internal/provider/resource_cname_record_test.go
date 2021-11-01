package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// TestAccCNAMERecord acceptance test for the CNAME record resource
func TestAccCNAMERecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCNAMERecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testLocalCNAMEResourceConfig("foo", "foo.com", "bar.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_cname_record.foo", "domain", "foo.com"),
					resource.TestCheckResourceAttr("pihole_cname_record.foo", "target", "bar.com"),
					testCheckLocalCNAMEResourceExists(t, "foo.com", "bar.com"),
				),
			},
			{
				Config: testLocalCNAMEResourceConfig("foo", "foo.com", "woz.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_cname_record.foo", "domain", "foo.com"),
					resource.TestCheckResourceAttr("pihole_cname_record.foo", "target", "woz.com"),
					testCheckLocalCNAMEResourceExists(t, "foo.com", "woz.com"),
				),
			},
		},
	})
}

// testLocalCNAMEResourceConfig returns HCL to configure a CNAME record
func testLocalCNAMEResourceConfig(name string, domain string, target string) string {
	return fmt.Sprintf(`
		resource "pihole_cname_record" %q {
			domain = %q
			target = %q
		}	
	`, name, domain, target)
}

// testCheckLocalCNAMEResourceExists checks that the CNAME record exists in Pi-hole
func testCheckLocalCNAMEResourceExists(t *testing.T, domain string, target string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client := testAccProvider.Meta().(*pihole.Client)

		record, err := client.GetCNAMERecord(context.Background(), domain)
		if err != nil {
			return err
		}

		if record.Target != target {
			return fmt.Errorf("requested %s:%s does not match: %s", domain, target, record.Target)
		}

		return nil
	}
}

// testAccCheckCNAMERecordDestroy checks that all resources have been deleted
func testAccCheckCNAMERecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pihole.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "pihole_cname_record" {
			continue
		}

		if _, err := client.GetCNAMERecord(context.Background(), r.Primary.ID); err != nil {
			if _, ok := err.(*pihole.NotFoundError); !ok {
				return err
			}
		}
	}
	return nil
}
