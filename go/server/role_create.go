package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) RoleCreate(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	serverID := c.Param("id")

	var server orm.Server
	if err := s.db.First(&server, "id = ?", serverID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "server not found")
	}

	mem, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !mem.CanUser("manage_roles", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	var req struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "role name is required")
	}

	// check that user has all permissions they're trying to assign
	if !mem.CanAssignRole(&orm.Role{Permissions: req.Permissions}, s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "cannot create role with permissions you don't have")
	}

	role := orm.Role{
		ServerID:    serverID,
		Name:        req.Name,
		Permissions: req.Permissions,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create role")
	}

	return c.JSON(http.StatusCreated, role)
}
