package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PIHOLE_PASSWORD", nil),
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PIHOLE_URL", nil),
				Default:     "http://pi.hole",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"pihole_dns_records": dataSourceDNSRecords(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"pihole_dns_record":        resourceDNSRecord(),
			"pihole_ad_blocker_status": resourceAdBlockerStatus(),
		},
	}

	provider.ConfigureContextFunc = providerConfigure

	return provider
}

// providerConfigure configures a pihole client to be used in terraform resource requests
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	client, err := Config{
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
	}.Client(ctx)

	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}
