package server

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"solace/middleware"
	"solace/orm"
)

// delete role
// @Summary delete role
// @Description delete an existing role from the server
// @Tags roles
// @Produce json
// @Param id path string true "server id"
// @Param role_id path string true "role id"
// @Success 200 {object} map[string]string
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /server/{id}/roles/{role_id} [delete]
func (s *Server) RoleDelete(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	serverID := c.Param("id")
	roleID := c.Param("role_id")

	mem, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !mem.CanUser("manage_roles", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	var role orm.Role
	if err := s.db.First(&role, "id = ? AND server_id = ?", roleID, serverID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "role not found")
	}

	// cannot delete roles with permissions user doesn't have
	if !mem.CanAssignRole(&role, s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "cannot delete role with permissions you don't have")
	}

	// remove all user-role assignments for this role
	if err := s.db.Where("role_id = ?", roleID).Delete(&orm.UserRole{}).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to remove role assignments")
	}

	if err := s.db.Delete(&role).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete role")
	}

	return c.JSON(http.StatusOK, map[string]any{"message": "role deleted"})
}
