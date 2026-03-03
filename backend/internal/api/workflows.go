// Package api contains the HTTP handlers for the payer service
package api

import (
	"net/http"

	"evolutionary-mcp/backend/internal/repository"
	"evolutionary-mcp/backend/pkg/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

// GetWorkflow returns a single workflow by ID
// (GET /api/v1/workflows/:id)
func (s *Server) GetWorkflow(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()
	workflow, err := s.Repo.GetWorkflow(ctx, id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Workflow not found: "+err.Error())
	}

	return c.JSON(http.StatusOK, workflow)
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

	// Logic for versioned vs non-versioned save
	if workflow.SaveAsNewVersion || workflow.WorkflowID == "" {
		// Create new version or new workflow concept
		if workflow.WorkflowID == "" {
			workflow.WorkflowID = uuid.New().String()
		}
		if err := s.Repo.CreateWorkflow(ctx, &workflow); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create workflow version: "+err.Error())
		}
	} else {
		// Update existing LATEST version for this concept (Draft mode)
		if err := s.Repo.UpdateWorkflow(ctx, &workflow); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update workflow: "+err.Error())
		}
	}

	return c.JSON(http.StatusOK, workflow)
}
