package server

import (
	"net/http"
	"solace/middleware"

	"github.com/labstack/echo/v4"
)

func (s *Server) SetServerAvatar(c echo.Context) error {
	serverID := c.Param("id")
	if serverID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid server id")
	}

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	mem, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !mem.CanUser("manage_server", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	fileID := c.QueryParam("id")
	if fileID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "file id is required")
	}

	if err := s.db.Model(&mem.Server).Update("avatar_id", fileID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update avatar")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id": fileID,
	})
}
