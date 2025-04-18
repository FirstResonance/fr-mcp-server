package firstresonance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SearchParts creates a tool to search for parts across First Resonance.
func SearchParts(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("search_parts",
			mcp.WithDescription(t("TOOL_SEARCH_PARTS_DESCRIPTION", "Search for parts across First Resonance")),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Search query"),
			),
			mcp.WithString("sort",
				mcp.Description("Sort field"),
			),
			mcp.WithString("order",
				mcp.Description("Sort order"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, err := requiredParam[string](request, "query")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			sort, err := OptionalParam[string](request, "sort")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			order, err := OptionalParam[string](request, "order")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &SearchOptions{
				Query: query,
				Sort:  sort,
				Order: order,
				ListOptions: ListOptions{
					Page:    pagination.page,
					PerPage: pagination.perPage,
				},
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			result, resp, err := client.Search.Parts(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to search parts: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to search parts: %s", string(body))), nil
			}

			r, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// SearchOrders creates a tool to search for orders.
func SearchOrders(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("search_orders",
			mcp.WithDescription(t("TOOL_SEARCH_ORDERS_DESCRIPTION", "Search for orders")),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Search query"),
			),
			mcp.WithString("sort",
				mcp.Description("Sort field"),
			),
			mcp.WithString("order",
				mcp.Description("Sort order"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, err := requiredParam[string](request, "query")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			sort, err := OptionalParam[string](request, "sort")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			order, err := OptionalParam[string](request, "order")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &SearchOptions{
				Query: query,
				Sort:  sort,
				Order: order,
				ListOptions: ListOptions{
					Page:    pagination.page,
					PerPage: pagination.perPage,
				},
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			result, resp, err := client.Search.Orders(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to search orders: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to search orders: %s", string(body))), nil
			}

			r, err := json.Marshal(result)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
} 