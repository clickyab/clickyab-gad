package routes

import (
	"assert"
	"fmt"
	"middlewares"
	"mr"
	"net/http"
	"net/url"
	"redis"
	"strconv"
	"transport"
	"utils"

	"bytes"
	"encoding/gob"

	"gopkg.in/labstack/echo.v3"
)

func (tc *selectController) nativeShow(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	var suspicious bool
	mega := c.Param("mega")
	native := c.Param("native")
	long := c.Request().URL.Query().Get("l")
	pos := c.Request().URL.Query().Get("pos")

	websiteID, err := strconv.ParseInt(c.Param("wid"), 10, 64)
	website, err := mr.NewManager().FetchWebsite(websiteID)
	assert.Nil(err)

	//TODO :validate Wid compare to us

	assert.Nil(err)
	if err != nil {
		// TODO : check error
		suspicious = true
	}

	// MGA-NATIVE
	megaImp, err := aredis.HGetAllString(transport.MEGA+transport.DELIMITER+mega, false, 0)
	assert.Nil(err)
	// native should be added
	asd, ok := megaImp[fmt.Sprintf("%s%s%s", transport.NATIVE, transport.DELIMITER, native)]
	if !ok {
		return c.String(http.StatusNotFound, "native ads not found")
	}

	// map[AdIDs]winnerBids
	var nativeAds map[int64]int64
	err = fetchToMap(asd, &nativeAds)
	assert.Nil(err)

	keys := make([]int64, 0, len(nativeAds))
	for i := range nativeAds {
		keys = append(keys, i)
	}

	ads, err := mr.NewManager().GetAds(false, keys...)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}

	// map[AdID]random
	rands := make(map[int64]string)

	// map[url]ad
	data := make(map[string]*mr.Ad, len(ads))
	for i := range ads {
		rnd := <-utils.ID
		u := url.URL{
			Scheme: rd.Scheme,
			Host:   rd.Host,
			Path:   fmt.Sprintf("/click/native/%d/%s/%d/%s", websiteID, mega, ads[i].AdID, rnd),
		}
		v := url.Values{}
		v.Set("tid", rd.TID)
		v.Set("ref", rd.Referrer)
		v.Set("parent", rd.Parent)
		u.RawQuery = v.Encode()

		rands[ads[i].AdID] = rnd
		data[u.String()] = ads[i]
	}

	res, err := tc.makeAdData(c, "native", long, pos, data, rd.Scheme != "http")
	if err != nil {
		return err
	}

	for i := range ads {
		slotID, err := strconv.ParseInt(megaImp[fmt.Sprintf("%s%s%d", transport.SLOT, transport.DELIMITER, ads[i].AdID)], 10, 64)
		assert.Nil(err)

		imp := tc.fillImp(rd, suspicious, ads[i], nativeAds[ads[i].AdID], website, slotID)
		go tc.callWebWorker(website, slotID, ads[i].AdID, mega, rands[ads[i].AdID], imp, rd)
	}

	return c.HTML(http.StatusOK, res)
}

func fetchToMap(in string, out *map[int64]int64) error {
	buf := bytes.NewBuffer([]byte(in))
	dc := gob.NewDecoder(buf)
	return dc.Decode(out)
}
