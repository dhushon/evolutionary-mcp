package models

import (
	"time"
)

// Workflow represents the evolutionary definition of a logic pipeline.
type Workflow struct {
	ID           string                 `json:"id"`          // Unique Version ID
	TenantID     string                 `json:"tenant_id"`   // Multi-tenancy isolation
	WorkflowID   string                 `json:"workflow_id"` // Stable Concept ID
	Version      int                    `json:"version"`
	IsLatest     bool                   `json:"is_latest"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Status       string                 `json:"status"`
	InputSchema  map[string]interface{} `json:"input_schema"`
	OutputSchema map[string]interface{} `json:"output_schema"`
	CreatedBy    string                 `json:"created_by"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}