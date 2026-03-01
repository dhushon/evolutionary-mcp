package repository

import (
	"context"
	"evolutionary-mcp/backend/pkg/models"
)

// Memory represents a single memory entry.
type Memory struct {
	ID         string
	Content    string
	Embedding  []float32
	Confidence float64
	Version    int
	// Provenance tracks the system state (model ver, rag ver) and user context (session)
	Provenance map[string]interface{}
	WorkflowID string // Links to the specific version of the workflow definition
	TenantID   string // Multi-tenancy isolation
}

// Repository is an interface for all data access operations.
type Repository interface {
	// Save saves a memory to the store.
	Save(ctx context.Context, memory *Memory) error
	// Get retrieves a memory by its ID.
	Get(ctx context.Context, id string) (*Memory, error)
	// Search searches for memories based on a query.
	Search(ctx context.Context, embedding []float32) ([]*Memory, error)
	// Update updates an existing memory.
	Update(ctx context.Context, memory *Memory) error
	// Ping checks the connection to the storage backend.
	Ping(ctx context.Context) error
	// CreateWorkflow creates a new workflow or evolves an existing one (append-only).
	CreateWorkflow(ctx context.Context, workflow *models.Workflow) error
	ListWorkflows(ctx context.Context) ([]*models.Workflow, error)
	// Tenant operations
	GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error)
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
}

// MemoryStore is an interface for storing and retrieving memories.
type MemoryStore interface {
	// Save saves a memory to the store.
	Save(ctx context.Context, memory *Memory) error
	// Get retrieves a memory by its ID.
	Get(ctx context.Context, id string) (*Memory, error)
	// Search searches for memories based on a query.
	Search(ctx context.Context, embedding []float32) ([]*Memory, error)
	// Update updates an existing memory.
	Update(ctx context.Context, memory *Memory) error
}
