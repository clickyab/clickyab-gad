package selectroute

import (
	"config"
	"middlewares"
	"statics"

	"github.com/labstack/echo"
)

// Routes function @todo
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Get("/select", tc.selectWebAd, middlewares.RequestCollector, middlewares.Header)
	e.Get("/show/:type/:mega/:wid/:ad", tc.show, middlewares.RequestCollector, middlewares.Header)
	e.Get("/click/:wid/:mega/:ad/:rand", tc.click, middlewares.RequestCollector, middlewares.Header)
	e.Get("/conversion/", tc.conversion, middlewares.RequestCollector, middlewares.Header)
	e.Get("/ads/vast/", tc.selectVastAd, middlewares.RequestCollector, middlewares.Header)
	e.Get("/apply", tc.applyAd, middlewares.RequestCollector, middlewares.Header)
	e.Get("/allads", tc.allAds, middlewares.RequestCollector, middlewares.Header)

	postfix := "-min.js"
	if config.Config.DevelMode {
		postfix = ".js"
	}
	e.Get("/show.js", tc.assetRoute("show"+postfix))
	e.Get("/conversion/clickyab-tracking.js", tc.assetRoute("clickyab-tracking"+postfix))
}

func (selectController) assetRoute(asset string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(200, string(statics.MustAsset(asset)))
	}
}
