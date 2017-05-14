package middlewares

import (
	"assert"
	"errors"
	"utils"

	"encoding/json"

	"mr"
	"net"

	"config"

	"github.com/mssola/user_agent"
	"gopkg.in/labstack/echo.v3"
)

// RequestDataExchange is the data for request
type RequestDataFromExchange struct {
	TrackID   string `json:"track_id"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`

	Source struct {
		Name         string                 `json:"name"`
		Supplier     string                 `json:"supplier"`
		FloorCPM     int                    `json:"floor_cpm"`
		SoftFloorCPM int                    `json:"soft_floor_cpm"`
		Attributes   map[string]interface{} `json:"attributes"`
	} `json:"source"`

	Location struct {
		Country struct {
			Valid bool   `json:"valid"`
			Name  string `json:"name"`
			ISO   string `json:"iso"`
		} `json:"country"`
		Province struct {
			Valid bool   `json:"valid"`
			Name  string `json:"name"`
		} `json:"province"`
		LatLon struct {
			Valid bool    `json:"valid"`
			Lat   float64 `json:"lat"`
			Long  float64 `json:"long"`
		} `json:"latlon"`
	} `json:"location"`

	Attributes map[string]interface{} `json:"attributes"`
	Slots      []Slot                 `json:"slots"`

	Category []string `json:"category"`

	Platform   string `json:"platform"`
	Underfloor bool   `json:"underfloor"`
}

type Slot struct {
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	TrackID string `json:"track_id"`
}
type Source struct {
	Website      string                 `json:"website"`
	Supplier     string                 `json:"supplier"`
	FloorCPM     int                    `json:"floor_cpm"`
	SoftFloorCPM int                    `json:"soft_floor_cpm"`
	Attributes   map[string]interface{} `json:"attributes"`
}
type Country struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
	ISO   string `json:"iso"`
}
type Province struct {
	Valid bool   `json:"valid"`
	Name  string `json:"name"`
}
type LatLon struct {
	Valid bool    `json:"valid"`
	Lat   float64 `json:"lat"`
	Long  float64 `json:"long"`
}

const requestDataTokenExchange = "__exchange__"

// RequestCollectorGenerator try to collect data from request
func RequestExchangeCollectorGenerator(copKey func(echo.Context, *RequestData, int) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			e := &RequestDataFromExchange{}

			dec := json.NewDecoder(ctx.Request().Body)
			defer ctx.Request().Body.Close()
			err := dec.Decode(e)
			assert.Nil(err)

			ctx.Set(requestDataTokenExchange, e)
			rde := RequestData{}
			rde.IP = net.ParseIP(e.IP)
			rde.UserAgent = e.UserAgent
			ua := user_agent.New(rde.UserAgent)
			browser, version := ua.Browser()
			rde.Browser = browser
			rde.BrowserVersion = version
			rde.OS = ua.OS()
			rde.Mobile = ua.Mobile()
			rde.Platform = ua.Platform()
			rde.PlatformID = config.FindOsID(ua.Platform())
			if v, ok := e.Attributes["referrer"]; ok {
				rde.Referrer = v.(string)
			}
			if vv, ok := e.Attributes["parent"]; ok {
				rde.Referrer = vv.(string)
			}
			rde.MegaImp = e.TrackID
			rde.TID = utils.CreateHash(config.Config.Clickyab.CopLen, []byte(rde.UserAgent), []byte(rde.IP))
			rde.CopID = mr.NewManager().CreateCookieProfile(rde.TID, rde.IP).ID
			ctx.Set(requestDataToken, rde)
			return next(ctx)
		}
	}
}

// GetExchangeRequestData is the helper function to extract request data from context
func GetExchangeRequestData(ctx echo.Context) (*RequestDataFromExchange, error) {
	rd, ok := ctx.Get(requestDataTokenExchange).(*RequestDataFromExchange)
	if !ok {
		return nil, errors.New("not valid data in context")
	}

	return rd, nil
}

// MustExchangeGetRequestData try to get request data, or panic if there is no request data
func MustExchangeGetRequestData(ctx echo.Context) *RequestDataFromExchange {
	rd, err := GetExchangeRequestData(ctx)
	assert.Nil(err)
	return rd
}
