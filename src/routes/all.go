package routes

import (
	"config"
	"filter"
	"middlewares"
	"mr"
	"net"
	"net/url"
	"selector"
	"sort"
	"strconv"

	echo "gopkg.in/labstack/echo.v3"

	"assert"
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

type allAdsPayload struct {
	TID   int64  `json:"tid"`
	IP    net.IP `json:"ip"`
	Slots []struct {
		Size   string `json:"size"`
		Repeat int    `json:"repeat"`
	} `json:"slots"`
}

func (tc *selectController) allAds(c echo.Context) error {
	params := c.QueryParams()
	adType := params.Get("type")

	var err error
	switch adType {
	case "web":
		err = tc.allWebAds(c)
	case "app":
		err = tc.allAppAds(params)
	case "vast":
		err = tc.allVastAds(params)
	}

	return err
}

func (tc *selectController) allWebAds(c echo.Context) error {
	var payload allAdsPayload
	err := c.Bind(&payload)
	assert.Nil(err)

	rd := middlewares.MustGetRequestData(c)
	rd.CopID = payload.TID
	wid := c.QueryParam("i")
	wpid, _ := strconv.ParseInt(wid, 10, 0)
	website, _ := mr.NewManager().FetchWebsite(wpid)

	province, err := tc.fetchProvince(rd.IP, c.Request().Header.Get("Cf-Ipcountry"))
	assert.Nil(err)

	slotSize, sizeNumSlice := tc.slotSizeWeb(c, *website, rd.Mobile, true)
	caf := c.QueryParam("capping")
	capping, _ := strconv.ParseBool(caf)

	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	winnerAds := tc.webBiding(rd, filteredAds, slotSize, sizeNumSlice, capping)

	// map[size]number
	var sizeNum = make(map[int]int)
	for i := range sizeNumSlice {
		sizeNum[sizeNumSlice[i]]++
	}

	for i := range sizeNum {
		winnerAds[i] = winnerAds[i][:sizeNum[i]]
	}

	return c.JSON(200, winnerAds)
}

func (tc *selectController) webBiding(rd *middlewares.RequestData, filteredAds map[int][]*mr.AdData, slotSize map[string]*slotData, sizeNumSlice map[string]int, capping bool) map[int][]*mr.AdData {
	if capping {
		filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds)
	} else {
		filteredAds = emptyCapping(filteredAds)
	}

	for i := range filteredAds {
		Ads := mr.ByMulti(filteredAds[i])
		sort.Sort(Ads)

		filteredAds[i] = []*mr.AdData(Ads)
	}

	return filteredAds
}

func (tc *selectController) allVastAds(p url.Values) error {
	return nil
}

func (tc *selectController) allAppAds(p url.Values) error {
	return nil
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
