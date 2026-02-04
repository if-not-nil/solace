package server

import (
	"net/http"
	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

func (s *Server) LinksJoin(c echo.Context) error {

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User
	inviteID := c.Param("invite_id")

	if inviteID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid invite id")
	}

	var invite orm.InviteLink
	if err := s.db.First(&invite, "id = ?", inviteID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "invite link not found")
	}

	if invite.JoinsLeft <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invite link has no uses left")
	}

	if user.InServer(s.db, invite.ServerID) {
		return echo.NewHTTPError(http.StatusBadRequest, "you are already a member of this server")
	}

	// Check if user was previously banned
	var existingMembership orm.Membership
	if err := s.db.Where("user_id = ? AND server_id = ?", user.ID, invite.ServerID).First(&existingMembership).Error; err == nil {
		if existingMembership.Banned {
			return echo.NewHTTPError(http.StatusForbidden, "you are banned from this server")
		}
	}

	tx := s.db.Begin()
	lockedInvite := invite

	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&lockedInvite, "id = ?", inviteID).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusNotFound, "invite link not found")
	}

	if lockedInvite.JoinsLeft <= 0 {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusBadRequest, "invite link has no uses left")
	}

	m := orm.Membership{
		ServerID: invite.ServerID,
		UserID:   user.ID,
	}
	if err := tx.Create(&m).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to join server")
	}
	lockedInvite.JoinsLeft -= 1
	if err := tx.Save(&lockedInvite).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update invite")
	}

	if err := tx.Commit().Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}
	return c.JSON(http.StatusCreated, map[string]any{"message": "joined server successfully"})
}
