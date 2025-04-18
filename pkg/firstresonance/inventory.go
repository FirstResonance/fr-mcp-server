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

// GetInventoryItem creates a tool to get details of a specific inventory item.
func GetInventoryItem(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_inventory_item",
			mcp.WithDescription(t("TOOL_GET_INVENTORY_ITEM_DESCRIPTION", "Get details of a specific inventory item")),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Inventory item ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			itemID, err := requiredParam[string](request, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			item, resp, err := client.Inventory.Get(ctx, itemID)
			if err != nil {
				return nil, fmt.Errorf("failed to get inventory item: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get inventory item: %s", string(body))), nil
			}

			r, err := json.Marshal(item)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal inventory item: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// ListInventoryItems creates a tool to list and filter inventory items.
func ListInventoryItems(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_inventory_items",
			mcp.WithDescription(t("TOOL_LIST_INVENTORY_ITEMS_DESCRIPTION", "List and filter inventory items")),
			mcp.WithString("location",
				mcp.Description("Filter by location"),
			),
			mcp.WithString("status",
				mcp.Description("Filter by status"),
			),
			mcp.WithString("sort",
				mcp.Description("Sort field"),
			),
			mcp.WithString("direction",
				mcp.Description("Sort direction"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			location, err := OptionalParam[string](request, "location")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			status, err := OptionalParam[string](request, "status")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			sort, err := OptionalParam[string](request, "sort")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			direction, err := OptionalParam[string](request, "direction")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &ListInventoryItemsOptions{
				Location:  location,
				Status:    status,
				Sort:     sort,
				Direction: direction,
				ListOptions: ListOptions{
					Page:    pagination.page,
					PerPage: pagination.perPage,
				},
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			items, resp, err := client.Inventory.List(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list inventory items: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list inventory items: %s", string(body))), nil
			}

			r, err := json.Marshal(items)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// UpdateInventoryItem creates a tool to update an existing inventory item.
func UpdateInventoryItem(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("update_inventory_item",
			mcp.WithDescription(t("TOOL_UPDATE_INVENTORY_ITEM_DESCRIPTION", "Update an existing inventory item")),
			mcp.WithString("item_id",
				mcp.Required(),
				mcp.Description("Inventory item ID to update"),
			),
			mcp.WithNumber("quantity",
				mcp.Description("New quantity"),
			),
			mcp.WithString("location",
				mcp.Description("New location"),
			),
			mcp.WithString("status",
				mcp.Description("New status"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			itemID, err := requiredParam[string](request, "item_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			update := &InventoryItemUpdateRequest{}
			updateNeeded := false

			if quantity, ok, err := OptionalParamOK[float64](request, "quantity"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				q := int(quantity)
				update.Quantity = &q
				updateNeeded = true
			}

			if location, ok, err := OptionalParamOK[string](request, "location"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Location = &location
				updateNeeded = true
			}

			if status, ok, err := OptionalParamOK[string](request, "status"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Status = &status
				updateNeeded = true
			}

			if !updateNeeded {
				return mcp.NewToolResultError("No update parameters provided."), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			updatedItem, resp, err := client.Inventory.Update(ctx, itemID, update)
			if err != nil {
				return nil, fmt.Errorf("failed to update inventory item: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to update inventory item: %s", string(body))), nil
			}

			r, err := json.Marshal(updatedItem)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
} 