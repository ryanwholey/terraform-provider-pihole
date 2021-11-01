package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// dataSourceCNAMERecords returns a schema resource for listing Pi-hole CNAME records
func dataSourceCNAMERecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCNAMERecordsRead,
		Schema: map[string]*schema.Schema{
			"records": {
				Description: "List of CNAME Pi-hole records",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Description: "CNAME record domain",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"target": {
							Description: "CNAME target value where traffic is routed to from the domain",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceCNAMERecordsRead lists all Pi-hole CNAME records
func dataSourceCNAMERecordsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	cnameList, err := client.ListCNAMERecords(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	list := make([]map[string]interface{}, len(cnameList))
	idRef := ""

	for i, r := range cnameList {
		idRef = fmt.Sprintf("%s%s%s", idRef, r.Domain, r.Target)

		list[i] = map[string]interface{}{
			"domain": r.Domain,
			"target": r.Target,
		}
	}

	if err := d.Set("records", list); err != nil {
		return diag.FromErr(err)
	}

	hash := sha256.Sum256([]byte(idRef))
	d.SetId(fmt.Sprintf("%x", hash[:]))

	return diags
}
