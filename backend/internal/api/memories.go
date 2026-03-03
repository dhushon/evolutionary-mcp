package api

import (
	"net/http"

	"evolutionary-mcp/backend/internal/repository"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ListMemories returns all memories for the tenant
// (GET /api/v1/memories)
func (s *Server) ListMemories(c echo.Context) error {
	ctx := c.Request().Context()
	memories, err := s.Repo.ListMemories(ctx, "") // TenantID extracted from ctx inside ListMemories
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, memories)
}

// SearchMemories performs semantic search
// (POST /api/v1/memories/search)
func (s *Server) SearchMemories(c echo.Context) error {
	var body struct {
		Query string `json:"query"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Inject MemoryService into Server to handle embedding + search logic
	return c.JSON(http.StatusOK, []repository.Memory{})
}

// GiveMemoryFeedback updates confidence
// (POST /api/v1/memories/:id/feedback)
func (s *Server) GiveMemoryFeedback(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()
	var feedback MemoryFeedback
	if err := c.Bind(&feedback); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	memory, err := s.Repo.Get(ctx, id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Memory not found")
	}

	memory.Confidence = float64(feedback.Confidence)
	memory.Version++

	if err := s.Repo.Update(ctx, memory); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
