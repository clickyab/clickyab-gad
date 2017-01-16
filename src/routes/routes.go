package routes

import (
	"config"
	"middlewares"
	"statics"

	"gopkg.in/labstack/echo.v3"
)

// Routes function register all routes in system
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.GET("/select", tc.selectWebAd, middlewares.RequestCollector, middlewares.Header)
	e.GET("/show/:type/:mega/:wid/:ad", tc.show, middlewares.RequestCollector, middlewares.Header)
	e.GET("/click/:wid/:mega/:ad/:rand", tc.click, middlewares.RequestCollector, middlewares.Header)
	e.GET("/conversion/", tc.conversion, middlewares.RequestCollector, middlewares.Header)
	e.GET("/ads/vast/", tc.selectVastAd, middlewares.RequestCollector, middlewares.Header)
	e.GET("/allads", tc.allAds, middlewares.RequestCollector, middlewares.Header)
	e.GET("/version", tc.version, middlewares.Header)

	postfix := "-min.js"
	if config.Config.DevelMode {
		postfix = ".js"
	}
	e.GET("/show.js", tc.assetRoute("show"+postfix))
	e.GET("/vastAD.js", tc.assetRoute("vastAD"+postfix))
	e.GET("/conversion/clickyab-tracking.js", tc.assetRoute("clickyab-tracking"+postfix))
}

func (selectController) assetRoute(asset string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(200, string(statics.MustAsset(asset)))
	}
}
