package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

// kick user
// @Summary kick user
// @Description remove a user from the server and ban them
// @Tags servers
// @Accept json
// @Produce json
// @Param id path string true "server id"
// @Param request body ServerKickRequest true "user to kick"
// @Success 200 {object} map[string]any
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /server/{id}/kick [post]
func (s *Server) ServerKick(c echo.Context) error {

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	serverID := c.Param("id")
	if serverID == "" {
		return echo.NewHTTPError(http.StatusNotFound, "server not found")
	}
	mem, err := user.GetMembership(s.db, serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "you are not a member of this server")
	}

	if !mem.CanUser("kick_users", s.db) {
		return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
	}

	var req struct {
		ID string `json:"id"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if req.ID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id is required")
	}

	var m orm.Membership
	if err := s.db.Preload("Server").First(&m, "user_id = ? AND server_id = ?", req.ID, serverID).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user not found in server")
	}

	if m.Server.OwnerID == req.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot kick the server owner")
	}

	m.Banned = true
	if err := s.db.Save(&m).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not ban user from server")
	}
	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
