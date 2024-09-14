package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// resourceGroup returns the Terraform resource management configuration for a Pi-hole group
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "A construct to associate clients with allow/deny lists and/or adlists",
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Group name",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					name := strings.TrimSpace(val.(string))
					validName := regexp.MustCompile(`^\S*$`)

					if !validName.MatchString(name) {
						errs = append(errs, fmt.Errorf("%s field cannot contain spaces: %q", key, name))
					}

					return
				},
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
			Enabled:     &enabled,
			Description: description,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.FormatInt(group.ID, 10))

	return diags
}

// resourceGroupRead reads a Pi-hole group resource
func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := client.GetGroupByID(ctx, id)
	if err != nil {
		if _, ok := err.(*pihole.NotFoundError); ok {
			d.SetId("")
			return nil
		}

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

	d.SetId(strconv.FormatInt(group.ID, 10))

	return diags
}

// resourceGroupUpdate handles updates of a Pi-hole group
func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	client, ok := meta.(*pihole.Client)
	if !ok {
		return diag.Errorf("Could not load client in resource request")
	}

	group, err := client.UpdateGroup(ctx, &pihole.GroupUpdateRequest{
		Name:        d.Get("name").(string),
		Enabled:     pihole.Bool(d.Get("enabled").(bool)),
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

	d.SetId(strconv.FormatInt(group.ID, 10))

	return resourceGroupRead(ctx, d, meta)
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
