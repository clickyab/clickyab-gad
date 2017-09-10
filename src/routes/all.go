package routes

import (
	"config"
	"filter"
	"middlewares"
	"mr"
	"net"
	"net/http"
	"selector"
	"sort"
	"strconv"

	echo "gopkg.in/labstack/echo.v3"

	"assert"
	"encoding/json"
	"fmt"
	"ip2location"
)

// AllData return all data required to render the all routes
// TODO : Rename this
type AllData struct {
	Website  []*mr.Website
	Province []*mr.Province
	//Campaign *[]mr.Campaign
	Size map[string]int
	Vast bool
	Data []*mr.AdData
	Len  int
}

var allFiter = map[string]selector.FilterFunc{
	"isWebNetwork":  filter.IsWebNetwork,
	"webSize":       filter.CheckWebSize,
	"appSize":       filter.CheckAppSize,
	"vastSize":      filter.CheckVastSize,
	"os":            filter.CheckOS,
	"whiteList":     filter.CheckWhiteList,
	"blackList":     filter.CheckWebBlackList,
	"webCategory":   filter.CheckWebCategory,
	"checkProvince": filter.CheckProvince,
	"isWebMobile":   filter.IsWebMobile,
	"notWebMobile":  filter.IsNotWebMobile,
	"checkCampaign": filter.CheckCampaign,
	"webMobileSize": filter.CheckWebMobileSize,
	"appBlackList":  filter.CheckAppBlackList,
	"appWhiteList":  filter.CheckAppWhiteList,
	"appCategory":   filter.CheckAppCategory,
	"appBrand":      filter.CheckAppBrand,
	"appHood":       filter.CheckAppHood,
	"appProvider":   filter.CheckProvder,
	"appAreaInGlob": filter.CheckAppAreaInGlob,
}

// Ints returns a unique subset of the int slice provided.
func UniqueStr(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

type allAdsWebPayload struct {
	TID   int64  `json:"tid"`
	IP    net.IP `json:"ip"`
	Slots []struct {
		Size  string `json:"size"`
		Count int    `json:"count"`
	} `json:"slots"`
}

type allAdsResponse struct {
	CTR    float64 `json:"ctr"`
	AdID   int64   `json:"ad_id"`
	CampID int64   `json:"camp_id"`
	AdImg  string  `json:"ad_img"`
	AdType int     `json:"ad_type"`
}

type allAdsNativePayload struct {
	TID   int64  `json:"tid"`
	IP    net.IP `json:"ip"`
	Count int    `json:"ad_count"`
}

type allAdsVastPayload struct {
	IP net.IP `json:"ip"`

	Start string `json:"start"`
	Mid   string `json:"mid"`
	End   string `json:"end"`
}

func (tc *selectController) allAds(c echo.Context) error {
	adType := c.QueryParams().Get("type")

	rd := middlewares.MustGetRequestData(c)

	wid := c.QueryParam("i")
	wpid, _ := strconv.ParseInt(wid, 10, 0)
	website, err := mr.NewManager().FetchWebsite(wpid)
	assert.Nil(err)

	switch adType {
	case "web":
		err = tc.allWebAds(c, rd, website)
	case "native":
		err = tc.allNativeAds(c, rd, website)
	case "vast":
		err = tc.allVastAds(c, rd, website)
	}

	return err
}

func (tc *selectController) allWebAds(c echo.Context, rd *middlewares.RequestData, website *mr.Website) error {
	payload := allAdsWebPayload{}
	dec := json.NewDecoder(c.Request().Body)
	defer c.Request().Body.Close()
	err := dec.Decode(&payload)
	assert.Nil(err)

	c.Set("payload", payload)

	rd.CopID = payload.TID
	slotSize, sizeNumSlice := tc.slotSizeWeb(c, *website, rd.Mobile, true)

	provinceID, _ := ip2location.GetProvinceISPByIP(payload.IP)

	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    provinceID,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	winnerAds := tc.webBiding(rd, filteredAds, slotSize, sizeNumSlice)

	// map[size]number
	var sizeNum = make(map[int]int)
	for i := range sizeNumSlice {
		sizeNum[sizeNumSlice[i]]++
	}

	for i := range sizeNum {
		winnerAds[i] = winnerAds[i][:sizeNum[i]]
	}

	//filling ads url
	resp := map[string][]allAdsResponse{}
	for i := range winnerAds {
		for _, j := range winnerAds[i] {
			// not sure bout the true part
			ad, err := mr.NewManager().GetAd(j.AdID, true)
			assert.Nil(err)

			temp := allAdsResponse{
				AdImg:  ad.AdImg.String,
				AdID:   j.AdID,
				AdType: j.AdType,
				CampID: j.CampaignAdID,
				CTR:    j.CTR,
			}

			assert.Nil(storeCapping(rd.CopID, j.AdID))
			resp[config.GetSizeByNumString(i)] = append(resp[config.GetSizeByNumString(i)], temp)
		}
	}

	println(fmt.Sprintf("%v", resp))

	return c.JSON(200, resp)
}

func (tc *selectController) allNativeAds(ctx echo.Context, rd *middlewares.RequestData, website *mr.Website) error {
	payload := allAdsNativePayload{}
	dec := json.NewDecoder(ctx.Request().Body)
	err := dec.Decode(&payload)
	ctx.Request().Body.Close()
	assert.Nil(err)
	ctx.Set("payload", payload)

	rd.CopID = payload.TID
	provinceID, _ := ip2location.GetProvinceISPByIP(payload.IP)

	slotSize, sizeNumSlice, _ := tc.slotSizeNative(ctx, *website, true)
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    provinceID,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), nativeSelector)
	winnerAds := tc.nativeBiding(rd, filteredAds, slotSize, sizeNumSlice)

	resp := []allAdsResponse{}
	for _, i := range winnerAds[20][:len(sizeNumSlice)] {
		// not sure bout the true part
		ad, err := mr.NewManager().GetAd(i.AdID, true)
		assert.Nil(err)

		resp = append(resp, allAdsResponse{
			AdID:   i.AdID,
			CTR:    i.CTR,
			AdType: i.AdType,
			CampID: i.CampaignAdID,
			AdImg:  ad.AdImg.String,
		})
		assert.Nil(storeCapping(rd.CopID, i.AdID))
	}

	return ctx.JSON(200, resp)
}

func (tc *selectController) allVastAds(ctx echo.Context, rd *middlewares.RequestData, website *mr.Website) error {
	payload := allAdsVastPayload{}
	dec := json.NewDecoder(ctx.Request().Body)
	err := dec.Decode(&payload)
	ctx.Request().Body.Close()
	assert.Nil(err)
	ctx.Set("payload", payload)

	provinceID, _ := ip2location.GetProvinceISPByIP(payload.IP)

	lenVast, vastCon := config.MakeVastLen(ctx.QueryParam("l"), payload.Start, payload.Mid, payload.End)
	vastSlot, pubs, pubSize := makeVastSlot(vastCon, website)
	SlotData, sizeNumSlice := tc.slotSizeNormal(pubs, website.WID, pubSize, true)

	// setting extra params
	for i := range vastSlot {
		SlotData[i].ExtraParam = map[string]string{
			"pos":  vastSlot[i].Offset,
			"type": vastSlot[i].Type,
			"l":    lenVast,
		}
	}

	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    provinceID,
	}

	filteredAds := selector.Apply(&m, selector.GetAdData(), vastSelector)
	_, allAds := tc.makeShow(ctx, "sync", rd, filteredAds, nil, sizeNumSlice, SlotData, nil, website, false, config.Config.Clickyab.MinCPCVast, config.Config.Clickyab.UnderFloor, true, config.Config.Clickyab.FloorDiv.Vast)

	response := map[string][]allAdsResponse{}
	for i := range allAds {
		response[config.GetSizeByNumString(allAds[i].AdSize)] = append(response[config.GetSizeByNumString(allAds[i].AdSize)], allAdsResponse{
			CTR:    .1,
			AdImg:  allAds[i].AdImg.String,
			AdType: allAds[i].AdType,
			CampID: allAds[i].CampaignAdID,
			AdID:   allAds[i].AdID,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (tc *selectController) webBiding(rd *middlewares.RequestData, filteredAds map[int][]*mr.AdData, slotSize map[string]*slotData, sizeNumSlice map[string]int) map[int][]*mr.AdData {
	filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds, "")

	for i := range filteredAds {
		Ads := mr.ByMulti{Ads: filteredAds[i]}
		sort.Sort(Ads)

		filteredAds[i] = Ads.Ads
	}

	return filteredAds
}

func (tc *selectController) nativeBiding(rd *middlewares.RequestData, filteredAds map[int][]*mr.AdData, slotSize map[string]*slotData, sizeNumSlice map[string]int) map[int][]*mr.AdData {
	filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds, "")

	for i := range filteredAds {
		Ads := mr.ByMulti{Ads: filteredAds[i]}
		sort.Sort(Ads)

		filteredAds[i] = Ads.Ads
	}

	return filteredAds
}

func makeVastSlot(length map[string][]string, website *mr.Website) (map[string]vastSlotData, []string, map[string]int) {
	var i int
	var sizeNumSlice = make(map[string]int)
	var slotPublic []string
	var vastSlot = make(map[string]vastSlotData)
	for m := range length {
		i++
		lenType := length[m][0]
		if lenType != "linear" {
			continue
		}
		pub := fmt.Sprintf("%d%s", website.WPubID, length[m][1])
		sizeNumSlice[pub] = config.VastNonLinearSize
		if lenType == "linear" {
			sizeNumSlice[pub] = config.VastLinearSize
		}
		slotPublic = append(slotPublic, pub)
		vastSlot[pub] = vastSlotData{
			Offset: m,
			Repeat: length[m][2],
			Type:   lenType,
		}
	}

	return vastSlot, slotPublic, sizeNumSlice
}

func allDate() AllData {
	/*c, err := mr.NewManager().FetchCampaignAll()
	if err != nil {
		c = nil
	}*/
	p, err := mr.NewManager().FetchProvinceAll()
	if err != nil {
		p = nil
	}
	w, err := mr.NewManager().FetchWebsiteAll()
	if err != nil {
		w = nil
	}
	s := config.GetAllSize()
	al := AllData{
		//Campaign: c,
		Province: p,
		Website:  w,
		Size:     s,
	}
	return al
}
