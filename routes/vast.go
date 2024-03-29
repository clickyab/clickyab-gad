package routes

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"clickyab.com/gad/filter"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/pin"
	"clickyab.com/gad/models/selector"
	"clickyab.com/gad/redis"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"

	"gopkg.in/labstack/echo.v3"

	"net/url"

	"clickyab.com/gad/ip2location"
	"clickyab.com/gad/store"
	"github.com/sirupsen/logrus"
)

var (
	vastSelector = selector.Mix(
		filter.IsWebNetwork,
		filter.IsWebMobile,
		filter.CheckDesktopNetwork,
		filter.CheckVastSize,
		filter.RemoveSlotPins,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.CheckWebCategory,
		filter.CheckProvince,
		filter.CheckMinBid,
		filter.CheckVastNetwork,
		filter.CheckISP,
	)
)

type vastAdTemplate struct {
	Link     template.HTML
	Repeat   string
	Offset   string
	Type     string
	PublicID string
	Len      string
}

// Select function is the route that the real biding happen
func (tc *selectController) selectVastAd(c echo.Context) error {

	rd, website, province, isp, lenType, length, err := tc.getVastDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	webPublicID := website.WPubID

	middlewares.SetData(c, "site_id", website.WID)
	middlewares.SetData(c, "site_domain", website.WDomain.String)

	var slotFixFound bool
	slotPins := pin.GetPinAdData()
	slotSize, sizeNumSlice, vastSlotData := tc.slotSizeVast(rd.Mobile, webPublicID, length, *website)

	middlewares.SetData(c, "video_len", length)
	middlewares.SetData(c, "ad_count", len(sizeNumSlice))
	var floorBids = make(map[string]int64)
	for i := range sizeNumSlice {
		floorBids[i] = minCPCVast.Int64()
	}
	// TODO : Move this to slotSizeVast func
	for i := range slotSize {
		slotSize[i].ExtraParam = map[string]string{
			"pos":  vastSlotData[i].Offset,
			"type": vastSlotData[i].Type,
			"l":    lenType,
		}
	}
	slotFixFound, slotSize, sizeNumSlice, slotPins, fixSlotSize, _ := checkForFixSlot(slotPins, slotSize, sizeNumSlice, "vast")
	//call context
	m := selector.Context{
		RequestData:      *rd,
		Website:          website,
		Size:             sizeNumSlice,
		Province:         province,
		ISP:              isp,
		SlotPins:         slotPins,
		MinBidPercentage: 1, // TODO : Hard coded shit.
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), vastSelector)
	var show map[string]string
	show, _ = tc.makeShow(c, "vast", rd, filteredAds, nil, sizeNumSlice, slotSize, nil, website, true, floorBids, allowUnderFloor.Bool(), true, floorDivVast.Int64(), false)
	var vTemp = make([]vastAdTemplate, 0)
	if slotFixFound {
		for _, val := range slotPins {
			reserve := make(map[string]string)
			tc.updateMegaKey(rd, val.AdID, val.Bid, val.SlotID, "", "", "")
			tmp := <-utils.ID
			reserve[val.SlotPublicID] = tmp
			store.Set(reserve[val.SlotPublicID], fmt.Sprintf("%d", val.AdID))
			u := url.URL{
				Scheme: rd.Scheme,
				Host:   rd.Host,
				Path:   fmt.Sprintf("/show/%s/%s/%d/%s", "vast", rd.MegaImp, website.GetID(), tmp),
			}
			v := url.Values{}
			v.Set("tid", rd.TID)
			v.Set("ref", rd.Referrer)
			v.Set("parent", rd.Parent)
			v.Set("s", fmt.Sprintf("%d", val.SlotID))
			for i, j := range fixSlotSize[val.SlotPublicID].ExtraParam {
				v.Set(i, j)
			}
			u.RawQuery = v.Encode()
			show[val.SlotPublicID] = u.String()

			vTemp = append(vTemp, vastAdTemplate{
				Link:   template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", show[val.SlotPublicID])),
				Offset: vastSlotData[val.SlotPublicID].Offset,
				Type:   vastSlotData[val.SlotPublicID].Type,
				Repeat: vastSlotData[val.SlotPublicID].Repeat,
			})
		}
	}

	for i := range sizeNumSlice {
		vTemp = append(vTemp, vastAdTemplate{
			Link:   template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", show[i])),
			Offset: vastSlotData[i].Offset,
			Type:   vastSlotData[i].Type,
			Repeat: vastSlotData[i].Repeat,
		})
	}
	result := &bytes.Buffer{}

	assert.Nil(vastIndex.Execute(result, vTemp))
	return c.XMLBlob(http.StatusOK, result.Bytes())
}

func (tc *selectController) slotSizeVast(mobile bool, websitePublicID int64, length map[string][]string, website models.Website, alladscase ...bool) (map[string]*slotData, map[string]int, map[string]vastSlotData) {
	var sizeNumSlice = make(map[string]int)
	var slotPublic []string
	var vastSlot = make(map[string]vastSlotData)
	var i int
	i = 0
	for m := range length {
		i++
		lenType := length[m][0]
		if lenType != "linear" && mobile {
			continue
		}
		pub := fmt.Sprintf("%d%s", websitePublicID, length[m][1])
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
	// add to prevent panic in the query
	if len(slotPublic) == 0 {
		return make(map[string]*slotData), make(map[string]int), make(map[string]vastSlotData)
	}
	finalSlotData, finalSizeNumSlice := tc.slotSizeNormal(slotPublic, website.WID, sizeNumSlice)
	return finalSlotData, finalSizeNumSlice, vastSlot

}

func (tc *selectController) getVastDataFromCtx(c echo.Context) (*middlewares.RequestData, *models.Website, int64, int64, string, map[string][]string, error) {
	rd := middlewares.MustGetRequestData(c)

	publicID, err := strconv.ParseInt(c.QueryParam("a"), 10, 0)
	if err != nil {
		return nil, nil, 0, 0, "", nil, errors.New("invalid request")
	}
	//fetch website and set in Context
	website, err := tc.fetchWebsite(publicID)
	if err != nil {
		return nil, nil, 0, 0, "", nil, errors.New("invalid request")
	}
	start := c.QueryParam("start")
	mid := c.QueryParam("mid")
	end := c.QueryParam("end")
	if !website.GetActive() {
		return nil, nil, 0, 0, "", nil, errors.New("web is not active")
	}

	if !models.NewManager().IsUserActive(website.UserID) {
		return nil, nil, 0, 0, "", nil, errors.New("user is banned")
	}

	province, isp, ll := ip2location.GetProvinceISPByIP(rd.IP)
	middlewares.SetData(c, "province", ll.Province)
	middlewares.SetData(c, "country", ll.Country)
	middlewares.SetData(c, "city", ll.City)
	middlewares.SetData(c, "isp", ll.ISP)
	lenVast, vastCon := utils.MakeVastLen(c.QueryParam("l"), start, mid, end)
	return rd, website, province, isp, lenVast, vastCon, nil
}

// TODO : Move this function to models and fix the cache problem
func (tc *selectController) slotSizeNormal(slotPublic []string, webID int64, sizeNumSlice map[string]int, alladscase ...bool) (map[string]*slotData, map[string]int) {
	if len(alladscase) == 1 && alladscase[0] {
		slotData2 := make(map[string]*slotData)

		for i := range sizeNumSlice {
			slotData2[i] = &slotData{
				PublicID: i,
				Ctr:      .1,
				SlotSize: sizeNumSlice[i],
			}
		}

		return slotData2, sizeNumSlice
	}

	slotPublicString := models.Build(slotPublic)
	res, err := models.NewManager().FetchWebSlots(slotPublicString, webID)
	assert.Nil(err)

	answer := make(map[string]*slotData)
	var (
		newSlots []int64
		newSize  []int
	)
	for i := range slotPublic {
		if _, ok := answer[slotPublic[i]]; ok {
			continue
		}
		for j := range res {
			if fmt.Sprintf("%d", res[j].PublicID) == slotPublic[i] {
				answer[slotPublic[i]] = &slotData{
					ID:       res[j].ID,
					PublicID: slotPublic[i],
					SlotSize: sizeNumSlice[slotPublic[i]],
				}
				break
			}
		}
		if _, ok := answer[slotPublic[i]]; !ok {
			s, err := strconv.ParseInt(slotPublic[i], 10, 0)
			if err == nil {
				newSlots = append(newSlots, s)
				newSize = append(newSize, sizeNumSlice[slotPublic[i]])
			}
		}
	}
	if len(newSlots) > 0 {
		// Expire the cache for the select
		key := utils.Hash(fmt.Sprintf("slot_%s_%d", slotPublicString, webID))
		err := aredis.RemoveKey(key)
		logrus.Debug(err)
	}
	insertedSlots := tc.insertNewSlots(webID, newSlots, newSize)
	for i := range insertedSlots {
		answer[i] = &slotData{
			ID:       insertedSlots[i],
			PublicID: i,
			SlotSize: sizeNumSlice[i],
			Ctr:      defaultCTR.Float64(),
		}
	}

	for i := range answer {
		result, err := aredis.SumHMGetField(transport.KeyGenDaily(transport.Slot, strconv.FormatInt(answer[i].ID, 10)), dailyClickDays.Int(), "i", "c")
		if err != nil || result["c"] == 0 || result["i"] < minImp.Int64() {
			answer[i].Ctr = defaultCTR.Float64()
		} else {
			answer[i].Ctr = utils.Ctr(result["i"], result["c"])
		}
	}

	return answer, sizeNumSlice
}
