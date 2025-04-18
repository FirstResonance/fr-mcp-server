package firstresonance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Get retrieves an inventory item by its ID
func (s *InventoryService) Get(ctx context.Context, id string) (*InventoryItem, error) {
	// GraphQL query to fetch an inventory item by ID
	query := `
		query GetInventoryItem($id: ID!) {
			inventoryItem(id: $id) {
				id
				part_id
				location
				quantity
				status
				last_updated
			}
		}
	`

	// Create variables for the query
	variables := map[string]interface{}{
		"id": id,
	}

	// Create the request body
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.client.baseURL+"/graphql", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.client.apiToken)

	// Send the request
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var result struct {
		Data struct {
			InventoryItem *InventoryItem `json:"inventoryItem"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Check for GraphQL errors
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	// Check if the inventory item was found
	if result.Data.InventoryItem == nil {
		return nil, fmt.Errorf("inventory item not found")
	}

	return result.Data.InventoryItem, nil
}

// List retrieves a list of inventory items
func (s *InventoryService) List(ctx context.Context, opts *ListInventoryItemsOptions) ([]*InventoryItem, error) {
	// GraphQL query to fetch a list of inventory items
	query := `
		query ListInventoryItems($part_id: ID, $location: String, $status: String, $page: Int, $perPage: Int) {
			inventoryItems(part_id: $part_id, location: $location, status: $status, page: $page, perPage: $perPage) {
				id
				part_id
				location
				quantity
				status
				last_updated
			}
		}
	`

	// Create variables for the query
	variables := map[string]interface{}{}
	if opts != nil {
		if opts.Location != "" {
			variables["location"] = opts.Location
		}
		if opts.Status != "" {
			variables["status"] = opts.Status
		}
		if opts.Page > 0 {
			variables["page"] = opts.Page
		}
		if opts.PerPage > 0 {
			variables["perPage"] = opts.PerPage
		}
	}

	// Create the request body
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.client.baseURL+"/graphql", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.client.apiToken)

	// Send the request
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var result struct {
		Data struct {
			InventoryItems []*InventoryItem `json:"inventoryItems"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Check for GraphQL errors
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	return result.Data.InventoryItems, nil
}

// Update updates an existing inventory item
func (s *InventoryService) Update(ctx context.Context, id string, update *InventoryItemUpdateRequest) (*InventoryItem, error) {
	// GraphQL mutation to update an existing inventory item
	query := `
		mutation UpdateInventoryItem($id: ID!, $input: UpdateInventoryItemInput!) {
			updateInventoryItem(id: $id, input: $input) {
				id
				part_id
				location
				quantity
				status
				last_updated
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"id": id,
		"input": map[string]interface{}{
			"location": update.Location,
			"quantity": update.Quantity,
			"status":   update.Status,
		},
	}

	// Create the request body
	requestBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", s.client.baseURL+"/graphql", strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.client.apiToken)

	// Send the request
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var result struct {
		Data struct {
			UpdateInventoryItem *InventoryItem `json:"updateInventoryItem"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Check for GraphQL errors
	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", result.Errors[0].Message)
	}

	return result.Data.UpdateInventoryItem, nil
}
