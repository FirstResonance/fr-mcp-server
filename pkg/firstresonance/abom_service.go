package firstresonance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Get retrieves an ABOM by its ID
func (s *ABomService) Get(ctx context.Context, id string) (*ABom, error) {
	// GraphQL query to fetch an ABOM by ID
	query := `
		query GetABom($id: ID!) {
			abom(id: $id) {
				id
				name
				description
				version
				status
				created_at
				updated_at
				items {
					id
					part_id
					quantity
					unit
					notes
				}
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
			ABom *ABom `json:"abom"`
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

	// Check if the ABOM was found
	if result.Data.ABom == nil {
		return nil, fmt.Errorf("ABOM not found")
	}

	return result.Data.ABom, nil
}

// List retrieves a list of ABOMs
func (s *ABomService) List(ctx context.Context, opts *ListABomsOptions) ([]*ABom, error) {
	// GraphQL query to fetch a list of ABOMs
	query := `
		query ListABoms($status: String, $sort: String, $direction: String, $page: Int, $perPage: Int) {
			aboms(status: $status, sort: $sort, direction: $direction, page: $page, perPage: $perPage) {
				id
				name
				description
				version
				status
				created_at
				updated_at
				items {
					id
					part_id
					quantity
					unit
					notes
				}
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
			ABoms []*ABom `json:"aboms"`
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

	return result.Data.ABoms, nil
}

// Create creates a new ABOM
func (s *ABomService) Create(ctx context.Context, abom *ABom) (*ABom, error) {
	// GraphQL mutation to create a new ABOM
	query := `
		mutation CreateABom($input: CreateABomInput!) {
			createABom(input: $input) {
				id
				name
				description
				version
				status
				created_at
				updated_at
				items {
					id
					part_id
					quantity
					unit
					notes
				}
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"name":        abom.Name,
			"description": abom.Description,
			"version":     abom.Version,
			"status":      abom.Status,
			"items":       abom.Items,
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
			CreateABom *ABom `json:"createABom"`
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

	return result.Data.CreateABom, nil
}

// Update updates an existing ABOM
func (s *ABomService) Update(ctx context.Context, id string, update *ABomUpdateRequest) (*ABom, error) {
	// GraphQL mutation to update an existing ABOM
	query := `
		mutation UpdateABom($id: ID!, $input: UpdateABomInput!) {
			updateABom(id: $id, input: $input) {
				id
				name
				description
				version
				status
				created_at
				updated_at
				items {
					id
					part_id
					quantity
					unit
					notes
				}
			}
		}
	`

	// Create variables for the mutation
	variables := map[string]interface{}{
		"id": id,
		"input": map[string]interface{}{
			"name":        update.Name,
			"description": update.Description,
			"version":     update.Version,
			"status":      update.Status,
			"items":       update.Items,
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
			UpdateABom *ABom `json:"updateABom"`
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

	return result.Data.UpdateABom, nil
}
