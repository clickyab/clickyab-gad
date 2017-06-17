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

	ads := make([]nativeAd, 0)
	var p protocol = httpScheme
	if rd.Scheme == httpsScheme {
		p = httpsScheme
	}

	for _, j := range h {
		if j == nil {
			continue
		}

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
		ads = append(ads, nativeAd{
			Image:   j.AdImg.String,
			URL:     u.String(),
			Lead:    j.AdAttribute["banner_description_text_type"].(string),
			More:    params.Get("more"),
			Title:   j.AdAttribute["banner_title_text_type"].(string),
			Corners: params.Get("corners"),
			Site:    j.AdURL.String,
		})

	}

	logrus.Debugf("%+v", ads)
	if len(ads) == 0 {
		return c.HTML(http.StatusBadRequest, "<div class=\"no-ads\"></div>")
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
