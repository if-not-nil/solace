package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

// function to revoke invite link (admin only)
func (s *Server) LinksRevoke(c echo.Context) error {

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	inviteID := c.Param("invite_id")
	if inviteID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid invite link ID")
	}

	var link orm.InviteLink
	if err := s.db.First(&link, "id = ?", inviteID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "invite link not found")
	}

	mem, err := user.GetMembership(s.db, link.ServerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "youre prolly not in the server")
	}

	if !mem.CanUser("manage_invites", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	if err := s.db.Delete(&link).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to revoke invite link")
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "deleted"})
}
