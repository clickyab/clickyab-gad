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
	e.Get("/select", tc.selectAd)
	e.Get("/show/:mega/:ad", tc.show)
	e.Get("/click/:imp/:data", tc.click)
}
