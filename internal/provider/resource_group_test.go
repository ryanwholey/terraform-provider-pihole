package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

func TestAccGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testGroupResourceConfig("foo", "mygroup", "description", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_group.foo", "name", "mygroup"),
					resource.TestCheckResourceAttr("pihole_group.foo", "description", "description"),
					resource.TestCheckResourceAttr("pihole_group.foo", "enabled", "true"),
					testCheckGroupResourceExists(t, "mygroup", "description", true),
				),
			},
			{
				Config: testGroupResourceConfig("foo", "mygroup", "updated", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("pihole_group.foo", "name", "mygroup"),
					resource.TestCheckResourceAttr("pihole_group.foo", "description", "updated"),
					resource.TestCheckResourceAttr("pihole_group.foo", "enabled", "false"),
					testCheckGroupResourceExists(t, "mygroup", "updated", false),
				),
			},
		},
	})
}

func testGroupResourceConfig(resourceName, name, description string, enabled bool) string {
	return fmt.Sprintf(`
		resource "pihole_group" %q {
			name        = %q
			description = %q
			enabled     = %v
		}	
	`, resourceName, name, description, enabled)
}

func testCheckGroupResourceExists(t *testing.T, name string, description string, enabled bool) resource.TestCheckFunc {
	return func(*terraform.State) error {
		client := testAccProvider.Meta().(*pihole.Client)

		group, err := client.GetGroup(context.Background(), name)
		if err != nil {
			return err
		}

		if group.Description != description {
			return fmt.Errorf("requested group %s:%s does not match description: %s", name, description, group.Description)
		}

		if group.Enabled != enabled {
			return fmt.Errorf("requested group %s:%v does not match enabled value: %v", name, enabled, group.Enabled)
		}

		return nil
	}
}

func testAccCheckGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*pihole.Client)

	for _, r := range s.RootModule().Resources {
		if r.Type != "pihole_group" {
			continue
		}

		name, ok := r.Primary.Attributes["name"]
		if !ok {
			return fmt.Errorf("group name not found on primary resource")
		}

		if _, err := client.GetGroup(context.Background(), name); err != nil {
			if _, ok := err.(*pihole.NotFoundError); !ok {
				return err
			}
		}
	}

	return nil
}
