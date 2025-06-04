package sdk

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	baseURL    string
	httpClient HTTPClient
	oauthCfg   *clientcredentials.Config
	tokenSrc   oauth2.TokenSource
	tokenMu    sync.Mutex
	token      *oauth2.Token
	retries    int
}

func NewClient(baseURL, tokenURL, clientID, clientSecret, scope string, retries int, httpClient HTTPClient) *Client {
	cfg := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       []string{scope},
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	tokenSrc := cfg.TokenSource(context.Background())

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		oauthCfg:   cfg,
		tokenSrc:   tokenSrc,
		retries:    retries,
	}
}

func (c *Client) getToken(ctx context.Context) (*oauth2.Token, error) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if c.token != nil && c.token.Valid() && time.Until(c.token.Expiry) > 60*time.Second {
		return c.token, nil
	}

	token, err := c.tokenSrc.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	c.token = token
	return token, nil
}

func (c *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	for i := 0; i <= c.retries; i++ {
		token, err := c.getToken(ctx)

		//log.Printf("token: %+v\n", token)

		if err != nil {
			return nil, fmt.Errorf("token error: %w", err)
		}

		req = req.Clone(ctx)
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}

		if resp.StatusCode == http.StatusUnauthorized && i < c.retries {
			c.tokenMu.Lock()
			c.token = nil
			c.tokenMu.Unlock()
			_ = resp.Body.Close()
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("all retries failed: %w", lastErr)
}
