// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"testing"
)

func TestNewClient(t *testing.T) {
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
			name: "invalid config - missing client ID",
			config: &Config{
				ClientSecret:   "test-secret",
				OrganizationID: "test-org",
				Region:         RegionUS,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
			if !tt.wantErr && client != nil {
				// Verify client has correct config
				if client.GetOrganizationID() != tt.config.OrganizationID {
					t.Errorf("NewClient() organization ID = %v, want %v", client.GetOrganizationID(), tt.config.OrganizationID)
				}
			}
		})
	}
}

func TestClientLogging(t *testing.T) {
	cfg := &Config{
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		OrganizationID: "test-org",
		Region:         RegionUS,
		MaxRetries:     3,
		MaxConcurrency: 10,
		PageSize:       100,
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test enabling logging
	client.EnableLogging()

	// Test setting custom logger
	customLogger := &NoOpLogger{}
	client.SetLogger(customLogger)
}
