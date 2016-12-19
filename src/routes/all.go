package routes

import (
	"selector"

	"gopkg.in/labstack/echo.v3"
)

func (tc *selectController) allAds(c echo.Context) error {
	return c.JSON(200, selector.GetAdData())
}
