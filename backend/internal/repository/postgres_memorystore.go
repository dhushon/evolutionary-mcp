package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresMemoryStore is a PostgreSQL implementation of the MemoryStore interface.
type PostgresMemoryStore struct {
	db *pgxpool.Pool
}

// NewPostgresMemoryStore creates a new PostgresMemoryStore.
func NewPostgresMemoryStore(db *pgxpool.Pool) *PostgresMemoryStore {
	return &PostgresMemoryStore{db: db}
}

// Save saves a memory to the store.
func (s *PostgresMemoryStore) Save(ctx context.Context, memory *Memory) error {
	_, err := s.db.Exec(ctx, "INSERT INTO memories (id, content, embedding, confidence, version) VALUES ($1, $2, $3, $4, $5)", memory.ID, memory.Content, memory.Embedding, memory.Confidence, memory.Version)
	return err
}

// Get retrieves a memory by its ID.
func (s *PostgresMemoryStore) Get(ctx context.Context, id string) (*Memory, error) {
	var memory Memory
	err := s.db.QueryRow(ctx, "SELECT id, content, embedding, confidence, version FROM memories WHERE id = $1", id).Scan(&memory.ID, &memory.Content, &memory.Embedding, &memory.Confidence, &memory.Version)
	if err != nil {
		return nil, err
	}
	return &memory, nil
}

// Search searches for memories based on a query.
func (s *PostgresMemoryStore) Search(ctx context.Context, embedding []float32) ([]*Memory, error) {
	rows, err := s.db.Query(ctx, "SELECT id, content, embedding, confidence, version FROM memories ORDER BY embedding <=> $1 LIMIT 10", embedding)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*Memory
	for rows.Next() {
		var memory Memory
		err := rows.Scan(&memory.ID, &memory.Content, &memory.Embedding, &memory.Confidence, &memory.Version)
		if err != nil {
			return nil, err
		}
		memories = append(memories, &memory)
	}

	return memories, nil
}

// Update updates an existing memory.
func (s *PostgresMemoryStore) Update(ctx context.Context, memory *Memory) error {
	_, err := s.db.Exec(ctx, "UPDATE memories SET content = $1, embedding = $2, confidence = $3, version = $4 WHERE id = $5", memory.Content, memory.Embedding, memory.Confidence, memory.Version, memory.ID)
	return err
}
