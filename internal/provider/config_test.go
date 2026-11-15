package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestClientSerializesRequests(t *testing.T) {
	var inFlight, maxInFlight int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&inFlight, 1)
		defer atomic.AddInt32(&inFlight, -1)

		for {
			m := atomic.LoadInt32(&maxInFlight)
			if n <= m || atomic.CompareAndSwapInt32(&maxInFlight, m, n) {
				break
			}
		}

		time.Sleep(20 * time.Millisecond)
		_, _ = w.Write([]byte(`{"config":{"dns":{"hosts":[]}}}`))
	}))
	defer server.Close()

	client, err := Config{URL: server.URL, SessionID: "external"}.Client(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.LocalDNS.Create(context.Background(), "foo.com", "127.0.0.1"); err == nil {
				t.Error("expected an error, the test server never answers 201")
			}
			if _, err := client.LocalDNS.List(context.Background()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	if maxInFlight != 1 {
		t.Fatalf("expected at most 1 request in flight, got %d", maxInFlight)
	}
}

func TestLogoutOnlyOwnSessions(t *testing.T) {
	var mu sync.Mutex
	var logouts []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/auth":
			_, _ = w.Write([]byte(`{"session":{"valid":true,"sid":"own","validity":1800}}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/api/auth/":
			mu.Lock()
			logouts = append(logouts, r.Header.Get("X-FTL-SID"))
			mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
		default:
			_, _ = w.Write([]byte(`{"config":{"dns":{"hosts":[]}}}`))
		}
	}))
	defer server.Close()

	ctx := context.Background()
	for _, sid := range []string{"", "external"} {
		client, err := Config{URL: server.URL, Password: "test", SessionID: sid}.Client(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := client.LocalDNS.List(ctx); err != nil {
			t.Fatal(err)
		}
	}

	Logout(ctx)

	if len(logouts) != 1 || logouts[0] != "own" {
		t.Fatalf("expected exactly one logout of session %q, got %v", "own", logouts)
	}
}
