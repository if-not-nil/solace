package server

import (
	"net/http"
	"strings"

	"solace/orm"
	"solace/util"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func (s *Server) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	name := strings.TrimSpace(req.Name)
	password := strings.TrimSpace(req.Password)

	if name == "" || password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and password are required")
	}
	if len(name) > 50 || len(password) > 128 {
		return echo.NewHTTPError(http.StatusBadRequest, "name or password too long")
	}

	var user orm.User
	if err := user.ByName(s.db, name); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "user doesn't exist")
	}

	if !util.VerifyPassword(password, user.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid password")
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	var userRes UserResponse
	copier.Copy(&userRes, &user)

	return c.JSON(http.StatusOK, LoginResponse{
		Message: "login successful",
		Token:   token,
		User:    userRes,
	})
}
