package sdk

import (
	"reflect"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		baseURL      string
		tokenURL     string
		clientID     string
		clientSecret string
		scope        string
		retries      int
		httpClient   HTTPClient
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.baseURL, tt.args.tokenURL, tt.args.clientID, tt.args.clientSecret, tt.args.scope, tt.args.retries, tt.args.httpClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient1(t *testing.T) {
	type args struct {
		baseURL      string
		tokenURL     string
		clientID     string
		clientSecret string
		scope        string
		retries      int
		httpClient   HTTPClient
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.baseURL, tt.args.tokenURL, tt.args.clientID, tt.args.clientSecret, tt.args.scope, tt.args.retries, tt.args.httpClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
