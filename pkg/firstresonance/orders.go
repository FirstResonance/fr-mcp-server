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

// GetOrder creates a tool to get details of a specific order.
func GetOrder(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_order",
			mcp.WithDescription(t("TOOL_GET_ORDER_DESCRIPTION", "Get details of a specific order")),
			mcp.WithString("order_id",
				mcp.Required(),
				mcp.Description("Order ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			orderID, err := requiredParam[string](request, "order_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
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
				return mcp.NewToolResultError(fmt.Sprintf("failed to get order: %s", string(body))), nil
			}

			r, err := json.Marshal(order)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal order: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// ListOrders creates a tool to list and filter orders.
func ListOrders(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_orders",
			mcp.WithDescription(t("TOOL_LIST_ORDERS_DESCRIPTION", "List and filter orders")),
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

			opts := &ListOrdersOptions{
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
			orders, resp, err := client.Orders.List(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list orders: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list orders: %s", string(body))), nil
			}

			r, err := json.Marshal(orders)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// CreateOrder creates a tool to create a new order.
func CreateOrder(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_order",
			mcp.WithDescription(t("TOOL_CREATE_ORDER_DESCRIPTION", "Create a new order")),
			mcp.WithString("customer_id",
				mcp.Required(),
				mcp.Description("Customer ID"),
			),
			mcp.WithArray("items",
				mcp.Required(),
				mcp.Description("Order items"),
			),
			mcp.WithString("priority",
				mcp.Description("Order priority"),
			),
			mcp.WithString("due_date",
				mcp.Description("Due date"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			customerID, err := requiredParam[string](request, "customer_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			items, err := requiredParam[[]interface{}](request, "items")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			priority, err := OptionalParam[string](request, "priority")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			dueDate, err := OptionalParam[string](request, "due_date")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			order := &Order{
				CustomerID: customerID,
				Items:     items,
				Priority:  priority,
				DueDate:   dueDate,
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			createdOrder, resp, err := client.Orders.Create(ctx, order)
			if err != nil {
				return nil, fmt.Errorf("failed to create order: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusCreated {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to create order: %s", string(body))), nil
			}

			r, err := json.Marshal(createdOrder)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// UpdateOrder creates a tool to update an existing order.
func UpdateOrder(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("update_order",
			mcp.WithDescription(t("TOOL_UPDATE_ORDER_DESCRIPTION", "Update an existing order")),
			mcp.WithString("order_id",
				mcp.Required(),
				mcp.Description("Order ID to update"),
			),
			mcp.WithString("status",
				mcp.Description("New status"),
			),
			mcp.WithString("priority",
				mcp.Description("New priority"),
			),
			mcp.WithString("due_date",
				mcp.Description("New due date"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			orderID, err := requiredParam[string](request, "order_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			update := &OrderUpdateRequest{}
			updateNeeded := false

			if status, ok, err := OptionalParamOK[string](request, "status"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Status = &status
				updateNeeded = true
			}

			if priority, ok, err := OptionalParamOK[string](request, "priority"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Priority = &priority
				updateNeeded = true
			}

			if dueDate, ok, err := OptionalParamOK[string](request, "due_date"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.DueDate = &dueDate
				updateNeeded = true
			}

			if !updateNeeded {
				return mcp.NewToolResultError("No update parameters provided."), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			updatedOrder, resp, err := client.Orders.Update(ctx, orderID, update)
			if err != nil {
				return nil, fmt.Errorf("failed to update order: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to update order: %s", string(body))), nil
			}

			r, err := json.Marshal(updatedOrder)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
} 