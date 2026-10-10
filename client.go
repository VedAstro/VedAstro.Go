package vedastro

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Version is the SDK release version.
const Version = "v1.0.0"

// DefaultBaseURL is the public calculator service.
const DefaultBaseURL = "https://vedastro.zaishi.net/api/Calculate"

const banner = "VedAstro : Easy To Use Advanced Astrology Engine"

var bannerOnce sync.Once

// Client owns immutable settings and is safe for concurrent calculation calls.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    *time.Duration
	ayanamsa   Ayanamsa
	configErr  error
	updateURL  string
}

// ClientOption configures a client during construction.
type ClientOption func(*Client)

// NewClient creates a customer API client. An empty key uses the service's free tier.
// Invalid options are reported by calculation methods before any HTTP request.
func NewClient(apiKey string, options ...ClientOption) *Client {
	c := &Client{
		apiKey: apiKey, baseURL: DefaultBaseURL, httpClient: &http.Client{},
		timeout:   nil,
		updateURL: "https://proxy.golang.org/github.com/!ved!astro/!ved!astro.!go/@latest",
	}
	for _, option := range options {
		if option != nil {
			option(c)
		}
	}
	bannerOnce.Do(func() {
		if os.Stdout != nil {
			_, _ = fmt.Fprintln(os.Stdout, banner)
		}
	})
	return c
}

// WithBaseURL selects a calculator service URL, for example for an HTTP test server.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		u, err := url.Parse(baseURL)
		if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") ||
			u.User != nil || u.RawQuery != "" || u.Fragment != "" {
			c.configErr = errors.New("vedastro: base URL must be an HTTP(S) URL without credentials, query, or fragment")
			return
		}
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient supplies a transport configuration. A shallow copy is retained;
// callers must not mutate a shared custom transport while requests are in flight.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		if client == nil {
			c.configErr = errors.New("vedastro: HTTP client cannot be nil")
			return
		}
		copy := *client
		c.httpClient = &copy
	}
}

// WithTimeout sets a calculation deadline, including reading the response body. An earlier
// deadline on the request context still takes precedence.
//
// No deadline is applied by default, and that is deliberate: a VedAstro calculation can take
// milliseconds or minutes, and this client cannot know what is acceptable for your workload, so a
// built-in deadline would be brittle logic that can truncate a valid answer. Use this only when
// your own code has decided a request has run too long.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout <= 0 {
			c.configErr = errors.New("vedastro: timeout must be positive")
			return
		}
		c.timeout = &timeout
	}
}

// WithDefaultAyanamsa sets the default sent with this client's requests.
// The empty value leaves the choice to the service (currently Lahiri).
func WithDefaultAyanamsa(value Ayanamsa) ClientOption {
	return func(c *Client) {
		name, err := normalizeAyanamsa(value)
		if err != nil {
			c.configErr = err
			return
		}
		c.ayanamsa = name
	}
}

type ayanamsaKey struct{}

// WithAyanamsa overrides ayanamsa for this context and its children.
// An empty value explicitly selects the server default, overriding the client default.
// Invalid values are reported when a calculation is called. ctx must be non-nil.
func WithAyanamsa(ctx context.Context, value Ayanamsa) context.Context {
	return context.WithValue(ctx, ayanamsaKey{}, value)
}

func normalizeAyanamsa(value Ayanamsa) (Ayanamsa, error) {
	name := strings.TrimSpace(string(value))
	if name == "" {
		return "", nil
	}
	for _, known := range knownAyanamsas {
		if strings.EqualFold(string(known), name) {
			return known, nil
		}
	}
	return "", errors.New("vedastro: unknown ayanamsa; use an Ayanamsa constant")
}

// APIError describes a remote or response-format failure without exposing credentials.
type APIError struct {
	Endpoint   string
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("vedastro: %s (HTTP %d): %s", e.Endpoint, e.StatusCode, e.Message)
}

func (c *Client) call(ctx context.Context, endpoint string, parameters map[string]any) (any, error) {
	if c == nil {
		return nil, errors.New("vedastro: nil client")
	}
	if ctx == nil {
		return nil, errors.New("vedastro: nil context")
	}
	if c.configErr != nil {
		return nil, c.configErr
	}
	ayanamsa := c.ayanamsa
	if scoped, ok := ctx.Value(ayanamsaKey{}).(Ayanamsa); ok {
		var err error
		ayanamsa, err = normalizeAyanamsa(scoped)
		if err != nil {
			return nil, err
		}
	}
	body := make(map[string]any, len(parameters)+2)
	for key, value := range parameters {
		body[key] = value
	}
	if c.apiKey != "" {
		body["APIKey"] = c.apiKey
	}
	if ayanamsa != "" {
		body["Ayanamsa"] = string(ayanamsa)
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("vedastro: parameters cannot be encoded as API JSON")
	}
	// Only apply a deadline when the caller asked for one: a calculation can take minutes, so
	// imposing a default would truncate a valid answer.
	if c.timeout != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, *c.timeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+url.PathEscape(endpoint), strings.NewReader(string(encoded)))
	if err != nil {
		return nil, errors.New("vedastro: cannot create API request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, image/svg+xml")
	response, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("vedastro: request interrupted: %w", ctx.Err())
		}
		return nil, errors.New("vedastro: HTTP request failed")
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("vedastro: response interrupted: %w", ctx.Err())
		}
		return nil, errors.New("vedastro: cannot read API response")
	}
	fail := func(message string) (any, error) {
		if c.apiKey != "" {
			message = strings.ReplaceAll(message, c.apiKey, "[redacted]")
		}
		return nil, &APIError{Endpoint: endpoint, StatusCode: response.StatusCode, Message: message}
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fail("API HTTP request failed")
	}
	contentType, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if contentType == "image/svg+xml" {
		return string(data), nil
	}
	var envelope struct {
		Status  string          `json:"Status"`
		Payload json.RawMessage `json:"Payload"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return fail("API returned invalid JSON")
	}
	if strings.EqualFold(envelope.Status, "Fail") {
		var message string
		if json.Unmarshal(envelope.Payload, &message) != nil || message == "" {
			message = "API reported failure"
		}
		return fail(message)
	}
	if envelope.Payload == nil {
		return fail("API response is missing its Payload")
	}
	var payload any
	if err := json.Unmarshal(envelope.Payload, &payload); err != nil {
		return fail("API returned an invalid Payload")
	}
	if object, ok := payload.(map[string]any); ok && len(object) == 1 {
		for _, value := range object {
			return value, nil
		}
	}
	return payload, nil
}

// Ptr constructs an optional argument pointer, including zero and false values.
func Ptr[T any](value T) *T { return &value }
