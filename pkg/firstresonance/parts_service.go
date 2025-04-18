package firstresonance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Get retrieves a part by its ID
func (s *PartsService) Get(ctx context.Context, id string) (*Part, error) {
	// GraphQL query to fetch a part by ID
	query := `
		query GetPart($id: ID!) {
			part(id: $id) {
				id
				name
				description
				type
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
			Part *Part `json:"part"`
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

	// Check if the part was found
	if result.Data.Part == nil {
		return nil, fmt.Errorf("part not found")
	}

	return result.Data.Part, nil
}

// List retrieves a list of parts
func (s *PartsService) List(ctx context.Context, opts *ListPartsOptions) ([]*Part, error) {
	// GraphQL query to fetch a list of parts
	query := `
		query ListParts($status: String, $type: String, $page: Int, $perPage: Int) {
			parts(status: $status, type: $type, page: $page, perPage: $perPage) {
				id
				name
				description
				type
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
		if opts.Type != "" {
			variables["type"] = opts.Type
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
			Parts []*Part `json:"parts"`
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

	return result.Data.Parts, nil
}

// Create creates a new part
func (s *PartsService) Create(ctx context.Context, part *Part) (*Part, error) {
	// GraphQL mutation to create a new part
	query := `
		mutation CreatePart($input: CreatePartInput!) {
			createPart(input: $input) {
				id
				name
				description
				type
				status
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"name":        part.Name,
			"description": part.Description,
			"type":        part.Type,
			"status":      part.Status,
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
			CreatePart *Part `json:"createPart"`
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

	return result.Data.CreatePart, nil
}

// Update updates an existing part
func (s *PartsService) Update(ctx context.Context, id string, update *PartUpdateRequest) (*Part, error) {
	// GraphQL mutation to update an existing part
	query := `
		mutation UpdatePart($id: ID!, $input: UpdatePartInput!) {
			updatePart(id: $id, input: $input) {
				id
				name
				description
				type
				status
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"id": id,
		"input": map[string]interface{}{
			"name":        update.Name,
			"description": update.Description,
			"type":        update.Type,
			"status":      update.Status,
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
			UpdatePart *Part `json:"updatePart"`
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

	return result.Data.UpdatePart, nil
}
