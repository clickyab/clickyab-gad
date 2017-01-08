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
	"redis"
	"regexp"
	"selector"
	"strconv"
	"time"
	"transport"
	"utils"

	"net/http/httputil"
	"rabbit"
	"sort"
	"store"

	"net/url"

	"strings"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

var (
	webSelector = selector.Mix(
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.IsWebNetwork,
		filter.CheckCategory,
		filter.CheckProvince,
	)

	vastSelector = selector.Mix(
		filter.CheckVastSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.IsWebNetwork,
		filter.CheckCategory,
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
	show := tc.makeShow(c, "web", rd, filteredAds, sizeNumSlice, slotSize, website, false)

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

func (tc *selectController) doBid(adData *mr.AdData, website *mr.Website, slot *slotData) bool {
	adData.CTR = tc.calculateCTR(
		adData,
		slot,
	)
	adData.CPM = utils.Cpm(adData.CampaignMaxBid, adData.CTR)
	//exceed cpm floor
	return adData.CPM >= website.WFloorCpm.Int64
}

func (tc *selectController) getSecondCPM(floorCPM int64, exceedFloor []*mr.AdData) int64 {
	var secondCPM = floorCPM
	if len(exceedFloor) > 1 && exceedFloor[0].Capping.GetSelected() == exceedFloor[1].Capping.GetSelected() {
		secondCPM = exceedFloor[1].CPM
	}

	return secondCPM
}

func (tc *selectController) createMegaKey(rd *middlewares.RequestData, website *mr.Website) error {
	tmp := map[string]string{
		"IP": rd.IP.String(),
		"UA": rd.UserAgent,
		"WS": fmt.Sprintf("%d", website.WID),
		"T":  fmt.Sprintf("%d", time.Now().Unix()),
	}
	assert.True(config.Config.Clickyab.MegaImpExpire > 1, "invalid config")
	return aredis.HMSet(
		fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, rd.MegaImp),
		config.Config.Clickyab.MegaImpExpire,
		tmp,
	)
}

func (tc *selectController) updateMegaKey(rd *middlewares.RequestData, adID int64, winnerBid int64, slotID string) {
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
		slotID,
		config.Config.Clickyab.MegaImpExpire,
	))
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
	if website.WFloorCpm.Int64 < config.Config.Clickyab.MinCPMFloor {
		website.WFloorCpm.Int64 = config.Config.Clickyab.MinCPMFloor
	}
	return website, err
}

//fetchIP2Location find ip
func (tc *selectController) fetchIP2Location(c net.IP) (*mr.IP2Location, error) {
	ip, err := mr.NewManager().GetLocation(c)
	if err != nil {
		return nil, errors.New("location not found")
	}

	return ip, nil

}

// The following code is not required anymore, but I keep it here as long as we are
// in development, remove it when we are done
//FetchCountry find country and set context
// func (tc *selectController) fetchCountry(c net.IP) (*mr.CountryInfo, error) {
// 	var country mr.CountryInfo
// 	ip, err := tc.fetchIP2Location(c)
// 	if err != nil || !ip.CountryName.Valid {
// 		return nil, errors.New("Country not found")
// 	}

// 	country, err = mr.NewManager().ConvertCountry2Info(ip.CountryCode.String)
// 	if err != nil {
// 		return nil, errors.New("Country not found")
// 	}
// 	return &country, nil
// }

//fetchProvince find province and set context
func (tc *selectController) fetchProvince(c net.IP, cfHeader string) (*mr.Province, error) {
	if strings.ToUpper(cfHeader) != "IR" {
		return nil, errors.New("not inside iran")
	}
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

func (selectController) insertNewSlots(wID int64, newSlots ...int64) map[string]int64 {
	result := make(map[string]int64)
	if len(newSlots) > 0 {
		insertedSlots, err := mr.NewManager().InsertSlots(wID, newSlots...)
		if err == nil {
			for i := range insertedSlots {
				p := fmt.Sprintf("%d", insertedSlots[i].PublicID)
				result[p] = insertedSlots[i].ID
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
	website *mr.Website,
	multipleVideo bool) map[string]string {
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
			Path:   fmt.Sprintf("/show/%s/%s/%d/%s", typ, rd.MegaImp, website.WID, tmp),
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
	assert.Nil(tc.createMegaKey(rd, website))
	middlewares.SafeGO(c, false, func() {
		filteredAds = getCapping(c, rd.CopID, sizeNumSlice, filteredAds)
		// TODO : must loop over this values, from lowest data to highest. the size with less ad count must be in higher priority
		for slotID := range slotSize {
			exceedFloor := &mr.CappingLocker{}
			cappedFloor := &mr.CappingLocker{}
			underFloor := &mr.CappingLocker{}
			for _, adData := range filteredAds[slotSize[slotID].SlotSize] {
				if adData.AdType == config.AdTypeVideo && noVideo {
					continue
				}
				if tc.doBid(adData, website, slotSize[slotID]) {
					if exceedFloor.Len() == 0 {
						exceedFloor.Set(adData.Capping.GetCapping())
					}
					//minimum capping
					if adData.Capping.GetCapping() <= exceedFloor.Get() && adData.WinnerBid == 0 {
						exceedFloor.Append(adData)
					} else if adData.WinnerBid == 0 { // the campaign is lost based on capping
						cappedFloor.Append(adData)
					}
				} else if adData.WinnerBid == 0 {
					underFloor.Append(adData)
				}
			}
			var sorted []*mr.AdData
			var (
				ef     mr.ByCPM
				secBid bool
			)

			// order is to get data from exceed flor, then capping passed and if the config allowed,
			// use the under floor. for under floor there is no second biding pricing
			if exceedFloor.Len() > 1 {
				ef = mr.ByCPM(exceedFloor.GetData())
				secBid = true
			} else if cappedFloor.Len() > 1 {
				ef = mr.ByCPM(cappedFloor.GetData())
				secBid = true
			} else if config.Config.Clickyab.UnderFloor && underFloor.Len() > 0 {
				ef = mr.ByCPM(underFloor.GetData())
				secBid = false
			}

			if len(ef) == 0 {
				middlewares.SafeGO(c, false, func() {
					w, h := config.GetSizeByNum(slotSize[slotID].SlotSize)
					warn := transport.Warning{
						Level: "warning",
						When:  time.Now(),
						Where: website.WDomain.String,
						Message: fmt.Sprintf(
							"no ad pass the bid \n "+
								"size was %sx%s \n"+
								"the floor was %d \n"+
								"all add count in this size %d \n "+
								"under floor is allowd? %v \n"+
								"under floor count %d",
							w, h,
							website.WFloorCpm.Int64,
							len(filteredAds[slotSize[slotID].SlotSize]),
							config.Config.Clickyab.UnderFloor,
							underFloor.Len(),
						),
					}
					warn.Request, _ = httputil.DumpRequest(c.Request(), false)
					err := rabbit.Publish("cy.warn", warn)
					if err != nil {
						logrus.Error(err)
					}
				})
				show[slotID] = ""
				store.Set(reserve[slotID], "no add")
				continue
			}

			sort.Sort(ef)
			sorted = []*mr.AdData(ef)
			// Do not do second biding pricing on this ads, they can not pass CPMFloor
			if secBid {
				secondCPM := tc.getSecondCPM(website.WFloorCpm.Int64, sorted)
				sorted[0].WinnerBid = utils.WinnerBid(secondCPM, sorted[0].CTR)
			} else {
				sorted[0].WinnerBid = sorted[0].CampaignMaxBid
			}
			sorted[0].Capping.IncView(1)
			winnerAd[slotID] = sorted[0]
			winnerAd[slotID].SlotID = slotSize[slotID].ID

			if !multipleVideo {
				noVideo = noVideo || sorted[0].AdType == config.AdTypeVideo
			}
			tc.updateMegaKey(rd, sorted[0].AdID, sorted[0].WinnerBid, slotID)
			store.Set(reserve[slotID], fmt.Sprintf("%d", sorted[0].AdID))
			//show[slotID] = fmt.Sprintf("%s://%s/show/%s/%s/%d/%d?tid=%s&ref=%s&s=%d", rd.Proto, rd.URL, typ, rd.MegaImp, website.WID, sorted[0].AdID, rd.TID, rd.Parent, slotSize[slotID].ID)
			assert.Nil(storeCapping(rd.CopID, sorted[0].CampaignID))
			// TODO {fzerorubigd} : Can we check for inner capping increase?

		}
	})
	return show
}

func init() {
	modules.Register(&selectController{})
}
