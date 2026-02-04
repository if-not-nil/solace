package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"solace/middleware"
	"solace/orm"
)

// get user roles
// @Summary get user roles
// @Description get roles assigned to the current user in the server
// @Tags roles
// @Produce json
// @Param id path string true "server id"
// @Success 200 {array} RoleResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /server/{id}/user/roles [get]
func (s *Server) UserRoles(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	serverID := c.Param("id")

	_, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	var userRoles []orm.UserRole
	if err := s.db.Where("user_id = ? AND server_id = ?", user.ID, serverID).Preload("Role").Find(&userRoles).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to fetch user roles")
	}

	var roles []orm.Role
	for _, ur := range userRoles {
		roles = append(roles, ur.Role)
	}

	return c.JSON(http.StatusOK, roles)
}

