package routes

import (
	"errors"
	"filter"
	"middlewares"
	"mr"
	"selector"
	"strconv"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

var (
	appSelector = selector.Mix(
		filter.IsAppNetwork,
		filter.CheckWebSize,
		filter.CheckOS,
		filter.CheckWhiteList,
		filter.CheckWebBlackList,
		filter.CheckCategory,
		filter.CheckProvince,
	)
)

func (tc *selectController) inApp(c echo.Context) error {
	// t := time.Now()
	// params := c.QueryParams()
	// rd, app, province, phone, cell, err := tc.getAppDataFromCtx(c)
	// if err != nil {
	// 	return c.HTML(http.StatusBadRequest, err.Error())
	// }
	// slotSize, sizeNumSlice := tc.slotSizeWeb(params, *website, rd.Mobile)
	// //call context
	// m := selector.Context{
	// 	RequestData:  *rd,
	// 	Website:      nil,
	// 	Size:         sizeNumSlice,
	// 	Province:     province,
	// 	App:          app,
	// 	PhoneData:    phone,
	// 	CellLocation: cell,
	// }
	// filteredAds := selector.Apply(&m, selector.GetAdData(), webSelector)
	// show := tc.makeShow(c, "web", rd, filteredAds, sizeNumSlice, slotSize, website, false)
	return nil
}

func (tc *selectController) slotSizeApp() {

}

func (tc *selectController) getAppDataFromCtx(c echo.Context) (*middlewares.RequestData, *mr.App, *mr.Province, *mr.PhoneData, *mr.CellLocation, error) {
	rd := middlewares.MustGetRequestData(c)

	token := c.Param("token")
	if len(token) < 1 {
		return nil, nil, nil, nil, nil, errors.New("invalid request")
	}
	m := mr.NewManager()
	app, err := m.GetApp(token)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	if !m.IsUserActive(app.UserID) {
		return nil, nil, nil, nil, nil, errors.New("user is banned")
	}

	province, err := tc.fetchProvince(rd.IP, c.Request().Header.Get("Cf-Ipcountry"))
	if err != nil {
		logrus.Debug(err)
	}

	phone, err := m.GetPhoneData(c.Param("brand"), c.Param("carrier"), c.Param("network"))
	if err != nil {
		logrus.Warn(err)
	}
	mcc, _ := strconv.ParseInt(c.Param("mcc"), 10, 0)
	mnc, _ := strconv.ParseInt(c.Param("mnc"), 10, 0)
	cid, _ := strconv.ParseInt(c.Param("cid"), 10, 0)
	lac, _ := strconv.ParseInt(c.Param("lac"), 10, 0)
	cell, err := m.GetCellLocation(mcc, mnc, lac, cid, phone.Carrier)
	if err != nil {
		logrus.Warn(err)
	}

	return rd, app, province, phone, cell, nil
}
