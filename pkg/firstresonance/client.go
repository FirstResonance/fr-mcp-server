package firstresonance

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

// cacheItem represents a cached item with a timestamp
type cacheItem struct {
	value     interface{}
	timestamp time.Time
}

// NewClient creates a new First Resonance API client
func NewClient(baseURL string, apiToken string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
			},
		}
	}

	// Ensure baseURL ends with a slash
	if baseURL != "" && baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	client := &Client{
		baseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: httpClient,
		cache:      &sync.Map{},
		cacheTTL:   5 * time.Minute,
	}

	// Initialize services
	client.Parts = &PartsService{client: client}
	client.Orders = &OrdersService{client: client}
	client.Suppliers = &SuppliersService{client: client}
	client.Inventory = &InventoryService{client: client}
	client.Search = &SearchService{client: client}
	client.ABom = &ABomService{client: client}

	return client
}

// Client represents the First Resonance API client
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	Parts      *PartsService
	Orders     *OrdersService
	Suppliers  *SuppliersService
	Inventory  *InventoryService
	Search     *SearchService
	ABom       *ABomService
	cache      *sync.Map
	cacheTTL   time.Duration
}

// SetCacheTTL sets the cache time-to-live
func (c *Client) SetCacheTTL(ttl time.Duration) {
	c.cacheTTL = ttl
}

// ClearCache clears the entire cache
func (c *Client) ClearCache() {
	c.cache = &sync.Map{}
}

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: %d - %s", e.StatusCode, e.Message)
}

// buildURL builds a URL from the base URL and path segments
func (c *Client) buildURL(pathSegments ...string) string {
	// Join path segments with slashes
	pathStr := path.Join(pathSegments...)

	// Ensure the path starts with a slash
	if !strings.HasPrefix(pathStr, "/") {
		pathStr = "/" + pathStr
	}

	// Combine with base URL
	return c.baseURL + strings.TrimPrefix(pathStr, "/")
}

// NewRequest creates a new HTTP request
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.buildURL(path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")
	if reqBody != nil {
		req.Header.Set("Content-Length", fmt.Sprintf("%d", reqBody.(*bytes.Buffer).Len()))
	}

	return req, nil
}

// Do sends an HTTP request and returns the response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}
