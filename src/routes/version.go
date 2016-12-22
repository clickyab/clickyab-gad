package routes

import (
	"net/http"
	"version"

	"gopkg.in/labstack/echo.v3"
)

func (tc *selectController) version(c echo.Context) error {
	return c.JSON(http.StatusOK, version.GetVersion())
}
