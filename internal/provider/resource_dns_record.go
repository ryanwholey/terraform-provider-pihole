package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// resourceDNSRecord returns the local DNS Terraform resource management configuration
func resourceDNSRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordCreate,
		ReadContext:   resourceDNSRecordRead,
		UpdateContext: resourceDNSRecordUpdate,
		DeleteContext: resourceDNSRecordDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "DNS record domain",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"ip": {
				Description: "IP address where traffic is routed to from the DNS record domain",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

// resourceDNSRecordCreate handles the creation a local DNS record via Terraform
func resourceDNSRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	domain := d.Get("domain").(string)
	ip := d.Get("ip").(string)

	_, err := client.CreateDNSRecord(ctx, &pihole.DNSRecord{
		Domain: domain,
		IP:     ip,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	return diags
}

// resourceDNSRecordRead finds a local DNS record based on the associated domain ID
func resourceDNSRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	record, err := client.GetDNSRecord(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*pihole.NotFoundError); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	if err = d.Set("domain", record.Domain); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("ip", record.IP); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceDNSRecordUpdate handles updates of a local DNS record via Terraform
func resourceDNSRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	_, err := client.UpdateDNSRecord(ctx, &pihole.DNSRecord{
		Domain: d.Get("domain").(string),
		IP:     d.Get("ip").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDNSRecordRead(ctx, d, meta)
}

// resourceDNSRecordDelete handles the deletion of a local DNS record via Terraform
func resourceDNSRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	if err := client.DeleteDNSRecord(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
