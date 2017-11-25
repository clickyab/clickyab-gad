package middlewares

import (
	"errors"
	"net"
	"strings"

	"clickyab.com/gad/models"
	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"

	"github.com/clickyab/services/config"
	"github.com/mssola/user_agent"
	"gopkg.in/labstack/echo.v3"
)

// RequestData is the data for request
type RequestData struct {
	IP             net.IP
	UserAgent      string
	IP2Location    *models.IP2Location
	Browser        string
	SupplierKey    string
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
	Alexa          bool
	Rate           float64

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

const (
	requestDataToken = "__request_data__"
	https            = "https"
)

var (
	copLen = config.RegisterInt("clickyab.cop_len", 10, "cop key len")
)

// RequestCollectorGenerator try to collect data from request
func RequestCollectorGenerator(copKey func(echo.Context, *RequestData, int) string, tag string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			SetData(ctx, "type", tag)
			e := &RequestData{}
			e.IP = net.ParseIP(ctx.RealIP())
			uaStr := ctx.Request().UserAgent()
			SetData(ctx, "ip", e.IP)
			e.UserAgent = ctx.Request().UserAgent()
			ua := user_agent.New(uaStr)
			SetData(ctx, "ua", uaStr)
			e.Host = ctx.Request().Host
			SetData(ctx, "host", e.Host)
			e.Scheme = ctx.Scheme()
			if xh := strings.ToLower(ctx.Request().Header.Get("X-Forwarded-Proto")); xh == https {
				e.Scheme = "https"
			}
			SetData(ctx, https, e.Scheme == https)
			name, version := ua.Browser()
			e.Browser = name
			e.BrowserVersion = version
			SetData(ctx, "version", version)
			e.OS = ua.OS()
			e.Mobile = ua.Mobile()
			e.Platform = ua.Platform()
			if e.Platform == "" && uaStr == "CLICKYAB" {
				e.Platform = "ClickyabSDK"
				e.OS = "Android"
				e.Mobile = true
				e.Browser = "AndroidSDK"
			}
			SetData(ctx, "platform", e.Platform)
			SetData(ctx, "os", e.OS)
			SetData(ctx, "mobile", e.Mobile)
			SetData(ctx, "browser", name)

			e.PlatformID = utils.FindOsID(ua.Platform())
			e.Referrer = ctx.Request().URL.Query().Get("ref")
			e.Method = ctx.Request().Method
			e.MegaImp = <-utils.ID
			e.Parent = ctx.Request().URL.Query().Get("parent")
			if e.Parent != "" {
				SetData(ctx, "parent", e.Parent)
			}
			if e.Referrer == "" {
				e.Referrer = ctx.Request().Referer()
			}
			if e.Referrer != "" {
				SetData(ctx, "ref", e.Referrer)
			}

			if e.TID = ctx.Request().URL.Query().Get("tid"); len(e.TID) < copLen.Int() {
				e.TID = copKey(ctx, e, copLen.Int())
			}
			e.CopID = models.NewManager().CreateCookieProfile(e.TID, e.IP).ID

			//extract app stuff
			e.GoogleID = ctx.Request().URL.Query().Get("GoogleAdvertisingId")
			e.AndroidID = ctx.Request().URL.Query().Get("androidid")
			e.AndroidDevice = ctx.Request().URL.Query().Get("deviceid")

			// In go headers are not case sensitive and ok with _ and -
			if strings.Contains(uaStr, "Alexa") || ctx.Request().Header.Get("ALEXATOOLBAR-ALX_NS_PH") != "" {
				e.Alexa = true
				SetData(ctx, "alexa", 1)
			}

			ctx.Set(requestDataToken, e)
			err := next(ctx)
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
