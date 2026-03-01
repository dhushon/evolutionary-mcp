package repository

import (
	"context"
	"evolutionary-mcp/backend/pkg/models"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Logger defines the logging interface compatible with the application logger.
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// DBTX is an interface that abstracts pgxpool.Pool and pgx.Tx to allow dependency injection for testing.
type DBTX interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
	Ping(context.Context) error
}

// PostgresMemoryStore is a PostgreSQL implementation of the MemoryStore interface.
type PostgresMemoryStore struct {
	db     DBTX
	logger Logger

	// Metrics
	memoriesStored   metric.Int64Counter
	memoriesUpdated  metric.Int64Counter
	memoriesSearched metric.Int64Counter
	workflowsCreated metric.Int64Counter
}

// NewPostgresMemoryStore creates a new PostgresMemoryStore.
func NewPostgresMemoryStore(db DBTX, logger Logger) *PostgresMemoryStore {
	meter := otel.Meter("evolutionary-mcp/backend/repository")

	memStored, err := meter.Int64Counter("memories_stored_total", metric.WithDescription("Total number of new memories stored"))
	if err != nil {
		logger.Error("failed to create memories_stored_total metric", "error", err)
	}
	memUpdated, err := meter.Int64Counter("memories_updated_total", metric.WithDescription("Total number of memories updated/versioned"))
	if err != nil {
		logger.Error("failed to create memories_updated_total metric", "error", err)
	}
	memSearched, err := meter.Int64Counter("memory_searches_total", metric.WithDescription("Total number of semantic searches performed"))
	if err != nil {
		logger.Error("failed to create memory_searches_total metric", "error", err)
	}
	wfCreated, err := meter.Int64Counter("workflows_created_total", metric.WithDescription("Total number of workflows created or evolved"))
	if err != nil {
		logger.Error("failed to create workflows_created_total metric", "error", err)
	}

	return &PostgresMemoryStore{
		db:               db,
		logger:           logger,
		memoriesStored:   memStored,
		memoriesUpdated:  memUpdated,
		memoriesSearched: memSearched,
		workflowsCreated: wfCreated,
	}
}

// Save saves a memory to the store.
func (s *PostgresMemoryStore) Save(ctx context.Context, memory *Memory) error {
	s.logger.Debug("Saving memory", "id", memory.ID, "version", memory.Version, "workflow_id", memory.WorkflowID)
	var workflowID interface{} = memory.WorkflowID
	if memory.WorkflowID == "" {
		workflowID = nil
	}

	_, err := s.db.Exec(ctx, "INSERT INTO memories (id, tenant_id, content, embedding, confidence, version, provenance, workflow_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", memory.ID, memory.TenantID, memory.Content, memory.Embedding, memory.Confidence, memory.Version, memory.Provenance, workflowID)
	if err == nil && s.memoriesStored != nil {
		s.memoriesStored.Add(ctx, 1, metric.WithAttributes(attribute.String("workflow_id", memory.WorkflowID)))
	}
	return err
}

// Get retrieves a memory by its ID.
func (s *PostgresMemoryStore) Get(ctx context.Context, id string) (*Memory, error) {
	s.logger.Debug("Getting memory", "id", id)
	var memory Memory
	var workflowID *string
	err := s.db.QueryRow(ctx, "SELECT id, tenant_id, content, embedding, confidence, version, provenance, workflow_id FROM memories WHERE id = $1", id).Scan(&memory.ID, &memory.TenantID, &memory.Content, &memory.Embedding, &memory.Confidence, &memory.Version, &memory.Provenance, &workflowID)
	if err != nil {
		return nil, err
	}
	if workflowID != nil {
		memory.WorkflowID = *workflowID
	}
	return &memory, nil
}

// Search searches for memories based on a query.
func (s *PostgresMemoryStore) Search(ctx context.Context, embedding []float32) ([]*Memory, error) {
	s.logger.Debug("Searching memories", "embedding_dim", len(embedding))

	tenantID, _ := ctx.Value("tenant_id").(string)
	if tenantID == "" {
		tenantID = "default"
	}

	rows, err := s.db.Query(ctx, "SELECT id, tenant_id, content, embedding, confidence, version, provenance, workflow_id FROM memories WHERE tenant_id = $1 ORDER BY embedding <=> $2 LIMIT 10", tenantID, embedding)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []*Memory
	for rows.Next() {
		var memory Memory
		var workflowID *string
		err := rows.Scan(&memory.ID, &memory.TenantID, &memory.Content, &memory.Embedding, &memory.Confidence, &memory.Version, &memory.Provenance, &workflowID)
		if err != nil {
			return nil, err
		}
		if workflowID != nil {
			memory.WorkflowID = *workflowID
		}
		memories = append(memories, &memory)
	}

	if s.memoriesSearched != nil {
		s.memoriesSearched.Add(ctx, 1)
	}
	s.logger.Debug("Search completed", "results", len(memories))
	return memories, nil
}

// Update updates an existing memory.
func (s *PostgresMemoryStore) Update(ctx context.Context, memory *Memory) error {
	s.logger.Debug("Updating memory", "id", memory.ID, "new_version", memory.Version)
	var workflowID interface{} = memory.WorkflowID
	if memory.WorkflowID == "" {
		workflowID = nil
	}

	_, err := s.db.Exec(ctx, "UPDATE memories SET content = $1, embedding = $2, confidence = $3, version = $4, provenance = $5, workflow_id = $6 WHERE id = $7", memory.Content, memory.Embedding, memory.Confidence, memory.Version, memory.Provenance, workflowID, memory.ID)
	if err == nil && s.memoriesUpdated != nil {
		s.memoriesUpdated.Add(ctx, 1)
	}
	return err
}

// Ping checks the database connection.
func (s *PostgresMemoryStore) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

// ListWorkflows retrieves all workflows from the database.
func (s *PostgresMemoryStore) ListWorkflows(ctx context.Context) ([]*models.Workflow, error) {
	tenantID, _ := ctx.Value("tenant_id").(string)
	if tenantID == "" {
		tenantID = "default"
	}

	s.logger.Debug("Listing active workflows")
	rows, err := s.db.Query(ctx, "SELECT id, workflow_id, tenant_id, version, is_latest, name, description, status, input_schema, output_schema, created_by, created_at, updated_at FROM workflows WHERE is_latest = true AND tenant_id = $1", tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	workflows := make([]*models.Workflow, 0)
	for rows.Next() {
		var workflow models.Workflow
		err := rows.Scan(&workflow.ID, &workflow.WorkflowID, &workflow.TenantID, &workflow.Version, &workflow.IsLatest, &workflow.Name, &workflow.Description, &workflow.Status, &workflow.InputSchema, &workflow.OutputSchema, &workflow.CreatedBy, &workflow.CreatedAt, &workflow.UpdatedAt)
		if err != nil {
			return nil, err
		}
		workflows = append(workflows, &workflow)
	}

	return workflows, nil
}

// CreateWorkflow creates a new workflow or evolves an existing one.
// It manages the is_latest flag and version incrementing transactionally.
func (s *PostgresMemoryStore) CreateWorkflow(ctx context.Context, workflow *models.Workflow) error {
	s.logger.Debug("Creating/Evolving workflow", "workflow_id", workflow.WorkflowID, "name", workflow.Name)
	if workflow.ID == "" {
		workflow.ID = uuid.New().String()
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// If WorkflowID is provided, we are evolving an existing workflow for this tenant.
	// We need to retire the current latest version.
	var nextVersion = 1
	if workflow.WorkflowID != "" {
		// 1. Retire the current latest version
		_, err := tx.Exec(ctx, "UPDATE workflows SET is_latest = false WHERE workflow_id = $1 AND tenant_id = $2 AND is_latest = true", workflow.WorkflowID, workflow.TenantID)
		if err != nil {
			return fmt.Errorf("failed to retire old workflow version: %w", err)
		}

		// 2. Determine the next version number
		var maxVer *int
		err = tx.QueryRow(ctx, "SELECT MAX(version) FROM workflows WHERE workflow_id = $1 AND tenant_id = $2", workflow.WorkflowID, workflow.TenantID).Scan(&maxVer)
		if err != nil {
			s.logger.Error("Failed to determine next version", "error", err)
			return fmt.Errorf("failed to determine next version: %w", err)
		}
		if maxVer != nil {
			nextVersion = *maxVer + 1
		}
	}

	s.logger.Debug("Setting new version", "version", nextVersion)
	workflow.Version = nextVersion
	workflow.IsLatest = true

	// 2. Insert the new version
	_, err = tx.Exec(ctx, `
		INSERT INTO workflows (id, tenant_id, workflow_id, version, is_latest, name, description, status, input_schema, output_schema, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
	`, workflow.ID, workflow.TenantID, workflow.WorkflowID, workflow.Version, workflow.IsLatest, workflow.Name, workflow.Description, workflow.Status, workflow.InputSchema, workflow.OutputSchema, workflow.CreatedBy)
	if err != nil {
		return fmt.Errorf("failed to insert workflow: %w", err)
	}

	if s.workflowsCreated != nil {
		s.workflowsCreated.Add(ctx, 1, metric.WithAttributes(attribute.Bool("is_evolution", workflow.Version > 1)))
	}
	return tx.Commit(ctx)
}

// GetTenantByDomain retrieves a tenant by their email domain.
func (s *PostgresMemoryStore) GetTenantByDomain(ctx context.Context, domain string) (*models.Tenant, error) {
	var t models.Tenant
	err := s.db.QueryRow(ctx, "SELECT id, name, domain, created_at, updated_at FROM tenants WHERE domain = $1", domain).Scan(&t.ID, &t.Name, &t.Domain, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// CreateTenant creates a new tenant.
func (s *PostgresMemoryStore) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	return s.db.QueryRow(ctx, `
		INSERT INTO tenants (name, domain, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at`, tenant.Name, tenant.Domain).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}
