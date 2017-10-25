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

	selector2 "pin"

	"math/rand"

	"store"

	"github.com/sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

// example request
// a.clickyab.com/ads/show.php?a=1941478513606&width=300&height=250&slot=48812001338&eventpage=995681655&ck=true&loc=http://myreal.ir/dara&ref=false&tid=2188781033

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
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    provinceID,
		ISP:         ispID,
		SlotPins:    slotPins,
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
			config.Config.Clickyab.MinCPCWeb,
			config.Config.Clickyab.UnderFloor,
			true,
			config.Config.Clickyab.FloorDiv.Web,
		)
	} else {
		res := &slotPins[0].AdData
		res.WinnerBid = slotPins[0].Bid
		res.SlotID = slotPins[0].SlotID
		pubsAds[slotReq] = res
		tc.updateMegaKey(rd, slotPins[0].AdID, slotPins[0].Bid, slotPins[0].SlotID, "", "", "")
		reserve := make(map[string]string)
		slotPubID := slotReq
		tmp := config.Config.MachineName + <-utils.ID
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
		if a.AdType == config.AdTypeDynamic {
			return false
		}
		if a.CampaignNetwork != 0 && a.CampaignNetwork != 2 {
			return false
		}
		return a.AdType == config.AdTypeVideo || config.InVastSize(a.AdSize)
	} else if typ == "banner" {
		if a.CampaignNetwork != 0 {
			return false
		}
		if a.AdType == config.AdTypeVideo {
			if config.InVideoSize(a.AdSize) {
				return true
			}
			return false
		}
		return a.SlotSize == a.AdSize
	} else if typ == "native" {
		return a.AdSize == 20 && a.CampaignNetwork == 3
	} else {
		panic("[BUG] unsupported type")
	}
}
