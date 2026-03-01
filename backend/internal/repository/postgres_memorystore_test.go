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

	// Initialize Schema
	schema := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE EXTENSION IF NOT EXISTS vector;
	
	CREATE TABLE IF NOT EXISTS tenants (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name TEXT NOT NULL,
		domain TEXT UNIQUE NOT NULL,
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
		input_schema JSONB,
		output_schema JSONB,
		created_by TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	CREATE UNIQUE INDEX idx_workflows_version ON workflows (tenant_id, workflow_id, version);
	CREATE UNIQUE INDEX idx_workflows_latest_active ON workflows (tenant_id, workflow_id) WHERE is_latest = TRUE;

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
		defer tx.Rollback(ctx) // Cleanup motion: Always rollback at end of test

		store := NewPostgresMemoryStore(TxWrapper{tx}, &TestLogger{t})
		fn(store)
	}

	t.Run("Memories: Save and Get", func(t *testing.T) {
		withTx(t, func(store *PostgresMemoryStore) {
			id := uuid.New().String()
			memory := &Memory{
				ID:         id,
				TenantID:   "tenant-a",
				Content:    "test content",
				Confidence: 0.9,
				Version:    1,
			}

			err := store.Save(ctx, memory)
			require.NoError(t, err)

			retrieved, err := store.Get(ctx, id)
			require.NoError(t, err)
			assert.Equal(t, memory.ID, retrieved.ID)
			assert.Equal(t, "tenant-a", retrieved.TenantID)
		})
	})

	t.Run("Workflows: Evolution Strategy", func(t *testing.T) {
		withTx(t, func(store *PostgresMemoryStore) {
			// 1. Create Initial Workflow
			wfID := uuid.New().String()
			wf := &models.Workflow{
				ID:          uuid.New().String(),
				TenantID:    "default",
				WorkflowID:  wfID,
				Name:        "Summarizer",
				Description: "v1 description",
				Status:      "active",
			}

			err := store.CreateWorkflow(ctx, wf)
			assert.NoError(t, err)
			assert.Equal(t, 1, wf.Version)
			assert.True(t, wf.IsLatest)

			// 2. Evolve Workflow (Update)
			wf2 := &models.Workflow{
				ID:          uuid.New().String(),
				TenantID:    "default",
				WorkflowID:  wfID,
				Name:        "Summarizer",
				Description: "v2 description",
			}
			err = store.CreateWorkflow(ctx, wf2)
			assert.NoError(t, err)
			assert.Equal(t, 2, wf2.Version)
			assert.True(t, wf2.IsLatest)

			// 3. Verify List only returns latest
			list, err := store.ListWorkflows(ctx)
			assert.NoError(t, err)
			assert.Len(t, list, 1)
			assert.Equal(t, "v2 description", list[0].Description)
			assert.Equal(t, 2, list[0].Version)
		})
	})

	t.Run("Tenants: Create and Lookup", func(t *testing.T) {
		withTx(t, func(store *PostgresMemoryStore) {
			tenant := &models.Tenant{
				Name:   "Acme Corp",
				Domain: "acme.com",
			}
			err := store.CreateTenant(ctx, tenant)
			assert.NoError(t, err)
			assert.NotEmpty(t, tenant.ID)

			found, err := store.GetTenantByDomain(ctx, "acme.com")
			assert.NoError(t, err)
			assert.Equal(t, tenant.ID, found.ID)
		})
	})

	t.Run("Multi-Tenancy: Isolation", func(t *testing.T) {
		withTx(t, func(store *PostgresMemoryStore) {
			// 1. Create Workflow for Tenant A
			wfA := &models.Workflow{
				ID:         uuid.New().String(),
				TenantID:   "tenant-a",
				WorkflowID: uuid.New().String(),
				Name:       "Workflow A",
				Version:    1,
				IsLatest:   true,
			}
			err := store.CreateWorkflow(ctx, wfA)
			require.NoError(t, err)

			// 2. Create Workflow for Tenant B
			wfB := &models.Workflow{
				ID:         uuid.New().String(),
				TenantID:   "tenant-b",
				WorkflowID: uuid.New().String(),
				Name:       "Workflow B",
				Version:    1,
				IsLatest:   true,
			}
			err = store.CreateWorkflow(ctx, wfB)
			require.NoError(t, err)

			// 3. List as Tenant A (Should only see A)
			ctxA := context.WithValue(ctx, "tenant_id", "tenant-a")
			listA, err := store.ListWorkflows(ctxA)
			require.NoError(t, err)
			assert.Len(t, listA, 1)
			assert.Equal(t, "Workflow A", listA[0].Name)

			// 4. List as Tenant B (Should only see B)
			ctxB := context.WithValue(ctx, "tenant_id", "tenant-b")
			listB, err := store.ListWorkflows(ctxB)
			require.NoError(t, err)
			assert.Len(t, listB, 1)
			assert.Equal(t, "Workflow B", listB[0].Name)

			// 5. List as Unknown Tenant (Should see nothing or default)
			ctxC := context.WithValue(ctx, "tenant_id", "tenant-c")
			listC, err := store.ListWorkflows(ctxC)
			require.NoError(t, err)
			assert.Len(t, listC, 0)
		})
	})
}
