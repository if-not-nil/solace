package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

// update role
// @Summary update role
// @Description update permissions for an existing role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "server id"
// @Param role_id path string true "role id"
// @Param request body RoleUpdateRequest true "role update data"
// @Success 200 {object} RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /server/{id}/roles/{role_id} [put]
func (s *Server) RoleUpdate(c echo.Context) error {
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

	var req map[string]bool
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	// build new permission array from booleans
	var newPerms []string
	if req["administrator"] {
		newPerms = append(newPerms, "administrator")
	}
	if req["manage_server"] {
		newPerms = append(newPerms, "manage_server")
	}
	if req["manage_roles"] {
		newPerms = append(newPerms, "manage_roles")
	}
	if req["manage_channels"] {
		newPerms = append(newPerms, "manage_channels")
	}
	if req["kick_users"] {
		newPerms = append(newPerms, "kick_users")
	}
	if req["ban_users"] {
		newPerms = append(newPerms, "ban_users")
	}
	if req["create_invites"] {
		newPerms = append(newPerms, "create_invites")
	}
	if req["manage_invites"] {
		newPerms = append(newPerms, "manage_invites")
	}
	if req["send_messages"] {
		newPerms = append(newPerms, "send_messages")
	}
	if req["manage_messages"] {
		newPerms = append(newPerms, "manage_messages")
	}
	if req["embed_links"] {
		newPerms = append(newPerms, "embed_links")
	}
	if req["attach_files"] {
		newPerms = append(newPerms, "attach_files")
	}
	if req["mention_everyone"] {
		newPerms = append(newPerms, "mention_everyone")
	}

	// check that user can modify this role
	if !mem.CanAssignRole(&orm.Role{Permissions: newPerms}, s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "cannot modify role with permissions you don't have")
	}

	role.Permissions = newPerms
	if err := s.db.Save(&role).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update role")
	}

	return c.JSON(http.StatusOK, role)
}
