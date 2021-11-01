package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// dataSourceDNSRecords returns a schema resource for listing pihole local DNS records
func dataSourceDNSRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSRecordsRead,
		Schema: map[string]*schema.Schema{
			"records": {
				Description: "List of Pi-hole DNS records",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Description: "DNS record domain",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"ip": {
							Description: "IP address where traffic is routed to from the DNS record domain",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceDNSRecordsRead lists all pihole local DNS records
func dataSourceDNSRecordsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	dnsList, err := client.ListDNSRecords(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	list := make([]map[string]interface{}, len(dnsList))
	idRef := ""

	for i, r := range dnsList {
		idRef = fmt.Sprintf("%s%s%s", idRef, r.Domain, r.IP)

		list[i] = map[string]interface{}{
			"domain": r.Domain,
			"ip":     r.IP,
		}
	}

	if err := d.Set("records", list); err != nil {
		return diag.FromErr(err)
	}

	hash := sha256.Sum256([]byte(idRef))
	d.SetId(fmt.Sprintf("%x", hash[:]))

	return diags
}
