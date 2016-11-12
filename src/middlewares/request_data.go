package middlewares

import (
	"assert"
	"config"
	"errors"
	"mr"

	"net"

	"utils"

	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
)

// RequestData is the data for request
type RequestData struct {
	CloudIP        string
	IP             net.IP
	UserAgent      string
	IP2Location    *mr.IP2Location
	Browser        string
	Os             string
	Platform       string
	PlatformID     int64
	BrowserVersion string
	Method         string
	Referrer       string
	Mobile         bool
	Url            string
	Proto            string

	MegaImp string
	CopID   string
}

const requestDataToken = "__request_data__"

// RequestCollector try to collect data from request
func RequestCollector(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		e := &RequestData{}
		e.IP = net.ParseIP(ctx.Request().RealIP())
		e.UserAgent = ctx.Request().UserAgent()
		ua := user_agent.New(ctx.Request().UserAgent())
		e.Url = ctx.Request().Host()
		e.Proto = ctx.Request().Scheme()
		name, version := ua.Browser()
		e.Browser = name
		e.UserAgent = ctx.Request().UserAgent()
		e.BrowserVersion = version
		e.Os = ua.OS()
		e.Mobile = ua.Mobile()
		e.Platform = ua.Platform()
		e.PlatformID = config.FindOsID(ua.Platform())
		e.Referrer = ctx.Request().Referer()
		e.Method = ctx.Request().Method()
		e.MegaImp = <-utils.ID
		e.CopID = utils.CreateCopID(e.UserAgent, e.IP)

		ctx.Set(requestDataToken, e)
		return next(ctx)
	}
}

// GetRequestData is the hgelper function to extract request data from context
func GetRequestData(ctx echo.Context) (*RequestData, error) {
	rd, ok := ctx.Get(requestDataToken).(*RequestData)
	if !ok {
		return nil, errors.New("not valid data in context")
	}

	return rd, nil
}

// MustGetRequestData try to get request data, or panic if there is no request data
func MustGetRequestData(ctx echo.Context) *RequestData {
	rd, err := GetRequestData(ctx)
	assert.Nil(err)
	return rd
}
