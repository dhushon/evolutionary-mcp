package main

import (
	"context"
	"fmt"
	"log"

	"evolutionary-mcp/backend/internal/config"
	"evolutionary-mcp/backend/internal/logging"
	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	logger := logging.NewLogger()

	// Load config
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to DB
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode,
	)
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	store := repository.NewPostgresMemoryStore(pool, logger)

	// 1. Ensure Tenant Exists
	domain := "localhost"
	tenant, err := store.GetTenantByDomain(ctx, domain)
	if err != nil {
		logger.Info("Creating default tenant", "domain", domain)
		tenant = &models.Tenant{
			Name:   "Local Dev Tenant",
			Domain: domain,
		}
		if err := store.CreateTenant(ctx, tenant); err != nil {
			log.Fatalf("Failed to create tenant: %v", err)
		}
	} else {
		logger.Info("Found existing tenant", "id", tenant.ID)
	}

	// Inject tenant_id into context for subsequent operations
	ctx = context.WithValue(ctx, "tenant_id", tenant.ID)

	// 2. Check for existing workflows to prevent duplicates
	existingWorkflows, err := store.ListWorkflows(ctx)
	if err != nil {
		log.Fatalf("Failed to list existing workflows: %v", err)
	}

	existingMap := make(map[string]bool)
	for _, w := range existingWorkflows {
		existingMap[w.Name] = true
	}

	// 3. Create Seed Workflows
	workflows := []struct {
		Name        string
		Description string
		Status      string
	}{
		{"Summarizer", "Summarizes long conversations into concise notes.", "active"},
		{"Fact Checker", "Verifies claims against stored long-term memories.", "active"},
		{"Code Reviewer", "Analyzes code snippets for style and bugs.", "draft"},
	}

	for _, w := range workflows {
		if existingMap[w.Name] {
			logger.Info("Skipping existing workflow", "name", w.Name)
			continue
		}

		wf := &models.Workflow{
			TenantID:    tenant.ID,
			WorkflowID:  uuid.New().String(),
			Name:        w.Name,
			Description: w.Description,
			Status:      w.Status,
			Version:     1,
			IsLatest:    true,
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{"type": "string"},
				},
			},
			CreatedBy: "seed-script",
		}

		if err := store.CreateWorkflow(ctx, wf); err != nil {
			log.Printf("Failed to create workflow %s: %v", w.Name, err)
		} else {
			logger.Info("Seeded workflow", "name", w.Name, "id", wf.WorkflowID)
		}
	}
	logger.Info("Seeding complete!")
}
