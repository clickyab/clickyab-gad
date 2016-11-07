package selectroute

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Use(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:   "/home/develop/gad/showjs",
			Browse: true,
		},
	))
	e.Get("/select", tc.selectAd)
	e.Get("/show/:mega/:ad", tc.show)
}
