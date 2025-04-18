package firstresonance

import (
	"github.com/mark3labs/mcp-go/server"
)

type TranslationHelperFunc func(key string, defaultValue string) string

// NewServer creates a new First Resonance MCP server with the specified client and logger.
func NewServer(getClient GetClientFn, version string, readOnly bool, t TranslationHelperFunc) *server.MCPServer {
	// Create a new MCP server
	s := server.NewMCPServer(
		"firstresonance-mcp-server",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithLogging())

	// Add First Resonance Resources
	s.AddResourceTemplate(GetPartContent(getClient, t))
	s.AddResourceTemplate(GetOrderContent(getClient, t))
	s.AddResourceTemplate(GetSupplierContent(getClient, t))
	s.AddResourceTemplate(GetInventoryItemContent(getClient, t))

	// Add First Resonance tools - Parts
	s.AddTool(GetPart(getClient, t))
	s.AddTool(ListParts(getClient, t))
	if !readOnly {
		s.AddTool(CreatePart(getClient, t))
		s.AddTool(UpdatePart(getClient, t))
	}

	// Add First Resonance tools - Orders
	s.AddTool(GetOrder(getClient, t))
	s.AddTool(ListOrders(getClient, t))
	if !readOnly {
		s.AddTool(CreateOrder(getClient, t))
		s.AddTool(UpdateOrder(getClient, t))
	}

	// Add First Resonance tools - Suppliers
	s.AddTool(GetSupplier(getClient, t))
	s.AddTool(ListSuppliers(getClient, t))
	if !readOnly {
		s.AddTool(CreateSupplier(getClient, t))
		s.AddTool(UpdateSupplier(getClient, t))
	}

	// Add First Resonance tools - Inventory
	s.AddTool(GetInventoryItem(getClient, t))
	s.AddTool(ListInventoryItems(getClient, t))
	if !readOnly {
		s.AddTool(UpdateInventoryItem(getClient, t))
	}

	// Add First Resonance tools - Search
	s.AddTool(SearchParts(getClient, t))
	s.AddTool(SearchOrders(getClient, t))

	return s
}
