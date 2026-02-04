package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) MessageDelete(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	messageID := c.Param("id")

	if messageID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid message id")
	}

	// find the message with channel info
	var message orm.Message
	if err := s.db.Preload("Channel").First(&message, "id = ?", messageID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}
	// check if user is the author or has manage messages permission

	if message.UserID != user.ID {
		mem, err := user.GetMembership(s.db, message.Channel.ServerID)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
		}

		if !mem.CanUser("manage_messages", s.db) {
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}

	// delete the message
	if err := s.db.Delete(&message).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete message")
	}

	// notify clients
	SendMessageDeletedToChannel(message.ChannelID, messageID)

	return c.JSON(http.StatusOK, map[string]any{"message": "message deleted"})
}
