package server

import (
	"net/http"
	"strings"

	"solace/middleware"
	"solace/orm"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func (s *Server) ChannelNew(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	serverID := c.Param("id")
	if serverID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid server ID")
	}

	var req ChannelNewRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	channelName := strings.TrimSpace(req.Name)
	if channelName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "channel name is required")
	}

	if len(channelName) < 1 || len(channelName) > 50 {
		return echo.NewHTTPError(http.StatusBadRequest, "channel name must be between 1 and 50 characters")
	}

	mem, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !mem.CanUser("manage_channels", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	server := mem.Server

	channel := orm.Channel{
		ServerID: serverID,
		Name:     channelName,
		Type:     "text",
	}

	if err := server.CreateChannel(s.db, &channel); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var res orm.ChannelResponse
	copier.Copy(&res, &channel)

	return c.JSON(http.StatusCreated, res)
}
