// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Logger interface for logging
type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// DefaultLogger is a simple logger that writes to standard output
type DefaultLogger struct{}

func (l *DefaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (l *DefaultLogger) Println(v ...interface{}) {
	log.Println(v...)
}

// NoOpLogger is a logger that doesn't log anything
type NoOpLogger struct{}

func (l *NoOpLogger) Printf(format string, v ...interface{}) {}
func (l *NoOpLogger) Println(v ...interface{})               {}

// Client is the main SDK client for interacting with the Upwind API
type Client struct {
	config      *Config
	httpClient  HTTPClient
	oauthCfg    *clientcredentials.Config
	tokenSrc    oauth2.TokenSource
	tokenMu     sync.Mutex
	token       *oauth2.Token
	rateLimiter *rate.Limiter
	logger      Logger
	clientMu    sync.Mutex // Protects httpClient during refresh
}

// NewClient creates a new Upwind API client with the provided configuration
func NewClient(cfg *Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	oauthCfg := &clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.GetTokenURL(),
		EndpointParams: map[string][]string{
			"audience": {cfg.GetAudience()},
		},
		AuthStyle: oauth2.AuthStyleInParams,
	}

	tokenSrc := oauthCfg.TokenSource(context.Background())

	// Create rate limiter
	var rateLimiter *rate.Limiter
	if cfg.RateLimitPerSecond > 0 {
		rateLimiter = rate.NewLimiter(rate.Limit(cfg.RateLimitPerSecond), cfg.RateLimitPerSecond)
	}

	client := &Client{
		config:      cfg,
		oauthCfg:    oauthCfg,
		tokenSrc:    tokenSrc,
		rateLimiter: rateLimiter,
		logger:      &NoOpLogger{}, // Default to no logging
	}

	// Create initial HTTP client
	client.httpClient = client.createHTTPClient()

	return client, nil
}

// createHTTPClient creates a new HTTP client with proper connection management
// This is used both at initialization and when refreshing connections
func (c *Client) createHTTPClient() HTTPClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     c.config.IdleConnTimeout,
		DisableKeepAlives:   false,
		// Force connection closure after response if needed for very large datasets
		// DisableKeepAlives can be set to true via environment if HTTP/2 issues persist
	}

	// Configure HTTP/2 or disable it based on config
	if c.config.DisableHTTP2 {
		// Force HTTP/1.1 by disabling HTTP/2
		// Setting TLSNextProto to an empty map disables HTTP/2
		transport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
		c.logger.Println("HTTP/2 disabled, using HTTP/1.1")
	} else {
		// HTTP/2 is enabled by default in Go 1.6+
		// The transport will automatically negotiate HTTP/2 with proper connection management
		c.logger.Println("HTTP/2 enabled with connection management")
	}

	return &http.Client{
		Timeout:   c.config.RequestTimeout,
		Transport: transport,
	}
}

// RefreshHTTPClient creates a new HTTP client, replacing the old one
// This is useful for long-running operations to avoid HTTP/2 GOAWAY issues
func (c *Client) RefreshHTTPClient() {
	c.clientMu.Lock()
	defer c.clientMu.Unlock()

	c.logger.Println("Refreshing HTTP client to avoid connection issues...")
	c.httpClient = c.createHTTPClient()
}

// NewClientFromEnv creates a new client from environment variables
func NewClientFromEnv() (*Client, error) {
	cfg, err := LoadConfigFromEnv()
	if err != nil {
		return nil, err
	}
	return NewClient(cfg)
}

// NewClientFromFile creates a new client from a configuration file
func NewClientFromFile(path string) (*Client, error) {
	cfg, err := LoadConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	return NewClient(cfg)
}

// SetLogger sets the logger for the client
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// EnableLogging enables logging with the default logger
func (c *Client) EnableLogging() {
	c.logger = &DefaultLogger{}
}

// GetOrganizationID returns the organization ID
func (c *Client) GetOrganizationID() string {
	return c.config.OrganizationID
}

// getToken retrieves a valid OAuth2 token, refreshing if necessary
func (c *Client) getToken(ctx context.Context) (*oauth2.Token, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// Return cached token if still valid with 60 second buffer
	if c.token != nil && c.token.Valid() && time.Until(c.token.Expiry) > 60*time.Second {
		return c.token, nil
	}

	c.logger.Println("Refreshing OAuth2 token...")
	token, err := c.oauthCfg.TokenSource(ctx).Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	c.token = token
	c.logger.Println("Token refreshed successfully")
	return token, nil
}

// doRequest executes an HTTP request with retries, rate limiting, and proper error handling
func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Apply rate limiting
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiter error: %w", err)
		}
	}

	var lastErr error
	backoff := 100 * time.Millisecond

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		// Get valid token
		token, err := c.getToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("token error: %w", err)
		}

		// Clone request and add auth header
		reqClone := req.Clone(ctx)
		reqClone.Header.Set("Authorization", "Bearer "+token.AccessToken)
		reqClone.Header.Set("User-Agent", UserAgent())
		reqClone.Header.Set("Accept", "application/json")

		c.logger.Printf("Request: %s %s (attempt %d/%d)", reqClone.Method, reqClone.URL.String(), attempt+1, c.config.MaxRetries+1)

		// Execute request with client mutex to handle concurrent refresh
		c.clientMu.Lock()
		httpClient := c.httpClient
		c.clientMu.Unlock()

		resp, err := httpClient.Do(reqClone)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			c.logger.Printf("Request error: %v", lastErr)

			if attempt < c.config.MaxRetries {
				time.Sleep(backoff)
				backoff *= 2 // Exponential backoff
			}
			continue
		}

		// Handle 401 Unauthorized - token might be invalid
		if resp.StatusCode == http.StatusUnauthorized && attempt < c.config.MaxRetries {
			c.logger.Println("Received 401 Unauthorized, invalidating token and retrying...")
			c.tokenMu.Lock()
			c.token = nil
			c.tokenMu.Unlock()
			_ = resp.Body.Close()
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		// Handle 429 Too Many Requests
		if resp.StatusCode == http.StatusTooManyRequests && attempt < c.config.MaxRetries {
			retryAfter := resp.Header.Get("Retry-After")
			waitDuration := backoff
			if retryAfter != "" {
				if duration, err := time.ParseDuration(retryAfter + "s"); err == nil {
					waitDuration = duration
				}
			}
			c.logger.Printf("Rate limited (429), waiting %v before retry...", waitDuration)
			_ = resp.Body.Close()
			time.Sleep(waitDuration)
			backoff *= 2
			continue
		}

		// Handle 5xx errors with retry
		if resp.StatusCode >= 500 && attempt < c.config.MaxRetries {
			c.logger.Printf("Server error (%d), retrying...", resp.StatusCode)
			_ = resp.Body.Close()
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("all retries exhausted: %w", lastErr)
}
