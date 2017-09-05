package routes

import (
	"config"
	"middlewares"
	"net"
	"statics"
	"utils"

	"gopkg.in/labstack/echo.v3"
)

func webCopCreateor(ctx echo.Context, e *middlewares.RequestData, len int) string {
	//return utils.CreateCopID(ctx.Request().UserAgent(), net.ParseIP(ctx.RealIP()), len)
	return utils.CreateHash(len, []byte(ctx.Request().UserAgent()), []byte(net.ParseIP(ctx.RealIP())))
}

func appCopCreateor(ctx echo.Context, e *middlewares.RequestData, len int) string {
	return utils.CreateHash(
		len,
		[]byte(ctx.Request().URL.Query().Get("androidid")),
		[]byte(ctx.Request().URL.Query().Get("deviceid")),
		[]byte(ctx.Request().URL.Query().Get("operator")),
		[]byte(ctx.Request().URL.Query().Get("model")),
	)
}

// Routes function register all routes in system
func (tc *selectController) Routes(e *echo.Echo, _ string) {
	e.Use(middlewares.ServerID)
	e.GET("/native", tc.selectNativeAd, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/select", tc.selectWebAd, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.POST("/demand", tc.selectDemandAd, middlewares.RequestExchangeCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/show/:type/:mega/:wid/:ad", tc.show, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/ads/", tc.showphp, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/click/:typ/:wid/:mega/:ad/:rand", tc.click, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/conversion/", tc.conversion, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/ads/vast/", tc.selectVastAd, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.POST("/allads", tc.allAds, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/version", tc.version, middlewares.Header)
	e.GET("/ads/inapp.php", tc.inApp, middlewares.RequestCollectorGenerator(appCopCreateor), middlewares.Header)
	e.GET("/ads/json-inapp.php", tc.inAppJson, middlewares.Header)
	e.GET("/healthz", tc.healthz, middlewares.Header)
	e.GET("/showjs", tc.showjs, middlewares.RequestCollectorGenerator(webCopCreateor), middlewares.Header)

	postfix := "-min.js"
	if config.Config.DevelMode {
		postfix = ".js"
	}
	e.GET("/show.js", tc.assetRoute("show"+postfix))
	e.GET("/vastAD.js", tc.assetRoute("vastAD"+postfix))
	e.GET("/conversion/clickyab-tracking.js", tc.assetRoute("clickyab-tracking"+postfix))

	//echo.NotFoundHandler = fcgi.NewPHPFastCGIHandler(config.Config.PHPCode.Root, "/", config.Config.PHPCode.FPM, 30*time.Second, 30*time.Second, 30*time.Second)
}

func (selectController) assetRoute(asset string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(200, string(statics.MustAsset(asset)))
	}
}
