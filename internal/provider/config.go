package provider

import (
	"context"

	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// Config defines the configuration options for the Pihole client
type Config struct {
	// The Pi-hole URL
	URL string

	// The Pi-hole admin password
	Password string

	// UserAgent for requests
	UserAgent string
}

func (c Config) Client(ctx context.Context) (*pihole.Client, error) {
	config := pihole.Config{
		URL:       c.URL,
		Password:  c.Password,
		UserAgent: c.UserAgent,
	}

	client := pihole.New(config)

	if err := client.Init(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
