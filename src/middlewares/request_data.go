package middlewares

import (
	"assert"
	"errors"
	"mr"

	"github.com/labstack/echo"
	"github.com/mssola/user_agent"
)

// RequestData is the data for request
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

const requestDataToken = "__request_data__"

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
