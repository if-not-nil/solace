package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) LinksNew(c echo.Context) error {

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

	var req InviteLinksNewRequest
	if err := c.Bind(&req); err != nil || req.MaxUsers == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "valid max_users required")
	}

	link := orm.InviteLink{
		ServerID:  serverID,
		JoinsLeft: req.MaxUsers,
	}

	if err := s.db.Create(&link).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create invite link")
	}

	return c.JSON(http.StatusCreated, link)
}
