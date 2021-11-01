package pihole

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Config struct {
	Password  string
	URL       string
	UserAgent string
}

type Client struct {
	URL         string
	UserAgent   string
	password    string
	sessionID   string
	token       string
	webPassword string
	client      *http.Client
}

// doubleHash256 takes a string, double hashes it using the sha256 algorithm and returns the value
func doubleHash256(data string) string {
	hash := sha256.Sum256([]byte(data))
	sha1 := fmt.Sprintf("%x", hash[:])

	hash2 := sha256.Sum256([]byte(sha1))
	return fmt.Sprintf("%x", hash2[:])
}

// New returns a new pihole client
func New(config *Config) (*Client, error) {
	client := &Client{
		URL:         config.URL,
		UserAgent:   config.UserAgent,
		password:    config.Password,
		client:      &http.Client{},
		webPassword: doubleHash256(config.Password),
	}

	return client, nil
}

// Init sets fields on the client which are a product of pihole network requests or other side effects
func (c *Client) Init(ctx context.Context) error {
	return c.login(ctx)
}

// Validate returns an error if any of the required fields on the client are empty
func (c *Client) Validate() error {
	if c.password == "" {
		return fmt.Errorf("failed to validate client: password is not set")
	}
	if c.webPassword == "" {
		return fmt.Errorf("failed to validate client: webPassword is not set")
	}
	if c.token == "" {
		return fmt.Errorf("failed to validate client: token is not set")
	}
	if c.sessionID == "" {
		return fmt.Errorf("failed to validate client: sessionID is not set")
	}

	return nil
}

// Request executes a basic unauthenticated http request
func (c *Client) Request(ctx context.Context, method string, path string, data *url.Values) (*http.Request, error) {
	d := data
	if d == nil {
		d = &url.Values{}
	}

	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.URL, path), strings.NewReader(d.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// mergeURLValues merges the passed URL values into a single url.Values object
func mergeURLValues(vs ...url.Values) url.Values {
	data := url.Values{}

	for _, val := range vs {
		for k, v := range val {
			data.Add(k, v[0])
		}
	}

	return data
}

// RequestWithSession executes a request with appropriate session authentication (login() must have been called)
func (c Client) RequestWithSession(ctx context.Context, method string, path string, data *url.Values) (*http.Request, error) {
	d := mergeURLValues(url.Values{
		"token": []string{c.token},
	}, *data)
	req, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", c.URL, path), strings.NewReader(d.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", fmt.Sprintf("PHPSESSID=%s", c.sessionID))

	return req, nil
}

// RequestWithAuth adds an auth token to the passed request
func (c Client) RequestWithAuth(ctx context.Context, method string, path string, data *url.Values) (*http.Request, error) {
	u, err := url.Parse(fmt.Sprintf("%s%s", c.URL, path))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Add("auth", c.webPassword)
	u.RawQuery = q.Encode()

	d := data
	if d == nil {
		d = &url.Values{}
	}

	return http.NewRequestWithContext(ctx, method, u.String(), strings.NewReader(d.Encode()))
}

// login sets a new sessionID and csrf token in the client to be used for logged in requests
func (c *Client) login(ctx context.Context) error {
	data := &url.Values{
		"pw": []string{c.password},
	}

	req, err := c.Request(ctx, "POST", "/admin/index.php?login=", data)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	sessionID := strings.Split(strings.Split(res.Header.Get("Set-Cookie"), "=")[1], ";")[0]
	c.sessionID = sessionID

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	c.token = doc.Find("#token").Text()

	return nil
}
