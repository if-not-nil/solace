package server

import (
	"net/http"
	"strings"

	"solace/middleware"
	"solace/orm"

	"github.com/labstack/echo/v4"
)

// create server
// @Summary create server
// @Description create a new server with the user as owner
// @Tags servers
// @Accept json
// @Produce json
// @Param request body ServerNewRequest true "server creation data"
// @Success 201 {object} ServerResponse
// @Failure 400 {object} ErrorResponse
// @Router /server [post]
func (s *Server) ServerNew(c echo.Context) error {

	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	var req ServerNewRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "server name is required")
	}

	if len(req.Name) < 3 || len(req.Name) > 50 {
		return echo.NewHTTPError(http.StatusBadRequest, "server name must be between 3 and 50 characters")
	}

	server := orm.Server{Name: req.Name, OwnerID: user.ID}
	if err := s.db.Create(&server).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	general := orm.Channel{Name: "general", ServerID: server.ID, Type: "text"}
	if err := server.CreateChannel(s.db, &general); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create general channel")
	}

	mem := orm.Membership{
		UserID:   user.ID,
		ServerID: server.ID,
	}

	if err := s.db.Create(&mem).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server created, but failed to assign owner")
	}

	res := orm.ServerResponse{
		ID:       server.ID,
		Name:     server.Name,
		OwnerID:  server.OwnerID,
		AvatarID: server.AvatarID,
		Channels: []orm.ChannelResponse{{
			ID:       general.ID,
			Name:     general.Name,
			ServerID: general.ServerID,
			Type:     general.Type,
		}},
	}

	return c.JSON(http.StatusCreated, res)
}
