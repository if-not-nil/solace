package server

import (
	"net/http"

	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) LinksList(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	serverID := c.Param("id")
	if serverID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid server ID")
	}

	mem, err := user.GetMembership(s.db, serverID)

	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "youre prolly not in the server")
	}
	if !mem.CanUser("manage_invites", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	var links []orm.InviteLink
	if err := s.db.Where("server_id = ?", serverID).Find(&links).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch invite links")
	}

	return c.JSON(http.StatusOK, links)
}
