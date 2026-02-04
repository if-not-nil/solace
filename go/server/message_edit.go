package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"
	"strings"

	"github.com/labstack/echo/v4"
)

func (s *Server) MessageEdit(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	messageID := c.Param("id")
	var req struct {
		Content string `json:"content"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if messageID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid message id")
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "content cannot be empty")
	}

	if len(req.Content) > 2000 {
		return echo.NewHTTPError(http.StatusBadRequest, "content too long")
	}

	var message orm.Message
	if err := s.db.Preload("Channel").First(&message, "id = ?", messageID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "message not found")
	}

	if message.UserID != user.ID {
		mem, err := user.GetMembership(s.db, message.Channel.ServerID)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
		}

		if !mem.CanUser("manage_messages", s.db) {
			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
	message.Content = req.Content

	if err := s.db.Save(&message).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to edit message")
	}

	SendMessageEditedToChannel(message.ChannelID, messageID, req.Content)

	return c.JSON(http.StatusOK, map[string]any{"message": "message edited"})
}
