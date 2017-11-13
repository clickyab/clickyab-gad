package routes

import (
	"net"
	"net/http"
	"sort"

	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
	"github.com/clickyab/services/assert"

	echo "gopkg.in/labstack/echo.v3"

	"encoding/json"
	"fmt"

	"clickyab.com/gad/ip2location"
	"clickyab.com/gad/utils"
)

// AllData return all data required to render the all routes
// TODO : Rename this
type AllData struct {
	Website  []*models.Website
	Province []*models.Province
	//Campaign *[]models.Campaign
	Size map[string]int
	Vast bool
	Data []*models.AdData
	Len  int
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
	AdID       int64   `json:"ad_id"`
	AdImg      string  `json:"ad_img"`
	AdType     int     `json:"ad_type"`
	CampID     int64   `json:"camp_id"`
	CampName   string  `json:"camp_name"`
	CampMail   string  `json:"camp_mail"`
	CampBudget int     `json:"camp_budget"`
	CTR        float64 `json:"ctr"`
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
	website, err := models.NewManager().FindWebsiteByDomain(wid)
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

func (tc *selectController) allWebAds(c echo.Context, rd *middlewares.RequestData, website *models.Website) error {
	payload := allAdsWebPayload{}
	dec := json.NewDecoder(c.Request().Body)
	defer c.Request().Body.Close()
	err := dec.Decode(&payload)
	assert.Nil(err)

	c.Set("payload", payload)

	rd.CopID = payload.TID
	slotSize, sizeNumSlice := tc.slotSizeWeb(c, *website, rd.Mobile, true)

	provinceID, _, _ := ip2location.GetProvinceISPByIP(payload.IP)

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
			temp := allAdsResponse{
				AdImg:      j.AdImg.String,
				AdID:       j.AdID,
				AdType:     j.AdType,
				CampID:     j.CampaignID,
				CTR:        j.CTR,
				CampName:   j.CampaignName.String,
				CampMail:   j.UserEmail,
				CampBudget: j.CampaignTotalBudget,
			}

			assert.Nil(storeCapping(rd.CopID, j.AdID))
			resp[utils.GetSizeByNumString(i)] = append(resp[utils.GetSizeByNumString(i)], temp)
		}
	}

	println(fmt.Sprintf("%v", resp))

	return c.JSON(http.StatusOK, resp)
}

func (tc *selectController) allNativeAds(ctx echo.Context, rd *middlewares.RequestData, website *models.Website) error {
	payload := allAdsNativePayload{}
	dec := json.NewDecoder(ctx.Request().Body)
	err := dec.Decode(&payload)
	ctx.Request().Body.Close()
	assert.Nil(err)
	ctx.Set("payload", payload)

	rd.CopID = payload.TID
	provinceID, _, _ := ip2location.GetProvinceISPByIP(payload.IP)

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

		resp = append(resp, allAdsResponse{
			AdID:       i.AdID,
			AdType:     i.AdType,
			AdImg:      i.AdImg.String,
			CampBudget: i.Campaign.CampaignTotalBudget,
			CampMail:   i.UserEmail,
			CampName:   i.CampaignName.String,
			CampID:     i.CampaignAdID,
			CTR:        i.CTR,
		})
		assert.Nil(storeCapping(rd.CopID, i.AdID))
	}

	return ctx.JSON(200, resp)
}

func (tc *selectController) allVastAds(ctx echo.Context, rd *middlewares.RequestData, website *models.Website) error {
	payload := allAdsVastPayload{}
	dec := json.NewDecoder(ctx.Request().Body)
	err := dec.Decode(&payload)
	ctx.Request().Body.Close()
	assert.Nil(err)
	ctx.Set("payload", payload)

	provinceID, _, _ := ip2location.GetProvinceISPByIP(payload.IP)

	lenVast, vastCon := utils.MakeVastLen(ctx.QueryParam("l"), payload.Start, payload.Mid, payload.End)
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
	_, allAds := tc.makeShow(ctx, "sync", rd, filteredAds, nil, sizeNumSlice, SlotData, nil, website, false, minCPCVast.Int64(), allowUnderFloor.Bool(), true, floorDivVast.Int64())

	response := map[string][]allAdsResponse{}
	for i := range allAds {
		response[utils.GetSizeByNumString(allAds[i].AdSize)] = append(response[utils.GetSizeByNumString(allAds[i].AdSize)], allAdsResponse{
			CTR:        .1,
			AdImg:      allAds[i].AdImg.String,
			AdType:     allAds[i].AdType,
			CampID:     allAds[i].CampaignAdID,
			AdID:       allAds[i].AdID,
			CampName:   allAds[i].CampaignName.String,
			CampMail:   allAds[i].UserEmail,
			CampBudget: allAds[i].CampaignTotalBudget,
		})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (tc *selectController) webBiding(rd *middlewares.RequestData, filteredAds map[int][]*models.AdData, slotSize map[string]*slotData, sizeNumSlice map[string]int) map[int][]*models.AdData {
	filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds, "")

	for i := range filteredAds {
		Ads := models.ByMulti{Ads: filteredAds[i]}
		sort.Sort(Ads)

		filteredAds[i] = Ads.Ads
	}

	return filteredAds
}

func (tc *selectController) nativeBiding(rd *middlewares.RequestData, filteredAds map[int][]*models.AdData, slotSize map[string]*slotData, sizeNumSlice map[string]int) map[int][]*models.AdData {
	filteredAds = getCapping(rd.CopID, sizeNumSlice, filteredAds, "")

	for i := range filteredAds {
		Ads := models.ByMulti{Ads: filteredAds[i]}
		sort.Sort(Ads)

		filteredAds[i] = Ads.Ads
	}

	return filteredAds
}

func makeVastSlot(length map[string][]string, website *models.Website) (map[string]vastSlotData, []string, map[string]int) {
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
		sizeNumSlice[pub] = utils.VastNonLinearSize
		if lenType == "linear" {
			sizeNumSlice[pub] = utils.VastLinearSize
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

func (tc *selectController) allAdsTemp(c echo.Context) error {
	return c.HTML(http.StatusOK, allAddTemplate)
}
