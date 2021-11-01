package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

func TestAccLocalDNS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLocalDNSDestroy,
		Steps: []resource.TestStep{
			{
				Config: testLocalDNSResourceConfig("foo", "foo.com", "127.0.0.1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_record.foo", "domain", "foo.com"),
					resource.TestCheckResourceAttr("pihole_dns_record.foo", "ip", "127.0.0.1"),
					testCheckLocalDNSResourceExists(t, "foo.com", "127.0.0.1"),
				),
			},
			{
				Config: testLocalDNSResourceConfig("foo", "foo.com", "127.0.0.2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_dns_record.foo", "domain", "foo.com"),
					resource.TestCheckResourceAttr("pihole_dns_record.foo", "ip", "127.0.0.2"),
					testCheckLocalDNSResourceExists(t, "foo.com", "127.0.0.2"),
				),
			},
		},
	})
}

func testLocalDNSResourceConfig(name string, domain string, ip string) string {
	return fmt.Sprintf(`
		resource "pihole_dns_record" %q {
			domain = %q
			ip     = %q
		}	
	`, name, domain, ip)
}

func testCheckLocalDNSResourceExists(t *testing.T, domain string, ip string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client := testAccProvider.Meta().(*pihole.Client)

		record, err := client.GetDNSRecord(context.Background(), domain)
		if err != nil {
			return err
		}

		if record.IP != ip {
			return fmt.Errorf("requested %s:%s does not match IP: %s", domain, ip, record.IP)
		}

		return nil
	}
}

func testAccCheckLocalDNSDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pihole.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "pihole_local_dns" {
			continue
		}

		if _, err := client.GetDNSRecord(context.Background(), r.Primary.ID); err != nil {
			if err.Error() != fmt.Sprintf("record %q not found", r.Primary.ID) {
				return err
			}
		}
	}
	return nil
}
