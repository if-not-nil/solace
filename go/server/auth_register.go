package server

import (
	"net/http"
	"strings"

	"solace/orm"
	"solace/util"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
)

func (s *Server) Register(c echo.Context) error {
	var req RegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)

	if req.Name == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and password are required")
	}

	if len(req.Name) < 3 || len(req.Name) > 50 {
		return echo.NewHTTPError(http.StatusBadRequest, "name must be between 3 and 50 characters")
	}

	if len(req.Password) < 8 || len(req.Password) > 128 {
		return echo.NewHTTPError(http.StatusBadRequest, "password must be between 8 and 128 characters")
	}

	var existing orm.User
	if err := s.db.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	user := orm.User{
		Name:     req.Name,
		Password: hashedPassword,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	token, err := util.GenerateJWT(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}
	var res orm.UserResponse
	copier.Copy(&res, &user)
	c.SetCookie(&http.Cookie{Name: "token", Value: token})
	return c.JSON(http.StatusCreated, map[string]any{
		"message": "registration successful",
		"token":   token,
		"user":    res,
	})
}
