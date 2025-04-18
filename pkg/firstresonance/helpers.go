package firstresonance

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// WithPagination returns a ToolOption that adds "page" and "perPage" parameters to the tool.
func WithPagination() mcp.ToolOption {
	return func(tool *mcp.Tool) {
		mcp.WithNumber("page",
			mcp.Description("Page number for pagination (min 1)"),
			mcp.Min(1),
		)(tool)

		mcp.WithNumber("perPage",
			mcp.Description("Results per page for pagination (min 1, max 100)"),
			mcp.Min(1),
			mcp.Max(100),
		)(tool)
	}
}

type paginationParams struct {
	page    int
	perPage int
}

// OptionalPaginationParams returns the "page" and "perPage" parameters from the request,
// or their default values if not present.
func OptionalPaginationParams(r mcp.CallToolRequest) (paginationParams, error) {
	page, err := OptionalIntParamWithDefault(r, "page", 1)
	if err != nil {
		return paginationParams{}, err
	}
	perPage, err := OptionalIntParamWithDefault(r, "perPage", 30)
	if err != nil {
		return paginationParams{}, err
	}
	return paginationParams{
		page:    page,
		perPage: perPage,
	}, nil
}

// requiredParam is a helper function that can be used to fetch a required parameter from the request.
func requiredParam[T comparable](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, fmt.Errorf("missing required parameter: %s", p)
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T", p, zero)
	}

	if r.Params.Arguments[p].(T) == zero {
		return zero, fmt.Errorf("missing required parameter: %s", p)
	}

	return r.Params.Arguments[p].(T), nil
}

// OptionalParam is a helper function that can be used to fetch an optional parameter from the request.
func OptionalParam[T any](r mcp.CallToolRequest, p string) (T, error) {
	var zero T

	// Check if the parameter is present in the request
	if _, ok := r.Params.Arguments[p]; !ok {
		return zero, nil
	}

	// Check if the parameter is of the expected type
	if _, ok := r.Params.Arguments[p].(T); !ok {
		return zero, fmt.Errorf("parameter %s is not of type %T", p, zero)
	}

	return r.Params.Arguments[p].(T), nil
}

// OptionalParamOK is a helper function that can be used to fetch an optional parameter from the request.
func OptionalParamOK[T any](r mcp.CallToolRequest, p string) (value T, ok bool, err error) {
	// Check if the parameter is present in the request
	val, exists := r.Params.Arguments[p]
	if !exists {
		// Not present, return zero value, false, no error
		return
	}

	// Check if the parameter is of the expected type
	value, ok = val.(T)
	if !ok {
		// Present but wrong type
		err = fmt.Errorf("parameter %s is not of type %T", p, value)
		ok = true // Set ok to true because the parameter *was* present, even if wrong type
		return
	}

	// Present and correct type
	ok = true
	return
}

// OptionalIntParam is a helper function that can be used to fetch an optional integer parameter from the request.
func OptionalIntParam(r mcp.CallToolRequest, p string) (int, error) {
	v, err := OptionalParam[float64](r, p)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// OptionalIntParamWithDefault is a helper function that can be used to fetch an optional integer parameter from the request
// with a default value.
func OptionalIntParamWithDefault(r mcp.CallToolRequest, p string, d int) (int, error) {
	v, err := OptionalIntParam(r, p)
	if err != nil {
		return 0, err
	}
	if v == 0 {
		return d, nil
	}
	return v, nil
}
