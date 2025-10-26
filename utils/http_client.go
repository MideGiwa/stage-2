package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is a wrapper for http.Client to make external API calls
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTPClient with a default timeout
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: &http.Client{
			Timeout: 30 * time.Second, // Default timeout
		},
	}
}

// Get performs a GET request to the specified URL and decodes the JSON response into the target interface
func (c *HTTPClient) Get(url string, target interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request to %s failed with status %d: %s", url, resp.StatusCode, string(bodyBytes))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response from %s: %w", url, err)
	}

	return nil
}
