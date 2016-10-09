package middlewares

import (
	_ "fmt"
	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
	"mr"
	"strconv"
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

// RequestCollector try to collect data from request
func RequestCollector(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		e := &RequestData{}
		ua := user_agent.New(ctx.Request().UserAgent())
		name, version := ua.Browser()
		e.Browser = name
		e.BrowserVersion = version
		e.Os = ua.OS()
		e.OsVersion = ua.Platform()
		e.RealIP = ctx.Request().RealIP()
		e.Referrer = ctx.Request().Referer()
		e.Method = ctx.Request().Method()

		ctx.Set("RequestData", e)

		params := ctx.QueryParams()
		public_id, _ := strconv.Atoi(params["i"][0])
		domain := params["d"][0]

		////fetch website and set in Context
		wd, err := mr.NewManager().FetchWebsite(public_id, domain)
		if err != nil {
			logrus.Fatal(err)
		}
		ctx.Set("WebsiteData", wd)
		//web:=new selector.Context{WebsiteData:}
		////fetch size and add to context

		return next(ctx)
	}
}
