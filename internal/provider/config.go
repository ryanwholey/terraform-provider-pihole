package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	retryablehttp "github.com/hashicorp/go-retryablehttp"
	pihole "github.com/ryanwholey/go-pihole"
)

// Config defines the configuration options for the Pi-hole client
type Config struct {
	// The Pi-hole URL
	URL string

	// The Pi-hole admin password
	Password string

	// UserAgent for requests
	UserAgent string

	// Custom CA file
	CAFile string

	// SessionID can be passed to reduce the number of requests against the /api/auth endpoint
	SessionID string

	// CSRFToken is needed for the secure session management
	CSRFToken string
}

func doAuthRequest(ctx context.Context, cfg pihole.Config) (string, string, error) {
	authURL := strings.TrimRight(cfg.BaseURL, "/") + "/api/auth"

	payload := map[string]string{"password": cfg.Password}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}

	req.Header = cfg.Headers.Clone()
	req.Header.Set("Content-Type", "application/json")

	resp, err := cfg.HttpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("auth failed: %s", b)
	}

	var parsed struct {
		Session struct {
			Sid  string `json:"sid"`
			Csrf string `json:"csrf"`
		} `json:"session"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", "", err
	}

	return parsed.Session.Sid, parsed.Session.Csrf, nil
}

func loginAndCreateClient(ctx context.Context, cfg pihole.Config) (*pihole.Client, error) {
	sid, csrf, err := doAuthRequest(ctx, cfg)
	if err != nil {
		return nil, err
	}

	cfg.SessionID = sid
	cfg.CSRFToken = csrf

	return pihole.New(cfg)
}

func (c Config) Client(ctx context.Context) (*pihole.Client, error) {
	httpClient := retryablehttp.NewClient().StandardClient()

	if c.CAFile != "" {
		ca, err := os.ReadFile(c.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA file %q: %w", c.CAFile, err)
		}

		rootCAs := x509.NewCertPool()
		rootCAs.AppendCertsFromPEM(ca)

		transport := &http.Transport{}
		transport.TLSClientConfig = &tls.Config{
			RootCAs: rootCAs,
		}

		httpClient.Transport = transport
	}

	headers := http.Header{}
	headers.Add("User-Agent", c.UserAgent)

	config := pihole.Config{
		BaseURL:    c.URL,
		Password:   c.Password,
		Headers:    headers,
		HttpClient: httpClient,
		SessionID:  c.SessionID,
		CSRFToken:  c.CSRFToken,
	}

	client, err := loginAndCreateClient(ctx, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}
