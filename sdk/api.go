package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) GetVulnerabilityFindings(ctx context.Context, query FindingsQuery) ([]VulnerabilityFinding, error) {
	var allResults []VulnerabilityFinding
	baseURL := fmt.Sprintf("%s/organizations/%s/vulnerability-findings", c.baseURL, c.orgID)

	// Manually construct query parameters
	var queryParams []string
	if query.ImageName != nil {
		queryParams = append(queryParams, "image-name="+*query.ImageName)
	}
	if query.InUse != nil {
		queryParams = append(queryParams, "in-use="+fmt.Sprintf("%v", *query.InUse))
	}
	if query.IngressActiveCommunication != nil {
		queryParams = append(queryParams, "ingress-active-communication="+fmt.Sprintf("%v", *query.IngressActiveCommunication))
	}
	if len(query.Severities) > 0 {
		queryParams = append(queryParams, "severity="+strings.Join(query.Severities, ","))
	}
	if query.PerPage != nil {
		queryParams = append(queryParams, "per-page="+fmt.Sprintf("%d", *query.PerPage))
	}

	queryString := strings.Join(queryParams, "&")
	urlWithQuery := baseURL

	if len(queryParams) > 0 {
		urlWithQuery += "?" + queryString
	}

	for {
		req, err := http.NewRequestWithContext(ctx, "GET", urlWithQuery, nil)

		//log.Printf("GET %s", urlWithQuery)
		//
		//log.Printf("Req %w", req)

		if err != nil {
			return nil, fmt.Errorf("creating request: %w", err)
		}

		resp, err := c.doRequest(ctx, req)

		//log.Printf("response: %+v\n", resp)

		if err != nil {
			return nil, fmt.Errorf("executing request: %w", err)
		}

		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				fmt.Printf("warning: closing response body: %v\n", cerr)
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API call failed: %s - %s", resp.Status, string(body))
		}

		var results []VulnerabilityFinding
		if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
			return nil, fmt.Errorf("decoding response: %w", err)
		}

		allResults = append(allResults, results...)

		// Check for next page
		nextURL, err := extractNextLink(resp.Header.Get("Link"))
		//
		//log.Printf("Next Link: %+v\n", nextURL)

		if err != nil {
			return nil, fmt.Errorf("extracting next link: %w", err)
		}
		if nextURL == "" {
			break
		}

		urlWithQuery = nextURL
	}

	return allResults, nil
}
