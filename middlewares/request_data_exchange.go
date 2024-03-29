package middlewares

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"

	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"

	"net/http"

	"bytes"
	"io/ioutil"

	"strings"

	"clickyab.com/gad/models"
	"github.com/clickyab/services/config"
	"github.com/clickyab/simple-rtb"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

const (
	usdCur                   = "USD"
	rialCur                  = "IRR"
	requestDataTokenExchange = "__exchange__"
)

var (
	// centExchangeRate change cent to rial
	centExchangeRate = config.RegisterFloat64("demand.cent.rate", 42000, "cent to rial change rate")
	domain           = regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
)

func validateCurr(imps []srtb.Impression) (bool, float64) {
	if imps[0].Currency != usdCur && imps[0].Currency != rialCur {
		return false, 0
	}
	var curr = imps[0].Currency
	var rate float64
	if curr == usdCur {
		rate = centExchangeRate.Float64()
	} else if curr == rialCur {
		rate = 1
	}
	for i := range imps {
		if imps[i].Currency != curr {
			return false, 0
		}
		imps[i].BidFloor = imps[i].BidFloor * float64(rate)
	}
	return true, rate
}

// RequestExchangeCollectorGenerator try to collect data from request
func RequestExchangeCollectorGenerator(copKey func(echo.Context, *RequestData, int) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			e := &srtb.BidRequest{}

			tmp, err := ioutil.ReadAll(ctx.Request().Body)
			assert.Nil(err)
			buf := bytes.NewBuffer(tmp)
			logrus.Debug(string(tmp))
			dec := json.NewDecoder(buf)
			defer func() {
				err := ctx.Request().Body.Close()
				assert.Nil(err)
			}()
			err = dec.Decode(e)
			assert.Nil(err)
			if len(e.Imp) == 0 {
				return ctx.HTML(http.StatusBadRequest, "no impression")
			}
			//validate cur and do the conversion
			ok, rate := validateCurr(e.Imp)
			if !ok {
				return ctx.HTML(http.StatusBadRequest, "Wrong Currency")
			}
			if e.Site != nil {
				if !domain.MatchString(e.Site.Domain) {
					return ctx.HTML(http.StatusBadRequest, fmt.Sprintf("invalid publisher site %s", e.Site.Domain))
				}
			}
			ctx.Set(requestDataTokenExchange, e)
			rde := &RequestData{}
			//get supplier from its key
			rde.SupplierKey = ctx.Param("key")
			rde.IP = net.ParseIP(e.Device.IP)
			rde.UserAgent = e.Device.UA
			ua := user_agent.New(rde.UserAgent)
			browser, version := ua.Browser()
			rde.Rate = rate
			rde.Browser = browser
			rde.BrowserVersion = version
			rde.OS = ua.OS()
			rde.Mobile = ua.Mobile()
			rde.Platform = ua.Platform()
			rde.PlatformID = utils.FindOsID(ua.Platform())
			rde.Parent = e.Site.Ref
			rde.Referrer = e.Site.Page
			rde.MegaImp = e.ID
			rde.TID = e.User.ID
			if rde.TID == "" {
				rde.TID = utils.CreateHash(copLen.Int(), []byte(rde.UserAgent), []byte(rde.IP))
			}
			rde.CopID = models.NewManager().CreateCookieProfile(rde.TID, rde.IP).ID
			rde.Host = ctx.Request().Host
			rde.Scheme = ctx.Scheme()
			if xh := strings.ToLower(ctx.Request().Header.Get("X-Forwarded-Proto")); xh == https {
				rde.Scheme = "https"
			}
			if e.App != nil {
				//if v, ok := e.Attributes["brand"]; ok {
				//	rde.Brand = v.(string)
				//}
				//if v, ok := e.Attributes["cid"]; ok {
				//	rde.CID = int64(v.(float64))
				//}
				//if v, ok := e.Attributes["lac"]; ok {
				//	rde.LAC = int64(v.(float64))
				//}
				//if v, ok := e.Attributes["language"]; ok {
				//	rde.Language = v.(string)
				//}
				//if v, ok := e.Attributes["mcc"]; ok {
				//	rde.MCC = int64(v.(float64))
				//}
				//if v, ok := e.Attributes["mnc"]; ok {
				//	rde.MNC = int64(v.(float64))
				//}
				//if v, ok := e.Attributes["model"]; ok {
				//	rde.Model = v.(string)
				//}
				//if v, ok := e.Attributes["operator"]; ok {
				//	rde.Operator = v.(string)
				//}
				//if v, ok := e.Attributes["os_identity"]; ok {
				//	rde.OSIdentity = v.(string)
				//}
				//if v, ok := e.Attributes["carrier"]; ok {
				//	rde.Carrier = v.(string)
				//}
			}
			ctx.Set(requestDataToken, rde)
			return next(ctx)
		}
	}
}

// GetExchangeRequestData is the helper function to extract request data from context
func GetExchangeRequestData(ctx echo.Context) (*srtb.BidRequest, error) {
	rd, ok := ctx.Get(requestDataTokenExchange).(*srtb.BidRequest)
	if !ok {
		return nil, errors.New("not valid data in context")
	}

	return rd, nil
}

// MustExchangeGetRequestData try to get request data, or panic if there is no request data
func MustExchangeGetRequestData(ctx echo.Context) *srtb.BidRequest {
	rd, err := GetExchangeRequestData(ctx)
	assert.Nil(err)
	return rd
}
