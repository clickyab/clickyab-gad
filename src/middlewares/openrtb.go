package middlewares

import (
	"bytes"
	"encoding/json"

	"assert"
	"io/ioutil"

	"errors"

	"config"
	"mr"
	"net"
	"utils"

	"net/http"

	"github.com/bsm/openrtb"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

const rtbDataToken = "__rtb_data__"

// RequestOpenRTBCollectorGenerator try to collect data from request
func RequestOpenRTBCollectorGenerator(copKey func(echo.Context, *RequestData, int) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			e := &openrtb.BidRequest{}
			tmp, err := ioutil.ReadAll(ctx.Request().Body)
			assert.Nil(err)
			defer ctx.Request().Body.Close()
			buf := bytes.NewBuffer(tmp)
			logrus.Debug(string(tmp))
			dec := json.NewDecoder(buf)
			err = dec.Decode(e)
			assert.Nil(err)

			err = e.Validate()
			if err != nil {
				return ctx.JSON(http.StatusBadRequest, nil)
			}

			ctx.Set(rtbDataToken, e)
			rde := &RequestData{}
			rde.IP = net.ParseIP(e.Device.IP)
			rde.UserAgent = e.Device.UA
			ua := user_agent.New(rde.UserAgent)
			browser, version := ua.Browser()
			rde.Browser = browser
			rde.BrowserVersion = version
			rde.OS = ua.OS()
			rde.Mobile = ua.Mobile()
			rde.Platform = ua.Platform()
			rde.PlatformID = config.FindOsID(ua.Platform())
			rde.MegaImp = e.ID
			rde.TID = e.User.ID
			if rde.TID == "" {
				rde.TID = utils.CreateHash(config.Config.Clickyab.CopLen, []byte(rde.UserAgent), []byte(rde.IP))
			}
			rde.CopID = mr.NewManager().CreateCookieProfile(rde.TID, rde.IP).ID
			rde.Host = ctx.Request().Host
			rde.Scheme = "http"
			ctx.Set(requestDataToken, rde)
			return next(ctx)
		}
	}
}

// GetRtbRequestData is the helper function to extract request data from context
func GetRtbRequestData(ctx echo.Context) (*openrtb.BidRequest, error) {
	rd, ok := ctx.Get(rtbDataToken).(*openrtb.BidRequest)
	if !ok {
		return nil, errors.New("not valid data in context")
	}

	return rd, nil
}

// MustRtbGetRequestData try to get request data, or panic if there is no request data
func MustRtbGetRequestData(ctx echo.Context) *openrtb.BidRequest {
	rd, err := GetRtbRequestData(ctx)
	assert.Nil(err)
	return rd
}
