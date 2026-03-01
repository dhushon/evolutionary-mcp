package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Server) GetHealth(ctx echo.Context) error {
	// Check database connectivity using the request context
	if err := h.Repo.Ping(ctx.Request().Context()); err != nil {
		status := "unavailable"
		service := "evolutionary-mcp"
		version := "1.0.0"
		now := time.Now()
		return ctx.JSON(http.StatusServiceUnavailable, HealthStatus{
			Status:    &status,
			Timestamp: &now,
			Service:   &service,
			Version:   &version,
		})
	}

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

func (h *Server) GetStatus(ctx echo.Context) error {
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

func (h *Server) PatchWorkflow(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}

func (h *Server) DeleteWorkflow(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNotImplemented)
}
