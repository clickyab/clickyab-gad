package routes

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"clickyab.com/gad/filter"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"

	"net/http"

	"clickyab.com/gad/redis"
	"clickyab.com/gad/transport"

	"bytes"
	"encoding/json"

	"clickyab.com/gad/ip2location"

	"github.com/sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

var (
	appSelector = selector.Mix(
		filter.IsAppNetwork,
		filter.CheckAppSize,
		filter.CheckAppBlackList,
		filter.CheckAppWhiteList,
		filter.CheckAppCategory,
		filter.CheckAppCarrier,
		filter.CheckProvince,
		filter.CheckAppBrand,
		filter.CheckProvder,
		filter.CheckAppHood,
		filter.CheckAppAreaInGlob,
		filter.CheckISP,
	)
)

const inAppJSON = `{"status":1,"apps":[{"name":"snapp","packaage":"cab.snapp.passenger"},{"name":"tap30","packaage":"taxi.tap30.passenger"},{"name":"ajancy","packaage":"com.mammutgroup.ajancy.passenger"},{"name":"digikala","packaage":"com.digikala"},{"name":"bamilo","packaage":"com.bamilo.android"},{"name":"pintapin","packaage":"com.pintapin.pintapin"},{"name":"alibaba","packaage":"ir.alibaba"}]}`

type appJSON struct {
	Status int64 `json:"status"`
	Apps   []struct {
		Name     string `json:"name"`
		Packaage string `json:"packaage"`
	} `json:"apps"`
}

func (tc *selectController) inApp(c echo.Context) error {
	//t := time.Now()
	sdkVers, _ := strconv.ParseInt(c.Request().URL.Query().Get("clickyabVersion"), 10, 0)
	rd, app, province, isp, phone, cell, err := tc.getAppDataFromCtx(c)
	if err != nil {
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	slotSize, sizeNumSlice, slotString, full := tc.slotSizeApp(c, app)
	//call context
	m := selector.Context{
		RequestData:  *rd,
		Website:      nil,
		Size:         sizeNumSlice,
		Province:     province,
		App:          app,
		PhoneData:    phone,
		CellLocation: cell,
		ISP:          isp,
	}

	middlewares.SetData(c, "app_name", app.AppName)
	middlewares.SetData(c, "app_pkg", app.AppPackage)
	middlewares.SetData(c, "app_id", app.ID)
	middlewares.SetData(c, "brand", phone.Brand)
	middlewares.SetData(c, "carrier", phone.Carrier)
	middlewares.SetData(c, "network", phone.Network)
	middlewares.SetData(c, "ad_count", 1) // Ad in app is always one

	filteredAds := selector.Apply(&m, selector.GetAdData(), appSelector)

	_, ads := tc.makeShow(c, "sync", rd, filteredAds, nil, sizeNumSlice, slotSize, nil, app, false, minCPCApp.Int64(), allowUnderFloor.Bool(), true, floorDivApp.Int64())
	assert.True(len(ads) == 1, "[BUG] why select no ad?")

	var (
		noAd = ads[slotString] == nil
		u    url.URL
		img  string
	)
	rnd := <-utils.ID
	if !noAd {
		ad, err := mr.NewManager().GetAd(ads[slotString].AdID, false)
		assert.Nil(err)
		imp := tc.fillImp(rd, false, ad, ads[slotString].WinnerBid, app, slotSize[slotString].ID)
		u = url.URL{
			Scheme: rd.Scheme,
			Host:   rd.Host,
			// mega in this case is the current request
			Path: fmt.Sprintf("/click/%s/%d/%s/%d/%s", "app", app.ID, rd.MegaImp, ad.AdID, rnd),
		}
		v := url.Values{}
		v.Set("tid", rd.TID)
		v.Set("ref", rd.Referrer)
		v.Set("parent", rd.Parent)
		u.RawQuery = v.Encode()
		// Pass it to worker
		img = ad.AdImg.String
		go tc.callWebWorker(app, slotSize[slotString].ID, ad.AdID, rd.MegaImp, rnd, imp, rd)
	}
	closeClass := "largeclose"
	if slotSize[slotString].SlotSize == 8 {
		closeClass = "close"
	}
	d, err := renderInApp(inappContext{
		FullScreen:    full != "",
		ExtraStyle:    "",
		BodyClass:     full,
		Dynamic:       false,
		DynamicBody:   "",
		FatFinger:     app.AppFatFinger.Bool,
		ClickURL:      u.String(),
		Src:           img,
		CloseClass:    closeClass,
		ImpID:         rand.Int(),
		SdkVersion:    sdkVers,
		RefreshMinute: 2, // TODO : Config,
		NoAd:          noAd,
	})
	assert.Nil(err)

	// This is the actual imp so call the imp
	return c.HTML(http.StatusOK, d)
}

func (tc *selectController) inAppJSON(c echo.Context) error {
	res := appJSON{}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(inAppJSON)))
	err := dec.Decode(&res)
	assert.Nil(err)
	return c.JSON(http.StatusOK, res)
}

func (tc *selectController) slotSizeApp(ctx echo.Context, app *mr.App) (map[string]*slotData, map[string]int, string, string) {
	adsMedia := ctx.Request().URL.Query().Get("adsMedia")
	var (
		bs   int
		full string
	)
	switch strings.ToLower(adsMedia) {
	case "banner":
		bs = 8
	case "largebanner":
		bs = 3
	case "xlargebannerportrait":
		bs = 16
	case "fullbannerportrait":
		bs = 16
		full = "portrait"
	case "xlargebannerlandscap":
		bs = 17
	case "fullbannerlandscape":
		bs = 17
		full = "landscape"
	}
	middlewares.SetData(ctx, "ad_size", bs)
	slotString := fmt.Sprintf("%d0%d0%d", app.ID, app.UserID, bs)
	slot, _ := strconv.ParseInt(slotString, 10, 0)
	s, err := mr.NewManager().FetchAppSlot(app.ID, slot)
	if err != nil {
		// no slot found
		s, err = mr.NewManager().InsertSlots(0, app.ID, slot, bs)
		assert.Nil(err)
	}
	data := map[string]*slotData{
		slotString: {
			SlotSize: bs,
			ID:       s.ID,
			PublicID: slotString,
			Ctr:      defaultCTR.Float64(),
		},
	}
	sizes := map[string]int{slotString: bs}

	for i := range data {
		result, err := aredis.SumHMGetField(transport.KeyGenDaily(transport.Slot, strconv.FormatInt(data[i].ID, 10)), dailyClickDays.Int(), "i", "c")
		if err != nil || result["c"] == 0 || result["i"] < minImp.Int64() {
			data[i].Ctr = defaultCTR.Float64()
		} else {
			data[i].Ctr = utils.Ctr(result["i"], result["c"])
		}
	}
	return data, sizes, slotString, full
}

func (tc *selectController) getAppDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.App, int64, int64, *mr.PhoneData, *mr.CellLocation, error) {
	rd := middlewares.MustGetRequestData(c)

	token := c.Request().URL.Query().Get("token")
	if len(token) < 1 {
		return nil, nil, 0, 0, nil, nil, errors.New("invalid request")
	}
	m := mr.NewManager()
	app, err := m.GetApp(token)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, err
	}

	if !app.GetActive() {
		return nil, nil, 0, 0, nil, nil, errors.New("app is disabled")
	}

	if !m.IsUserActive(app.UserID) {
		return nil, nil, 0, 0, nil, nil, errors.New("user is banned")
	}

	province, isp, ll := ip2location.GetProvinceISPByIP(rd.IP)
	middlewares.SetData(c, "province", ll.Province)
	middlewares.SetData(c, "country", ll.Country)
	middlewares.SetData(c, "city", ll.City)
	middlewares.SetData(c, "isp", ll.ISP)

	phone := m.GetPhoneData(c.Request().URL.Query().Get("brand"), strings.Trim(c.Request().URL.Query().Get("carrier"), "# \n\t"), c.Request().URL.Query().Get("network"))
	mcc, _ := strconv.ParseInt(c.Request().URL.Query().Get("mcc"), 10, 0)
	mnc, _ := strconv.ParseInt(c.Request().URL.Query().Get("mnc"), 10, 0)
	cid, _ := strconv.ParseFloat(c.Request().URL.Query().Get("cid"), 64)
	lac, _ := strconv.ParseFloat(c.Request().URL.Query().Get("lac"), 64)
	cell, err := m.GetCellLocation(mcc, mnc, int64(lac), int64(cid), phone.Carrier)
	if err != nil {
		logrus.Debug(err)
	}

	return rd, app, province, isp, phone, cell, nil
}
