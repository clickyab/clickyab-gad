package middlewares

import (
	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
	"mr"
)

type RequestData struct {
	CloudIP        string
	RealIP         string
	IP2Location    *mr.IP2Location
	Browser        string
	Os             string
	Platform       string
	BrowserVersion string
	Method         string
	Referrer       string
	Mobile         bool
}

type Size []int

// RequestCollector try to collect data from request
func RequestCollector(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		e := &RequestData{}
		ua := user_agent.New(ctx.Request().UserAgent())
		name, version := ua.Browser()
		e.Browser = name
		e.BrowserVersion = version
		e.Os = ua.OS()
		e.Mobile = ua.Mobile()
		e.Platform = ua.Platform()
		e.RealIP = ctx.Request().RealIP()
		e.Referrer = ctx.Request().Referer()
		e.Method = ctx.Request().Method()

		ctx.Set("RequestData", e)
		return next(ctx)
	}
}
