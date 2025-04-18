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

// GetPartContent creates a resource template to get part content.
func GetPartContent(getClient GetClientFn, t TranslationHelperFunc) server.ResourceTemplate {
	return server.ResourceTemplate{
		Template: "part://{part_id}",
		Description: t("RESOURCE_GET_PART_CONTENT_DESCRIPTION",
			"Retrieves the content of a part"),
		Parameters: map[string]server.ResourceParameter{
			"part_id": {
				Description: "Part ID",
				Required:   true,
			},
		},
		Handler: func(ctx context.Context, params map[string]string) (*mcp.ResourceResult, error) {
			partID := params["part_id"]
			if partID == "" {
				return nil, fmt.Errorf("part_id is required")
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			part, resp, err := client.Parts.Get(ctx, partID)
			if err != nil {
				return nil, fmt.Errorf("failed to get part: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return nil, fmt.Errorf("failed to get part: %s", string(body))
			}

			r, err := json.Marshal(part)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal part: %w", err)
			}

			return &mcp.ResourceResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: string(r),
					},
				},
			}, nil
		},
	}
}

// GetOrderContent creates a resource template to get order content.
func GetOrderContent(getClient GetClientFn, t TranslationHelperFunc) server.ResourceTemplate {
	return server.ResourceTemplate{
		Template: "order://{order_id}",
		Description: t("RESOURCE_GET_ORDER_CONTENT_DESCRIPTION",
			"Retrieves the content of an order"),
		Parameters: map[string]server.ResourceParameter{
			"order_id": {
				Description: "Order ID",
				Required:   true,
			},
		},
		Handler: func(ctx context.Context, params map[string]string) (*mcp.ResourceResult, error) {
			orderID := params["order_id"]
			if orderID == "" {
				return nil, fmt.Errorf("order_id is required")
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			order, resp, err := client.Orders.Get(ctx, orderID)
			if err != nil {
				return nil, fmt.Errorf("failed to get order: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return nil, fmt.Errorf("failed to get order: %s", string(body))
			}

			r, err := json.Marshal(order)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal order: %w", err)
			}

			return &mcp.ResourceResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: string(r),
					},
				},
			}, nil
		},
	}
}

// GetSupplierContent creates a resource template to get supplier content.
func GetSupplierContent(getClient GetClientFn, t TranslationHelperFunc) server.ResourceTemplate {
	return server.ResourceTemplate{
		Template: "supplier://{supplier_id}",
		Description: t("RESOURCE_GET_SUPPLIER_CONTENT_DESCRIPTION",
			"Retrieves the content of a supplier"),
		Parameters: map[string]server.ResourceParameter{
			"supplier_id": {
				Description: "Supplier ID",
				Required:   true,
			},
		},
		Handler: func(ctx context.Context, params map[string]string) (*mcp.ResourceResult, error) {
			supplierID := params["supplier_id"]
			if supplierID == "" {
				return nil, fmt.Errorf("supplier_id is required")
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			supplier, resp, err := client.Suppliers.Get(ctx, supplierID)
			if err != nil {
				return nil, fmt.Errorf("failed to get supplier: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return nil, fmt.Errorf("failed to get supplier: %s", string(body))
			}

			r, err := json.Marshal(supplier)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal supplier: %w", err)
			}

			return &mcp.ResourceResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: string(r),
					},
				},
			}, nil
		},
	}
}

// GetInventoryItemContent creates a resource template to get inventory item content.
func GetInventoryItemContent(getClient GetClientFn, t TranslationHelperFunc) server.ResourceTemplate {
	return server.ResourceTemplate{
		Template: "inventory://{item_id}",
		Description: t("RESOURCE_GET_INVENTORY_ITEM_CONTENT_DESCRIPTION",
			"Retrieves the content of an inventory item"),
		Parameters: map[string]server.ResourceParameter{
			"item_id": {
				Description: "Inventory item ID",
				Required:   true,
			},
		},
		Handler: func(ctx context.Context, params map[string]string) (*mcp.ResourceResult, error) {
			itemID := params["item_id"]
			if itemID == "" {
				return nil, fmt.Errorf("item_id is required")
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
				return nil, fmt.Errorf("failed to get inventory item: %s", string(body))
			}

			r, err := json.Marshal(item)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal inventory item: %w", err)
			}

			return &mcp.ResourceResult{
				Content: []mcp.Content{
					mcp.TextContent{
						Type: "text",
						Text: string(r),
					},
				},
			}, nil
		},
	}
} 