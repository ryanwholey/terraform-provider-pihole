package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// resourceCNAMERecord returns the CNAME Terraform resource management configuration
func resourceCNAMERecord() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a Pi-hole CNAME record",
		CreateContext: resourceCNAMERecordCreate,
		ReadContext:   resourceCNAMERecordRead,
		UpdateContext: resourceCNAMERecordUpdate,
		DeleteContext: resourceCNAMERecordDelete,
		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Domain to create a CNAME record for",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"target": {
				Description: "Value of the CNAME record where traffic will be directed to from the configured domain value",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

// resourceCNAMERecordCreate handles the creation a CNAME record via Terraform
func resourceCNAMERecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	domain := d.Get("domain").(string)
	target := d.Get("target").(string)

	_, err := client.CreateCNAMERecord(ctx, &pihole.CNAMERecord{
		Domain: domain,
		Target: target,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	return diags
}

// resourceCNAMERecordRead retrieves the CNAME record of the associated domain ID
func resourceCNAMERecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	record, err := client.GetCNAMERecord(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("domain", record.Domain); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("target", record.Target); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

// resourceCNAMERecordUpdate handles CNAME record updates via Terraform
func resourceCNAMERecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	_, err := client.UpdateCNAMERecord(ctx, &pihole.CNAMERecord{
		Domain: d.Get("domain").(string),
		Target: d.Get("target").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCNAMERecordRead(ctx, d, meta)
}

// resourceCNAMERecordDelete handles the deletion of a CNAME record via Terraform
func resourceCNAMERecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	if err := client.DeleteCNAMERecord(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
