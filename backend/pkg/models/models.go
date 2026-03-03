package models

import (
	"time"
)

// GroundingRule represents a high-priority context pointer or foundational truth 
// that helps stabilize AI reasoning and provides factual grounding.
// These are often stored as vectors to allow for semantic lookup during the RAG process.
type GroundingRule struct {
	ID         string    `json:"id"`
	TenantID   string    `json:"tenant_id"`
	WorkflowID *string   `json:"workflow_id,omitempty"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"-"` // Not exposed in JSON, used for vector search
	IsGlobal   bool      `json:"is_global"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// HealthStatus represents service health
type HealthStatus struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// ProblemDetails represents RFC 7807 Problem Details
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	TraceID  string `json:"trace_id,omitempty"`
}
