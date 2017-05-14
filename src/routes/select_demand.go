package routes

import (
	"config"
	"encoding/json"
	"errors"
	"fmt"
	"middlewares"
	"mr"
	"net/http"
	"selector"

	"strings"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

type Demand struct {
	ID      int64  `json:"id"`
	CPM     int64  `json:"cpm"`
	With    string `json:"with"`
	Height  string `json:"height"`
	URL     string `json:"url"`
	Landing string `json:"landing"`
}

// Select function is the route that the real biding happen
func (tc *selectController) selectDemandWebAd(c echo.Context) error {
	//t := time.Now()
	rd, e, website, province, err := tc.getWebDataExchangeFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice := tc.slotSizeWebExchange(e.Slots, *website)
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	show, ads := tc.makeShow(c, "web", rd, filteredAds, sizeNumSlice, slotSize, website, false, config.Config.Clickyab.MinCPCWeb)

	//substitute the webMobile slot if exists
	wm := fmt.Sprintf("%d%s", website.WPubID, webMobile)
	val, ok := show[wm]
	if ok {
		show["web-mobile"] = val
		delete(show, wm)
	}

	for i := range ads {
		c := []Demand{}
		d := Demand{
			ID:      ads[i].AdID,
			Height:  config.GetSizeByNumStringWith(ads[i].AdSize),
			With:    config.GetSizeByNumStringWith(ads[i].AdSize),
			URL:     show[i],
			CPM:     ads[i].CPM,
			Landing: stripURLParts(ads[i].AdURL.String),
		}
		c = append(c, d)
	}

	/*b, _ := json.MarshalIndent(c, "\t", "\t")
	result := "renderFarm(" + string(b) + "); \n//" + time.Since(t).String()*/
	result, err := json.Marshal(c)

	return c.JSON(200, result)
}
func stripURLParts(url string) string {
	//Lower case the url
	url = strings.ToLower(url)

	//Strip protocol
	if index := strings.Index(url, "://"); index > -1 {
		url = url[index+3:]
	}

	//Strip path (and query with it)
	if index := strings.Index(url, "/"); index > -1 {
		url = url[:index]
	} else if index := strings.Index(url, "?"); index > -1 { //Strip query if path is not found
		url = url[:index]
	}

	//Return domain
	return url
}

//func (tc *selectController) getWebDataExchangeFromCtx(c echo.Context) (*middlewares.RequestDataExchange, *mr.Website, string, error) {
func (tc *selectController) getWebDataExchangeFromCtx(c echo.Context) (*middlewares.RequestData, *middlewares.RequestDataFromExchange, *mr.Website, *mr.Province, error) {
	rd := middlewares.MustGetRequestData(c)
	e := middlewares.MustExchangeGetRequestData(c)
	u_id, err := config.GetSupplier(e.Source.Supplier)
	if err != nil {
		return nil, nil, nil, nil, errors.New("invalid request")
	}
	website, err := tc.fetchWebsiteDomain(fmt.Sprintf("%s/%s", e.Source.Supplier, e.Source.Website), u_id)
	if err != nil {
		return nil, nil, nil, nil, errors.New("invalid request")
	}

	if !website.GetActive() {
		return nil, nil, nil, nil, errors.New("web is not active")
	}

	if !mr.NewManager().IsUserActive(website.UserID) {
		return nil, nil, nil, nil, errors.New("user is banned")
	}

	//province := rd.Province.Name
	province, err := tc.fetchProvince(rd.IP, c.Request().Header.Get("Cf-Ipcountry"))
	if err != nil {
		logrus.Debug(err)
	}
	return rd, e, website, province, nil
}

//fetchWebsiteDomain website and check if the minimum floor is applied
func (tc *selectController) fetchWebsiteDomain(domain string, user int64) (*mr.Website, error) {
	website, err := mr.NewManager().FetchWebsiteByDomain(domain)
	if err != nil {
		_, err := mr.NewManager().InsertWebsite(domain, user)
		if err != nil {
			return nil, err
		}
		website, err = mr.NewManager().FetchWebsiteByDomain(domain)
		if err != nil {
			return nil, err
		}
	}
	if website.WFloorCpm.Int64 < config.Config.Clickyab.MinCPMFloorWeb {
		website.WFloorCpm.Int64 = config.Config.Clickyab.MinCPMFloorWeb
	}
	return website, err
}

func (tc selectController) slotSizeWebExchange(slots []middlewares.Slot, website mr.Website) (map[string]*slotData, map[string]int) {
	var sizeNumSlice = make(map[string]int)
	var slotPublics []string
	for slot := range slots {
		size, err := config.GetSize(fmt.Sprintf("%dx%d", slots[slot].Width, slots[slot].Height))
		slotPublic := fmt.Sprintf("%d%d%d", website.WPubID, size, slot)
		if err != nil {
			sizeNumSlice[slotPublic] = size
			slotPublics = append(slotPublics, slotPublic)
		}
	}
	return tc.slotSizeNormal(slotPublics, website.WID, sizeNumSlice)
}
