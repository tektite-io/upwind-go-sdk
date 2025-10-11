// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListThreatDetections retrieves all threat detections (no pagination for this endpoint)
func (c *Client) ListThreatDetections(ctx context.Context, query *ThreatDetectionsQuery) ([]ThreatDetection, error) {
	if query == nil {
		query = &ThreatDetectionsQuery{}
	}

	urlPath := fmt.Sprintf("%s/organizations/%s/threat-detections", c.config.GetBaseURL(), c.config.OrganizationID)
	queryParams := c.buildThreatDetectionsQueryParams(query)

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

	var detections []ThreatDetection
	if err := json.NewDecoder(resp.Body).Decode(&detections); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return detections, nil
}

// GetThreatDetection retrieves a specific threat detection by ID
func (c *Client) GetThreatDetection(ctx context.Context, detectionID string) (*ThreatDetection, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/threat-detections/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, detectionID)

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
		return nil, fmt.Errorf("threat detection not found: %s", detectionID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var detection ThreatDetection
	if err := json.NewDecoder(resp.Body).Decode(&detection); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &detection, nil
}

// UpdateThreatDetection updates a threat detection (e.g., to archive it)
func (c *Client) UpdateThreatDetection(ctx context.Context, detectionID string, update map[string]interface{}) (*ThreatDetection, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/threat-detections/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, detectionID)

	body, err := json.Marshal(update)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", urlPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("threat detection not found: %s", detectionID)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var detection ThreatDetection
	if err := json.NewDecoder(resp.Body).Decode(&detection); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &detection, nil
}

// ArchiveThreatDetection archives a threat detection
func (c *Client) ArchiveThreatDetection(ctx context.Context, detectionID string) (*ThreatDetection, error) {
	return c.UpdateThreatDetection(ctx, detectionID, map[string]interface{}{
		"status": "ARCHIVED",
	})
}

// ListThreatEvents retrieves threat events with page-based pagination
func (c *Client) ListThreatEvents(ctx context.Context, query *ThreatEventsQuery) ([]ThreatEvent, error) {
	var allEvents []ThreatEvent

	if query == nil {
		query = &ThreatEventsQuery{}
	}

	// Set defaults
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PerPage == 0 {
		query.PerPage = c.config.PageSize
	}

	page := query.Page
	for {
		query.Page = page
		events, hasMore, err := c.listThreatEventsPage(ctx, query)
		if err != nil {
			return nil, err
		}

		allEvents = append(allEvents, events...)

		if !hasMore || len(events) == 0 {
			break
		}
		page++
	}

	return allEvents, nil
}

// listThreatEventsPage retrieves a single page of threat events
func (c *Client) listThreatEventsPage(ctx context.Context, query *ThreatEventsQuery) ([]ThreatEvent, bool, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/threat-events", c.config.GetBaseURL(), c.config.OrganizationID)
	queryParams := c.buildThreatEventsQueryParams(query)

	if len(queryParams) > 0 {
		urlPath += "?" + queryParams
	}

	req, err := http.NewRequestWithContext(ctx, "GET", urlPath, nil)
	if err != nil {
		return nil, false, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, false, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, false, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var events []ThreatEvent
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, false, fmt.Errorf("decoding response: %w", err)
	}

	// If we received a full page, there might be more
	hasMore := len(events) == query.PerPage

	return events, hasMore, nil
}

// ListThreatPolicies retrieves all threat policies
func (c *Client) ListThreatPolicies(ctx context.Context, managedBy string) ([]ThreatPolicy, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/threat-policies", c.config.GetBaseURL(), c.config.OrganizationID)

	if managedBy != "" {
		urlPath += "?managed-by=" + url.QueryEscape(managedBy)
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

	var policies []ThreatPolicy
	if err := json.NewDecoder(resp.Body).Decode(&policies); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return policies, nil
}

// UpdateThreatPolicy updates a threat policy (e.g., to enable/disable it)
func (c *Client) UpdateThreatPolicy(ctx context.Context, policyID string, update map[string]interface{}) (*ThreatPolicy, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/threat-policies/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, policyID)

	body, err := json.Marshal(update)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PATCH", urlPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("threat policy not found: %s", policyID)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var policy ThreatPolicy
	if err := json.NewDecoder(resp.Body).Decode(&policy); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &policy, nil
}

// buildThreatDetectionsQueryParams constructs URL query parameters for threat detections
func (c *Client) buildThreatDetectionsQueryParams(query *ThreatDetectionsQuery) string {
	params := url.Values{}

	if query.Severity != "" {
		params.Add("severity", query.Severity)
	}
	if query.Type != "" {
		params.Add("type", query.Type)
	}
	if query.Category != "" {
		params.Add("category", query.Category)
	}
	if query.MinFirstSeenTime != "" {
		params.Add("min-first-seen-time", query.MinFirstSeenTime)
	}
	if query.MaxFirstSeenTime != "" {
		params.Add("max-first-seen-time", query.MaxFirstSeenTime)
	}
	if query.MinLastSeenTime != "" {
		params.Add("min-last-seen-time", query.MinLastSeenTime)
	}
	if query.MaxLastSeenTime != "" {
		params.Add("max-last-seen-time", query.MaxLastSeenTime)
	}

	return params.Encode()
}

// buildThreatEventsQueryParams constructs URL query parameters for threat events
func (c *Client) buildThreatEventsQueryParams(query *ThreatEventsQuery) string {
	params := url.Values{}

	if query.CloudAccountID != "" {
		params.Add("cloud-account-id", query.CloudAccountID)
	}
	if query.Severity != "" {
		params.Add("severity", query.Severity)
	}
	if query.Category != "" {
		params.Add("category", query.Category)
	}
	if query.MinFirstSeenTime != "" {
		params.Add("min-first-seen-time", query.MinFirstSeenTime)
	}
	if query.MaxFirstSeenTime != "" {
		params.Add("max-first-seen-time", query.MaxFirstSeenTime)
	}
	if query.MinLastSeenTime != "" {
		params.Add("min-last-seen-time", query.MinLastSeenTime)
	}
	if query.MaxLastSeenTime != "" {
		params.Add("max-last-seen-time", query.MaxLastSeenTime)
	}
	if query.Page > 0 {
		params.Add("page", fmt.Sprintf("%d", query.Page))
	}
	if query.PerPage > 0 {
		params.Add("per-page", fmt.Sprintf("%d", query.PerPage))
	}

	return params.Encode()
}
