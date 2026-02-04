package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

// get server
// @Summary get server
// @Description get server information and channels
// @Tags servers
// @Produce json
// @Param id path string true "server id"
// @Success 200 {object} ServerResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /server/{id} [get]
func (s *Server) ServerGet(c echo.Context) error {
	serverID := c.Param("id")
	var server orm.Server

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	if !user.InServer(s.db, serverID) {
		return echo.NewHTTPError(http.StatusUnauthorized, "not a member")
	}

	if err := s.db.First(&server, "id = ?", serverID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "server not found")
	}

	var channels []orm.Channel
	if err := s.db.Where("server_id = ?", serverID).Find(&channels).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not retrieve channels")
	}

	var channelRes []orm.ChannelResponse
	copier.Copy(&channelRes, &channels)

	var res orm.ServerResponse
	copier.Copy(&res, &server)
	res.Channels = channelRes

	return c.JSON(http.StatusOK, res)
}
