package middlewares

import (
	_ "fmt"
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
	OsVersion      string
	BrowserVersion string
	Method         string
	Referrer       string
}


// Recovery is the middleware to prevent the panic to crash the app
func Test(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		e := &RequestData{}
		ua := user_agent.New(ctx.Request().Header().Get("User-Agent"))
		name, version := ua.Browser()
		e.Browser=name
		e.BrowserVersion=version
		e.Os=ua.OS()
		e.OsVersion=ua.Platform()
		e.RealIP = ctx.Request().RealIP()
		e.Referrer = ctx.Request().Referer()
		e.Method = ctx.Request().Method()


		ctx.Set("RequestData", e)

		return next(ctx)
	}
}
