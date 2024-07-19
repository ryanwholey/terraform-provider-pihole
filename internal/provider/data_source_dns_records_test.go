package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDNSRecordsData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "pihole_dns_record" "record" {
					  domain = "foo.com"
					  ip     = "127.0.0.1"
					}

					data "pihole_dns_records" "records" {
					  depends_on = [pihole_dns_record.record]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pihole_dns_records.records", "records.#", "1"),

					resource.TestCheckResourceAttr("data.pihole_dns_records.records", "records.0.domain", "foo.com"),
					resource.TestCheckResourceAttr("data.pihole_dns_records.records", "records.0.ip", "127.0.0.1"),
				),
			},
		},
	})
}
