// Package models defines the domain models for the payer service
package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
)

// PayerType represents the type of payer organization
type PayerType string

const (
	PayerTypeInsurance   PayerType = "insurance"
	PayerTypeMedicare    PayerType = "medicare"
	PayerTypeMedicaid    PayerType = "medicaid"
	PayerTypeCommercial  PayerType = "commercial"
	PayerTypeGovernment  PayerType = "government"
	PayerTypeSelfPay     PayerType = "self_pay"
	PayerTypeOther       PayerType = "other"
)

// PayerStatus represents the operational status of a payer
type PayerStatus string

const (
	PayerStatusActive    PayerStatus = "active"
	PayerStatusInactive  PayerStatus = "inactive"
	PayerStatusPending   PayerStatus = "pending"
	PayerStatusSuspended PayerStatus = "suspended"
)

// SearchType represents the search strategy to use
type SearchType string

const (
	SearchTypeFuzzy    SearchType = "fuzzy"
	SearchTypeVector   SearchType = "vector"
	SearchTypeHybrid   SearchType = "hybrid"
	SearchTypeFulltext SearchType = "fulltext"
)

// MatchType represents how a search result matched
type MatchType string

const (
	MatchTypeExact    MatchType = "exact"
	MatchTypeFuzzy    MatchType = "fuzzy"
	MatchTypeVector   MatchType = "vector"
	MatchTypeFulltext MatchType = "fulltext"
)

// Payer represents a healthcare payer/insurance provider
type Payer struct {
	ID          string         `json:"id" db:"id"`
	Name        string         `json:"name" db:"name"`
	DisplayName *string        `json:"display_name,omitempty" db:"display_name"`
	PayerID     string         `json:"payer_id" db:"payer_id"`
	PayerType   PayerType      `json:"payer_type" db:"payer_type"`
	Status      PayerStatus    `json:"status" db:"status"`
	
	// Contact information
	Contact *ContactInfo `json:"contact,omitempty"`
	
	// Address information
	Address *Address `json:"address,omitempty"`
	
	// Geolocation
	GeoLocation *GeoLocation `json:"geolocation,omitempty"`
	
	// Metadata
	Description *string  `json:"description,omitempty" db:"description"`
	Notes       *string  `json:"notes,omitempty" db:"notes"`
	Tags        []string `json:"tags,omitempty" db:"tags"`
	
	// Vector embedding (not exposed in JSON)
	Embedding pgvector.Vector `json:"-" db:"embedding"`
	
	// Audit fields
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	CreatedBy *string    `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy *string    `json:"updated_by,omitempty" db:"updated_by"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	DeletedBy *string    `json:"-" db:"deleted_by"`
}

// ContactInfo represents contact information
type ContactInfo struct {
	Website *string `json:"website,omitempty" db:"website"`
	Phone   *string `json:"phone,omitempty" db:"phone"`
	Email   *string `json:"email,omitempty" db:"email"`
	Fax     *string `json:"fax,omitempty" db:"fax"`
}

// Address represents a physical address
type Address struct {
	Line1   *string `json:"line1,omitempty" db:"address_line1"`
	Line2   *string `json:"line2,omitempty" db:"address_line2"`
	City    *string `json:"city,omitempty" db:"city"`
	State   *string `json:"state,omitempty" db:"state"`
	ZipCode *string `json:"zip_code,omitempty" db:"zip_code"`
	Country *string `json:"country,omitempty" db:"country"`
}

// GeoLocation represents geographic coordinates
type GeoLocation struct {
	Latitude  *float64 `json:"latitude,omitempty" db:"latitude"`
	Longitude *float64 `json:"longitude,omitempty" db:"longitude"`
}

// PayerSearchResult extends Payer with search-specific metadata
type PayerSearchResult struct {
	*Payer
	SimilarityScore *float64   `json:"similarity_score,omitempty"`
	MatchType       *MatchType `json:"match_type,omitempty"`
}

// SearchOptions contains parameters for searching payers
type SearchOptions struct {
	Query               string
	State               *string
	City                *string
	PayerType           *PayerType
	Status              *PayerStatus
	SearchType          SearchType
	SimilarityThreshold float64
	Limit               int
	Offset              int
}

// SearchResponse contains search results and metadata
type SearchResponse struct {
	Results   []*PayerSearchResult `json:"results"`
	Total     int                  `json:"total"`
	Limit     int                  `json:"limit"`
	Offset    int                  `json:"offset"`
	QueryInfo *QueryInfo           `json:"query_info,omitempty"`
}

// QueryInfo contains metadata about the executed query
type QueryInfo struct {
	Query           string     `json:"query"`
	SearchType      SearchType `json:"search_type"`
	ExecutionTimeMs int64      `json:"execution_time_ms"`
	CacheHit        bool       `json:"cache_hit"`
}

// Vendor represents a healthcare vendor
type Vendor struct {
	ID        string          `json:"id" db:"id"`
	Name      string          `json:"name" db:"name"`
	VendorID  string          `json:"vendor_id" db:"vendor_id"`
	Contact   *ContactInfo    `json:"contact,omitempty"`
	Embedding pgvector.Vector `json:"-" db:"embedding"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time      `json:"-" db:"deleted_at"`
}

// VendorSearchResult extends Vendor with search metadata
type VendorSearchResult struct {
	*Vendor
	SimilarityScore *float64 `json:"similarity_score,omitempty"`
}

// PayerVendorRelationship represents a relationship between a payer and vendor
type PayerVendorRelationship struct {
	ID                string     `json:"id" db:"id"`
	PayerID           string     `json:"payer_id" db:"payer_id"`
	VendorID          string     `json:"vendor_id" db:"vendor_id"`
	RelationshipType  *string    `json:"relationship_type,omitempty" db:"relationship_type"`
	ContractStartDate *time.Time `json:"contract_start_date,omitempty" db:"contract_start_date"`
	ContractEndDate   *time.Time `json:"contract_end_date,omitempty" db:"contract_end_date"`
	Status            string     `json:"status" db:"status"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// SearchHistory tracks search queries for analytics
type SearchHistory struct {
	ID              string     `json:"id" db:"id"`
	Query           string     `json:"query" db:"query"`
	SearchType      *string    `json:"search_type,omitempty" db:"search_type"`
	Filters         []byte     `json:"filters,omitempty" db:"filters"` // JSONB
	ResultCount     *int       `json:"result_count,omitempty" db:"result_count"`
	ExecutionTimeMs *int       `json:"execution_time_ms,omitempty" db:"execution_time_ms"`
	CacheHit        bool       `json:"cache_hit" db:"cache_hit"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UserID          *string    `json:"user_id,omitempty" db:"user_id"`
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

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID          string    `json:"id" db:"id"`
	WorkflowID  string    `json:"workflow_id" db:"workflow_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	Action      string    `json:"action" db:"action"`
	Config      []byte    `json:"config,omitempty" db:"config"` // JSONB
	Order       int       `json:"order" db:"order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// WorkflowExecution represents an instance of a running workflow
type WorkflowExecution struct {
	ID         string    `json:"id" db:"id"`
	WorkflowID string    `json:"workflow_id" db:"workflow_id"`
	Status     string    `json:"status" db:"status"`
	Input      []byte    `json:"input,omitempty" db:"input"`     // JSONB
	Output     []byte    `json:"output,omitempty" db:"output"`   // JSONB
	StartedAt  time.Time `json:"started_at" db:"started_at"`
	EndedAt    *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	CreatedBy  *string   `json:"created_by,omitempty" db:"created_by"`
}
