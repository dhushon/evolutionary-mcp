package services

import (
	"context"
	"testing"

	"evolutionary-mcp/backend/internal/contextutil"
	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMemoryStore satisfies repository.Repository
type MockMemoryStore struct {
	mock.Mock
}

func (m *MockMemoryStore) Save(ctx context.Context, memory *repository.Memory) error {
	args := m.Called(ctx, memory)
	return args.Error(0)
}

func (m *MockMemoryStore) Get(ctx context.Context, id string) (*repository.Memory, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.Memory), args.Error(1)
}

func (m *MockMemoryStore) Search(ctx context.Context, embedding []float32) ([]*repository.Memory, error) {
	args := m.Called(ctx, embedding)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*repository.Memory), args.Error(1)
}

func (m *MockMemoryStore) Update(ctx context.Context, memory *repository.Memory) error {
	args := m.Called(ctx, memory)
	return args.Error(0)
}

func (m *MockMemoryStore) Ping(ctx context.Context) error { return nil }
func (m *MockMemoryStore) CreateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	return nil
}
func (m *MockMemoryStore) UpdateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	return nil
}
func (m *MockMemoryStore) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
	return nil, nil
}
func (m *MockMemoryStore) ListWorkflows(ctx context.Context) ([]*models.Workflow, error) {
	return nil, nil
}
func (m *MockMemoryStore) GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	return nil, nil
}
func (m *MockMemoryStore) GetTenantByID(ctx context.Context, id string) (*models.Tenant, error) {
	return nil, nil
}
func (m *MockMemoryStore) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return nil
}
func (m *MockMemoryStore) CreateGroundingRule(ctx context.Context, rule *models.GroundingRule) error {
	return nil
}
func (m *MockMemoryStore) GetGroundingRule(ctx context.Context, id string) (*models.GroundingRule, error) {
	return nil, nil
}
func (m *MockMemoryStore) ListGroundingRules(ctx context.Context, tenantID string) ([]*models.GroundingRule, error) {
	return nil, nil
}
func (m *MockMemoryStore) UpdateGroundingRule(ctx context.Context, rule *models.GroundingRule) error {
	return nil
}
func (m *MockMemoryStore) DeleteGroundingRule(ctx context.Context, id string) error {
	return nil
}
func (m *MockMemoryStore) SearchGroundingRules(ctx context.Context, tenantID string, embedding []float32) ([]*models.GroundingRule, error) {
	return nil, nil
}

func (m *MockMemoryStore) ListMemories(ctx context.Context, tenantID string) ([]*repository.Memory, error) {
	return nil, nil
}

// MockMLClient satisfies MLClient interface
type MockMLClient struct {
	mock.Mock
}

func (m *MockMLClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	args := m.Called(ctx, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]float32), args.Error(1)
}

func TestMemoryService_Remember(t *testing.T) {
	mockStore := new(MockMemoryStore)
	mockML := new(MockMLClient)
	svc := NewMemoryService(mockStore, mockML)

	tenantID := "test-tenant"
	ctx := contextutil.WithTenant(context.Background(), tenantID)
	content := "test memory"
	fakeEmbedding := []float32{0.1, 0.2, 0.3}

	mockML.On("GetEmbedding", ctx, content).Return(fakeEmbedding, nil)
	mockStore.On("Save", ctx, mock.MatchedBy(func(m *repository.Memory) bool {
		return m.Content == content && 
			m.TenantID == tenantID && 
			assert.ObjectsAreEqual(m.Embedding, fakeEmbedding)
	})).Return(nil)

	memory, err := svc.Remember(ctx, content)

	assert.NoError(t, err)
	assert.NotNil(t, memory)
	assert.Equal(t, content, memory.Content)
	assert.Equal(t, tenantID, memory.TenantID)
	mockML.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}

func TestMemoryService_Recall(t *testing.T) {
	mockStore := new(MockMemoryStore)
	mockML := new(MockMLClient)
	svc := NewMemoryService(mockStore, mockML)

	tenantID := "test-tenant"
	ctx := contextutil.WithTenant(context.Background(), tenantID)
	query := "search query"
	fakeEmbedding := []float32{0.1, 0.2, 0.3}
	expectedResults := []*repository.Memory{
		{ID: "1", Content: "result 1", TenantID: tenantID},
	}

	mockML.On("GetEmbedding", ctx, query).Return(fakeEmbedding, nil)
	mockStore.On("Search", ctx, fakeEmbedding).Return(expectedResults, nil)

	results, err := svc.Recall(ctx, query)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "result 1", results[0].Content)
	mockML.AssertExpectations(t)
	mockStore.AssertExpectations(t)
}
