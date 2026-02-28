package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Handler contains HTTP handlers for the payer service REST API
type Handler struct {
}

// NewHandler creates a new Handler with required dependencies
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetHealth(ctx echo.Context) error {
	status := "ok"
	service := "evolutionary-mcp"
	version := "1.0.0"
	now := time.Now()
	return ctx.JSON(http.StatusOK, HealthStatus{
		Status:    &status,
		Timestamp: &now,
		Service:   &service,
		Version:   &version,
	})
}

func (h *Handler) GetStatus(ctx echo.Context) error {
	status := "ok"
	service := "evolutionary-mcp"
	version := "1.0.0"
	now := time.Now()
	return ctx.JSON(http.StatusOK, HealthStatus{
		Status:    &status,
		Timestamp: &now,
		Service:   &service,
		Version:   &version,
	})
}

func (h *Handler) ListWorkflows(ctx echo.Context) error {
	id := uuid.MustParse("0c1a4b6e-8e5e-4b1d-8c1a-4b6e8e5e4b1d")
	name := "test"
	description := "test workflow"
	workflows := []Workflow{
		{
			Id:          (*openapi_types.UUID)(&id),
			Name:        &name,
			Description: &description,
		},
	}
	return ctx.JSON(http.StatusOK, workflows)
}

func (h *Handler) PutWorkflow(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *Handler) PatchWorkflow(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *Handler) DeleteWorkflow(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
