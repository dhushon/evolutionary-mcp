package api

import (
	"net/http"

	"evolutionary-mcp/backend/pkg/models"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// ListGroundingRules returns all rules for the tenant
// (GET /api/v1/grounding)
func (s *Server) ListGroundingRules(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, _ := ctx.Value("tenant_id").(string)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Tenant ID missing")
	}

	rules, err := s.Repo.ListGroundingRules(ctx, tenantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rules)
}

// CreateGroundingRule creates a new rule
// (POST /api/v1/grounding)
func (s *Server) CreateGroundingRule(c echo.Context) error {
	ctx := c.Request().Context()
	tenantID, _ := ctx.Value("tenant_id").(string)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Tenant ID missing")
	}

	var rule models.GroundingRule
	if err := c.Bind(&rule); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	rule.TenantID = tenantID

	if err := s.Repo.CreateGroundingRule(ctx, &rule); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, rule)
}

// GetGroundingRule returns a single rule
// (GET /api/v1/grounding/:id)
func (s *Server) GetGroundingRule(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()

	rule, err := s.Repo.GetGroundingRule(ctx, id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Rule not found")
	}

	return c.JSON(http.StatusOK, rule)
}

// UpdateGroundingRule updates an existing rule
// (PUT /api/v1/grounding/:id)
func (s *Server) UpdateGroundingRule(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()
	tenantID, _ := ctx.Value("tenant_id").(string)

	var rule models.GroundingRule
	if err := c.Bind(&rule); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	rule.ID = id.String()
	rule.TenantID = tenantID

	if err := s.Repo.UpdateGroundingRule(ctx, &rule); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rule)
}

// DeleteGroundingRule removes a rule
// (DELETE /api/v1/grounding/:id)
func (s *Server) DeleteGroundingRule(c echo.Context, id openapi_types.UUID) error {
	ctx := c.Request().Context()

	if err := s.Repo.DeleteGroundingRule(ctx, id.String()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
