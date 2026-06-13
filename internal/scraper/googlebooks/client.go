package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://www.googleapis.com/books/v1/volumes"

type Client struct {
	http   *http.Client
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		http:   &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

type Response struct {
	Items []Item `json:"items"`
}

type Item struct {
	SelfLink   string     `json:"selfLink"`
	VolumeInfo VolumeInfo `json:"volumeInfo"`
}

type VolumeInfo struct {
	Title               string   `json:"title"`
	Subtitle            string   `json:"subtitle"`
	Authors             []string `json:"authors"`
	Publisher           string   `json:"publisher"`
	PublishedDate       string   `json:"publishedDate"`
	Description         string   `json:"description"`
	PageCount           int      `json:"pageCount"`
	Language            string   `json:"language"`
	PrintType           string   `json:"printType"`
	IndustryIdentifiers []struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"industryIdentifiers"`
	ImageLinks struct {
		Thumbnail      string `json:"thumbnail"`
		SmallThumbnail string `json:"smallThumbnail"`
	} `json:"imageLinks"`
}

func (c *Client) FetchByISBN(ctx context.Context, isbn string) (*VolumeInfo, error) {
	url := baseURL + "?q=isbn:" + isbn
	if c.apiKey != "" {
		url += "&key=" + c.apiKey
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}

	if len(data.Items) == 0 {
		return nil, fmt.Errorf("isbn %s not found", isbn)
	}

	// Exit if no self link exist
	if data.Items[0].SelfLink == "" {
		return &data.Items[0].VolumeInfo, nil
	}

	// If we have the response, we do a second request on the self link for more accurate data
	url = data.Items[0].SelfLink + "?key=" + c.apiKey
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err = c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var refinedData Item
	if err := json.NewDecoder(resp.Body).Decode(&refinedData); err != nil {
		return nil, fmt.Errorf("decoding: %w", err)
	}
	defer resp.Body.Close()

	return &refinedData.VolumeInfo, nil
}
