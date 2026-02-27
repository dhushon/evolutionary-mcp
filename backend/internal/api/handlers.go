package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// Handler contains HTTP handlers for the payer service REST API
type Handler struct {
}

// NewHandler creates a new Handler with required dependencies
func NewHandler() *Handler {
	return &Handler{}
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
	Version   string    `json:"version"`
}

// HandleHealth returns basic health status (always returns 200 OK)
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "evolutionary-mcp",
		Version:   "1.0.0",
	}
	writeJSON(w, http.StatusOK, status)
}

// writeJSON writes a JSON response with the given status code
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but can't change response at this point
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ProblemDetails represents an RFC 7807 Problem Details response
type ProblemDetails struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

// writeError writes an RFC 7807 Problem Details JSON error response
func writeError(w http.ResponseWriter, status int, title, detail string) {
	problem := ProblemDetails{
		Type:   "about:blank",
		Title:  title,
		Status: status,
		Detail: detail,
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(problem)
}
