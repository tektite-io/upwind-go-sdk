// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Region != RegionUS {
		t.Errorf("Expected default region US, got %s", cfg.Region)
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("Expected default max retries 3, got %d", cfg.MaxRetries)
	}
	if cfg.MaxConcurrency != 10 {
		t.Errorf("Expected default max concurrency 10, got %d", cfg.MaxConcurrency)
	}
	if cfg.PageSize != 100 {
		t.Errorf("Expected default page size 100, got %d", cfg.PageSize)
	}
	if cfg.RateLimitPerSecond != 10 {
		t.Errorf("Expected default rate limit 10, got %d", cfg.RateLimitPerSecond)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				ClientID:       "test-client",
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
				MaxRetries:     3,
				MaxConcurrency: 10,
				PageSize:       100,
			},
			wantErr: false,
		},
		{
			name: "missing client ID",
			config: &Config{
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
			},
			wantErr: true,
		},
		{
			name: "missing client secret",
			config: &Config{
				ClientID:       "test-client",
				OrganizationID: "test-org",
				Region:         RegionUS,
			},
			wantErr: true,
		},
		{
			name: "missing organization ID",
			config: &Config{
				ClientID:     "test-client",
				ClientSecret: "test-secret",
				Region:       RegionUS,
			},
			wantErr: true,
		},
		{
			name: "invalid region",
			config: &Config{
				ClientID:       "test-client",
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         "INVALID",
			},
			wantErr: true,
		},
		{
			name: "negative max retries",
			config: &Config{
				ClientID:       "test-client",
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
				MaxRetries:     -1,
			},
			wantErr: true,
		},
		{
			name: "invalid max concurrency",
			config: &Config{
				ClientID:       "test-client",
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
				MaxConcurrency: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid page size - too small",
			config: &Config{
				ClientID:       "test-client",
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
				PageSize:       0,
			},
			wantErr: true,
		},
		{
			name: "invalid page size - too large",
			config: &Config{
				ClientID:       "test-client",
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
				PageSize:       1001,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		region   Region
		baseURL  string
		expected string
	}{
		{"US region", RegionUS, "", "https://api.upwind.io/v1"},
		{"EU region", RegionEU, "", "https://api.eu.upwind.io/v1"},
		{"ME region", RegionME, "", "https://api.me.upwind.io/v1"},
		{"Custom URL", RegionUS, "https://custom.api.com", "https://custom.api.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Region:  tt.region,
				BaseURL: tt.baseURL,
			}
			got := cfg.GetBaseURL()
			if got != tt.expected {
				t.Errorf("GetBaseURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetTokenURL(t *testing.T) {
	cfg := &Config{}
	if got := cfg.GetTokenURL(); got != "https://auth.upwind.io/oauth/token" {
		t.Errorf("GetTokenURL() = %v, want https://auth.upwind.io/oauth/token", got)
	}

	cfg.TokenURL = "https://custom.auth.com/token"
	if got := cfg.GetTokenURL(); got != "https://custom.auth.com/token" {
		t.Errorf("GetTokenURL() with custom = %v, want https://custom.auth.com/token", got)
	}
}

func TestGetAudience(t *testing.T) {
	tests := []struct {
		region   Region
		expected string
	}{
		{RegionUS, "https://api.upwind.io"},
		{RegionEU, "https://api.eu.upwind.io"},
		{RegionME, "https://api.me.upwind.io"},
	}

	for _, tt := range tests {
		t.Run(string(tt.region), func(t *testing.T) {
			cfg := &Config{Region: tt.region}
			got := cfg.GetAudience()
			if got != tt.expected {
				t.Errorf("GetAudience() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original environment
	originalEnv := make(map[string]string)
	envVars := []string{
		"UPWIND_CLIENT_ID",
		"UPWIND_CLIENT_SECRET",
		"UPWIND_ORGANIZATION_ID",
		"UPWIND_REGION",
		"UPWIND_MAX_RETRIES",
	}
	for _, key := range envVars {
		originalEnv[key] = os.Getenv(key)
	}

	// Restore environment after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	// Set test environment variables
	os.Setenv("UPWIND_CLIENT_ID", "test-client")
	os.Setenv("UPWIND_CLIENT_SECRET", "test-secret")
	os.Setenv("UPWIND_ORGANIZATION_ID", "test-org")
	os.Setenv("UPWIND_REGION", "EU")
	os.Setenv("UPWIND_MAX_RETRIES", "5")

	cfg, err := LoadConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadConfigFromEnv() error = %v", err)
	}

	if cfg.ClientID != "test-client" {
		t.Errorf("ClientID = %v, want test-client", cfg.ClientID)
	}
	if cfg.ClientSecret != "test-secret" {
		t.Errorf("ClientSecret = %v, want test-secret", cfg.ClientSecret)
	}
	if cfg.OrganizationID != "test-org" {
		t.Errorf("OrganizationID = %v, want test-org", cfg.OrganizationID)
	}
	if cfg.Region != RegionEU {
		t.Errorf("Region = %v, want EU", cfg.Region)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %v, want 5", cfg.MaxRetries)
	}
}
