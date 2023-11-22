package pihole

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Run("Fail login if URL is not set", func(t *testing.T) {
		t.Parallel()

		client := New(Config{
			Password: "test",
		})

		err := client.Init(context.Background())
		require.ErrorIs(t, err, ErrClientValidationFailed)
		require.Contains(t, err.Error(), "Pi-hole URL is not set")
	})

	t.Run("Fail login if login http request fails", func(t *testing.T) {
		t.Parallel()

		client := New(Config{
			Password: "test",
			URL:      "fake-url",
		})

		err := client.Init(context.Background())
		require.NoError(t, err)

		err = client.Login(context.Background())
		require.ErrorIs(t, err, ErrLoginFailed)
		require.Contains(t, err.Error(), "request failed")
	})

	t.Run("Fail login if no session ID is found", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()

		mux.HandleFunc("/admin/index.php", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<div id="token">token</div>`)) //nolint:errcheck
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client := New(Config{
			Password: "test",
			URL:      server.URL,
		})

		err := client.Init(context.Background())
		require.NoError(t, err)

		err = client.Login(context.Background())
		require.ErrorIs(t, err, ErrLoginFailed)
		require.Contains(t, err.Error(), "session ID not found")
	})

	t.Run("Fail login if session ID is malformed", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()

		mux.HandleFunc("/admin/index.php", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Set-Cookie", "malformed")
			w.Write([]byte(`<div id="token">token</div>`)) //nolint:errcheck
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client := New(Config{
			Password: "test",
			URL:      server.URL,
		})

		err := client.Init(context.Background())
		require.NoError(t, err)

		err = client.Login(context.Background())
		require.ErrorIs(t, err, ErrLoginFailed)
		require.Contains(t, err.Error(), "malformed session cookie")
	})

	t.Run("Fail login if token is not found in response", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()

		mux.HandleFunc("/admin/index.php", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Set-Cookie", "session-id=ID;")
			w.Write([]byte(``)) //nolint:errcheck
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client := New(Config{
			Password: "test",
			URL:      server.URL,
		})

		err := client.Init(context.Background())
		require.NoError(t, err)

		err = client.Login(context.Background())
		require.ErrorIs(t, err, ErrLoginFailed)
		require.Contains(t, err.Error(), "invalid password")
	})

	t.Run("Initialize a client", func(t *testing.T) {
		t.Parallel()

		mux := http.NewServeMux()

		mux.HandleFunc("/admin/index.php", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Set-Cookie", "session-id=ID;")
			w.Write([]byte(`<div id="token">token</div>`)) //nolint:errcheck
		})

		server := httptest.NewServer(mux)
		defer server.Close()

		client := New(Config{
			Password: "test",
			URL:      server.URL,
		})

		require.NoError(t, client.Init(context.Background()))
		require.Equal(t, client.password, "test")
		require.Equal(t, client.webPassword, doubleHash256("test"))

		require.NoError(t, client.Login(context.Background()))
		require.Equal(t, client.sessionID, "ID")
		require.Equal(t, client.sessionToken, "token")
	})
}
