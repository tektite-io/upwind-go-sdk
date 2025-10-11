// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"testing"
)

func TestExtractNextLink(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		want       string
		wantErr    bool
	}{
		{
			name:       "valid next link",
			linkHeader: `<https://api.upwind.io/v1/organizations/org_123/findings?page-token=abc123>; rel="next"`,
			want:       "https://api.upwind.io/v1/organizations/org_123/findings?page-token=abc123",
			wantErr:    false,
		},
		{
			name:       "multiple links",
			linkHeader: `<https://api.upwind.io/v1/organizations/org_123/findings?page-token=abc123>; rel="first", <https://api.upwind.io/v1/organizations/org_123/findings?page-token=xyz789>; rel="next"`,
			want:       "https://api.upwind.io/v1/organizations/org_123/findings?page-token=xyz789",
			wantErr:    false,
		},
		{
			name:       "no next link",
			linkHeader: `<https://api.upwind.io/v1/organizations/org_123/findings>; rel="first"`,
			want:       "",
			wantErr:    false,
		},
		{
			name:       "empty link header",
			linkHeader: "",
			want:       "",
			wantErr:    false,
		},
		{
			name:       "invalid format",
			linkHeader: "not a valid link header",
			want:       "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractNextLink(tt.linkHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractNextLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractNextLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
