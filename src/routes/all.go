package routes

import (
	"selector"

	"github.com/labstack/echo"
)

func (tc *selectController) allAds(c echo.Context) error {
	return c.JSON(200, selector.GetAdData())
}
