package routes

import (
	"config"
	"net/http"
	selector2 "pin"
	"selector"
	"strings"
	"utils"

	"fmt"
	"net/url"

	"mr"
	"store"

	"middlewares"

	"strconv"

	"github.com/sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

// Select function is the route that the real biding happen
func (tc *selectController) selectNativeAd(c echo.Context) error {
	logrus.Debug("select native ad")
	params := c.QueryParams()
	rd, website, province, isp, err := tc.getWebDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}

	middlewares.SetData(c, "site_id", website.WID)
	middlewares.SetData(c, "site_domain", website.WDomain.String)

	var slotFixFound bool
	slotPins := selector2.GetPinAdData()
	slotSize, sizeNumSlice, order := tc.slotSizeNative(c, *website)
	slotFixFound, slotSize, sizeNumSlice, slotPins, _, _ = checkForFixSlot(slotPins, slotSize, sizeNumSlice, "native")

	middlewares.SetData(c, "ad_count", len(sizeNumSlice))
	//call context
	m := selector.Context{
		RequestData:      *rd,
		Website:          website,
		Size:             sizeNumSlice,
		Province:         province,
		ISP:              isp,
		SlotPins:         slotPins,
		MinBidPercentage: 0.5, // TODO : Hard coded :) make it some how calculated
	}
	// remove fucking order
	var resOrder = []string{}
	for i := range order {
		var appending bool = true
		for _, v := range slotPins {
			if v.SlotPublicID == order[i] {
				appending = false
			}
		}
		if appending {
			resOrder = append(resOrder, order[i])
		}
	}
	var h = make(map[string]*mr.AdData)
	filteredAds := selector.Apply(&m, selector.GetAdData(), nativeSelector)
	// TODO : Currently underfloor is always true
	_, h = tc.makeShow(c, "sync", rd, filteredAds, resOrder, sizeNumSlice, slotSize, nil, website, false, config.Config.Clickyab.MinCPCNative, true, true, config.Config.Clickyab.FloorDiv.Native)

	if slotFixFound {
		for _, val := range slotPins {
			res := &val.AdData
			res.WinnerBid = val.Bid
			res.SlotID = val.SlotID
			h[val.SlotPublicID] = res
			reserve := make(map[string]string)
			tc.updateMegaKey(rd, val.AdID, val.Bid, val.SlotID, "", "", "")
			tmp := config.Config.MachineName + <-utils.ID
			reserve[val.SlotPublicID] = tmp
			store.Set(reserve[val.SlotPublicID], fmt.Sprintf("%d", val.AdID))
		}
	}

	ads := make([]nativeAd, 0)
	var p protocol = httpScheme
	if rd.Scheme == httpsScheme {
		p = httpsScheme
	}
	var count int
	for i := len(order) - 1; i >= 0; i-- {
		j := h[order[i]]
		if j == nil {
			continue
		}
		count++

		rnd := <-utils.ID
		u := url.URL{
			Scheme: rd.Scheme,
			Host:   rd.Host,
			Path:   fmt.Sprintf("/click/%s/%d/%s/%d/%s", "native", website.WID, m.RequestData.MegaImp, j.AdID, rnd),
		}
		v := url.Values{}
		v.Set("tid", rd.TID)
		v.Set("ref", rd.Referrer)
		v.Set("parent", rd.Parent)
		u.RawQuery = v.Encode()
		//middlewares.SafeGO(c, false, false, func() {
		imp := tc.fillNativeImp(rd, false, j, j.WinnerBid, website, j.SlotID)
		tc.callWebWorker(website, j.SlotID, j.AdID, m.RequestData.MegaImp, rnd, imp, rd)
		//})

		if v == nil {
			continue
		}
		if p == httpsScheme {
			j.AdImg.String = strings.Replace(j.AdImg.String, "http://", "https://", -1)
		}
		fixTitle := utils.LimitCharacter(j.AdName.String, 50)
		nAd := nativeAd{
			Image:   j.AdImg.String,
			URL:     u.String(),
			Lead:    j.AdAttribute["banner_description_text_type"].(string),
			More:    params.Get("more"),
			Title:   fixTitle,
			Corners: params.Get("corners"),
			Site:    j.AdURL.String,
			Extra:   fmt.Sprintf("CTR=%f, CPM=%d, Winner=%d == %s", j.CTR, j.CPM, j.WinnerBid, j.Extra),
		}
		var fixSlot bool
		var direct bool
		for n := range slotPins {
			if slotPins[n].SlotID == j.SlotID {
				fixSlot = true
				direct = slotPins[n].Direct
				break
			}
		}

		if fixSlot && direct {
			nAd.URL = j.AdURL.String
		}
		ads = append(ads, nAd)

	}
	// TODO : handle this in select
	if count > 1 && count%2 == 1 {
		ads = ads[:len(ads)-1]
	}

	if len(ads) == 0 {
		return c.HTML(http.StatusBadRequest, "<div class=\"no-ads\"></div>")
	}

	//check min size
	if params.Get("minsize") != "" {
		minInt, err := strconv.ParseInt(params.Get("minsize"), 10, 64)
		if err != nil || minInt < 90 || minInt > 150 {
			minInt = 141
		}
	}

	n := nativeContainer{
		Ads:        ads,
		Title:      params.Get("title"),
		FontSize:   params.Get("fontSize"),
		FontFamily: params.Get("fontFamily"),
		Position:   params.Get("position"),
		MinSize:    params.Get("minsize"),
		IsVertical: params.Get("orientation") == "vertical",
	}

	return c.HTML(200, renderNative(n))

}
