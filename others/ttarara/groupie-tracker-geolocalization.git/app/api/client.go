package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// baseURL is the root endpoint for the Groupie Tracker API
	baseURL = "https://groupietrackers.herokuapp.com/api"
	timeout = 5 * time.Second
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient initializes and returns a configured API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
	}
}

// get performs a GET request to the given endpoint and unmarshals JSON into v
func (c *Client) get(endpoint string, v interface{}) error {
	url := c.baseURL + endpoint

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-200 status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return nil
}
