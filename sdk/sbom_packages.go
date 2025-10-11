// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// SbomPackagesQuery represents query parameters for SBOM packages
type SbomPackagesQuery struct {
	CloudAccountID string
	Framework      string
	ImageName      string
	PackageName    string
	PackageManager string
	PackageLicense string
}

// ListSbomPackages retrieves all SBOM packages
func (c *Client) ListSbomPackages(ctx context.Context, query *SbomPackagesQuery) ([]SbomPackage, error) {
	if query == nil {
		query = &SbomPackagesQuery{}
	}

	urlPath := fmt.Sprintf("%s/organizations/%s/sbom-packages", c.config.GetBaseURL(), c.config.OrganizationID)
	queryParams := c.buildSbomPackagesQueryParams(query)

	if len(queryParams) > 0 {
		urlPath += "?" + queryParams
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var packages []SbomPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return packages, nil
}

// GetSbomPackageDetails retrieves detailed information about a specific SBOM package
func (c *Client) GetSbomPackageDetails(ctx context.Context, packageName, version string) (*SbomPackage, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/sbom-packages/%s/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, url.PathEscape(packageName), url.PathEscape(version))

	req, err := http.NewRequestWithContext(ctx, "GET", urlPath, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("SBOM package not found: %s@%s", packageName, version)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var pkg SbomPackage
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &pkg, nil
}

// buildSbomPackagesQueryParams constructs URL query parameters for SBOM packages
func (c *Client) buildSbomPackagesQueryParams(query *SbomPackagesQuery) string {
	params := url.Values{}

	if query.CloudAccountID != "" {
		params.Add("cloud-account-id", query.CloudAccountID)
	}
	if query.Framework != "" {
		params.Add("framework", query.Framework)
	}
	if query.ImageName != "" {
		params.Add("image-name", query.ImageName)
	}
	if query.PackageName != "" {
		params.Add("package-name", query.PackageName)
	}
	if query.PackageManager != "" {
		params.Add("package-manager", query.PackageManager)
	}
	if query.PackageLicense != "" {
		params.Add("package-license", query.PackageLicense)
	}

	return params.Encode()
}
