package services

import (
	"context"
	"evolutionary-mcp/backend/internal/repository"
	"github.com/google/uuid"
)

// MemoryService is a service for managing memories.
type MemoryService struct {
	store    repository.MemoryStore
	mlClient MLClient
}

// NewMemoryService creates a new MemoryService.
func NewMemoryService(store repository.MemoryStore, mlClient MLClient) *MemoryService {
	return &MemoryService{
		store:    store,
		mlClient: mlClient,
	}
}

// Remember creates a new memory.
func (s *MemoryService) Remember(ctx context.Context, content string) (*repository.Memory, error) {
	embedding, err := s.mlClient.GetEmbedding(ctx, content)
	if err != nil {
		return nil, err
	}

	memory := &repository.Memory{
		ID:         uuid.New().String(),
		Content:    content,
		Embedding:  embedding,
		Confidence: 0.5, // Initial confidence
		Version:    1,
	}

	err = s.store.Save(ctx, memory)
	if err != nil {
		return nil, err
	}

	return memory, nil
}

// Recall searches for memories.
func (s *MemoryService) Recall(ctx context.Context, query string) ([]*repository.Memory, error) {
	embedding, err := s.mlClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, err
	}

	return s.store.Search(ctx, embedding)
}

// GiveFeedback updates a memory's confidence.
func (s *MemoryService) GiveFeedback(ctx context.Context, id string, confidence float64) error {
	memory, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}

	memory.Confidence = confidence
	memory.Version++

	return s.store.Update(ctx, memory)
}
