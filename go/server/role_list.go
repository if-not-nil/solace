package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

// list roles
// @Summary list roles
// @Description get all roles for the server
// @Tags roles
// @Produce json
// @Param id path string true "server id"
// @Success 200 {array} RoleResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /server/{id}/roles [get]
func (s *Server) RoleList(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	serverID := c.Param("id")

	_, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	var roles []orm.Role
	if err := s.db.Where("server_id = ?", serverID).Find(&roles).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch roles")
	}

	return c.JSON(http.StatusOK, roles)
}
