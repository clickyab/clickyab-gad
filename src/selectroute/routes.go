package selectroute

import (
	"config"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"middlewares"
)

// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Use(middleware.StaticWithConfig(
		middleware.StaticConfig{
			Root:   config.Config.StaticRoot,
			Browse: true,
		},
	))
	e.Get("/select", tc.selectWebAd,middlewares.RequestCollector,middlewares.Header)
	e.Get("/show/:type/:mega/:wid/:ad", tc.show,middlewares.RequestCollector,middlewares.Header)
	e.Get("/click/:wid/:mega/:ad/:rand", tc.click,middlewares.RequestCollector,middlewares.Header)
	e.Get("/conversion/", tc.conversion,middlewares.RequestCollector,middlewares.Header)
	e.Get("/ads/vast/", tc.selectVastAd,middlewares.RequestCollector,middlewares.Header)
	e.Get("/apply", tc.applyAd,middlewares.RequestCollector,middlewares.Header)
	e.Get("/allads", tc.allAds,middlewares.RequestCollector,middlewares.Header)
}
