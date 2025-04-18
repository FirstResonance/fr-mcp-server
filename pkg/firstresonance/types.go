package firstresonance

// ListOptions represents pagination parameters
type ListOptions struct {
	Page    int
	PerPage int
}

// Part represents a part in First Resonance
type Part struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Status      string `json:"status,omitempty"`
}

// PartUpdateRequest represents a request to update a part
type PartUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"`
	Status      *string `json:"status,omitempty"`
}

// ListPartsOptions represents options for listing parts
type ListPartsOptions struct {
	Status string
	Type   string
	ListOptions
}

// Order represents an order in First Resonance
type Order struct {
	ID         string        `json:"id"`
	CustomerID string        `json:"customer_id"`
	Items      []interface{} `json:"items"`
	Priority   string        `json:"priority,omitempty"`
	DueDate    string        `json:"due_date,omitempty"`
	Status     string        `json:"status,omitempty"`
}

// OrderUpdateRequest represents a request to update an order
type OrderUpdateRequest struct {
	Status   *string `json:"status,omitempty"`
	Priority *string `json:"priority,omitempty"`
	DueDate  *string `json:"due_date,omitempty"`
}

// ListOrdersOptions represents options for listing orders
type ListOrdersOptions struct {
	Status    string
	Sort      string
	Direction string
	ListOptions
}

// Supplier represents a supplier in First Resonance
type Supplier struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	ContactInfo map[string]interface{} `json:"contact_info,omitempty"`
	Status      string                 `json:"status,omitempty"`
}

// SupplierUpdateRequest represents a request to update a supplier
type SupplierUpdateRequest struct {
	Name        *string                 `json:"name,omitempty"`
	ContactInfo *map[string]interface{} `json:"contact_info,omitempty"`
	Status      *string                 `json:"status,omitempty"`
}

// ListSuppliersOptions represents options for listing suppliers
type ListSuppliersOptions struct {
	Status    string
	Sort      string
	Direction string
	ListOptions
}

// InventoryItem represents an inventory item in First Resonance
type InventoryItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
	Location string `json:"location,omitempty"`
	Status   string `json:"status,omitempty"`
}

// InventoryItemUpdateRequest represents a request to update an inventory item
type InventoryItemUpdateRequest struct {
	Quantity *int    `json:"quantity,omitempty"`
	Location *string `json:"location,omitempty"`
	Status   *string `json:"status,omitempty"`
}

// ListInventoryItemsOptions represents options for listing inventory items
type ListInventoryItemsOptions struct {
	Location  string
	Status    string
	Sort      string
	Direction string
	ListOptions
}

// SearchOptions represents options for search operations
type SearchOptions struct {
	Query string
	Sort  string
	Order string
	ListOptions
}

// ABomItem represents an item in an ABOM
type ABomItem struct {
	ID       string `json:"id"`
	PartID   string `json:"part_id"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit,omitempty"`
	Notes    string `json:"notes,omitempty"`
}

// ABom represents an ABOM in First Resonance
type ABom struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Version     string     `json:"version,omitempty"`
	Status      string     `json:"status,omitempty"`
	CreatedAt   string     `json:"created_at,omitempty"`
	UpdatedAt   string     `json:"updated_at,omitempty"`
	Items       []ABomItem `json:"items,omitempty"`
}

// ABomUpdateRequest represents a request to update an ABOM
type ABomUpdateRequest struct {
	Name        *string     `json:"name,omitempty"`
	Description *string     `json:"description,omitempty"`
	Version     *string     `json:"version,omitempty"`
	Status      *string     `json:"status,omitempty"`
	Items       *[]ABomItem `json:"items,omitempty"`
}

// ListABomsOptions represents options for listing ABOMs
type ListABomsOptions struct {
	Status    string
	Sort      string
	Direction string
	ListOptions
}

// PartsService handles operations on parts
type PartsService struct {
	client *Client
}

// OrdersService handles operations on orders
type OrdersService struct {
	client *Client
}

// SuppliersService handles operations on suppliers
type SuppliersService struct {
	client *Client
}

// InventoryService handles operations on inventory
type InventoryService struct {
	client *Client
}

// SearchService handles search operations
type SearchService struct {
	client *Client
}

// ABomService handles operations on ABOMs
type ABomService struct {
	client *Client
}
