package provider

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// dataSourceGroups returns all Pi-hole groups, including the default group
func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupsRead,
		Schema: map[string]*schema.Schema{
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Group ID",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"enabled": {
							Description: "Whether the group is enabled",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"name": {
							Description: "Name of the group",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "Group description",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceGroupsRead returns all Pi-hole groups
func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	groupList, err := client.ListGroups(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	list := make([]map[string]interface{}, len(groupList))

	for i, g := range groupList {
		list[i] = map[string]interface{}{
			"id":          g.ID,
			"enabled":     g.Enabled,
			"name":        g.Name,
			"description": g.Description,
		}
	}
	listString, err := json.Marshal(list)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("groups", list); err != nil {
		return diag.FromErr(err)
	}

	hash := sha256.Sum256([]byte(listString))
	d.SetId(fmt.Sprintf("%x", hash[:]))

	return diags
}
