package routes

import (
	"assert"
	"config"
	"errors"
	"fmt"
	"middlewares"
	"mr"
	"net/http"
	"net/url"
	"redis"
	"selector"
	"time"

	"redlock"

	"strconv"
	"transport"
	"utils"

	"github.com/Sirupsen/logrus"
	echo "gopkg.in/labstack/echo.v3"
)

type Demand struct {
	ID          string `json:"id"`
	CPM         int64  `json:"max_cpm"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	URL         string `json:"url"`
	Landing     string `json:"landing"`
	SlotTrackID string `json:"slot_track_id"`
}

func (tc *selectController) selectDemandAppAd(c echo.Context, rd *middlewares.RequestData, e *middlewares.RequestDataFromExchange) error {
	rd, e, app, province, phone, cell, err := tc.getDemandAppDataFromCtx(c, rd, e)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice, trackIDs, attr, err := tc.slotSizeAppExchange(e.Slots, *app)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	m := selector.Context{
		RequestData:  *rd,
		App:          app,
		Size:         sizeNumSlice,
		Province:     province,
		PhoneData:    phone,
		CellLocation: cell,
	}
	var sel selector.FilterFunc
	if e.Platform == "app" {
		sel = appSelector
	} else {
		return c.HTML(http.StatusBadRequest, "not supported platform")
	}
	lockSession := "DRD_SESS_" + e.SessionKey
	redlock.Lock(lockSession, lockSession, 100*time.Millisecond)
	defer redlock.Unlock([]string{lockSession}, lockSession)
	var sessionAds []int64
	// This is when the supplier is not support grouping
	if e.SessionKey != "" {
		e.SessionKey = "EXC_SESS_" + e.SessionKey
		sessionAds = aredis.SMembersInt(e.SessionKey)
		if len(sessionAds) > 0 {
			sel = selector.Mix(sel, func(_ *selector.Context, a mr.AdData) bool {
				for _, i := range sessionAds {
					if i == a.AdID {
						return false
					}
				}
				return true
			})
		}
	}

	filteredAds := selector.Apply(&m, selector.GetAdData(), sel)
	show, ads := tc.makeShow(c, "sync", rd, filteredAds, sizeNumSlice, slotSize, attr, app, false, config.Config.Clickyab.MinCPCApp, config.Config.Clickyab.UnderFloor, true)

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
			CPM:         int64(float64(ads[i].WinnerBid) * ads[i].CTR * 10),
			Landing:     stripURLParts(ads[i].AdURL.String),
			SlotTrackID: trackIDs[ads[i].SlotPublicID],
		}
		assert.False(d.SlotTrackID == "", "[BUG] invalid track id")
		dm = append(dm, d)
		sessionAds = append(sessionAds, ads[i].AdID)
	}
	if len(dm) < 1 {
		return c.NoContent(http.StatusNoContent)
	}

	if e.SessionKey != "" {
		err := aredis.SAddInt(e.SessionKey, true, time.Minute, sessionAds...)
		if err != nil {
			logrus.Warn(err)
		}
	}

	return c.JSON(http.StatusOK, dm)
}

func (tc *selectController) selectDemandWebAd(c echo.Context, rd *middlewares.RequestData, e *middlewares.RequestDataFromExchange) error {
	rd, e, website, province, err := tc.getWebDataExchangeFromCtx(c, rd, e)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice, trackIDs, attr, err := tc.slotSizeWebExchange(e.Slots, *website)
	if err != nil {
		return c.HTML(http.StatusBadRequest, "slot size was wrong, reason : "+err.Error())
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
	lockSession := "DRD_SESS_" + e.SessionKey
	redlock.Lock(lockSession, lockSession, 100*time.Millisecond)
	defer redlock.Unlock([]string{lockSession}, lockSession)
	var sessionAds []int64
	// This is when the supplier is not support grouping
	if e.SessionKey != "" {
		e.SessionKey = "EXC_SESS_" + e.SessionKey
		sessionAds = aredis.SMembersInt(e.SessionKey)
		if len(sessionAds) > 0 {
			sel = selector.Mix(sel, func(_ *selector.Context, a mr.AdData) bool {
				for _, i := range sessionAds {
					if i == a.AdID {
						return false
					}
				}
				return true
			})
		}
	}

	filteredAds := selector.Apply(&m, selector.GetAdData(), sel)
	show, ads := tc.makeShow(c, "sync", rd, filteredAds, sizeNumSlice, slotSize, attr, website, false, e.Source.FloorCPM, e.Underfloor, true)

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
			CPM:         int64(float64(ads[i].WinnerBid) * ads[i].CTR * 10),
			Landing:     stripURLParts(ads[i].AdURL.String),
			SlotTrackID: trackIDs[ads[i].SlotPublicID],
		}
		assert.False(d.SlotTrackID == "", "[BUG] invalid track id")
		dm = append(dm, d)
		sessionAds = append(sessionAds, ads[i].AdID)
	}
	if len(dm) < 1 {
		return c.NoContent(http.StatusNoContent)
	}

	if e.SessionKey != "" {
		err := aredis.SAddInt(e.SessionKey, true, time.Minute, sessionAds...)
		if err != nil {
			logrus.Warn(err)
		}
	}

	return c.JSON(http.StatusOK, dm)
}

// selectDemandWebAd function is the route that the real biding happens
func (tc *selectController) selectDemandAd(c echo.Context) error {
	//t := time.Now()
	rd := middlewares.MustGetRequestData(c)
	e := middlewares.MustExchangeGetRequestData(c)
	if e.Platform != "app" && e.Platform != "web" {
		return c.HTML(http.StatusBadRequest, "wrong platform")
	}
	if e.Platform == "web" {
		return tc.selectDemandWebAd(c, rd, e)
	} // app platform selected
	return tc.selectDemandAppAd(c, rd, e)

}

func (tc *selectController) getDemandAppDataFromCtx(c echo.Context, rd *middlewares.RequestData, e *middlewares.RequestDataFromExchange) (*middlewares.RequestData, *middlewares.RequestDataFromExchange, *mr.App, *mr.Province, *mr.PhoneData, *mr.CellLocation, error) {
	name, userID, err := config.GetSupplier(e.Source.Supplier)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, fmt.Errorf("can not accept from %s demand", e.Source.Supplier)
	}
	e.Source.Supplier = name
	app, err := tc.fetchAppPackage(e.Source.Name, e.Source.Supplier, userID)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, errors.New("invalid request")
	}
	var province *mr.Province
	if e.Location.Province.Valid {
		province, err = tc.fetchProvinceDemand(e.Location.Province.Name)
		if err != nil {
			logrus.Debug(err)
		}
	}
	m := mr.NewManager()
	phone := m.GetPhoneData(rd.Brand, rd.Carrier, rd.Network)
	cell, err := m.GetCellLocation(rd.MCC, rd.MNC, rd.LAC, rd.CID, rd.Carrier)
	if err != nil {
		logrus.Debug(err)
	}
	return rd, e, app, province, phone, cell, nil

}

func stripURLParts(in string) string {
	u, err := url.Parse(in)
	if err != nil {
		return ""
	}

	return u.Host
}

func (tc *selectController) getWebDataExchangeFromCtx(c echo.Context, rd *middlewares.RequestData, e *middlewares.RequestDataFromExchange) (*middlewares.RequestData, *middlewares.RequestDataFromExchange, *mr.Website, *mr.Province, error) {
	name, userID, err := config.GetSupplier(e.Source.Supplier)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("can not accept from %s supplier", e.Source.Supplier)
	}
	e.Source.Supplier = name
	website, err := tc.fetchWebsiteDomain(e.Source.Name, e.Source.Supplier, userID)
	if err != nil {
		return nil, nil, nil, nil, err //errors.New("invalid request")
	}
	// Set the floor here. its related to the demand request not our data
	website.WFloorCpm.Int64, website.WFloorCpm.Valid = int64(e.Source.FloorCPM), true
	if !website.GetActive() {
		return nil, nil, nil, nil, errors.New("website is not active")
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
func (tc *selectController) fetchWebsiteDomain(domain, supplier string, user int64) (*mr.Website, error) {
	website, err := mr.NewManager().FetchWebsiteByDomain(domain, supplier)
	if err != nil {
		website, err = mr.NewManager().InsertWebsite(domain, supplier, user)
		if err != nil {
			return nil, err
		}
	}
	if website.WFloorCpm.Int64 < config.Config.Clickyab.MinCPMFloorWeb {
		website.WFloorCpm.Int64 = config.Config.Clickyab.MinCPMFloorWeb
	}
	return website, err
}

// fetchAppPackage fetch app if not exists insert it
func (tc *selectController) fetchAppPackage(pack, supplier string, userID int64) (*mr.App, error) {
	app, err := mr.NewManager().FetchAppByPack(pack, supplier)
	if err != nil {
		app, err = mr.NewManager().InsertApp(pack, supplier, userID)
		if err != nil {
			return nil, err
		}
	}
	return app, nil
}

func (tc selectController) slotSizeWebExchange(slots []middlewares.Slot, website mr.Website) (map[string]*slotData, map[string]int, map[string]string, map[string]map[string]string, error) {
	var sizeNumSlice = make(map[string]int)
	var slotPublics []string
	var trackIDs = make(map[string]string)
	var attr = make(map[string]map[string]string)
	for slot := range slots {
		size, err := config.GetSize(fmt.Sprintf("%dx%d", slots[slot].Width, slots[slot].Height))
		slotPublic := fmt.Sprintf("%d%d%d", website.WPubID, size, slot)
		attr[slotPublic] = make(map[string]string)
		for k, v := range slots[slot].Attributes {
			attr[slotPublic][k] = v
		}
		if err == nil {
			sizeNumSlice[slotPublic] = size
			slotPublics = append(slotPublics, slotPublic)
			trackIDs[slotPublic] = slots[slot].TrackID
		}
	}
	if len(slotPublics) == 0 {
		return nil, nil, nil, nil, errors.New("no supported slot size")
	}
	all, size := tc.slotSizeNormal(slotPublics, website.WID, sizeNumSlice)

	return all, size, trackIDs, attr, nil
}

func (tc selectController) slotSizeAppExchange(slots []middlewares.Slot, app mr.App) (map[string]*slotData, map[string]int, map[string]string, map[string]map[string]string, error) {
	var sizeNumSlice = make(map[string]int)
	var slotPublics []string
	var trackIDs = make(map[string]string)
	var attr = make(map[string]map[string]string)

	for slot := range slots {
		size, err := config.GetSize(fmt.Sprintf("%dx%d", slots[slot].Width, slots[slot].Height))
		slotPublic := fmt.Sprintf("%d0%d0%d", app.ID, app.UserID, size)
		for k, v := range slots[slot].Attributes {
			attr[slotPublic][k] = v
		}
		if err == nil {
			sizeNumSlice[slotPublic] = size
			slotPublics = append(slotPublics, slotPublic)
			trackIDs[slotPublic] = slots[slot].TrackID
		}
	}
	if len(slotPublics) == 0 {
		return nil, nil, nil, nil, errors.New("no supported slot size")
	}
	all, size := tc.slotSizeAppExchangeNormal(slotPublics, app.ID, sizeNumSlice)
	return all, size, trackIDs, attr, nil
}

func (tc selectController) slotSizeAppExchangeNormal(slotPublic []string, appID int64, sizeNumSlice map[string]int) (map[string]*slotData, map[string]int) {
	slotPublicString := mr.Build(slotPublic)
	res, err := mr.NewManager().FetchAppSlots(slotPublicString, appID)
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
		key := utils.Sha1(fmt.Sprintf("slot_%s_%d", slotPublicString, appID))
		aredis.RemoveKey(key)
	}
	insertedSlots := tc.insertNewAppSlots(appID, newSlots, newSize)
	for i := range insertedSlots {
		answer[i] = &slotData{
			ID:       insertedSlots[i],
			PublicID: i,
			SlotSize: sizeNumSlice[i],
		}
	}

	for i := range answer {
		result, err := aredis.SumHMGetField(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(answer[i].ID, 10)), config.Config.Redis.Days, "i", "c")
		if err != nil || result["c"] == 0 || result["i"] < config.Config.Clickyab.MinImp {
			answer[i].Ctr = config.Config.Clickyab.DefaultCTR
		} else {
			answer[i].Ctr = utils.Ctr(result["i"], result["c"])
		}
	}

	return answer, sizeNumSlice
}
