package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"sync"

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

	httpClient.Transport = serialRequests{next: httpClient.Transport}

	headers := http.Header{}
	headers.Add("User-Agent", c.UserAgent)

	config := pihole.Config{
		BaseURL:    c.URL,
		Password:   c.Password,
		Headers:    headers,
		HttpClient: httpClient,
		SessionID:  c.SessionID,
	}

	client, err := pihole.New(config)
	if err != nil {
		return nil, err
	}

	// A session passed in from outside belongs to the caller, so it is not logged out.
	if c.SessionID == "" {
		ownSessions.Lock()
		ownSessions.clients = append(ownSessions.clients, client)
		ownSessions.Unlock()
	}

	return client, nil
}

// requestLock is shared by all clients in the process, so provider aliases pointing at
// the same Pi-hole are serialized as well.
var requestLock sync.Mutex

// serialRequests sends one request at a time. FTL applies each config change as copy,
// modify, replace: concurrent writes overwrite each other and records get lost (#88), and
// a read during a write can return a half-written dns.hosts entry.
type serialRequests struct {
	next http.RoundTripper
}

func (t serialRequests) RoundTrip(req *http.Request) (*http.Response, error) {
	requestLock.Lock()
	defer requestLock.Unlock()

	return t.next.RoundTrip(req)
}

var ownSessions struct {
	sync.Mutex
	clients []*pihole.Client
}

// Logout ends the sessions the provider opened. Pi-hole keeps a session for
// webserver.session.timeout (30 minutes by default) and refuses new logins with 429 once
// webserver.api.max_sessions (16 by default) are taken (#97).
func Logout(ctx context.Context) {
	ownSessions.Lock()
	defer ownSessions.Unlock()

	for _, client := range ownSessions.clients {
		_ = client.SessionAPI.Logout(ctx)
	}
	ownSessions.clients = nil
}
