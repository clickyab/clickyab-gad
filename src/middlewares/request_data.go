package middlewares

import (
	"assert"
	"config"
	"errors"
	"mr"
	"net"
	"strings"
	"utils"

	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

// RequestData is the data for request
type RequestData struct {
	IP             net.IP
	UserAgent      string
	IP2Location    *mr.IP2Location
	Browser        string
	SuppliersName  string
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
	ISP            string

	//App part
	Network    string
	Brand      string
	CID        int64
	LAC        int64
	MCC        int64
	MNC        int64
	Language   string
	Model      string
	Operator   string
	OSIdentity string
	Carrier    string
	//other app stuff
	GoogleID      string
	AndroidID     string
	AndroidDevice string
}

const requestDataToken = "__request_data__"

// RequestCollectorGenerator try to collect data from request
func RequestCollectorGenerator(copKey func(echo.Context, *RequestData, int) string, tag string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ll := logrus.Fields{
				"tag": tag,
			}
			e := &RequestData{}
			e.IP = net.ParseIP(ctx.RealIP())
			ll["ip"] = e.IP.String()
			e.UserAgent = ctx.Request().UserAgent()
			ua := user_agent.New(ctx.Request().UserAgent())
			ll["ua"] = ctx.Request().UserAgent()
			e.Host = ctx.Request().Host
			ll["host"] = e.Host
			e.Scheme = ctx.Scheme()
			if xh := strings.ToLower(ctx.Request().Header.Get("X-Forwarded-Proto")); xh == "https" {
				e.Scheme = "https"
			}
			ll["https"] = e.Scheme == "https"
			name, version := ua.Browser()
			e.Browser = name
			ll["browser"] = name
			e.BrowserVersion = version
			ll["version"] = version
			e.OS = ua.OS()
			ll["os"] = e.OS
			e.Mobile = ua.Mobile()
			ll["mobile"] = e.Mobile
			e.Platform = ua.Platform()
			ll["platform"] = e.Platform
			e.PlatformID = config.FindOsID(ua.Platform())
			e.Referrer = ctx.Request().URL.Query().Get("ref")
			e.Method = ctx.Request().Method
			e.MegaImp = <-utils.ID
			e.Parent = ctx.Request().URL.Query().Get("parent")
			ll["parent"] = e.Parent
			if e.Referrer == "" {
				e.Referrer = ctx.Request().Referer()
			}
			ll["ref"] = e.Referrer

			if e.TID = ctx.Request().URL.Query().Get("tid"); len(e.TID) < config.Config.Clickyab.CopLen {
				e.TID = copKey(ctx, e, config.Config.Clickyab.CopLen)
			}
			e.CopID = mr.NewManager().CreateCookieProfile(e.TID, e.IP).ID

			//extract app stuff
			e.GoogleID = ctx.Request().URL.Query().Get("GoogleAdvertisingId")
			e.AndroidID = ctx.Request().URL.Query().Get("androidid")
			e.AndroidDevice = ctx.Request().URL.Query().Get("deviceid")

			ctx.Set(requestDataToken, e)
			ctx.Set("LOG", &ll)
			err := next(ctx)
			logrus.WithFields(ll).Info("RequestCompleted")
			return err
		}
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
