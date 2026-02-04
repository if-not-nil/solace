package server

import (
	"net/http"
	"slices"
	"solace/middleware"
	"solace/orm"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Server) ChannelHistory(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	channelID := c.Param("id")

	if channelID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid channel ID")
	}

	var req ChannelHistoryRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	var channel orm.Channel
	if err := s.db.First(&channel, "id = ?", channelID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "channel not found")
	}

	var dbUser orm.User
	if err := dbUser.ByID(s.db, user.ID); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
	}

	mem, err := dbUser.GetMembership(s.db, channel.ServerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !mem.CanUser("view_channels", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	count := req.Count
	if count <= 0 {
		count = 20
	} else if count > 100 {
		count = 100
	}

	untilTime := time.Now()
	if req.Until != nil && *req.Until > 0 {
		untilTime = time.Unix(*req.Until, 0)
	}

	var messages []orm.Message
	err = s.db.Preload("User").
		Where("channel_id = ? AND created_at <= ?", channelID, untilTime).
		Order("created_at DESC").
		Limit(count).
		Find(&messages).Error

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch messages")
	}

	slices.Reverse(messages)

	userMap := make(map[string]UserResponse)
	plainMessages := make([]MessageHistoryResponse, len(messages))

	for i, msg := range messages {
		// unique users
		if _, exists := userMap[msg.UserID]; !exists {
			userMap[msg.UserID] = UserResponse{
				ID:       msg.User.ID,
				Name:     msg.User.Name,
				AvatarID: msg.User.AvatarID,
			}
		}

		plainMessages[i] = MessageHistoryResponse{
			ID:        msg.ID,
			Content:   msg.Content,
			UserID:    msg.UserID,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		}
	}

	users := make([]UserResponse, 0, len(userMap))
	for _, u := range userMap {
		users = append(users, u)
	}

	return c.JSON(http.StatusOK, ChannelHistoryResponse{
		Messages: plainMessages,
		Users:    users,
	})
}
