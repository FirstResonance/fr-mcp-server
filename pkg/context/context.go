package context

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Context represents a contextual information container
type Context struct {
	ID          string                 `json:"id"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]string      `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
	ParentID    *string                `json:"parent_id,omitempty"`
	ChildrenIDs []string               `json:"children_ids,omitempty"`
	Source      string                 `json:"source"`
}

// ContextManager manages contexts and their relationships
type ContextManager struct {
	contexts map[string]*Context
	mu       sync.RWMutex
}

// NewContextManager creates a new context manager
func NewContextManager() *ContextManager {
	return &ContextManager{
		contexts: make(map[string]*Context),
	}
}

// CreateContext creates a new context
func (cm *ContextManager) CreateContext(id string, data map[string]interface{}, metadata map[string]string, source string) *Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	ctx := &Context{
		ID:        id,
		Data:      data,
		Metadata:  metadata,
		CreatedAt: now,
		UpdatedAt: now,
		Source:    source,
	}

	cm.contexts[id] = ctx
	return ctx
}

// GetContext retrieves a context by ID
func (cm *ContextManager) GetContext(id string) (*Context, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	ctx, exists := cm.contexts[id]
	return ctx, exists
}

// UpdateContext updates an existing context
func (cm *ContextManager) UpdateContext(id string, data map[string]interface{}, metadata map[string]string) (*Context, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ctx, exists := cm.contexts[id]
	if !exists {
		return nil, false
	}

	if data != nil {
		ctx.Data = data
	}

	if metadata != nil {
		ctx.Metadata = metadata
	}

	ctx.UpdatedAt = time.Now()
	return ctx, true
}

// DeleteContext deletes a context
func (cm *ContextManager) DeleteContext(id string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.contexts[id]; !exists {
		return false
	}

	delete(cm.contexts, id)
	return true
}

// GetContextsBySource retrieves all contexts from a specific source
func (cm *ContextManager) GetContextsBySource(source string) []*Context {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var result []*Context
	for _, ctx := range cm.contexts {
		if ctx.Source == source {
			result = append(result, ctx)
		}
	}

	return result
}

// SetParentContext sets a parent-child relationship between contexts
func (cm *ContextManager) SetParentContext(childID, parentID string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	child, childExists := cm.contexts[childID]
	parent, parentExists := cm.contexts[parentID]

	if !childExists || !parentExists {
		return false
	}

	child.ParentID = &parentID
	parent.ChildrenIDs = append(parent.ChildrenIDs, childID)

	return true
}

// GetInheritedContext retrieves a context with inherited data from parent contexts
func (cm *ContextManager) GetInheritedContext(id string) map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make(map[string]interface{})

	// Get the context
	ctx, exists := cm.contexts[id]
	if !exists {
		return result
	}

	// Add current context data
	for k, v := range ctx.Data {
		result[k] = v
	}

	// Recursively add parent context data
	if ctx.ParentID != nil {
		parentData := cm.GetInheritedContext(*ctx.ParentID)
		for k, v := range parentData {
			// Only add parent data if not already in child
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}

	return result
}

// SaveContexts saves all contexts to a file
func (cm *ContextManager) SaveContexts(filePath string) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.MarshalIndent(cm.contexts, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadContexts loads contexts from a file
func (cm *ContextManager) LoadContexts(filePath string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &cm.contexts)
}
