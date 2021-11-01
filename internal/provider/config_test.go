package provider

import (
	"context"
	"testing"
)

func TestConfigEmptyPassword(t *testing.T) {
	config := Config{
		URL: "pihole.foo.com",
	}

	if _, err := config.Client(context.Background()); err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
