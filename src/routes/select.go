package routes

import (
	"assert"
	"config"
	"encoding/json"
	"errors"
	"filter"
	"fmt"
	"middlewares"
	"modules"
	"mr"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rabbit"
	aredis "redis"
	"regexp"
	"selector"
	"sort"
	"store"
	"strconv"
	"time"
	"transport"
	"utils"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

var (
	webSelector = selector.Mix(
		filter.IsWebNetwork,
		filter.IsWebMobile,
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.CheckWebCategory,
		filter.CheckProvince,
		filter.CheckVastOtherNetwork,
	)

	nativeSelector = selector.Mix(
		filter.IsNativeNetwork,
		filter.IsNativeAd,
		filter.IsWebMobile,
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.CheckWebCategory,
		filter.CheckProvince,
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
	params := c.QueryParams()
	rd, website, province, err := tc.getWebDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice := tc.slotSizeWeb(params, *website, rd.Mobile)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	show, _ := tc.makeShow(c, "web", rd, filteredAds, sizeNumSlice, slotSize, nil, website, false, config.Config.Clickyab.MinCPCWeb, config.Config.Clickyab.UnderFloor, true)

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

func (tc *selectController) doBid(adData *mr.AdData, website Publisher, slot *slotData) bool {
	adData.CTR = tc.calculateCTR(
		adData,
		slot,
	)
	adData.CPM = utils.Cpm(adData.CampaignMaxBid, adData.CTR)
	//exceed cpm floor
	return adData.CPM >= website.FloorCPM()
}

func (tc *selectController) getSecondCPM(floorCPM int64, exceedFloor []*mr.AdData) int64 {
	var secondCPM = floorCPM
	if len(exceedFloor) > 1 && exceedFloor[0].Capping.GetSelected() == exceedFloor[1].Capping.GetSelected() {
		secondCPM = exceedFloor[1].CPM
	}

	return secondCPM
}

func (tc *selectController) createMegaKey(rd *middlewares.RequestData, website Publisher) error {
	tmp := map[string]string{
		"IP": rd.IP.String(),
		"UA": rd.UserAgent,
		"WS": fmt.Sprintf("%d", website.GetID()),
		"T":  fmt.Sprintf("%d", time.Now().Unix()),
	}
	assert.True(config.Config.Clickyab.MegaImpExpire > 1, "invalid config")
	return aredis.HMSet(
		fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp),
		config.Config.Clickyab.MegaImpExpire,
		tmp,
	)
}

func (tc *selectController) updateMegaKey(rd *middlewares.RequestData, adID int64, winnerBid int64, slotID int64, clickURL, clickParam string) {
	assert.Nil(aredis.StoreHashKey(
		fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp),
		fmt.Sprintf(
			"%s%s%d",
			transport.ADVERTISE,
			transport.DELIMITER,
			adID),
		fmt.Sprintf("%d", winnerBid),
		config.Config.Clickyab.MegaImpExpire,
	))
	assert.Nil(aredis.StoreHashKey(
		fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp),
		fmt.Sprintf(
			"%s%s%d",
			transport.SLOT,
			transport.DELIMITER,
			adID),
		strconv.FormatInt(slotID, 10),
		config.Config.Clickyab.MegaImpExpire,
	))
	if clickURL != "" && clickParam != "" {
		assert.Nil(aredis.StoreHashKey(
			fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp),
			fmt.Sprintf(
				"%s%s%d",
				transport.CUSTOM_CLICK_URL,
				transport.DELIMITER,
				slotID),
			clickURL,
			config.Config.Clickyab.MegaImpExpire,
		))

		assert.Nil(aredis.StoreHashKey(
			fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp),
			fmt.Sprintf(
				"%s%s%d",
				transport.CUSTOM_CLICK_PARAM,
				transport.DELIMITER,
				slotID),
			clickParam,
			config.Config.Clickyab.MegaImpExpire,
		))
	}

}

func (tc *selectController) getWebDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.Website, *mr.Province, error) {
	rd := middlewares.MustGetRequestData(c)
	params := c.QueryParams()
	publicParams, ok := params["i"]
	if !ok {
		return nil, nil, nil, errors.New("invalid request")
	}
	publicID, err := strconv.ParseInt(publicParams[0], 10, 0)
	if err != nil {
		return nil, nil, nil, errors.New("invalid request")
	}
	domain, ok := params["d"]
	if !ok {
		return nil, nil, nil, errors.New("invalid request")
	}
	//fetch website and set in Context
	website, err := tc.fetchWebsite(publicID)
	if err != nil {
		return nil, nil, nil, errors.New("invalid request")
	}

	if !website.GetActive() {
		return nil, nil, nil, errors.New("web is not active")
	}

	if !mr.NewManager().IsUserActive(website.UserID) {
		return nil, nil, nil, errors.New("user is banned")
	}

	province, err := tc.fetchProvince(rd.IP, c.Request().Header.Get("Cf-Ipcountry"))
	if err != nil {
		logrus.Debug(err)
	}
	//check if the website domain is valid
	if website.WDomain.Valid && website.WDomain.String != domain[0] {
		return nil, nil, nil, errors.New("domain and public id mismatch")
	}

	return rd, website, province, nil
}

//FetchWebsite website and check if the minimum floor is applied
func (tc *selectController) fetchWebsite(publicID int64) (*mr.Website, error) {
	website, err := mr.NewManager().FetchWebsiteByPublicID(publicID)
	if err != nil {
		return nil, err
	}
	if website.WFloorCpm.Int64 < config.Config.Clickyab.MinCPMFloorWeb {
		website.WFloorCpm.Int64 = config.Config.Clickyab.MinCPMFloorWeb
	}
	return website, err
}

//fetchIP2Location find ip
func (tc *selectController) fetchIP2Location(c net.IP) (*mr.IP2Location, error) {
	if config.Config.DevelMode {
		// change the local ip to tehran ip
		if c.String() == net.IPv4(172, 17, 0, 1).String() {
			c = net.IPv4(5, 116, 150, 199) // An Irancell IP in iran
		}
	}
	ip, err := mr.NewManager().GetLocation(c)
	if err != nil {
		return nil, errors.New("location not found")
	}

	return ip, nil

}

//fetchProvince find province and set context
func (tc *selectController) fetchProvince(c net.IP, cfHeader string) (*mr.Province, error) {
	// if strings.ToUpper(cfHeader) != "IR" {
	// 	return nil, errors.New("not inside iran")
	// }
	var province mr.Province
	ip, err := tc.fetchIP2Location(c)
	if err != nil || !ip.RegionName.Valid {
		return nil, errors.New("province not found")
	}

	province, err = mr.NewManager().ConvertProvince2Info(ip.RegionName.String)
	if err != nil {
		return nil, errors.New("province not found")
	}
	return &province, nil

}

//fetchProvince find province and set context
func (tc *selectController) fetchProvinceDemand(r string) (*mr.Province, error) {
	// if strings.ToUpper(cfHeader) != "IR" {
	// 	return nil, errors.New("not inside iran")
	// }
	var province mr.Province
	province, err := mr.NewManager().ConvertProvince2Info(r)
	if err != nil {
		return nil, errors.New("province not found")
	}
	return &province, nil

}

func (tc selectController) slotSizeWeb(params map[string][]string, website mr.Website, mobile bool) (map[string]*slotData, map[string]int) {
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
			SizeNum, _ := config.GetSize(size[slice[1]])
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

func (tc selectController) slotSizeNative(params map[string][]string, website mr.Website) (map[string]*slotData, map[string]int) {
	var sizeNumSlice = make(map[string]int)
	var slotPublic []string

	count, ok := params["count"]
	if !ok {
		return make(map[string]*slotData), make(map[string]int)
	}
	countString := count[0]
	countInt, err := strconv.Atoi(countString)
	if err != nil || countInt < 1 {
		return make(map[string]*slotData), make(map[string]int)
	}
	if countInt > 10 {
		countInt = 10
	}

	for i := 1; i <= countInt; i++ { //range  over slots
		pub := fmt.Sprintf("%d%s%d", website.WID, "470", i)
		sizeNumSlice[pub] = 20
		slotPublic = append(slotPublic, pub)
	}

	return tc.slotSizeNormal(slotPublic, website.WID, sizeNumSlice)
}

func (selectController) insertNewSlots(wID int64, newSlots []int64, newSize []int) map[string]int64 {
	assert.True(len(newSlots) == len(newSize), "[BUG] slot public and count is not matched")
	result := make(map[string]int64)
	if len(newSlots) > 0 {
		for i := range newSlots {
			insertedSlots, err := mr.NewManager().InsertSlots(wID, 0, newSlots[i], newSize[i])
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
			insertedSlots, err := mr.NewManager().InsertSlots(0, appID, newSlots[i], newSize[i])
			if err == nil {
				p := fmt.Sprintf("%d", insertedSlots.PublicID)
				result[p] = insertedSlots.ID
			}
		}
	}

	return result
}

// CalculateCtr calculate ctr
func (selectController) calculateCTR(ad *mr.AdData, slot *slotData) float64 {
	//fmt.Println(ad.AdCTR*float64(config.Config.Clickyab.AdCTREffect),slot.Ctr*float64(config.Config.Clickyab.SlotCTREffect),(ad.AdCTR*float64(config.Config.Clickyab.AdCTREffect) + slot.Ctr*float64(config.Config.Clickyab.SlotCTREffect)) / float64(100))
	return (ad.AdCTR*float64(config.Config.Clickyab.AdCTREffect) + slot.Ctr*float64(config.Config.Clickyab.SlotCTREffect)) / float64(100)
}

func (tc *selectController) makeShow(
	c echo.Context,
	typ string,
	rd *middlewares.RequestData,
	filteredAds map[int][]*mr.AdData,
	sizeNumSlice map[string]int,
	slotSize map[string]*slotData,
	attr map[string]map[string]string,
	publisher Publisher,
	multipleVideo bool,
	minCPC int64,
	allowUnderFloor bool,
	capping bool) (map[string]string, map[string]*mr.AdData) {
	var (
		winnerAd = make(map[string]*mr.AdData)
		show     = make(map[string]string)
		noVideo  bool // once set, never unset it again
	)
	reserve := make(map[string]string)
	for slotID := range slotSize {
		tmp := config.Config.MachineName + <-utils.ID
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

	var wait chan map[string]*mr.AdData
	if typ == "sync" {
		wait = make(chan map[string]*mr.AdData)
	}
	assert.Nil(tc.createMegaKey(rd, publisher))
	middlewares.SafeGO(c, false, false, func() {
		ads := make(map[string]*mr.AdData)
		defer func() {
			if typ == "sync" {
				wait <- ads
			}
		}()

		if capping {
			filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds)
		} else {
			filteredAds = emptyCapping(filteredAds)
		}
		// TODO : must loop over this values, from lowest data to highest. the size with less ad count must be in higher priority
		selected := make(map[int]int)
		total := make(map[int]int)
		for slotID := range slotSize {
			exceedFloor := []*mr.AdData{}
			underFloor := []*mr.AdData{}
			for _, adData := range filteredAds[slotSize[slotID].SlotSize] {
				total[slotSize[slotID].SlotSize]++
				if adData.AdType == config.AdTypeVideo && noVideo {
					continue
				}
				if adData.WinnerBid == 0 && tc.doBid(adData, publisher, slotSize[slotID]) {
					exceedFloor = append(exceedFloor, adData)
					logrus.Debug("append to exceed => ", adData.AdName)
				} else if adData.WinnerBid == 0 {
					underFloor = append(underFloor, adData)
					logrus.Debug("append to under => ", adData.AdName)
				} else {
					logrus.Debug("already selected => ", adData.AdName)
				}
			}

			var sorted []*mr.AdData
			var (
				ef     mr.ByCapping
				secBid bool
			)

			// order is to get data from exceed flor, then capping passed and if the config allowed,
			// use the under floor. for under floor there is no second biding pricing
			if len(exceedFloor) > 0 {
				ef = mr.ByCapping(exceedFloor)
				secBid = true
			} else if allowUnderFloor && len(underFloor) > 0 {
				ef = mr.ByCapping(underFloor)
				secBid = false
			}
			if len(ef) == 0 {
				logrus.Debug("No ad????")
				middlewares.SafeGO(c, false, false, func() {
					w, h := config.GetSizeByNum(slotSize[slotID].SlotSize)
					warn := transport.Warning{
						Level: "warning",
						When:  time.Now(),
						Where: publisher.GetName(),
						Message: fmt.Sprintf(
							"no ad pass the bid \n"+"size was %sx%s \n"+"the floor was %d \n"+"all add count in this size %d \n"+"pass the floor %d \n"+"under floor is allowd? %v \n"+"under floor count %d \n"+"currently %d item of %d in this request is filled",
							w, h,
							publisher.FloorCPM(),
							len(filteredAds[slotSize[slotID].SlotSize]),
							len(exceedFloor),
							allowUnderFloor,
							len(underFloor),
							selected[slotSize[slotID].SlotSize], total[slotSize[slotID].SlotSize],
						),
					}
					warn.Request, _ = httputil.DumpRequest(c.Request(), false)
					err := rabbit.Publish(warn)
					if err != nil {
						logrus.Error(err)
					}
				})
				ads[slotID] = nil
				store.Set(reserve[slotID], "no add")
				continue
			}

			sort.Sort(ef)
			sorted = []*mr.AdData(ef)
			// Do not do second biding pricing on this ads, they can not pass CPMFloor
			if secBid {
				secondCPM := tc.getSecondCPM(publisher.FloorCPM(), sorted)
				sorted[0].WinnerBid = utils.WinnerBid(secondCPM, sorted[0].CTR)
			} else {
				sorted[0].WinnerBid = sorted[0].CampaignMaxBid
			}

			if sorted[0].WinnerBid > sorted[0].CampaignMaxBid {
				// TODO : must not happen, but it happen some how. check it later
				sorted[0].WinnerBid = sorted[0].CampaignMaxBid
			}

			if sorted[0].WinnerBid < minCPC {
				sorted[0].WinnerBid = minCPC
			}

			sorted[0].Capping.IncView(sorted[0].AdID, 1, true)
			winnerAd[slotID] = sorted[0]
			winnerAd[slotID].SlotID = slotSize[slotID].ID
			winnerAd[slotID].SlotPublicID = slotSize[slotID].PublicID
			ads[slotID] = sorted[0]

			if !multipleVideo {
				noVideo = noVideo || sorted[0].AdType == config.AdTypeVideo
			}
			var clu, clp string
			if sa, ok := attr[slotID]; ok {
				clu = sa["click_url"]
				clp = sa["click_parameter"]
			}
			tc.updateMegaKey(rd, sorted[0].AdID, sorted[0].WinnerBid, slotSize[slotID].ID, clu, clp)
			store.Set(reserve[slotID], fmt.Sprintf("%d", sorted[0].AdID))
			assert.Nil(storeCapping(rd.CopID, sorted[0].AdID))
			selected[slotSize[slotID].SlotSize]++
			// TODO {fzerorubigd} : Can we check for inner capping increase?
		}
	})
	var allAds map[string]*mr.AdData
	if typ == "sync" {
		allAds = <-wait
	}
	return show, allAds
}

func init() {
	modules.Register(&selectController{})
}
