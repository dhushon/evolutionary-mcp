package services

import (
	"context"
	"evolutionary-mcp/backend/internal/contextutil"
	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"
	"fmt"

	"github.com/google/uuid"
)

// MemoryService is a service for managing memories and grounding rules.
type MemoryService struct {
	store    repository.Repository
	mlClient MLClient
}

// NewMemoryService creates a new MemoryService.
func NewMemoryService(store repository.Repository, mlClient MLClient) *MemoryService {
	return &MemoryService{
		store:    store,
		mlClient: mlClient,
	}
}

// Remember creates a new memory with semantic embedding and tenant isolation.
func (s *MemoryService) Remember(ctx context.Context, content string) (*repository.Memory, error) {
	tenantID := contextutil.GetTenant(ctx)
	if tenantID == "" {
		return nil, fmt.Errorf("unauthorized: tenant_id missing from context")
	}

	embedding, err := s.mlClient.GetEmbedding(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	memory := &repository.Memory{
		ID:         uuid.New().String(),
		TenantID:   tenantID,
		Content:    content,
		Embedding:  embedding,
		Confidence: 1.0,
		Version:    1,
		Provenance: map[string]interface{}{
			"source": "mcp-tool",
		},
	}

	if err := s.store.Save(ctx, memory); err != nil {
		return nil, err
	}

	return memory, nil
}

// Recall retrieves memories similar to the query within the tenant's scope.
func (s *MemoryService) Recall(ctx context.Context, query string) ([]*repository.Memory, error) {
	tenantID := contextutil.GetTenant(ctx)
	if tenantID == "" {
		return nil, fmt.Errorf("unauthorized: tenant_id missing from context")
	}

	embedding, err := s.mlClient.GetEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Repository Search already extracts tenantID from context using GetTenant()
	return s.store.Search(ctx, embedding)
}

// GiveFeedback updates the confidence of a memory based on user/AI feedback.
func (s *MemoryService) GiveFeedback(ctx context.Context, id string, confidence float64) error {
	memory, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}

	tenantID := contextutil.GetTenant(ctx)
	if memory.TenantID != tenantID {
		return fmt.Errorf("unauthorized: memory belongs to another tenant")
	}

	memory.Confidence = confidence
	memory.Version++

	return s.store.Update(ctx, memory)
}

// GetGroundingRules returns the foundational rules for the current tenant.
func (s *MemoryService) GetGroundingRules(ctx context.Context) ([]*models.GroundingRule, error) {
	tenantID := contextutil.GetTenant(ctx)
	if tenantID == "" {
		return nil, fmt.Errorf("unauthorized: tenant_id missing from context")
	}

	return s.store.ListGroundingRules(ctx, tenantID)
}
