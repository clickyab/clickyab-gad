package selectroute

import (
	"github.com/labstack/echo"
	"selector"
)

func (tc *selectController) allAds(c echo.Context) error {
	return c.JSON(200, selector.GetAdData())
}
