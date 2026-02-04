package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"solace/middleware"
	"solace/orm"
	"solace/util"

	"github.com/labstack/echo/v4"
)

func parseFileUpload(c echo.Context, field string) (bytes []byte, contentType string, err error) {
	fileHeader, err := c.FormFile(field)
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusBadRequest, "invalid file")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusInternalServerError, "could not open file")
	}
	defer file.Close()

	bytes, err = io.ReadAll(file)
	if err != nil {
		return nil, "", echo.NewHTTPError(http.StatusInternalServerError, "could not read file")
	}

	contentType = fileHeader.Header.Get("Content-Type")
	return bytes, contentType, nil
}

func (s *Server) FilePost(c echo.Context) error {
	auth := c.Get("auth").(middleware.Auth)
	user := auth.User

	if user.ID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	bytes, contentType, err := parseFileUpload(c, "file")
	if err != nil {
		return err
	}

	id, err := util.UploadFile(s.db, bytes, contentType)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not save file")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":   id,
		"type": contentType,
	})
}

func (s *Server) FileGet(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing file id")
	}

	var file orm.File
	if err := s.db.First(&file, "id = ?", id).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "file not found")
	}

	path := fmt.Sprintf("%s/%s", util.UploadDir, id)
	data, err := os.ReadFile(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not read file")
	}

	return c.Blob(http.StatusOK, file.Type, data)
}
