package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/version"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PIHOLE_PASSWORD", nil),
				Description: "The admin password used to login to the admin dashboard",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PIHOLE_URL", "http://pi.hole"),
				Description: "URL where Pi-hole is deployed",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"pihole_dns_records":   dataSourceDNSRecords(),
			"pihole_cname_records": dataSourceCNAMERecords(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"pihole_dns_record":        resourceDNSRecord(),
			"pihole_ad_blocker_status": resourceAdBlockerStatus(),
		},
	}

	provider.ConfigureContextFunc = configure(version.ProviderVersion, provider)

	return provider
}

// configure configures a Pi-hole client to be used for terraform resource requests
func configure(version string, provider *schema.Provider) func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (client interface{}, diags diag.Diagnostics) {
		client, err := Config{
			Password:  d.Get("password").(string),
			URL:       d.Get("url").(string),
			UserAgent: provider.UserAgent("terraform-provider-pihole", version),
		}.Client(ctx)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
