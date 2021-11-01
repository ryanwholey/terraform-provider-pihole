package provider

import (
	"context"
	"fmt"

	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
	"github.com/ryanwholey/terraform-provider-pihole/internal/version"
)

// Config defines the configuration options for the Pihole client
type Config struct {
	// The Pihole URL
	URL string

	// The Pihole admin password
	Password string

	// UserAgent for requests
	UserAgent string
}

func (c Config) Client(ctx context.Context) (*pihole.Client, error) {
	config := &pihole.Config{
		URL:      c.URL,
		Password: c.Password,
		// TODO: Use more conventional UserAgent input
		UserAgent: fmt.Sprintf("%s-%s", "terraform-provider-pihole", version.ProviderVersion),
	}

	if config.Password == "" {
		return nil, fmt.Errorf("password must be set in the provider")
	}

	client, err := pihole.New(config)
	if err != nil {
		return nil, err
	}

	if err = client.Init(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
