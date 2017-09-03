package routes

import (
	"bytes"
	"config"
	"mr"
	"net/http"
	"selector"
	"utils"

	"fmt"

	"strconv"

	"time"

	"net/url"

	"assert"
	"middlewares"
	"redlock"

	"net"

	"errors"

	echo "gopkg.in/labstack/echo.v3"
	"ip2location"
)

const redLockTTL = time.Second

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

	}

	slotSize, sizeNumSlice := tc.slotSizeNormal([]string{slotReq}, website.WPubID, map[string]int{slotReq: size})

	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    provinceID,
		ISP:         ispID,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)

	redlock.Lock(eventpage, "", redLockTTL)
	_, pubsAds := tc.makeShow(c, typ, rd, filteredAds, nil, sizeNumSlice, slotSize, nil, website, false, config.Config.Clickyab.MinCPCWeb, config.Config.Clickyab.UnderFloor, true, config.Config.Clickyab.FloorDiv.Web)
	redlock.Unlock([]string{eventpage}, "")

	targetedAd, ok := pubsAds[slotReq]
	assert.True(len(pubsAds) == 1 && ok, "more than one or no ads where found")

	ad, err := mr.NewManager().GetAd(targetedAd.AdID, false)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}

	buf := &bytes.Buffer{}
	adURL := makeAdURL(rd, website.WPubID, ad.AdID, eventpage)
	res := tc.makeSingleAdData(ad, adURL, loc[:5] == "https")

	err = singleAdTemplate.Execute(buf, res)
	assert.Nil(err)

	return c.HTML(http.StatusOK, buf.String())
}

func makeAdURL(rd *middlewares.RequestData, wpid, adID int64, mega string) string {
	u := url.URL{
		Scheme: rd.Scheme,
		Host:   rd.Host,
		Path:   fmt.Sprintf("/click/%s/%d/%s/%d/%s", "web", wpid, mega, adID, <-utils.ID),
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

	provinceID := ip2location.GetProvinceIDByIP(ip)
	ispID, _ := ip2location.GetProvinceISPByIP(ip)

	return website, provinceID, ispID, nil
}
