package fandom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(wikiHost string) *Client {
	return &Client{
		http:    &http.Client{Timeout: 15 * time.Second},
		baseURL: fmt.Sprintf("https://%s/api.php", wikiHost),
	}
}

type ParseResponse struct {
	Parse struct {
		Title      string            `json:"title"`
		Wikitext   map[string]string `json:"wikitext"`
		Properties []struct {
			Name string `json:"name"`
			Data string `json:"*"`
		} `json:"properties"`
	} `json:"parse"`
}

type FandomData struct {
	Title     string
	Wikitext  map[string]string
	Infoboxes map[string]any
}

func (c *Client) FetchPage(ctx context.Context, pageTitle string) (*ParseResponse, error) {
	params := url.Values{
		"action": {"parse"},
		"page":   {pageTitle},
		"prop":   {"wikitext|properties"},
		"format": {"json"},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("User-Agent", "MyScraper/1.0 (contact@example.com)")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result ParseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &result, nil
}
