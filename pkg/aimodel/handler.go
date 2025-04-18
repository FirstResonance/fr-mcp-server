package aimodel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	ctx "github.com/firstresonance/fr-mcp-server/pkg/context"
	"github.com/firstresonance/fr-mcp-server/pkg/firstresonance"
)

// ModelRequest represents a request from an AI model
type ModelRequest struct {
	ModelID   string                 `json:"model_id"`
	Action    string                 `json:"action"`
	ContextID string                 `json:"context_id"`
	Params    map[string]interface{} `json:"params"`
}

// ModelResponse represents a response to an AI model
type ModelResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// ModelHandler handles requests from AI models
type ModelHandler struct {
	contextManager *ctx.ContextManager
	frClient       *firstresonance.Client
	models         map[string]bool
	mu             sync.RWMutex
}

// NewModelHandler creates a new model handler
func NewModelHandler(contextManager *ctx.ContextManager, frClient *firstresonance.Client) *ModelHandler {
	return &ModelHandler{
		contextManager: contextManager,
		frClient:       frClient,
		models:         make(map[string]bool),
	}
}

// RegisterModel registers a new AI model
func (h *ModelHandler) RegisterModel(modelID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.models[modelID] = true
}

// HandleRequest handles a request from an AI model
func (h *ModelHandler) HandleRequest(ctx context.Context, req *ModelRequest) (*ModelResponse, error) {
	// Verify model is registered
	h.mu.RLock()
	if !h.models[req.ModelID] {
		h.mu.RUnlock()
		return nil, fmt.Errorf("model %s not registered", req.ModelID)
	}
	h.mu.RUnlock()

	// Get context if provided
	var ctxData map[string]interface{}
	if req.ContextID != "" {
		ctxData = h.contextManager.GetInheritedContext(req.ContextID)
	}

	// Handle different actions
	switch req.Action {
	case "get_part":
		return h.handleGetPart(ctx, req, ctxData)
	case "create_order":
		return h.handleCreateOrder(ctx, req, ctxData)
	default:
		return &ModelResponse{
			Success: false,
			Error:   fmt.Sprintf("unknown action: %s", req.Action),
		}, nil
	}
}

// handleGetPart handles the get_part action
func (h *ModelHandler) handleGetPart(ctx context.Context, req *ModelRequest, ctxData map[string]interface{}) (*ModelResponse, error) {
	partID, ok := req.Params["part_id"].(string)
	if !ok {
		return &ModelResponse{
			Success: false,
			Error:   "part_id is required",
		}, nil
	}

	part, _, err := h.frClient.Parts.Get(ctx, partID)
	if err != nil {
		return &ModelResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get part: %v", err),
		}, nil
	}

	return &ModelResponse{
		Success: true,
		Data: map[string]interface{}{
			"part": part,
		},
	}, nil
}

// handleCreateOrder handles the create_order action
func (h *ModelHandler) handleCreateOrder(ctx context.Context, req *ModelRequest, ctxData map[string]interface{}) (*ModelResponse, error) {
	// Extract order data from params
	orderData, ok := req.Params["order"].(map[string]interface{})
	if !ok {
		return &ModelResponse{
			Success: false,
			Error:   "order data is required",
		}, nil
	}

	// Convert map to Order struct
	order := &firstresonance.Order{}

	// Required fields
	if customerID, ok := orderData["customer_id"].(string); ok {
		order.CustomerID = customerID
	} else {
		return &ModelResponse{
			Success: false,
			Error:   "customer_id is required and must be a string",
		}, nil
	}

	if items, ok := orderData["items"].([]interface{}); ok {
		order.Items = items
	} else {
		return &ModelResponse{
			Success: false,
			Error:   "items is required and must be an array",
		}, nil
	}

	// Optional fields
	if priority, ok := orderData["priority"].(string); ok {
		order.Priority = priority
	}

	if dueDate, ok := orderData["due_date"].(string); ok {
		order.DueDate = dueDate
	}

	if status, ok := orderData["status"].(string); ok {
		order.Status = status
	}

	// Create the order
	createdOrder, _, err := h.frClient.Orders.Create(ctx, order)
	if err != nil {
		return &ModelResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create order: %v", err),
		}, nil
	}

	return &ModelResponse{
		Success: true,
		Data: map[string]interface{}{
			"order": createdOrder,
		},
	}, nil
}

// ServeHTTP implements the http.Handler interface
func (h *ModelHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ModelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.HandleRequest(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
