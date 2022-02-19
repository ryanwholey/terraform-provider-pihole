package provider

import (
	"context"

	"github.com/ryanwholey/terraform-provider-pihole/internal/pihole"
)

// Config defines the configuration options for the Pi-hole client
type Config struct {
	// The Pi-hole URL
	URL string

	// The Pi-hole admin password
	Password string

	// UserAgent for requests
	UserAgent string

	// Pi-hole API token
	APIToken string
}

// Client initializes a new pihole client from the passed configuration
func (c Config) Client(ctx context.Context) (*pihole.Client, error) {
	config := pihole.Config{
		URL:       c.URL,
		Password:  c.Password,
		UserAgent: c.UserAgent,
		APIToken:  c.APIToken,
	}

	client := pihole.New(config)

	if err := client.Init(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
