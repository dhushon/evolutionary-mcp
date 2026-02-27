package repository

import "context"

// Memory represents a single memory entry.
type Memory struct {
	ID         string
	Content    string
	Embedding  []float32
	Confidence float64
	Version    int
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
