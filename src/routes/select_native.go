package routes

import (
	"config"
	"net/http"
	"selector"
	"strings"
	"utils"

	"fmt"
	"net/url"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

// Select function is the route that the real biding happen
func (tc *selectController) selectNativeAd(c echo.Context) error {
	logrus.Debug("select native ad")
	params := c.QueryParams()
	rd, website, province, err := tc.getWebDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice, order := tc.slotSizeNative(params, *website)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), nativeSelector)
	// TODO : Currently underfloor is always true
	_, h := tc.makeShow(c, "sync", rd, filteredAds, order, sizeNumSlice, slotSize, nil, website, false, config.Config.Clickyab.MinCPCNative, true, true, config.Config.Clickyab.FloorDiv.Native)
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
		ads = append(ads, nativeAd{
			Image:   j.AdImg.String,
			URL:     u.String(),
			Lead:    j.AdAttribute["banner_description_text_type"].(string),
			More:    params.Get("more"),
			Title:   fixTitle,
			Corners: params.Get("corners"),
			Site:    j.AdURL.String,
			Extra:   fmt.Sprintf("CTR=%f, CPM=%d, Winner=%d == %s", j.CTR, j.CPM, j.WinnerBid, j.Extra),
		})

	}
	// TODO : handle this in select
	if count > 1 && count%2 == 1 {
		ads = ads[:len(ads)-1]
	}

	if len(ads) == 0 {
		return c.HTML(http.StatusBadRequest, "<div class=\"no-ads\"></div>")
	}

	n := nativeContainer{
		Ads:        ads,
		Title:      params.Get("title"),
		FontSize:   params.Get("fontSize"),
		Position:   params.Get("position"),
		IsVertical: params.Get("orientation") == "vertical",
	}

	return c.HTML(200, renderNative(n))

}
