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

func dataSourceDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainsRead,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Description: "Filter on allowed or denied domains. Must be either 'allow' or 'deny'.",
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					domainType := val.(string)

					if domainType != pihole.DomainOptionsAllow && domainType != pihole.DomainOptionsDeny {
						errs = append(errs, fmt.Errorf("%s field must be one of %v: %q", key, domainType, []string{pihole.DomainOptionsAllow, pihole.DomainOptionsDeny}))
					}

					return
				},
			},
			"domains": {
				Type:        schema.TypeSet,
				Description: "Domains ",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Domain ID",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"type": {
							Description: "Whether the doamin is on the allow or deny list",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"enabled": {
							Description: "Whether the domain rule is enabled",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"domain": {
							Description: "Domain",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"comment": {
							Description: "Comments associated with the domain",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"wildcard": {
							Description: "Whether the domain should be interpreted using a wildcard parser",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"group_ids": {
							Description: "Groups to which the domain is associated",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

// dataSourceDomainsRead returns all Pi-hole domains
func dataSourceDomainsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	opts := pihole.ListDomainsOptions{}

	domainType := d.Get("type").(string)

	if domainType != "" {
		opts.Type = domainType
	}

	domainList, err := client.ListDomains(ctx, opts)
	if err != nil {
		return diag.FromErr(err)
	}

	list := make([]map[string]interface{}, len(domainList))

	for i, d := range domainList {
		list[i] = map[string]interface{}{
			"id":        d.ID,
			"type":      d.Type,
			"enabled":   d.Enabled,
			"domain":    d.Domain,
			"comment":   d.Comment,
			"wildcard":  d.Wildcard,
			"group_ids": d.GroupIDs,
		}
	}

	listString, err := json.Marshal(list)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("domains", list); err != nil {
		return diag.FromErr(err)
	}

	hash := sha256.Sum256([]byte(listString))
	d.SetId(fmt.Sprintf("%x", hash[:]))

	return diags
}
