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

type GetClientFn func(context.Context) (*Client, error)

// GetPart creates a tool to get details of a specific part.
func GetPart(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_part",
			mcp.WithDescription(t("TOOL_GET_PART_DESCRIPTION", "Get details of a specific part")),
			mcp.WithString("part_id",
				mcp.Required(),
				mcp.Description("Part ID"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			partID, err := requiredParam[string](request, "part_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
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
				return mcp.NewToolResultError(fmt.Sprintf("failed to get part: %s", string(body))), nil
			}

			r, err := json.Marshal(part)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal part: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// ListParts creates a tool to list and filter parts.
func ListParts(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("list_parts",
			mcp.WithDescription(t("TOOL_LIST_PARTS_DESCRIPTION", "List and filter parts")),
			mcp.WithString("status",
				mcp.Description("Filter by status"),
			),
			mcp.WithString("type",
				mcp.Description("Filter by type"),
			),
			WithPagination(),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			status, err := OptionalParam[string](request, "status")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			partType, err := OptionalParam[string](request, "type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			pagination, err := OptionalPaginationParams(request)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			opts := &ListPartsOptions{
				Status: status,
				Type:   partType,
				ListOptions: ListOptions{
					Page:    pagination.page,
					PerPage: pagination.perPage,
				},
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			parts, resp, err := client.Parts.List(ctx, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to list parts: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to list parts: %s", string(body))), nil
			}

			r, err := json.Marshal(parts)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// CreatePart creates a tool to create a new part.
func CreatePart(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("create_part",
			mcp.WithDescription(t("TOOL_CREATE_PART_DESCRIPTION", "Create a new part")),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Part name"),
			),
			mcp.WithString("description",
				mcp.Description("Part description"),
			),
			mcp.WithString("type",
				mcp.Required(),
				mcp.Description("Part type"),
			),
			mcp.WithString("status",
				mcp.Description("Part status"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := requiredParam[string](request, "name")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			description, err := OptionalParam[string](request, "description")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			partType, err := requiredParam[string](request, "type")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			status, err := OptionalParam[string](request, "status")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			part := &Part{
				Name:        name,
				Description: description,
				Type:       partType,
				Status:     status,
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get First Resonance client: %w", err)
			}
			createdPart, resp, err := client.Parts.Create(ctx, part)
			if err != nil {
				return nil, fmt.Errorf("failed to create part: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusCreated {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to create part: %s", string(body))), nil
			}

			r, err := json.Marshal(createdPart)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

// UpdatePart creates a tool to update an existing part.
func UpdatePart(getClient GetClientFn, t TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("update_part",
			mcp.WithDescription(t("TOOL_UPDATE_PART_DESCRIPTION", "Update an existing part")),
			mcp.WithString("part_id",
				mcp.Required(),
				mcp.Description("Part ID to update"),
			),
			mcp.WithString("name",
				mcp.Description("New name"),
			),
			mcp.WithString("description",
				mcp.Description("New description"),
			),
			mcp.WithString("type",
				mcp.Description("New type"),
			),
			mcp.WithString("status",
				mcp.Description("New status"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			partID, err := requiredParam[string](request, "part_id")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			update := &PartUpdateRequest{}
			updateNeeded := false

			if name, ok, err := OptionalParamOK[string](request, "name"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Name = &name
				updateNeeded = true
			}

			if description, ok, err := OptionalParamOK[string](request, "description"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Description = &description
				updateNeeded = true
			}

			if partType, ok, err := OptionalParamOK[string](request, "type"); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			} else if ok {
				update.Type = &partType
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
			updatedPart, resp, err := client.Parts.Update(ctx, partID, update)
			if err != nil {
				return nil, fmt.Errorf("failed to update part: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to update part: %s", string(body))), nil
			}

			r, err := json.Marshal(updatedPart)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
} 