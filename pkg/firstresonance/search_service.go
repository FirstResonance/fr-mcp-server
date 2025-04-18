package firstresonance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Search performs a search across all entities
func (s *SearchService) Search(ctx context.Context, opts *SearchOptions) ([]interface{}, error) {
	// GraphQL query to search across all entities
	query := `
		query Search($query: String!, $sort: String, $order: String, $page: Int, $perPage: Int) {
			search(query: $query, sort: $sort, order: $order, page: $page, perPage: $perPage) {
				id
				type
				name
				description
				status
				# Additional fields based on type
				... on Part {
					type
				}
				... on Order {
					customer_id
					priority
					due_date
				}
				... on Supplier {
					contact_info
				}
				... on InventoryItem {
					quantity
					location
				}
			}
		}
	`

	// Create variables for the query
	variables := map[string]interface{}{
		"query": opts.Query,
	}

	if opts.Sort != "" {
		variables["sort"] = opts.Sort
	}
	if opts.Order != "" {
		variables["order"] = opts.Order
	}
	if opts.Page > 0 {
		variables["page"] = opts.Page
	}
	if opts.PerPage > 0 {
		variables["perPage"] = opts.PerPage
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
			Search []map[string]interface{} `json:"search"`
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

	// Convert the search results to the appropriate types based on the "type" field
	var results []interface{}
	for _, item := range result.Data.Search {
		itemType, ok := item["type"].(string)
		if !ok {
			continue
		}

		switch itemType {
		case "part":
			var part Part
			if err := mapToStruct(item, &part); err != nil {
				return nil, fmt.Errorf("error converting search result to Part: %w", err)
			}
			results = append(results, &part)
		case "order":
			var order Order
			if err := mapToStruct(item, &order); err != nil {
				return nil, fmt.Errorf("error converting search result to Order: %w", err)
			}
			results = append(results, &order)
		case "supplier":
			var supplier Supplier
			if err := mapToStruct(item, &supplier); err != nil {
				return nil, fmt.Errorf("error converting search result to Supplier: %w", err)
			}
			results = append(results, &supplier)
		case "inventory_item":
			var inventoryItem InventoryItem
			if err := mapToStruct(item, &inventoryItem); err != nil {
				return nil, fmt.Errorf("error converting search result to InventoryItem: %w", err)
			}
			results = append(results, &inventoryItem)
		default:
			// For unknown types, just add the raw map
			results = append(results, item)
		}
	}

	return results, nil
}

// mapToStruct converts a map to a struct using JSON marshaling/unmarshaling
func mapToStruct(input map[string]interface{}, output interface{}) error {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, output)
}
