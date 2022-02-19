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
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("PIHOLE_PASSWORD", nil),
				Description:  "The admin password used to login to the admin dashboard. Conflicts with `api_token`.",
				ExactlyOneOf: []string{"api_token", "password"},
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PIHOLE_URL", "http://pi.hole"),
				Description: "URL where Pi-hole is deployed",
			},
			"api_token": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("PIHOLE_API_TOKEN", nil),
				Description:  "Experimental: Pi-hole API token. Conflicts with `password`.",
				ExactlyOneOf: []string{"api_token", "password"},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"pihole_cname_records": dataSourceCNAMERecords(),
			"pihole_dns_records":   dataSourceDNSRecords(),
			"pihole_domains":       dataSourceDomains(),
			"pihole_groups":        dataSourceGroups(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"pihole_ad_blocker_status": resourceAdBlockerStatus(),
			"pihole_cname_record":      resourceCNAMERecord(),
			"pihole_dns_record":        resourceDNSRecord(),
			"pihole_group":             resourceGroup(),
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
			APIToken:  d.Get("api_token").(string),
		}.Client(ctx)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return client, diags
	}
}
