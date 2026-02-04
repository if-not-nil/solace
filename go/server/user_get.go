package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

// get user info
// @Summary get user info
// @Description get current user information and servers
// @Tags users
// @Produce json
// @Success 200 {object} UserMeResponse
// @Failure 500 {object} ErrorResponse
// @Router /user/me [get]
func (s *Server) UserGet(c echo.Context) error {

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	var memberships []orm.Membership
	if err := s.db.Preload("Server").Where("user_id = ? AND banned = ?", user.ID, false).Find(&memberships).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch user data")
	}

	res := struct {
		User    orm.UserResponse `json:"user"`
		Servers []orm.ServerResponse
	}{
		User: orm.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			AvatarID: user.AvatarID,
		},
		Servers: make([]orm.ServerResponse, 0, len(memberships)),
	}

	for _, m := range memberships {
		var channels []orm.Channel
		s.db.Where("server_id = ?", m.ServerID).Find(&channels)

		channelRes := make([]orm.ChannelResponse, len(channels))
		for i, ch := range channels {
			channelRes[i] = orm.ChannelResponse{
				ID:       ch.ID,
				Name:     ch.Name,
				ServerID: ch.ServerID,
				Type:     ch.Type,
			}
		}

		res.Servers = append(res.Servers, orm.ServerResponse{
			ID:       m.Server.ID,
			Name:     m.Server.Name,
			OwnerID:  m.Server.OwnerID,
			AvatarID: m.Server.AvatarID,
			Channels: channelRes,
		})
	}

	return c.JSON(http.StatusOK, res)
}
