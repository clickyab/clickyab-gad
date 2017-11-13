package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"time"

	"clickyab.com/gad/filter"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/models"
	"clickyab.com/gad/modules"
	"clickyab.com/gad/redis"
	"clickyab.com/gad/selector"
	"clickyab.com/gad/store"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"

	"clickyab.com/gad/ip2location"

	"clickyab.com/gad/js"

	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

var (
	webSelector = selector.Mix(
		filter.IsWebNetwork,
		filter.IsWebMobile,
		filter.CheckDesktopNetwork,
		filter.CheckMinBid,
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.CheckWebCategory,
		filter.CheckProvince,
		filter.RemoveSlotPins,
		filter.CheckVastOtherNetwork,
		filter.CheckISP,
	)

	nativeSelector = selector.Mix(
		filter.IsNativeNetwork,
		filter.IsNativeAd,
		filter.IsWebMobile,
		filter.CheckDesktopNetwork,
		filter.CheckMinBid,
		filter.CheckWebSize,
		filter.RemoveSlotPins,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.CheckWebCategory,
		filter.CheckProvince,
		filter.CheckISP,
	)

	slotReg = regexp.MustCompile(`s\[(\d*)\]`)
)

const webMobile string = "1000"

type selectController struct {
}

// SlotData is the single slot data in database
type slotData struct {
	SlotSize   int
	ID         int64
	PublicID   string
	Ctr        float64
	ExtraParam map[string]string
}

type vastSlotData struct {
	Type   string
	Offset string
	Repeat string
}

// Select function is the route that the real biding happen
func (tc *selectController) selectWebAd(c echo.Context) error {
	t := time.Now()
	rd, website, province, isp, err := tc.getWebDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice := tc.slotSizeWeb(c, *website, rd.Mobile)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
		ISP:         isp,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	show, _ := tc.makeShow(c, "web", rd, filteredAds, nil, sizeNumSlice, slotSize, nil, website, false, minCPCWeb.Int64(), allowUnderFloor.Bool(), true, floorDivWeb.Int64())

	//substitute the webMobile slot if exists
	wm := fmt.Sprintf("%d%s", website.WPubID, webMobile)
	val, ok := show[wm]
	if ok {
		show["web-mobile"] = val
		delete(show, wm)
	}

	b, _ := json.MarshalIndent(show, "\t", "\t")
	result := "renderFarm(" + string(b) + "); \n//" + time.Since(t).String()

	return c.HTML(200, result)
}

func (tc *selectController) doBid(adData *models.AdData, website Publisher, slot *slotData, floorDiv int64) bool {
	adData.CTR = tc.calculateCTR(
		adData,
		slot,
	)
	adData.CPM = utils.Cpm(adData.CampaignMaxBid, adData.CTR)
	//exceed cpm floor
	if floorDiv < 1 {
		floorDiv = 1
	}
	//logrus.Debugf("%d / %f ==> %d >= %d", adData.CampaignMaxBid, adData.CTR, adData.CPM, website.FloorCPM()/floorDiv)
	return adData.CPM >= website.FloorCPM()/floorDiv
}

func (tc *selectController) getSecondCPM(floorCPM int64, exceedFloor []*models.AdData) int64 {
	var secondCPM = floorCPM
	if len(exceedFloor) > 1 && exceedFloor[0].Capping.GetSelected() == exceedFloor[1].Capping.GetSelected() {
		secondCPM = exceedFloor[1].CPM
	}

	return secondCPM
}

func (tc *selectController) createMegaKey(rd *middlewares.RequestData, website Publisher) error {
	tmp := map[string]interface{}{
		"IP": rd.IP.String(),
		"UA": rd.UserAgent,
		"WS": fmt.Sprintf("%d", website.GetID()),
		"T":  fmt.Sprintf("%d", time.Now().Unix()),
	}
	assert.True(megaImpExpire.Int64() > 1, "invalid config")
	return aredis.HMSet(
		fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, rd.MegaImp),
		megaImpExpire.Duration(),
		tmp,
	)
}

func (tc *selectController) updateMegaKey(rd *middlewares.RequestData, adID int64, winnerBid int64, slotID int64, clickURL, clickParam, clickType string) {
	assert.Nil(aredis.StoreHashKey(
		fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, rd.MegaImp),
		fmt.Sprintf(
			"%s%s%d",
			transport.Advertise,
			transport.Delimiter,
			adID),
		fmt.Sprintf("%d", winnerBid),
		megaImpExpire.Duration(),
	))
	assert.Nil(aredis.StoreHashKey(
		fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, rd.MegaImp),
		fmt.Sprintf(
			"%s%s%d",
			transport.Slot,
			transport.Delimiter,
			adID),
		strconv.FormatInt(slotID, 10),
		megaImpExpire.Duration(),
	))
	if clickURL != "" && clickParam != "" {
		assert.Nil(aredis.StoreHashKey(
			fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, rd.MegaImp),
			fmt.Sprintf(
				"%s%s%d",
				transport.CustomClickURL,
				transport.Delimiter,
				slotID),
			clickURL,
			megaImpExpire.Duration(),
		))

		assert.Nil(aredis.StoreHashKey(
			fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, rd.MegaImp),
			fmt.Sprintf(
				"%s%s%d",
				transport.CustomClickParam,
				transport.Delimiter,
				slotID),
			clickParam,
			megaImpExpire.Duration(),
		))
		assert.Nil(aredis.StoreHashKey(
			fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, rd.MegaImp),
			fmt.Sprintf(
				"%s%s%d",
				transport.CustomClickType,
				transport.Delimiter,
				slotID),
			clickType,
			megaImpExpire.Duration(),
		))
	}

}

func (tc *selectController) getWebDataFromCtx(c echo.Context) (*middlewares.RequestData, *models.Website, int64, int64, error) {
	rd := middlewares.MustGetRequestData(c)
	params := c.QueryParams()
	publicParams, ok := params["i"]
	if !ok {
		return nil, nil, 0, 0, errors.New("invalid request")
	}
	publicID, err := strconv.ParseInt(publicParams[0], 10, 0)
	if err != nil {
		return nil, nil, 0, 0, errors.New("invalid request")
	}
	domain, ok := params["d"]
	if !ok {
		return nil, nil, 0, 0, errors.New("invalid request")
	}
	//fetch website and set in Context
	website, err := tc.fetchWebsite(publicID)
	if err != nil {
		return nil, nil, 0, 0, errors.New("invalid request")
	}

	if !website.GetActive() {
		return nil, nil, 0, 0, errors.New("web is not active")
	}

	if !models.NewManager().IsUserActive(website.UserID) {
		return nil, nil, 0, 0, errors.New("user is banned")
	}
	province, isp, ll := ip2location.GetProvinceISPByIP(rd.IP)
	middlewares.SetData(c, "province", ll.Province)
	middlewares.SetData(c, "country", ll.Country)
	middlewares.SetData(c, "city", ll.City)
	middlewares.SetData(c, "isp", ll.ISP)
	//check if the website domain is valid
	if website.WDomain.Valid && website.WDomain.String != domain[0] {
		return nil, nil, 0, 0, errors.New("domain and public id mismatch")
	}

	return rd, website, province, isp, nil
}

//FetchWebsite website and check if the minimum floor is applied
func (tc *selectController) fetchWebsite(publicID int64) (*models.Website, error) {
	website, err := models.NewManager().FetchWebsiteByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if website.WFloorCpm.Int64 < minCPMFloorWeb.Int64() {
		website.WFloorCpm.Int64 = minCPMFloorWeb.Int64()
	}
	return website, err
}

//fetchProvince find province and set context
func (tc *selectController) fetchProvinceDemand(r string) (*models.Province, error) {
	// if strings.ToUpper(cfHeader) != "IR" {
	// 	return nil, errors.New("not inside iran")
	// }
	var province models.Province
	province, err := models.NewManager().ConvertProvince2Info(r)
	if err != nil {
		return nil, errors.New("province not found")
	}
	return &province, nil

}

func (tc selectController) slotSizeWeb(c echo.Context, website models.Website, mobile bool, allAdsCase ...bool) (map[string]*slotData, map[string]int) {
	// main if for all ads data
	if len(allAdsCase) == 1 && allAdsCase[0] {
		var pubAd = make(map[string]*slotData)
		var pubSize = make(map[string]int)
		payload := c.Get("payload").(allAdsWebPayload)

		for _, pid := range payload.Slots {
			for i := 0; i < pid.Count; i++ {
				s, _ := utils.GetSize(pid.Size)
				r := fmt.Sprintf("%d", rand.Intn(10000))

				pubAd[r] = &slotData{
					Ctr:      .1,
					PublicID: r,
					SlotSize: s,
				}

				pubSize[r] = s
			}
		}

		return pubAd, pubSize
	}

	params := c.QueryParams()

	var size = make(map[string]string)
	var sizeNumSlice = make(map[string]int)
	var slotPublic []string

	for key := range params {
		slice := slotReg.FindStringSubmatch(key)

		//fmt.Println(slice,len(slice))
		if len(slice) == 2 {

			slotPublic = append(slotPublic, slice[1])
			size[slice[1]] = params[key][0]
			//check for size
			//size[slice[1]] = strings.Trim(size[slice[1]], "a")
			SizeNum, _ := utils.GetSize(size[slice[1]])
			sizeNumSlice[slice[1]] = SizeNum
		}
	}

	if mobile {
		slotPub := fmt.Sprintf("%d%s", website.WPubID, webMobile)
		slotPublic = append(slotPublic, slotPub)
		sizeNumSlice[slotPub] = 8
	}
	return tc.slotSizeNormal(slotPublic, website.WID, sizeNumSlice)
}

func (tc selectController) slotSizeNative(c echo.Context, website models.Website, test ...bool) (map[string]*slotData, map[string]int, []string) {
	if len(test) == 1 && test[0] {
		var pubAd = make(map[string]*slotData)
		var pubSize = make(map[string]int)

		payload := c.Get("payload").(allAdsNativePayload)

		for i := 0; i < payload.Count; i++ {
			pid := strconv.FormatInt(int64(i), 10)
			pubAd[pid] = &slotData{
				PublicID: string(i),
				ID:       int64(i),
				SlotSize: 20,
				Ctr:      .1,
			}

			pubSize[pid] = 20
		}

		return pubAd, pubSize, nil
	}
	params := c.QueryParams()
	var sizeNumSlice = make(map[string]int)
	var slotPublic []string

	count, ok := params["count"]
	if !ok {
		return make(map[string]*slotData), make(map[string]int), nil
	}
	countString := count[0]
	countInt, err := strconv.Atoi(countString)
	if err != nil || countInt < 1 {
		return make(map[string]*slotData), make(map[string]int), nil
	}
	if countInt > nativeMaxCount.Int() {
		countInt = nativeMaxCount.Int()
	}
	order := []string{}
	for i := 1; i <= countInt; i++ { //range  over slots
		pub := fmt.Sprintf("%d%s%d", website.WID, "470", i)
		order = append(order, pub)
		sizeNumSlice[pub] = 20
		slotPublic = append(slotPublic, pub)
	}

	res, sizes := tc.slotSizeNormal(slotPublic, website.WID, sizeNumSlice)
	return res, sizes, order
}

func (selectController) insertNewSlots(wID int64, newSlots []int64, newSize []int) map[string]int64 {
	assert.True(len(newSlots) == len(newSize), "[BUG] slot public and count is not matched")
	result := make(map[string]int64)
	if len(newSlots) > 0 {
		for i := range newSlots {
			insertedSlots, err := models.NewManager().InsertSlots(wID, 0, newSlots[i], newSize[i])
			if err == nil {
				p := fmt.Sprintf("%d", insertedSlots.PublicID)
				result[p] = insertedSlots.ID
			}
		}
	}

	return result
}

func (selectController) insertNewAppSlots(appID int64, newSlots []int64, newSize []int) map[string]int64 {
	assert.True(len(newSlots) == len(newSize), "[BUG] slot public and count is not matched")
	result := make(map[string]int64)
	if len(newSlots) > 0 {
		for i := range newSlots {
			insertedSlots, err := models.NewManager().InsertSlots(0, appID, newSlots[i], newSize[i])
			if err == nil {
				p := fmt.Sprintf("%d", insertedSlots.PublicID)
				result[p] = insertedSlots.ID
			}
		}
	}

	return result
}

// CalculateCtr calculate ctr
func (selectController) calculateCTR(ad *models.AdData, slot *slotData) float64 {
	return (ad.AdCTR*float64(adCTREffect.Int()) + slot.Ctr*float64(slotCTREffect.Int())) / float64(100)
}

func (tc *selectController) makeShow(
	c echo.Context,
	typ string,
	rd *middlewares.RequestData,
	filteredAds map[int][]*models.AdData,
	order []string,
	sizeNumSlice map[string]int,
	slotSize map[string]*slotData,
	attr map[string]map[string]string,
	publisher Publisher,
	multipleVideo bool,
	minCPC int64,
	allowUnderFloor bool,
	capping bool,
	floorDiv int64, // I hate add parameter to this function :/ TODO : implement the function option pattern
) (map[string]string, map[string]*models.AdData) {
	//var dum []*models.AdData
	var (
		winnerAd = make(map[string]*models.AdData)
		show     = make(map[string]string)
		noVideo  bool // once set, never unset it again
	)

	sel := true // for campaign with order (native) do not handle selected and duplicates
	if order == nil {
		sel = false
		for i := range slotSize {
			order = append(order, i)
		}
	}

	reserve := make(map[string]string)
	for o := range order {
		slotID := order[o]
		tmp := <-utils.ID
		reserve[slotID] = tmp
		u := url.URL{
			Scheme: rd.Scheme,
			Host:   rd.Host,
			Path:   fmt.Sprintf("/show/%s/%s/%d/%s", typ, rd.MegaImp, publisher.GetID(), tmp),
		}
		v := url.Values{}
		v.Set("tid", rd.TID)
		v.Set("ref", rd.Referrer)
		v.Set("parent", rd.Parent)
		v.Set("s", fmt.Sprintf("%d", slotSize[slotID].ID))

		for i, j := range slotSize[slotID].ExtraParam {
			v.Set(i, j)
		}
		u.RawQuery = v.Encode()
		show[slotID] = u.String()
	}

	var wait chan map[string]*models.AdData
	if typ == "sync" {
		wait = make(chan map[string]*models.AdData)
	}
	assert.Nil(tc.createMegaKey(rd, publisher))
	safe.GoRoutine(func() {
		ads := make(map[string]*models.AdData)
		defer func() {
			if typ == "sync" {
				wait <- ads
			}
		}()

		if capping {
			eventPage := ""
			if ep, ok := c.Get("EVENT_PAGE").(string); ok {
				eventPage = ep
			}

			filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds, eventPage)
		} else {
			filteredAds = emptyCapping(filteredAds)
		}
		// TODO : must loop over this values, from lowest data to highest. the size with less ad count must be in higher priority
		selected := make(map[int]int)
		total := make(map[int]int)

		for o := range order {
			slotID := order[o]
			exceedFloor := []*models.AdData{}
			underFloor := []*models.AdData{}

			for _, adData := range filteredAds[slotSize[slotID].SlotSize] {
				total[slotSize[slotID].SlotSize]++
				if adData.AdType == utils.AdTypeVideo && noVideo {
					continue
				}
				if adData.WinnerBid == 0 && tc.doBid(adData, publisher, slotSize[slotID], floorDiv) {
					exceedFloor = append(exceedFloor, adData)
				} else if adData.WinnerBid == 0 {
					underFloor = append(underFloor, adData)
				}
			}

			extra := fmt.Sprintf("For Slot %s", slotID)

			var sorted []*models.AdData
			var (
				ef     models.ByMulti
				secBid bool
			)

			// order is to get data from exceed flor, then capping passed and if the config allowed,
			// use the under floor. for under floor there is no second biding pricing
			if len(exceedFloor) > 0 {
				ef = models.ByMulti{
					Ads:   exceedFloor,
					Video: multipleVideo,
				}
				secBid = true
				extra += " From Exceed, SecBID "
			} else if allowUnderFloor && len(underFloor) > 0 {
				ef = models.ByMulti{
					Ads:   underFloor,
					Video: multipleVideo,
				}
				secBid = false
				extra += " From Under, FirstBID "
			}
			if len(ef.Ads) == 0 {
				w, h := utils.GetSizeByNum(slotSize[slotID].SlotSize)
				req, _ := httputil.DumpRequest(c.Request(), false)
				logrus.WithFields(logrus.Fields{
					"publisher":     publisher.GetName(),
					"width":         w,
					"height":        h,
					"floor":         publisher.FloorCPM(),
					"filtered":      len(filteredAds[slotSize[slotID].SlotSize]),
					"exceed":        len(exceedFloor),
					"allowed_under": allowUnderFloor,
					"under":         len(underFloor),
					"selected":      selected[slotSize[slotID].SlotSize],
					"total":         total[slotSize[slotID].SlotSize],
					"request":       string(req),
				}).Warn("no ad")
				ads[slotID] = nil
				store.Set(reserve[slotID], "no add")
				continue
			}

			sort.Sort(ef)
			sorted = ef.Ads

			// Do not do second biding pricing on this ads, they can not pass CPMFloor
			if secBid {
				secondCPM := tc.getSecondCPM(publisher.FloorCPM(), sorted)
				sorted[0].WinnerBid = utils.WinnerBid(secondCPM, sorted[0].CTR)
				extra += fmt.Sprintf(" WinnerCPM = %d, SecCPM = %d , CTR = %f, WinnerBID: %d", sorted[0].CPM, secondCPM, sorted[0].CTR, sorted[0].WinnerBid)
				//if len(sorted) > 1 {
				//	extra += fmt.Sprintf("%s == %d", sorted[1].AdName.String, sorted[1].CampaignID)
				//
				//	//if sorted[0].CPM == secondCPM {
				//	//	sorted[0], sorted[1] = sorted[1], sorted[0]
				//	//}
				//}
			} else {
				sorted[0].WinnerBid = sorted[0].CampaignMaxBid
				extra += fmt.Sprintf(" WinnerCPM = %d, MaxBID is requested", sorted[0].CPM)
			}

			if sorted[0].WinnerBid > sorted[0].CampaignMaxBid {
				// TODO : must not happen, but it happen some how. check it later
				sorted[0].WinnerBid = sorted[0].CampaignMaxBid
				extra += " WTF? the winner bid is greater than max bid? fixing. "
			}

			if sorted[0].WinnerBid < minCPC {
				sorted[0].WinnerBid = minCPC
				extra += " Winner bid is less than min CPC? fixing. "
			}

			extra += fmt.Sprintf("CPID = %d", sorted[0].CampaignID)
			sorted[0].Capping.IncView(sorted[0].AdID, 1, sel)
			winnerAd[slotID] = sorted[0]
			winnerAd[slotID].SlotID = slotSize[slotID].ID
			winnerAd[slotID].SlotPublicID = slotSize[slotID].PublicID
			sorted[0].Extra = extra
			ads[slotID] = sorted[0]

			if !multipleVideo {
				noVideo = noVideo || sorted[0].AdType == utils.AdTypeVideo
			}
			var clu, clp, clt string
			if sa, ok := attr[slotID]; ok {
				clu = sa["click_url"]
				clp = sa["click_parameter"]
				clt = sa["type"]
			}
			tc.updateMegaKey(rd, sorted[0].AdID, sorted[0].WinnerBid, slotSize[slotID].ID, clu, clp, clt)
			store.Set(reserve[slotID], fmt.Sprintf("%d", sorted[0].AdID))
			assert.Nil(storeCapping(rd.CopID, sorted[0].AdID))
			selected[slotSize[slotID].SlotSize]++
			// TODO {fzerorubigd} : Can we check for inner capping increase?
		}
	})
	var allAds map[string]*models.AdData
	if typ == "sync" {
		allAds = <-wait
	}
	//t, _ := json.MarshalIndent(dum, "\t", "\t")
	//fmt.Println(string(t))
	return show, allAds
}

func (tc *selectController) vastJS(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/javascript")
	_, err := c.Response().Write([]byte(js.Vastjs))
	return err
}

func (tc *selectController) videoJS(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/javascript")
	_, err := c.Response().Write([]byte(js.Videojs))
	return err
}

func (tc *selectController) nativeJS(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/javascript")
	_, err := c.Response().Write([]byte(js.Nativejs))
	return err
}

func init() {
	modules.Register(&selectController{})
}
