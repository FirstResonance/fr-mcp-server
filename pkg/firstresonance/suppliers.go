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

// GetSupplier creates a tool to get details of a specific supplier.
func GetSupplier(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_supplier",
			mcp.WithDescription(t("TOOL_GET_SUPPLIER_DESCRIPTION", "Get details of a specific supplier")),
			mcp.WithString("supplier_id",
				mcp.Required(),
				mcp.Description("Supplier ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			supplierID, err := requiredParam[string](request, "supplier_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
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
				return mcp.NewToolResultError(fmt.Sprintf("failed to get supplier: %s", string(body))), nil
			}

			r, err := json.Marshal(supplier)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal supplier: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// ListSuppliers creates a tool to list and filter suppliers.
func ListSuppliers(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_suppliers",
			mcp.WithDescription(t("TOOL_LIST_SUPPLIERS_DESCRIPTION", "List and filter suppliers")),
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

			opts := &ListSuppliersOptions{
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
			suppliers, resp, err := client.Suppliers.List(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list suppliers: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list suppliers: %s", string(body))), nil
			}

			r, err := json.Marshal(suppliers)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// CreateSupplier creates a tool to create a new supplier.
func CreateSupplier(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_supplier",
			mcp.WithDescription(t("TOOL_CREATE_SUPPLIER_DESCRIPTION", "Create a new supplier")),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Supplier name"),
			),
			mcp.WithObject("contact_info",
				mcp.Description("Contact information"),
			),
			mcp.WithString("status",
				mcp.Description("Supplier status"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := requiredParam[string](request, "name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			contactInfo, err := OptionalParam[map[string]interface{}](request, "contact_info")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			status, err := OptionalParam[string](request, "status")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			supplier := &Supplier{
				Name:        name,
				ContactInfo: contactInfo,
				Status:     status,
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			createdSupplier, resp, err := client.Suppliers.Create(ctx, supplier)
			if err != nil {
				return nil, fmt.Errorf("failed to create supplier: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusCreated {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to create supplier: %s", string(body))), nil
			}

			r, err := json.Marshal(createdSupplier)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// UpdateSupplier creates a tool to update an existing supplier.
func UpdateSupplier(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("update_supplier",
			mcp.WithDescription(t("TOOL_UPDATE_SUPPLIER_DESCRIPTION", "Update an existing supplier")),
			mcp.WithString("supplier_id",
				mcp.Required(),
				mcp.Description("Supplier ID to update"),
			),
			mcp.WithString("name",
				mcp.Description("New name"),
			),
			mcp.WithObject("contact_info",
				mcp.Description("New contact information"),
			),
			mcp.WithString("status",
				mcp.Description("New status"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			supplierID, err := requiredParam[string](request, "supplier_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			update := &SupplierUpdateRequest{}
			updateNeeded := false

			if name, ok, err := OptionalParamOK[string](request, "name"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Name = &name
				updateNeeded = true
			}

			if contactInfo, ok, err := OptionalParamOK[map[string]interface{}](request, "contact_info"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.ContactInfo = &contactInfo
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
			updatedSupplier, resp, err := client.Suppliers.Update(ctx, supplierID, update)
			if err != nil {
				return nil, fmt.Errorf("failed to update supplier: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to update supplier: %s", string(body))), nil
			}

			r, err := json.Marshal(updatedSupplier)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
} 