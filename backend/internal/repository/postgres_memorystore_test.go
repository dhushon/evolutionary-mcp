package repository

import (
	"context"
	"testing"

	"evolutionary-mcp/backend/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// NoOpLogger for testing
type NoOpLogger struct{}

func (l *NoOpLogger) Debug(msg string, args ...any) {}
func (l *NoOpLogger) Info(msg string, args ...any)  {}
func (l *NoOpLogger) Error(msg string, args ...any) {}

// TestLogger logs to testing.T so we can see output on failure
type TestLogger struct {
	t *testing.T
}

func (l *TestLogger) Debug(msg string, args ...any) { l.t.Logf("DEBUG: %s %v", msg, args) }
func (l *TestLogger) Info(msg string, args ...any)  { l.t.Logf("INFO: %s %v", msg, args) }
func (l *TestLogger) Error(msg string, args ...any) { l.t.Logf("ERROR: %s %v", msg, args) }

// TxWrapper wraps pgx.Tx to satisfy the DBTX interface (adding Ping)
type TxWrapper struct {
	pgx.Tx
}

func (t TxWrapper) Ping(ctx context.Context) error {
	return t.Tx.Conn().Ping(ctx)
}

func TestPostgresMemoryStore(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"pgvector/pgvector:pg16",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	// Initialize Schema (matching all migrations)
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS vector;
	
	CREATE TABLE IF NOT EXISTS tenants (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		domain TEXT UNIQUE NOT NULL,
		logo_svg TEXT,
		brand_title TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	
	CREATE TABLE IF NOT EXISTS workflows (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id TEXT NOT NULL DEFAULT 'default',
		workflow_id UUID NOT NULL,
		version INT NOT NULL DEFAULT 1,
		is_latest BOOLEAN NOT NULL DEFAULT TRUE,
		name TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'draft',
		parent_id UUID,
		element_type TEXT NOT NULL DEFAULT 'workflow',
		input_schema JSONB,
		output_schema JSONB,
		created_by TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_workflows_version ON workflows (tenant_id, workflow_id, version);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_workflows_latest_active ON workflows (tenant_id, workflow_id) WHERE is_latest = TRUE;

	CREATE TABLE IF NOT EXISTS memories (
		id UUID PRIMARY KEY,
		tenant_id TEXT NOT NULL DEFAULT 'default',
		content TEXT NOT NULL,
		embedding VECTOR(384),
		confidence FLOAT NOT NULL,
		version INT NOT NULL,
		provenance JSONB DEFAULT '{}',
		workflow_id UUID
	);

	CREATE TABLE IF NOT EXISTS grounding_rules (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		tenant_id UUID NOT NULL REFERENCES tenants(id),
		workflow_id UUID REFERENCES workflows(id),
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		embedding VECTOR(384),
		is_global BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`
	_, err = pool.Exec(ctx, schema)
	if err != nil {
		t.Fatal(err)
	}

	// Helper to create a transactional store for a test
	withTx := func(t *testing.T, fn func(store *PostgresMemoryStore)) {
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatalf("failed to begin transaction: %v", err)
		}
		defer tx.Rollback(ctx)

		store := NewPostgresMemoryStore(TxWrapper{tx}, &TestLogger{t})
		fn(store)
	}

	t.Run("GroundingRules: CRUD", func(t *testing.T) {
		withTx(t, func(store *PostgresMemoryStore) {
			// Setup tenant first
			tenant := &models.Tenant{Name: "Test Tenant", Domain: "test.com"}
			err := store.CreateTenant(ctx, tenant)
			require.NoError(t, err)

			rule := &models.GroundingRule{
				Name:     "Test Rule",
				Content:  "Grounding content",
				TenantID: tenant.ID,
				IsGlobal: false,
			}

			// Create
			err = store.CreateGroundingRule(ctx, rule)
			assert.NoError(t, err)
			assert.NotEmpty(t, rule.ID)

			// Get
			retrieved, err := store.GetGroundingRule(ctx, rule.ID)
			assert.NoError(t, err)
			assert.Equal(t, rule.Name, retrieved.Name)

			// Update
			rule.Name = "Updated Rule"
			err = store.UpdateGroundingRule(ctx, rule)
			assert.NoError(t, err)

			// List
			list, err := store.ListGroundingRules(ctx, tenant.ID)
			assert.NoError(t, err)
			assert.NotEmpty(t, list)
			assert.Equal(t, "Updated Rule", list[0].Name)

			// Delete
			err = store.DeleteGroundingRule(ctx, rule.ID)
			assert.NoError(t, err)
		})
	})

	t.Run("Workflows: Hierarchical support", func(t *testing.T) {
		withTx(t, func(store *PostgresMemoryStore) {
			parent := &models.Workflow{
				WorkflowID:  uuid.New().String(),
				TenantID:    "tenant-1",
				Name:        "Parent",
				ElementType: "workflow",
			}
			err := store.CreateWorkflow(ctx, parent)
			require.NoError(t, err)

			child := &models.Workflow{
				WorkflowID:  uuid.New().String(),
				TenantID:    "tenant-1",
				Name:        "Child",
				ParentID:    &parent.ID,
				ElementType: "element",
			}
			err = store.CreateWorkflow(ctx, child)
			assert.NoError(t, err)

			retrieved, err := store.GetWorkflow(ctx, child.ID)
			assert.NoError(t, err)
			assert.Equal(t, parent.ID, *retrieved.ParentID)
		})
	})
}
