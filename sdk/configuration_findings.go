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
	"strings"
)

// ListConfigurationFindings streams configuration findings page by page via a channel.
// This is memory-efficient for large datasets. The channel will be closed when done.
// Returns an error channel that will receive any error that occurs during streaming.
//
// Example - streaming (memory efficient):
//
//	findings, errCh := client.ListConfigurationFindings(ctx, query)
//	for finding := range findings {
//	    process(finding)
//	}
//	if err := <-errCh; err != nil {
//	    log.Fatal(err)
//	}
//
// Example - collect all (loads everything in memory):
//
//	findingsCh, errCh := client.ListConfigurationFindings(ctx, query)
//	allFindings, err := sdk.CollectAll(ctx, findingsCh, errCh)
func (c *Client) ListConfigurationFindings(ctx context.Context, query *ConfigurationFindingsQuery) (<-chan ConfigurationFinding, <-chan error) {
	findingsCh := make(chan ConfigurationFinding, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(findingsCh)
		defer close(errCh)

		if query == nil {
			query = &ConfigurationFindingsQuery{}
		}

		pageToken := ""
		for {
			findings, nextToken, err := c.listConfigurationFindingsPage(ctx, query, pageToken)
			if err != nil {
				errCh <- err
				return
			}

			for _, finding := range findings {
				select {
				case findingsCh <- finding:
				case <-ctx.Done():
					errCh <- ctx.Err()
					return
				}
			}

			if nextToken == "" {
				break
			}
			pageToken = nextToken
		}
	}()

	return findingsCh, errCh
}

// listConfigurationFindingsPage retrieves a single page of configuration findings
func (c *Client) listConfigurationFindingsPage(ctx context.Context, query *ConfigurationFindingsQuery, pageToken string) ([]ConfigurationFinding, string, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/configuration-findings", c.config.GetBaseURL(), c.config.OrganizationID)
	queryParams := c.buildConfigurationFindingsQueryParams(query, pageToken)

	if len(queryParams) > 0 {
		urlPath += "?" + queryParams
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlPath, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var findings []ConfigurationFinding
	if err := json.NewDecoder(resp.Body).Decode(&findings); err != nil {
		return nil, "", fmt.Errorf("decoding response: %w", err)
	}

	// Extract next page token from Link header
	nextToken, err := extractNextLink(resp.Header.Get("Link"))
	if err != nil {
		return nil, "", fmt.Errorf("parsing pagination link: %w", err)
	}

	// If nextToken is a full URL, extract just the page-token parameter
	if nextToken != "" && strings.Contains(nextToken, "page-token=") {
		parsedURL, err := url.Parse(nextToken)
		if err == nil {
			nextToken = parsedURL.Query().Get("page-token")
		}
	}

	return findings, nextToken, nil
}

// GetConfigurationFinding retrieves a specific configuration finding by ID
func (c *Client) GetConfigurationFinding(ctx context.Context, findingID string, includeCloudAccountTags bool) (*ConfigurationFinding, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/configuration-findings/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, findingID)

	if includeCloudAccountTags {
		urlPath += "?include-cloud-account-tags=true"
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("configuration finding not found: %s", findingID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var finding ConfigurationFinding
	if err := json.NewDecoder(resp.Body).Decode(&finding); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &finding, nil
}

// buildConfigurationFindingsQueryParams constructs URL query parameters for configuration findings
func (c *Client) buildConfigurationFindingsQueryParams(query *ConfigurationFindingsQuery, pageToken string) string {
	params := url.Values{}

	if pageToken != "" {
		params.Add("page-token", pageToken)
	}

	if query.MinLastSeenTime != "" {
		params.Add("min-last-seen-time", query.MinLastSeenTime)
	}
	if query.MaxLastSeenTime != "" {
		params.Add("max-last-seen-time", query.MaxLastSeenTime)
	}
	if query.Status != "" {
		params.Add("status", query.Status)
	}
	if query.Severity != "" {
		params.Add("severity", query.Severity)
	}
	if query.ResourceName != "" {
		params.Add("resource-name", query.ResourceName)
	}
	if query.CheckTitle != "" {
		params.Add("check-title", query.CheckTitle)
	}
	if query.CheckID != "" {
		params.Add("check-id", query.CheckID)
	}
	if query.FrameworkID != "" {
		params.Add("framework-id", query.FrameworkID)
	}
	if query.FrameworkTitle != "" {
		params.Add("framework-title", query.FrameworkTitle)
	}
	if len(query.CloudAccountTags) > 0 {
		params.Add("cloud-account-tags", strings.Join(query.CloudAccountTags, ","))
	}
	if query.IncludeCloudAccountTags {
		params.Add("include-cloud-account-tags", "true")
	}

	return params.Encode()
}
