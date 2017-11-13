package routes

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"clickyab.com/gad/ip2location"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/mr"
	selector2 "clickyab.com/gad/pin"
	"clickyab.com/gad/redis"
	"clickyab.com/gad/redlock"
	"clickyab.com/gad/selector"
	"clickyab.com/gad/store"
	"clickyab.com/gad/utils"

	"github.com/sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

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
	size, err := utils.GetSize(fmt.Sprintf("%sx%s", width, height))
	if err != nil {
		return c.HTML(http.StatusBadRequest, "wrong slot size")
	}
	website, provinceID, ispID, err := locationStatus(c, tc, wpidReq, rd.IP)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	middlewares.SetData(c, "site_id", website.WID)
	middlewares.SetData(c, "site_domain", website.WDomain.String)
	middlewares.SetData(c, "ad_count", 1)
	middlewares.SetData(c, "ad_size", size)

	slotSize, sizeNumSlice := tc.slotSizeNormal([]string{slotReq}, website.WID, map[string]int{slotReq: size})

	lockSession := "DMDS_SESS_" + eventpage
	t := redlock.NewRedisDistributedLock(lockSession, 3000*time.Millisecond)
	t.Lock()
	defer t.Unlock()
	var sel selector.FilterFunc
	var sessionAds []int64

	sel = webSelector
	var slotFixFound bool
	slotPins := selector2.GetPinAdData()

	slotFixFound, slotSize, sizeNumSlice, slotPins, _, _ = checkForFixSlot(slotPins, slotSize, sizeNumSlice, "banner")
	m := selector.Context{
		RequestData:      *rd,
		Website:          website,
		Size:             sizeNumSlice,
		Province:         provinceID,
		ISP:              ispID,
		SlotPins:         slotPins,
		MinBidPercentage: 1, // TODO : hard coded shit
	}
	// TODO remove slot fix ads from normal pool

	c.Set("EVENT_PAGE", eventpage)
	var pubsAds = make(map[string]*mr.AdData)
	if !slotFixFound {
		filteredAds := selector.Apply(&m, selector.GetAdData(), sel)
		_, pubsAds = tc.makeShow(c,
			"sync",
			rd,
			filteredAds,
			nil,
			sizeNumSlice,
			slotSize,
			nil,
			website,
			false,
			minCPCWeb.Int64(),
			allowUnderFloor.Bool(),
			true,
			floorDivWeb.Int64(),
		)
	} else {
		res := &slotPins[0].AdData
		res.WinnerBid = slotPins[0].Bid
		res.SlotID = slotPins[0].SlotID
		pubsAds[slotReq] = res
		tc.updateMegaKey(rd, slotPins[0].AdID, slotPins[0].Bid, slotPins[0].SlotID, "", "", "")
		reserve := make(map[string]string)
		slotPubID := slotReq
		tmp := <-utils.ID
		reserve[slotPubID] = tmp
		store.Set(reserve[slotPubID], fmt.Sprintf("%d", slotPins[0].AdID))
	}

	targetedAd := pubsAds[slotReq]
	if targetedAd == nil {
		return c.String(http.StatusNotFound, "not found")
	}
	ad, err := mr.NewManager().GetAd(targetedAd.AdID, false)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	ad.RawSlotSize = &mr.RawSlotDimensions{
		Width:  width,
		Height: height,
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
	if slotFixFound && slotPins[0].Direct {
		adURL = ad.AdURL.String
	}
	showT := false
	if rd.Mobile && provinceID > 0 && rand.Intn(chanceShowT.Int()) == 1 {
		showT = true
		middlewares.SetData(c, "show_t", 1)
	}
	res, err := tc.makeWebTemplate(c, "", ad, adURL, "", "", loc[:5] == "https", showT)
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

func locationStatus(c echo.Context, tc *selectController, wpid string, ip net.IP) (*mr.Website, int64, int64, error) {
	wpID, _ := strconv.ParseInt(wpid, 10, 0)
	website, err := tc.fetchWebsite(wpID)
	if err != nil {
		return nil, 0, 0, errors.New("wrong website")
	}

	provinceID, ispID, ll := ip2location.GetProvinceISPByIP(ip)
	middlewares.SetData(c, "province", ll.Province)
	middlewares.SetData(c, "country", ll.Country)
	middlewares.SetData(c, "city", ll.City)
	middlewares.SetData(c, "isp", ll.ISP)

	return website, provinceID, ispID, nil
}

func checkForFixSlot(pins []mr.SlotPinData, a map[string]*slotData, b map[string]int, typ string) (bool, map[string]*slotData, map[string]int, []mr.SlotPinData, map[string]*slotData, map[string]int) {
	var fix []mr.SlotPinData
	var fixSizeNumSlice = make(map[string]int)
	var fixSlotSize = make(map[string]*slotData)
	var found bool
	for i := range pins {
		for j := range a {
			if pins[i].SlotID == a[j].ID {
				//check chance
				ok1 := checkFixChance(pins[i].Chance)
				//check size
				ok2 := checkFixSlotSize(pins[i], typ)
				if ok1 && ok2 {
					pins[i].SlotPublicID = j
					fix = append(fix, pins[i])
					found = true
					fixSlotSize[j] = a[j]
					fixSizeNumSlice[j] = a[j].SlotSize
					delete(a, j)
					delete(b, j)
				}
			}
		}
	}
	return found, a, b, fix, fixSlotSize, fixSizeNumSlice
}

func checkFixChance(a int) bool {

	return a >= rand.Intn(100)
}

func checkFixSlotSize(a mr.SlotPinData, typ string) bool {
	if typ == "vast" {
		if a.CampaignNetwork != 2 && a.AdType == 3 {
			return false
		}
		if a.AdType == utils.AdTypeDynamic {
			return false
		}
		if a.CampaignNetwork != 0 && a.CampaignNetwork != 2 {
			return false
		}
		return a.AdType == utils.AdTypeVideo || utils.InVastSize(a.AdSize)
	} else if typ == "banner" {
		if a.CampaignNetwork != 0 {
			return false
		}
		if a.AdType == utils.AdTypeVideo {
			if utils.InVideoSize(a.AdSize) {
				return true
			}
			return false
		}
		return a.SlotSize == a.AdSize
	} else if typ == "native" {
		return a.AdSize == 20 && a.CampaignNetwork == 3
	}
	panic("[BUG] unsupported type")
}
