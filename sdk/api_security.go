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

// ListApiEndpoints streams API endpoints page by page via a channel.
// This is memory-efficient for large datasets. The channel will be closed when done.
// Returns an error channel that will receive any error that occurs during streaming.
//
// Example - streaming (memory efficient):
//
//	endpoints, errCh := client.ListApiEndpoints(ctx, query)
//	for endpoint := range endpoints {
//	    process(endpoint)
//	}
//	if err := <-errCh; err != nil {
//	    log.Fatal(err)
//	}
//
// Example - collect all (loads everything in memory):
//
//	endpointsCh, errCh := client.ListApiEndpoints(ctx, query)
//	allEndpoints, err := sdk.CollectAll(ctx, endpointsCh, errCh)
func (c *Client) ListApiEndpoints(ctx context.Context, query *ApiEndpointsQuery) (<-chan ApiEndpoint, <-chan error) {
	endpointsCh := make(chan ApiEndpoint, 100)
	errCh := make(chan error, 1)

	go func() {
		defer close(endpointsCh)
		defer close(errCh)

		if query == nil {
			query = &ApiEndpointsQuery{}
		}

		if query.PerPage == 0 {
			query.PerPage = c.config.PageSize
		}

		pageToken := ""
		for {
			endpoints, nextToken, err := c.listApiEndpointsPage(ctx, query, pageToken)
			if err != nil {
				errCh <- err
				return
			}

			for _, endpoint := range endpoints {
				select {
				case endpointsCh <- endpoint:
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

	return endpointsCh, errCh
}

// listApiEndpointsPage retrieves a single page of API endpoints
func (c *Client) listApiEndpointsPage(ctx context.Context, query *ApiEndpointsQuery, pageToken string) ([]ApiEndpoint, string, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/apisecurity-endpoints", c.config.GetBaseURL(), c.config.OrganizationID)
	queryParams := c.buildApiEndpointsQueryParams(query, pageToken)

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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		body, _ := io.ReadAll(resp.Body)
		return nil, "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var endpoints []ApiEndpoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
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

	return endpoints, nextToken, nil
}

// buildApiEndpointsQueryParams constructs URL query parameters for API endpoints
func (c *Client) buildApiEndpointsQueryParams(query *ApiEndpointsQuery, pageToken string) string {
	params := url.Values{}

	if pageToken != "" {
		params.Add("page-token", pageToken)
	} else if query.PageToken != "" {
		params.Add("page-token", query.PageToken)
	}

	if query.PerPage > 0 {
		params.Add("per-page", fmt.Sprintf("%d", query.PerPage))
	}
	if query.Method != "" {
		params.Add("method", query.Method)
	}
	if query.AuthenticationState != "" {
		params.Add("authentication-state", query.AuthenticationState)
	}
	if query.HasInternetIngress != nil {
		params.Add("has-internet-ingress", fmt.Sprintf("%t", *query.HasInternetIngress))
	}
	if query.HasVulnerability != nil {
		params.Add("has-vulnerability", fmt.Sprintf("%t", *query.HasVulnerability))
	}
	if query.HasSensitiveData != nil {
		params.Add("has-sensitive-data", fmt.Sprintf("%t", *query.HasSensitiveData))
	}
	if query.CloudAccountID != "" {
		params.Add("cloud-account-id", query.CloudAccountID)
	}
	if query.CloudProvider != "" {
		params.Add("cloud-provider", query.CloudProvider)
	}
	if query.ResourceType != "" {
		params.Add("resource-type", query.ResourceType)
	}
	if query.CloudOrganizationID != "" {
		params.Add("cloud-organization-id", query.CloudOrganizationID)
	}
	if query.CloudOrganizationUnitID != "" {
		params.Add("cloud-organization-unit-id", query.CloudOrganizationUnitID)
	}
	if query.Domain != "" {
		params.Add("domain", query.Domain)
	}
	if query.ClusterID != "" {
		params.Add("cluster-id", query.ClusterID)
	}
	if query.Namespace != "" {
		params.Add("namespace", query.Namespace)
	}

	return params.Encode()
}
