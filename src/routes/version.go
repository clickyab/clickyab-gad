package routes

import (
	"version"

	"net/http"

	"gopkg.in/labstack/echo.v3"
)

func (tc *selectController) version(c echo.Context) error {
	return c.JSON(http.StatusOK, version.GetVersion())
}
