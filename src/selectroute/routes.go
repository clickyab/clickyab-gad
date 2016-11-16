package selectroute

import (
	"config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Use(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:   config.Config.StaticRoot,
			Browse: true,
		},
	))
	e.Get("/select", tc.selectWebAd)
	e.Get("/show/:type/:mega/:wid/:ad", tc.show)
	e.Get("/click/:wid/:mega/:ad/:rand", tc.click)
	e.Get("/conversion/", tc.conversion)
	e.Get("/ads/vast/", tc.selectVastAd)
}
