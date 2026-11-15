package main

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/ryanwholey/terraform-provider-pihole/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})

	// Serve returns once Terraform shuts the plugin down.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	provider.Logout(ctx)
}
