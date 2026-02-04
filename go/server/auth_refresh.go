package server

import (
	"net/http"
	"solace/middleware"
	"solace/util"

	"github.com/labstack/echo/v4"
)

// refresh token
// @Summary refresh token
// @Description generate a new JWT token
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [get]
func (s *Server) Refresh(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)

	token, err := util.GenerateJWT(auth.User.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]any{"token": token})
}
