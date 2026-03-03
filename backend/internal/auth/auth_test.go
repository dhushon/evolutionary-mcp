package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"evolutionary-mcp/backend/internal/config"
	"evolutionary-mcp/backend/internal/contextutil"
	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"

	"github.com/coreos/go-oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// NoOpLogger for testing
type NoOpLogger struct{}

func (l *NoOpLogger) Debug(msg string, args ...any) {}
func (l *NoOpLogger) Info(msg string, args ...any)  {}
func (l *NoOpLogger) Error(msg string, args ...any) {}

// MockKeySet satisfies oidc.KeySet to bypass signature verification
type MockKeySet struct {
	payload []byte
}

func (m *MockKeySet) VerifySignature(ctx context.Context, jwtToken string) ([]byte, error) {
	parts := strings.Split(jwtToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed jwt")
	}
	return base64.RawURLEncoding.DecodeString(parts[1])
}

// MockRepository satisfies repository.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	args := m.Called(ctx, domain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockRepository) GetTenantByID(ctx context.Context, id string) (*models.Tenant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tenant), args.Error(1)
}

func (m *MockRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	args := m.Called(ctx, tenant)
	return args.Error(0)
}

func (m *MockRepository) Save(ctx context.Context, memory *repository.Memory) error { return nil }
func (m *MockRepository) Get(ctx context.Context, id string) (*repository.Memory, error) {
	return nil, nil
}
func (m *MockRepository) Search(ctx context.Context, embedding []float32) ([]*repository.Memory, error) {
	return nil, nil
}
func (m *MockRepository) Update(ctx context.Context, memory *repository.Memory) error { return nil }
func (m *MockRepository) Ping(ctx context.Context) error                              { return nil }
func (m *MockRepository) CreateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	return nil
}
func (m *MockRepository) UpdateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	return nil
}
func (m *MockRepository) GetWorkflow(ctx context.Context, id string) (*models.Workflow, error) {
	return nil, nil
}
func (m *MockRepository) ListWorkflows(ctx context.Context) ([]*models.Workflow, error) {
	return nil, nil
}
func (m *MockRepository) CreateGroundingRule(ctx context.Context, rule *models.GroundingRule) error {
	return nil
}
func (m *MockRepository) GetGroundingRule(ctx context.Context, id string) (*models.GroundingRule, error) {
	return nil, nil
}
func (m *MockRepository) ListGroundingRules(ctx context.Context, tenantID string) ([]*models.GroundingRule, error) {
	return nil, nil
}
func (m *MockRepository) UpdateGroundingRule(ctx context.Context, rule *models.GroundingRule) error {
	return nil
}
func (m *MockRepository) DeleteGroundingRule(ctx context.Context, id string) error {
	return nil
}
func (m *MockRepository) SearchGroundingRules(ctx context.Context, tenantID string, embedding []float32) ([]*models.GroundingRule, error) {
	return nil, nil
}

func (m *MockRepository) ListMemories(ctx context.Context, tenantID string) ([]*repository.Memory, error) {
	return nil, nil
}

func TestRequireAuth_BearerToken_ExtractsTenant(t *testing.T) {
	mockRepo := new(MockRepository)
	expectedTenant := &models.Tenant{
		ID:     "tenant-123",
		Name:   "acme.com",
		Domain: "acme.com",
	}
	mockRepo.On("GetTenantByDomain", mock.Anything, "acme.com").Return(expectedTenant, nil)

	issuer := "https://test-issuer.com"
	clientID := "test-client"

	claims := map[string]interface{}{
		"iss":   issuer,
		"aud":   clientID,
		"sub":   "test-user",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Add(-1 * time.Minute).Unix(),
		"email": "user@acme.com",
	}
	headerData := map[string]interface{}{"alg": "RS256", "typ": "JWT", "kid": "test-key"}
	headerBytes, _ := json.Marshal(headerData)
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	payload, _ := json.Marshal(claims)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	encodedSignature := base64.RawURLEncoding.EncodeToString([]byte("fakesignature"))
	fakeToken := encodedHeader + "." + encodedPayload + "." + encodedSignature

	keySet := &MockKeySet{payload: payload}
	verifier := oidc.NewVerifier(issuer, keySet, &oidc.Config{ClientID: clientID, SkipClientIDCheck: true})

	a := &Auth{apiVerifier: verifier, repo: mockRepo}
	req := httptest.NewRequest("GET", "/api/v1/workflows", nil)
	req.Header.Set("Authorization", "Bearer "+fakeToken)
	rec := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := contextutil.GetTenant(r.Context())
		assert.Equal(t, "tenant-123", tenantID)
		w.WriteHeader(http.StatusOK)
	})

	a.RequireAuth(nextHandler).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRequireAuth_BypassMode(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("GetTenantByDomain", mock.Anything, "localhost").Return(nil, fmt.Errorf("not found"))
	mockRepo.On("CreateTenant", mock.Anything, mock.MatchedBy(func(tenant *models.Tenant) bool {
		return tenant.Domain == "localhost"
	})).Run(func(args mock.Arguments) {
		argTenant := args.Get(1).(*models.Tenant)
		argTenant.ID = "dev-tenant-id"
	}).Return(nil)

	cfg := &config.Config{
		Environment:   "DEV",
		DevModeBypass: true,
	}
	a, err := New(context.Background(), cfg, mockRepo, &NoOpLogger{})
	assert.NoError(t, err)

	req := httptest.NewRequest("GET", "/api/v1/workflows", nil)
	rec := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := contextutil.GetTenant(r.Context())
		assert.Equal(t, "dev-tenant-id", tenantID)
		w.WriteHeader(http.StatusOK)
	})

	a.RequireAuth(nextHandler).ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
