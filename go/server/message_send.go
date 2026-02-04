package server

import (
	"net/http"
	"strings"

	"solace/middleware"
	"solace/orm"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

// send message
// @Summary send message
// @Router /channel/{id}/send [post]
func (s *Server) MessageSend(c echo.Context) error {
	user := c.Get("auth").(middleware.Auth).User
	channelID := c.Param("id")

	var req ChannelSendRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "content is required")
	}
	if len(req.Content) > 2000 {
		return echo.NewHTTPError(http.StatusBadRequest, "content too long")
	}
	var channel orm.Channel
	if err := s.db.First(&channel, "id = ?", channelID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "channel not found")
	}
	if err := user.ByID(s.db, user.ID); err != nil { // check if user is in server
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}
	mem, err := user.GetMembership(s.db, channel.ServerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}
	if !mem.CanUser("send_messages", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}
	message := orm.Message{Content: req.Content, ChannelID: channelID, UserID: user.ID}
	if err := s.db.Create(&message).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to send message")
	}
	if err := s.db.Preload("User").First(&message, "id = ?", message.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to load message user")
	}

	var messageRes orm.MessageResponse
	copier.Copy(&messageRes, message)

	// send to websocket clients
	SendToChannel(message.ChannelID, messageRes)

	response := MessageSendResponse{
		Message: "message sent successfully",
	}

	return c.JSON(http.StatusOK, response)
}
