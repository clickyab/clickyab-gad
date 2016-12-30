package middlewares

import (
	"assert"
	"config"
	"errors"
	"mr"
	"net"
	"utils"

	"github.com/mssola/user_agent"
	"gopkg.in/labstack/echo.v3"
)

// RequestData is the data for request
type RequestData struct {
	IP             net.IP
	UserAgent      string
	IP2Location    *mr.IP2Location
	Browser        string
	OS             string
	Platform       string
	PlatformID     int64
	BrowserVersion string
	Method         string
	Referrer       string
	Mobile         bool
	Host           string
	Scheme         string
	MegaImp        string
	CopID          int64
	TID            string
	Parent         string
}

const requestDataToken = "__request_data__"

// RequestCollector try to collect data from request
func RequestCollector(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		e := &RequestData{}
		e.IP = net.ParseIP(ctx.RealIP())
		e.UserAgent = ctx.Request().UserAgent()
		ua := user_agent.New(ctx.Request().UserAgent())
		e.Host = ctx.Request().Host
		e.Scheme = ctx.Scheme()
		name, version := ua.Browser()
		e.Browser = name
		e.BrowserVersion = version
		e.OS = ua.OS()
		e.Mobile = ua.Mobile()
		e.Platform = ua.Platform()
		e.PlatformID = config.FindOsID(ua.Platform())
		e.Referrer = ctx.Request().URL.Query().Get("ref")
		e.Method = ctx.Request().Method
		e.MegaImp = <-utils.ID
		e.Parent = ctx.Request().URL.Query().Get("parent")
		if e.Referrer == "" {
			e.Referrer = ctx.Request().Referer()
		}

		if e.TID = ctx.Request().URL.Query().Get("tid"); len(e.TID) < config.Config.Clickyab.CopLen {
			e.TID = utils.CreateCopID(e.UserAgent, e.IP, config.Config.Clickyab.CopLen)
		}
		e.CopID = mr.NewManager().CreateCookieProfile(e.TID, e.IP).ID
		ctx.Set(requestDataToken, e)
		return next(ctx)
	}
}

// GetRequestData is the helper function to extract request data from context
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
