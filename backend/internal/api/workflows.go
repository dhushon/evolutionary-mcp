// Package api contains the HTTP handlers for the payer service
package api

import (
	"net/http"

	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Server holds the dependencies for the API server.
type Server struct {
	Repo repository.Repository
}

// NewServer creates a new Server.
func NewServer(repo repository.Repository) *Server {
	return &Server{Repo: repo}
}

// ListWorkflows returns a list of all workflows
// (GET /api/v1/workflows)
func (s *Server) ListWorkflows(c echo.Context) error {
	ctx := c.Request().Context()

	workflows, err := s.Repo.ListWorkflows(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, workflows)
}

// PutWorkflow creates or updates a workflow
// (PUT /api/v1/workflows)
func (s *Server) PutWorkflow(c echo.Context) error {
	ctx := c.Request().Context()

	var workflow models.Workflow
	if err := c.Bind(&workflow); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Tenant ID not found in context")
	}
	workflow.TenantID = tenantID

	// If this is a new workflow concept (no WorkflowID), generate one.
	// If WorkflowID is present, the repo will treat it as an evolution of that workflow.
	if workflow.WorkflowID == "" {
		workflow.WorkflowID = uuid.New().String()
	}

	if err := s.Repo.CreateWorkflow(ctx, &workflow); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to save workflow: "+err.Error())
	}

	return c.JSON(http.StatusOK, workflow)
}
