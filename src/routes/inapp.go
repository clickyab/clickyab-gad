package routes

import (
	"assert"
	"errors"
	"filter"
	"middlewares"
	"mr"
	"net/http"
	"selector"
	"strconv"

	"strings"

	"fmt"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

var (
	appSelector = selector.Mix(
		filter.IsAppNetwork,
		filter.CheckWebSize,
		filter.CheckAppBlackList,
		filter.CheckAppWhiteList,
		filter.CheckAppCategory,
		filter.CheckProvince,
		filter.CheckAppBrand,
		filter.CheckProvder,
		filter.CheckAppHood,
		filter.CheckAppAreaInGlob,
	)
)

func (tc *selectController) inApp(c echo.Context) error {
	//t := time.Now()
	rd, app, province, phone, cell, err := tc.getAppDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice := tc.slotSizeApp(c, app)
	//call context
	m := selector.Context{
		RequestData:  *rd,
		Website:      nil,
		Size:         sizeNumSlice,
		Province:     province,
		App:          app,
		PhoneData:    phone,
		CellLocation: cell,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), appSelector)
	_, ads := tc.makeShow(c, "sync", rd, filteredAds, sizeNumSlice, slotSize, app, false)
	assert.True(len(ads) == 1, "[BUG] why select no ad?")

	return nil
}

func (tc *selectController) slotSizeApp(ctx echo.Context, app *mr.App) (map[string]*slotData, map[string]int) {
	adsMedia := ctx.Request().URL.Query().Get("adsMedia")
	var bs int
	switch strings.ToLower(adsMedia) {
	case "banner":
		bs = 8
	case "largebanner":
		bs = 3
	case "xlargebannerportrait", "fullbannerportrait":
		bs = 16
	case "xlargebannerlandscap", "fullbannerlandscape":
		bs = 17
	}
	slotString := fmt.Sprintf("%d0%d0%d", app.ID, app.UserID, bs)
	slot, _ := strconv.ParseInt(slotString, 10, 0)
	s, err := mr.NewManager().FetchAppSlot(app.ID, slot)
	if err != nil {
		// no slot found
		slots, err := mr.NewManager().InsertSlots(0, app.ID, slot)
		assert.Nil(err)
		s = &slots[0]
	}
	data := map[string]*slotData{
		slotString: &slotData{
			SlotSize: bs,
			ID:       s.ID,
			PublicID: slotString,
		},
	}
	sizes := map[string]int{slotString: bs}

	return data, sizes
}

func (tc *selectController) getAppDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.App, *mr.Province, *mr.PhoneData, *mr.CellLocation, error) {
	rd := middlewares.MustGetRequestData(c)

	token := c.Request().URL.Query().Get("token")
	if len(token) < 1 {
		return nil, nil, nil, nil, nil, errors.New("invalid request")
	}
	m := mr.NewManager()
	app, err := m.GetApp(token)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if !app.GetActive() {
		return nil, nil, nil, nil, nil, errors.New("app is disabled")
	}

	if !m.IsUserActive(app.UserID) {
		return nil, nil, nil, nil, nil, errors.New("user is banned")
	}

	province, err := tc.fetchProvince(rd.IP, c.Request().Header.Get("Cf-Ipcountry"))
	if err != nil {
		logrus.Debug(err)
	}

	phone := m.GetPhoneData(c.Request().URL.Query().Get("brand"), c.Request().URL.Query().Get("carrier"), c.Request().URL.Query().Get("network"))
	mcc, _ := strconv.ParseInt(c.Request().URL.Query().Get("mcc"), 10, 0)
	mnc, _ := strconv.ParseInt(c.Request().URL.Query().Get("mnc"), 10, 0)
	cid, _ := strconv.ParseInt(c.Request().URL.Query().Get("cid"), 10, 0)
	lac, _ := strconv.ParseInt(c.Request().URL.Query().Get("lac"), 10, 0)
	cell, err := m.GetCellLocation(mcc, mnc, lac, cid, phone.Carrier)
	if err != nil {
		logrus.Warn(err)
	}

	return rd, app, province, phone, cell, nil
}
