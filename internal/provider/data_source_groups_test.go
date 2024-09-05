package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccGroupsData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "pihole_group" "enabled_test_group" {
					  name        = "enabled_test_group"
					  description = "Sample description"
					  enabled     = true
					}

					data "pihole_groups" "groups" {
					  depends_on = [pihole_group.enabled_test_group]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.#", "2"),

					resource.TestCheckResourceAttrSet("data.pihole_groups.groups", "groups.0.id"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.0.name", "enabled_test_group"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.0.description", "Sample description"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.0.enabled", "true"),

					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.id", "0"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.name", "Default"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.description", "The default group"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.enabled", "true"),
				),
			},
			{
				Config: `
					resource "pihole_group" "enabled_test_group" {
					  name        = "enabled_test_group"
					  description = "Sample description"
					  enabled 	  = true
					}

					resource "pihole_group" "disabled_test_group" {
					  name        = "disabled_test_group"
					  description = "Sample description"
					  enabled 	  = false
					}

					data "pihole_groups" "groups" {
					  depends_on = [pihole_group.enabled_test_group, pihole_group.disabled_test_group]
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.#", "3"),

					resource.TestCheckResourceAttrSet("data.pihole_groups.groups", "groups.0.id"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.0.name", "disabled_test_group"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.0.description", "Sample description"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.0.enabled", "false"),

					resource.TestCheckResourceAttrSet("data.pihole_groups.groups", "groups.1.id"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.name", "enabled_test_group"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.description", "Sample description"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.1.enabled", "true"),

					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.2.id", "0"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.2.name", "Default"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.2.description", "The default group"),
					resource.TestCheckResourceAttr("data.pihole_groups.groups", "groups.2.enabled", "true"),
				),
			},
		},
	})
}
