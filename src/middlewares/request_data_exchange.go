package middlewares

import (
	"assert"
	"config"
	"encoding/json"
	"errors"
	"fmt"
	"mr"
	"net"
	"regexp"
	"utils"

	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	"github.com/mssola/user_agent"
	"gopkg.in/labstack/echo.v3"
)

// RequestDataExchange is the data for request
type RequestDataFromExchange struct {
	TrackID   string `json:"track_id"`
	IP        string `json:"ip"`
	Scheme    string `json:"scheme"`
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

	Platform    string `json:"platform"`
	Underfloor  bool   `json:"underfloor"`
	SessionKey  string `json:"page_track_id"`
	UserTrackID string `json:"user_track_id"`
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

var domain = regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)

// RequestCollectorGenerator try to collect data from request
func RequestExchangeCollectorGenerator(copKey func(echo.Context, *RequestData, int) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			e := &RequestDataFromExchange{}

			tmp, err := ioutil.ReadAll(ctx.Request().Body)
			assert.Nil(err)
			buf := bytes.NewBuffer(tmp)
			logrus.Debug(string(tmp))
			dec := json.NewDecoder(buf)
			defer ctx.Request().Body.Close()
			err = dec.Decode(e)
			assert.Nil(err)

			if e.Platform == "web" {
				if !domain.MatchString(e.Source.Name) {
					return ctx.JSON(http.StatusBadRequest, fmt.Errorf("invalid publisher site %s", e.Source.Name))
				}
			}
			ctx.Set(requestDataTokenExchange, e)
			rde := &RequestData{}
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
			rde.TID = e.UserTrackID
			if rde.TID == "" {
				rde.TID = utils.CreateHash(config.Config.Clickyab.CopLen, []byte(rde.UserAgent), []byte(rde.IP))
			}
			rde.CopID = mr.NewManager().CreateCookieProfile(rde.TID, rde.IP).ID
			rde.Host = ctx.Request().Host
			rde.Scheme = "http"
			if e.Scheme == "https" {
				rde.Scheme = e.Scheme
			}

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
