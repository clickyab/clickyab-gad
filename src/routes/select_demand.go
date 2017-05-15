package routes

import (
	"config"
	"errors"
	"fmt"
	"middlewares"
	"mr"
	"net/http"
	"selector"

	"assert"

	"net/url"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

type Demand struct {
	ID          string `json:"id"`
	CPM         int64  `json:"cpm"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	URL         string `json:"url"`
	Landing     string `json:"landing"`
	SlotTrackID string `json:"slot_track_id"`
}

// Select function is the route that the real biding happen
func (tc *selectController) selectDemandWebAd(c echo.Context) error {
	//t := time.Now()
	rd, e, website, province, err := tc.getWebDataExchangeFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice, trackIDs, err := tc.slotSizeWebExchange(e.Slots, *website)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//call context
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    province,
	}

	var sel selector.FilterFunc
	if e.Platform == "web" {
		sel = webSelector
	} else {
		return c.HTML(http.StatusBadRequest, "not supported platform")
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), sel)
	show, ads := tc.makeShow(c, "sync", rd, filteredAds, sizeNumSlice, slotSize, website, false, config.Config.Clickyab.MinCPCWeb, e.Underfloor)

	//substitute the webMobile slot if exists
	dm := []Demand{}
	for i := range ads {
		if ads[i] == nil {
			continue
		}

		d := Demand{
			ID:          fmt.Sprint(ads[i].AdID),
			Height:      config.GetSizeByNumStringHeight(ads[i].AdSize),
			Width:       config.GetSizeByNumStringWith(ads[i].AdSize),
			URL:         show[i],
			CPM:         ads[i].CPM,
			Landing:     stripURLParts(ads[i].AdURL.String),
			SlotTrackID: trackIDs[ads[i].SlotPublicID],
		}
		assert.False(d.SlotTrackID == "", "[BUG] invalid track id")
		dm = append(dm, d)
	}
	if len(dm) < 1 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, dm)
}
func stripURLParts(in string) string {
	u, err := url.Parse(in)
	if err != nil {
		return ""
	}

	return u.Host
}

func (tc *selectController) getWebDataExchangeFromCtx(c echo.Context) (*middlewares.RequestData, *middlewares.RequestDataFromExchange, *mr.Website, *mr.Province, error) {
	rd := middlewares.MustGetRequestData(c)
	e := middlewares.MustExchangeGetRequestData(c)
	name, userID, err := config.GetSupplier(e.Source.Supplier)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("can not accept from %s demand", e.Source.Supplier)
	}
	e.Source.Supplier = name
	website, err := tc.fetchWebsiteDomain(fmt.Sprintf("%s/%s", e.Source.Supplier, e.Source.Name), userID)
	if err != nil {
		logrus.Warn(err)
		return nil, nil, nil, nil, errors.New("invalid request")
	}
	// Set the floor here. its related to the demand request not our data
	website.WFloorCpm.Int64, website.WFloorCpm.Valid = int64(e.Source.FloorCPM), true
	if !website.GetActive() {
		return nil, nil, nil, nil, errors.New("web is not active")
	}

	if !mr.NewManager().IsUserActive(website.UserID) {
		return nil, nil, nil, nil, errors.New("user is banned")
	}

	var province *mr.Province
	if e.Location.Province.Valid {
		province, err = tc.fetchProvinceDemand(e.Location.Province.Name)
		if err != nil {
			logrus.Debug(err)
		}
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

func (tc selectController) slotSizeWebExchange(slots []middlewares.Slot, website mr.Website) (map[string]*slotData, map[string]int, map[string]string, error) {
	var sizeNumSlice = make(map[string]int)
	var slotPublics []string
	var trackIDs = make(map[string]string)
	for slot := range slots {
		size, err := config.GetSize(fmt.Sprintf("%dx%d", slots[slot].Width, slots[slot].Height))
		slotPublic := fmt.Sprintf("%d%d%d", website.WPubID, size, slot)
		if err == nil {
			sizeNumSlice[slotPublic] = size
			slotPublics = append(slotPublics, slotPublic)
			trackIDs[slotPublic] = slots[slot].TrackID
		}
	}
	if len(slotPublics) == 0 {
		return nil, nil, nil, errors.New("no supported slot size")
	}
	all, size := tc.slotSizeNormal(slotPublics, website.WID, sizeNumSlice)
	return all, size, trackIDs, nil
}
