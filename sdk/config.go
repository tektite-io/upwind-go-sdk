// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Region represents the Upwind API region
type Region string

const (
	// RegionUS represents the US region
	RegionUS Region = "US"
	// RegionEU represents the EU region
	RegionEU Region = "EU"
	// RegionME represents the ME region
	RegionME Region = "ME"
)

// Config holds the configuration for the Upwind SDK
type Config struct {
	// ClientID is the OAuth2 client ID
	ClientID string `json:"client_id"`
	// ClientSecret is the OAuth2 client secret
	ClientSecret string `json:"client_secret"`
	// OrganizationID is the Upwind organization ID
	OrganizationID string `json:"organization_id"`
	// Region is the API region (US, EU, or ME)
	Region Region `json:"region"`
	// BaseURL is the API base URL (optional, derived from region if not provided)
	BaseURL string `json:"base_url,omitempty"`
	// TokenURL is the OAuth2 token endpoint (optional, defaults to Upwind auth endpoint)
	TokenURL string `json:"token_url,omitempty"`
	// MaxRetries is the maximum number of retry attempts for failed requests
	MaxRetries int `json:"max_retries"`
	// MaxConcurrency is the maximum number of concurrent API requests
	MaxConcurrency int `json:"max_concurrency"`
	// PageSize is the default page size for paginated requests
	PageSize int `json:"page_size"`
	// RateLimitPerSecond is the maximum number of requests per second (0 = no limit)
	RateLimitPerSecond int `json:"rate_limit_per_second"`

	// HTTP/2 Connection Management Settings (to avoid long-lived connection issues)

	// RequestTimeout is the timeout for individual HTTP requests
	RequestTimeout time.Duration `json:"request_timeout"`
	// IdleConnTimeout is how long idle connections are kept open
	IdleConnTimeout time.Duration `json:"idle_conn_timeout"`
	// DisableHTTP2 forces HTTP/1.1 to avoid HTTP/2 GOAWAY issues with long-running connections
	DisableHTTP2 bool `json:"disable_http2"`
	// ConnectionRefreshPages is the number of pages to fetch before refreshing the HTTP client
	// This helps avoid HTTP/2 connection issues with very large datasets (0 = never refresh)
	ConnectionRefreshPages int `json:"connection_refresh_pages"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Region:                 RegionUS,
		MaxRetries:             3,
		MaxConcurrency:         10,
		PageSize:               100,
		RateLimitPerSecond:     10,
		RequestTimeout:         30 * time.Second,
		IdleConnTimeout:        30 * time.Second, // Shorter than default 90s to avoid HTTP/2 GOAWAY
		DisableHTTP2:           false,            // Keep HTTP/2 by default, but with proper connection management
		ConnectionRefreshPages: 100,              // Refresh HTTP client every 100 pages for large datasets
	}
}

// LoadConfigFromEnv loads configuration from environment variables
// Supported environment variables:
//   - UPWIND_CLIENT_ID: OAuth2 client ID
//   - UPWIND_CLIENT_SECRET: OAuth2 client secret
//   - UPWIND_ORGANIZATION_ID: Organization ID
//   - UPWIND_REGION: API region (US, EU, ME)
//   - UPWIND_BASE_URL: Custom base URL (optional)
//   - UPWIND_TOKEN_URL: Custom token URL (optional)
//   - UPWIND_MAX_RETRIES: Maximum retry attempts (default: 3)
//   - UPWIND_MAX_CONCURRENCY: Maximum concurrent requests (default: 10)
//   - UPWIND_PAGE_SIZE: Default page size (default: 100)
//   - UPWIND_RATE_LIMIT: Requests per second limit (default: 10)
//   - UPWIND_REQUEST_TIMEOUT: Request timeout in seconds (default: 30)
//   - UPWIND_IDLE_CONN_TIMEOUT: Idle connection timeout in seconds (default: 30)
//   - UPWIND_DISABLE_HTTP2: Set to "true" to disable HTTP/2 (default: false)
//   - UPWIND_CONNECTION_REFRESH_PAGES: Number of pages before refreshing connection (default: 100)
func LoadConfigFromEnv() (*Config, error) {
	cfg := DefaultConfig()

	if clientID := os.Getenv("UPWIND_CLIENT_ID"); clientID != "" {
		cfg.ClientID = clientID
	}

	if clientSecret := os.Getenv("UPWIND_CLIENT_SECRET"); clientSecret != "" {
		cfg.ClientSecret = clientSecret
	}

	if orgID := os.Getenv("UPWIND_ORGANIZATION_ID"); orgID != "" {
		cfg.OrganizationID = orgID
	}

	if region := os.Getenv("UPWIND_REGION"); region != "" {
		cfg.Region = Region(strings.ToUpper(region))
	}

	if baseURL := os.Getenv("UPWIND_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	if tokenURL := os.Getenv("UPWIND_TOKEN_URL"); tokenURL != "" {
		cfg.TokenURL = tokenURL
	}

	// Parse integer values with defaults
	if maxRetries := os.Getenv("UPWIND_MAX_RETRIES"); maxRetries != "" {
		var retries int
		if _, err := fmt.Sscanf(maxRetries, "%d", &retries); err == nil {
			cfg.MaxRetries = retries
		}
	}

	if maxConcurrency := os.Getenv("UPWIND_MAX_CONCURRENCY"); maxConcurrency != "" {
		var concurrency int
		if _, err := fmt.Sscanf(maxConcurrency, "%d", &concurrency); err == nil {
			cfg.MaxConcurrency = concurrency
		}
	}

	if pageSize := os.Getenv("UPWIND_PAGE_SIZE"); pageSize != "" {
		var size int
		if _, err := fmt.Sscanf(pageSize, "%d", &size); err == nil {
			cfg.PageSize = size
		}
	}

	if rateLimit := os.Getenv("UPWIND_RATE_LIMIT"); rateLimit != "" {
		var limit int
		if _, err := fmt.Sscanf(rateLimit, "%d", &limit); err == nil {
			cfg.RateLimitPerSecond = limit
		}
	}

	// Parse duration values
	if requestTimeout := os.Getenv("UPWIND_REQUEST_TIMEOUT"); requestTimeout != "" {
		var seconds int
		if _, err := fmt.Sscanf(requestTimeout, "%d", &seconds); err == nil {
			cfg.RequestTimeout = time.Duration(seconds) * time.Second
		}
	}

	if idleConnTimeout := os.Getenv("UPWIND_IDLE_CONN_TIMEOUT"); idleConnTimeout != "" {
		var seconds int
		if _, err := fmt.Sscanf(idleConnTimeout, "%d", &seconds); err == nil {
			cfg.IdleConnTimeout = time.Duration(seconds) * time.Second
		}
	}

	// Parse boolean values
	if disableHTTP2 := os.Getenv("UPWIND_DISABLE_HTTP2"); disableHTTP2 != "" {
		cfg.DisableHTTP2 = strings.ToLower(disableHTTP2) == "true"
	}

	if connectionRefreshPages := os.Getenv("UPWIND_CONNECTION_REFRESH_PAGES"); connectionRefreshPages != "" {
		var pages int
		if _, err := fmt.Sscanf(connectionRefreshPages, "%d", &pages); err == nil {
			cfg.ConnectionRefreshPages = pages
		}
	}

	return cfg, cfg.Validate()
}

// LoadConfigFromFile loads configuration from a JSON file
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	cfg := DefaultConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return cfg, cfg.Validate()
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}
	if c.ClientSecret == "" {
		return fmt.Errorf("client_secret is required")
	}
	if c.OrganizationID == "" {
		return fmt.Errorf("organization_id is required")
	}

	if c.Region != RegionUS && c.Region != RegionEU && c.Region != RegionME {
		return fmt.Errorf("invalid region: %s (must be US, EU, or ME)", c.Region)
	}

	if c.MaxRetries < 0 {
		return fmt.Errorf("max_retries must be >= 0")
	}

	if c.MaxConcurrency < 1 {
		return fmt.Errorf("max_concurrency must be >= 1")
	}

	if c.PageSize < 1 || c.PageSize > 10000 {
		return fmt.Errorf("page_size must be between 1 and 10000")
	}

	return nil
}

// GetBaseURL returns the base URL for the API based on the region
func (c *Config) GetBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}

	switch c.Region {
	case RegionEU:
		return "https://api.eu.upwind.io/v1"
	case RegionME:
		return "https://api.me.upwind.io/v1"
	default:
		return "https://api.upwind.io/v1"
	}
}

// GetTokenURL returns the OAuth2 token URL
func (c *Config) GetTokenURL() string {
	if c.TokenURL != "" {
		return c.TokenURL
	}
	return "https://auth.upwind.io/oauth/token"
}

// GetAudience returns the OAuth2 audience based on the region
func (c *Config) GetAudience() string {
	switch c.Region {
	case RegionEU:
		return "https://api.eu.upwind.io"
	case RegionME:
		return "https://api.me.upwind.io"
	default:
		return "https://api.upwind.io"
	}
}
