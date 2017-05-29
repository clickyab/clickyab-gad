package routes

import (
	"config"
	"middlewares"
	"mr"
	"net/http"
	"selector"
	"strings"
	"utils"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

// Select function is the route that the real biding happen
func (tc *selectController) selectNativeAd(c echo.Context) error {
	logrus.Debug("select native ad")
	params := c.QueryParams()
	rd, website, province, err := tc.getWebDataFromCtx(c)
	if err != nil {
		logrus.Warn(1)
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice := tc.slotSizeNative(params, *website)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), nativeSelector)
	logrus.Debug("Pass filters => ", len(filteredAds[20]))
	// TODO : Currently underfloor is always true
	_, h := tc.makeShow(c, "sync", rd, filteredAds, sizeNumSlice, slotSize, website, false, config.Config.Clickyab.MinCPCWeb, true, false)
	logrus.Debugf("%+v", h)
	middlewares.SafeGO(c, false, false, func() {
		for _, j := range h {
			if j == nil {
				continue
			}
			ads, err := mr.NewManager().GetAd(j.AdID, false)
			if err != nil {

			}
			imp := tc.fillImp(rd, false, ads, j.WinnerBid, website, j.SlotID)
			tc.callWebWorker(website, j.SlotID, j.AdID, m.RequestData.MegaImp, <-utils.ID, imp, rd)
		}

	})

	ads := make([]nativeAd, 0)
	var p protocol = httpScheme
	if rd.Scheme == httpsScheme {
		p = httpsScheme
	}
	// check more param
	if params.Get("more") == "" {
		return c.HTML(http.StatusBadRequest, "more not found")
	}
	for _, v := range h {
		if v == nil {
			continue
		}
		if p == httpsScheme {
			v.AdImg.String = strings.Replace(v.AdImg.String, "http://", "https://", -1)
		}
		ads = append(ads, nativeAd{
			Image:   v.AdImg.String,
			URL:     v.AdURL.String,
			Lead:    v.AdAttribute["banner_title_text_type"].(string),
			More:    params.Get("more"),
			Title:   v.AdAttribute["banner_description_text_type"].(string),
			Corners: params.Get("corners"),
			Site:    v.AdURL.String,
		})
	}
	logrus.Debugf("%+v", ads)
	if len(ads) == 0 {
		return c.HTML(http.StatusBadRequest, "no ads")
	}
	var layout layout
	switch params.Get("position") {
	case "left":
		layout = layoutTitleRight
	case "right":
		layout = layoutImageRight
	case "top":
		layout = layoutImageFirst
	case "bottom":
		layout = layoutImageLast
	case "middle":
		layout = layoutTitleFirst
	default:
		return c.HTML(http.StatusBadRequest, "position not valid")
	}
	n := nativeContainer{
		Ads:      ads,
		Layout:   layout,
		Title:    params.Get("title"),
		Font:     params.Get("font"),
		FontSize: params.Get("fontsize"),
	}

	return c.HTML(200, renderNative(n))
}
