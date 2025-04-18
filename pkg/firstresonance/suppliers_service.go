package firstresonance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Get retrieves a supplier by its ID
func (s *SuppliersService) Get(ctx context.Context, id string) (*Supplier, error) {
	// GraphQL query to fetch a supplier by ID
	query := `
		query GetSupplier($id: ID!) {
			supplier(id: $id) {
				id
				name
				contact_info
				status
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
			Supplier *Supplier `json:"supplier"`
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

	// Check if the supplier was found
	if result.Data.Supplier == nil {
		return nil, fmt.Errorf("supplier not found")
	}

	return result.Data.Supplier, nil
}

// List retrieves a list of suppliers
func (s *SuppliersService) List(ctx context.Context, opts *ListSuppliersOptions) ([]*Supplier, error) {
	// GraphQL query to fetch a list of suppliers
	query := `
		query ListSuppliers($status: String, $sort: String, $direction: String, $page: Int, $perPage: Int) {
			suppliers(status: $status, sort: $sort, direction: $direction, page: $page, perPage: $perPage) {
				id
				name
				contact_info
				status
			}
		}
	`

	// Create variables for the query
	variables := map[string]interface{}{}
	if opts != nil {
		if opts.Status != "" {
			variables["status"] = opts.Status
		}
		if opts.Sort != "" {
			variables["sort"] = opts.Sort
		}
		if opts.Direction != "" {
			variables["direction"] = opts.Direction
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
			Suppliers []*Supplier `json:"suppliers"`
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

	return result.Data.Suppliers, nil
}

// Create creates a new supplier
func (s *SuppliersService) Create(ctx context.Context, supplier *Supplier) (*Supplier, error) {
	// GraphQL mutation to create a new supplier
	query := `
		mutation CreateSupplier($input: CreateSupplierInput!) {
			createSupplier(input: $input) {
				id
				name
				contact_info
				status
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"name":         supplier.Name,
			"contact_info": supplier.ContactInfo,
			"status":       supplier.Status,
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
			CreateSupplier *Supplier `json:"createSupplier"`
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

	return result.Data.CreateSupplier, nil
}

// Update updates an existing supplier
func (s *SuppliersService) Update(ctx context.Context, id string, update *SupplierUpdateRequest) (*Supplier, error) {
	// GraphQL mutation to update an existing supplier
	query := `
		mutation UpdateSupplier($id: ID!, $input: UpdateSupplierInput!) {
			updateSupplier(id: $id, input: $input) {
				id
				name
				contact_info
				status
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"id": id,
		"input": map[string]interface{}{
			"name":         update.Name,
			"contact_info": update.ContactInfo,
			"status":       update.Status,
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
			UpdateSupplier *Supplier `json:"updateSupplier"`
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

	return result.Data.UpdateSupplier, nil
}
