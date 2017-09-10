package routes

import (
	"config"
	"mr"
	"net/http"
	"selector"
	"utils"

	"fmt"

	"strconv"

	"time"

	"net/url"

	"middlewares"

	"net"

	"errors"

	"ip2location"

	"redis"

	"redlock"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

// example request
// a.clickyab.com/ads/show.php?a=1941478513606&width=300&height=250&slot=48812001338&eventpage=995681655&ck=true&loc=http://myreal.ir/dara&ref=false&tid=2188781033

const typ = "sync"

func (tc *selectController) showphp(c echo.Context) error {

	wpidReq := c.Request().URL.Query().Get("a")
	width := c.Request().URL.Query().Get("width")
	height := c.Request().URL.Query().Get("height")
	slotReq := c.Request().URL.Query().Get("slot")
	eventpage := c.Request().URL.Query().Get("eventpage")
	//ckReq := c.Request().URL.Query().Get("ck")
	loc := c.Request().URL.Query().Get("loc")
	//refReq := c.Request().URL.Query().Get("ref")
	tid := c.Request().URL.Query().Get("tid")

	rd := middlewares.MustGetRequestData(c)
	rd.TID = tid
	size, err := config.GetSize(fmt.Sprintf("%sx%s", width, height))
	if err != nil {
		return c.HTML(http.StatusBadRequest, "wrong slot size")
	}
	website, provinceID, ispID, err := locationStatus(tc, wpidReq, rd.IP)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	slotSize, sizeNumSlice := tc.slotSizeNormal([]string{slotReq}, website.WID, map[string]int{slotReq: size})
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    provinceID,
		ISP:         ispID,
	}
	lockSession := "DMDS_SESS_" + eventpage
	t := redlock.NewRedisDistributedLock(lockSession, 3000*time.Millisecond)
	t.Lock()
	defer t.Unlock()
	var sel selector.FilterFunc
	var sessionAds []int64
	sel = webSelector

	filteredAds := selector.Apply(&m, selector.GetAdData(), sel)
	c.Set("EVENT_PAGE", eventpage)
	_, pubsAds := tc.makeShow(c, typ, rd, filteredAds, nil, sizeNumSlice, slotSize, nil, website, false, config.Config.Clickyab.MinCPCWeb, config.Config.Clickyab.UnderFloor, true, config.Config.Clickyab.FloorDiv.Web)
	targetedAd := pubsAds[slotReq]
	if targetedAd == nil {
		return c.String(http.StatusNotFound, "not found")
	}
	ad, err := mr.NewManager().GetAd(targetedAd.AdID, false)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	var rnd string
	for i := range pubsAds {
		sessionAds = append(sessionAds, pubsAds[i].AdID)
		rnd = <-utils.ID
		imp := tc.fillNativeImp(rd, false, pubsAds[i], pubsAds[i].WinnerBid, website, pubsAds[i].SlotID)
		tc.callWebWorker(website, pubsAds[i].SlotID, pubsAds[i].AdID, m.RequestData.MegaImp, rnd, imp, rd)
		if eventpage != "" {
			err := aredis.SAddInt(eventpage, true, time.Minute, sessionAds...)
			if err != nil {
				logrus.Debug(err)
			}
		}
	}

	adURL := makeAdURL(rd, website.WID, ad.AdID, m.RequestData.MegaImp, rnd)
	res, err := tc.makeWebTemplate(c, "", ad, adURL, "", "", loc[:5] == "https")
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	return c.HTML(http.StatusOK, res)
}

func makeAdURL(rd *middlewares.RequestData, wID, adID int64, mega string, rnd string) string {
	u := url.URL{
		Scheme: rd.Scheme,
		Host:   rd.Host,
		Path:   fmt.Sprintf("/click/%s/%d/%s/%d/%s", "web", wID, mega, adID, rnd),
	}

	v := url.Values{}
	v.Set("tid", rd.TID)
	v.Set("ref", rd.Referrer)
	v.Set("parent", rd.Parent)
	u.RawQuery = v.Encode()

	return u.String()
}

func locationStatus(tc *selectController, wpid string, ip net.IP) (*mr.Website, int64, int64, error) {
	wpID, err := strconv.ParseInt(wpid, 10, 0)

	website, err := tc.fetchWebsite(wpID)
	if err != nil {
		return nil, 0, 0, errors.New("wrong website")
	}

	provinceID, ispID := ip2location.GetProvinceISPByIP(ip)

	return website, provinceID, ispID, nil
}
