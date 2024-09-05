package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccCNAMERecordsData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "pihole_cname_record" "record" {
					  domain = "foo.com"
					  target = "bar.com"
					}

					data "pihole_cname_records" "records" {
					  depends_on = [pihole_cname_record.record]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pihole_cname_records.records", "records.#", "1"),

					resource.TestCheckResourceAttr("data.pihole_cname_records.records", "records.0.domain", "foo.com"),
					resource.TestCheckResourceAttr("data.pihole_cname_records.records", "records.0.target", "bar.com"),
				),
			},
		},
	})
}
