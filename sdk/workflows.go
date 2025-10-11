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
)

// ListWorkflows retrieves all workflows
func (c *Client) ListWorkflows(ctx context.Context) ([]Workflow, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/workflows", c.config.GetBaseURL(), c.config.OrganizationID)

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

	var workflows []Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflows); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return workflows, nil
}

// GetWorkflow retrieves a specific workflow by ID
func (c *Client) GetWorkflow(ctx context.Context, workflowID string) (*Workflow, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/workflows/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, workflowID)

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
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &workflow, nil
}

// CreateWorkflow creates a new workflow
func (c *Client) CreateWorkflow(ctx context.Context, workflow map[string]interface{}) (*Workflow, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/workflows", c.config.GetBaseURL(), c.config.OrganizationID)

	body, err := json.Marshal(workflow)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", urlPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var createdWorkflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&createdWorkflow); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &createdWorkflow, nil
}

// UpdateWorkflow updates an existing workflow
func (c *Client) UpdateWorkflow(ctx context.Context, workflowID string, update map[string]interface{}) (*Workflow, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/workflows/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, workflowID)

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
		return nil, fmt.Errorf("workflow not found: %s", workflowID)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &workflow, nil
}

// DeleteWorkflow deletes a workflow
func (c *Client) DeleteWorkflow(ctx context.Context, workflowID string) error {
	urlPath := fmt.Sprintf("%s/organizations/%s/workflows/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, workflowID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", urlPath, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("workflow not found: %s", workflowID)
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ListIntegrationWebhooks retrieves all integration webhooks
func (c *Client) ListIntegrationWebhooks(ctx context.Context, vendor string) ([]IntegrationWebhook, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/integration-webhooks", c.config.GetBaseURL(), c.config.OrganizationID)

	if vendor != "" {
		urlPath += "?vendor=" + vendor
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

	var webhooks []IntegrationWebhook
	if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return webhooks, nil
}

// CreateIntegrationWebhook creates a new integration webhook
func (c *Client) CreateIntegrationWebhook(ctx context.Context, webhook map[string]interface{}) (*IntegrationWebhook, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/integration-webhooks", c.config.GetBaseURL(), c.config.OrganizationID)

	body, err := json.Marshal(webhook)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", urlPath, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var createdWebhook IntegrationWebhook
	if err := json.NewDecoder(resp.Body).Decode(&createdWebhook); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &createdWebhook, nil
}

// UpdateIntegrationWebhook updates an existing integration webhook
func (c *Client) UpdateIntegrationWebhook(ctx context.Context, webhookID string, update map[string]interface{}) (*IntegrationWebhook, error) {
	urlPath := fmt.Sprintf("%s/organizations/%s/integration-webhooks/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, webhookID)

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
		return nil, fmt.Errorf("integration webhook not found: %s", webhookID)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBody))
	}

	var webhook IntegrationWebhook
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &webhook, nil
}

// DeleteIntegrationWebhook deletes an integration webhook
func (c *Client) DeleteIntegrationWebhook(ctx context.Context, webhookID string) error {
	urlPath := fmt.Sprintf("%s/organizations/%s/integration-webhooks/%s",
		c.config.GetBaseURL(), c.config.OrganizationID, webhookID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", urlPath, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.doRequest(ctx, req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("integration webhook not found: %s", webhookID)
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
