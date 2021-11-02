package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// resorceGroup returns the Terraform resource management configuration for a Pi-hole group
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Group name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Group description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether to enable the group",
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
		},
	}
}

// resourceGroupCreate handles the creation a Pi-hole group
func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	enabled := d.Get("enabled").(bool)

	group, err := client.CreateGroup(ctx, &pihole.GroupCreateRequest{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if !enabled {
		_, err := client.UpdateGroup(ctx, &pihole.GroupUpdateRequest{
			Name:        name,
			Enabled:     enabled,
			Description: description,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.FormatUint(uint64(group.ID), 10))

	return diags
}

// resourceGroupRead reads a Pi-hole group resource
func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	name := d.Get("name").(string)

	group, err := client.GetGroup(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("description", group.Description); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("enabled", group.Enabled); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatUint(uint64(group.ID), 10))

	return diags
}

// resourceDNSRecordUpdate handles updates of a Pi-hole group
func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	group, err := client.UpdateGroup(ctx, &pihole.GroupUpdateRequest{
		Name:        d.Get("name").(string),
		Enabled:     d.Get("enabled").(bool),
		Description: d.Get("description").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("description", group.Description); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("enabled", group.Enabled); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatUint(uint64(group.ID), 10))

	return resourceDNSRecordRead(ctx, d, meta)
}

// resourceGroupDelete handles the deletion of a Pi-hole group
func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	if err := client.DeleteGroup(ctx, d.Get("name").(string)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
