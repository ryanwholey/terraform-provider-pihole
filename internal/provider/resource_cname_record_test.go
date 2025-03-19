package provider

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	pihole "github.com/ryanwholey/go-pihole"
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
			{
				Config: testLocalCNAMEResourceWithDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pihole_cname_records.records", "records.#", "20"),

					resource.TestCheckResourceAttr("data.pihole_cname_records.records", "records.0.domain", "aa.com"),
					resource.TestCheckResourceAttr("data.pihole_cname_records.records", "records.0.target", "ingress.example.local"),
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

func testLocalCNAMEResourceWithDataConfig() string {
	return `
		locals {
		  all_cnames = [
			"aa.com",
			"bb.com",
			"cc.com",
			"dd.com",
			"ee.com",
			"ff.com",
			"gg.com",
			"hh.com",
			"ii.com",
			"jj.com",
			"kk.com",
			"ll.com",
			"mm.com",
			"nn.com",
			"oo.com",
			"pp.com",
			"qq.com",
			"rr.com",
			"ss.com",
			"tt.com",
		  ]
		}

		resource "pihole_cname_record" "cname_records" {
		  count  = length(local.all_cnames)
		  domain = local.all_cnames[count.index]
		  target = "ingress.example.local"
		}

		data "pihole_cname_records" "records" {
		  depends_on = [pihole_cname_record.cname_records]
		}
    `
}

// testCheckLocalCNAMEResourceExists checks that the CNAME record exists in Pi-hole
func testCheckLocalCNAMEResourceExists(t *testing.T, domain string, target string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client := testAccProvider.Meta().(*pihole.Client)

		record, err := client.LocalCNAME.Get(context.Background(), domain)
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

		if _, err := client.LocalCNAME.Get(context.Background(), r.Primary.ID); err != nil {
			if errors.Is(err, pihole.ErrorLocalCNAMENotFound) {
				return err
			}
		}
	}
	return nil
}
