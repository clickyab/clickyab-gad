package routes

import (
	"net"

	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/statics_src"
	"clickyab.com/gad/utils"

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
	e.GET("/native", tc.selectNativeAd, middlewares.RequestCollectorGenerator(webCopCreateor, "native"), middlewares.Header)
	e.POST("/demand", tc.selectDemandAd, middlewares.RequestExchangeCollectorGenerator(webCopCreateor), middlewares.Header)
	e.GET("/ads/", tc.showphp, middlewares.RequestCollectorGenerator(webCopCreateor, "show"), middlewares.Header)
	e.GET("/click/:typ/:wid/:mega/:ad/:rand", tc.click, middlewares.RequestCollectorGenerator(webCopCreateor, "click"), middlewares.Header)
	e.GET("/conversion/", tc.conversion, middlewares.RequestCollectorGenerator(webCopCreateor, "conv"), middlewares.Header)
	e.GET("/ads/vast/", tc.selectVastAd, middlewares.RequestCollectorGenerator(webCopCreateor, "vast"), middlewares.Header)
	e.POST("/allads", tc.allAds, middlewares.RequestCollectorGenerator(webCopCreateor, "debug"), middlewares.Header)
	e.GET("/allads", tc.allAdsTemp, middlewares.RequestCollectorGenerator(webCopCreateor, "debug"), middlewares.Header)
	e.GET("/version", tc.version, middlewares.Header)
	e.GET("/ads/inapp.php", tc.inApp, middlewares.RequestCollectorGenerator(appCopCreateor, "app"), middlewares.Header)
	e.GET("/ads/json-inapp.php", tc.inAppJson, middlewares.Header)
	e.GET("/healthz", tc.healthz, middlewares.Header)
	e.GET("/show.js", tc.showjs, middlewares.RequestCollectorGenerator(webCopCreateor, "showjs"), middlewares.Header)
	e.GET("/js/jwvast.js", tc.vastJS)
	e.GET("/js/videovast.js", tc.videoJS)
	e.GET("/js/native.js", tc.nativeJS)

	postfix := "-min.js"
	if develMode.Bool() {
		postfix = ".js"
	}
	//e.GET("/show.js", tc.assetRoute("show"+postfix))
	e.GET("/vastAD.js", tc.assetRoute("vastAD"+postfix))
	e.GET("/conversion/clickyab-tracking.js", tc.assetRoute("clickyab-tracking"+postfix))

	e.GET("/select", tc.selectWebAd, middlewares.RequestCollectorGenerator(webCopCreateor, "old"), middlewares.Header)
	e.GET("/show/:type/:mega/:wid/:ad", tc.show, middlewares.RequestCollectorGenerator(webCopCreateor, "multi"), middlewares.Header)

	//echo.NotFoundHandler = fcgi.NewPHPFastCGIHandler(config.Config.PHPCode.Root, "/", config.Config.PHPCode.FPM, 30*time.Second, 30*time.Second, 30*time.Second)
}

func (selectController) assetRoute(asset string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.HTML(200, string(statics_src.MustAsset(asset)))
	}
}
