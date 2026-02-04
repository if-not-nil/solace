package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) ChannelDelete(c echo.Context) error {

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	channelID := c.Param("id")

	if channelID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing server or channel ID")
	}

	var channel orm.Channel
	if err := s.db.First(&channel, "id = ?", channelID).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "channel not found")
	}

	mem, err := user.GetMembership(s.db, channel.ServerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "youre prolly not in the server")
	}

	if !mem.CanUser("manage_channels", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	if err := s.db.Delete(&channel).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not delete channel")
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
