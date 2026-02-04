package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) RoleAssign(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	actingUser := auth.User

	serverID := c.Param("id")
	targetUserID := c.Param("user_id")
	roleID := c.Param("role_id")

	var server orm.Server
	if err := s.db.First(&server, "id = ?", serverID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "server not found")
	}

	actingUserMembership, err := actingUser.GetMembership(s.db, server.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !actingUserMembership.CanUser("manage_roles", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	var targetMembership orm.Membership
	if err := s.db.First(&targetMembership, "user_id = ? AND server_id = ?", targetUserID, server.ID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "target user is not a member of this server")
	}

	var role orm.Role
	if err := s.db.First(&role, "id = ? AND server_id = ?", roleID, serverID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "role not found")
	}

	if !actingUserMembership.CanAssignRole(&role, s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "cannot assign this role")
	}

	userRole := orm.UserRole{
		UserID:   targetUserID,
		RoleID:   roleID,
		ServerID: serverID,
	}

	if err := s.db.Create(&userRole).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to assign role")
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "role assigned successfully"})
}
