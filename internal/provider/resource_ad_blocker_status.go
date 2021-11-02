package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// resourceAdBlockerStatus returns the DNS Terraform resource management configuration
func resourceAdBlockerStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdBlockerStatusCreate,
		ReadContext:   resourceAdBlockerStatusRead,
		UpdateContext: resourceAdBlockerStatusUpdate,
		DeleteContext: resourceAdBlockerStatusDelete,
		Schema: map[string]*schema.Schema{
			"enabled": {
				Description: "Whether to enable the Pi-hole ad blocker",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

// rresourceAdBlockerStatusCreate handles the creation a DNS record via Terraform
func resourceAdBlockerStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	_, err := client.SetAdBlockEnabled(ctx, d.Get("enabled").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ad-block-enabled")

	return diags
}

// resourceAdBlockerStatusRead finds a DNS record based on the associated domain ID
func resourceAdBlockerStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	res, err := client.GetAdBlockerStatus(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", res.Enabled); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceAdBlockerStatusUpdate handles updates of a DNS record via Terraform
func resourceAdBlockerStatusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	res, err := client.SetAdBlockEnabled(ctx, d.Get("enabled").(bool))
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("enabled", res.Enabled); err != nil {
		return diag.FromErr(err)
	}

	return resourceAdBlockerStatusRead(ctx, d, meta)
}

// resourceAdBlockerStatusDelete handles the deletion of a DNS record via Terraform
func resourceAdBlockerStatusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	d.SetId("")

	return diags
}
