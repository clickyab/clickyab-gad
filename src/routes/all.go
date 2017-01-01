package routes

import (
	"bytes"
	"filter"
	"fmt"
	"math/rand"
	"middlewares"
	"mr"
	"selector"
	"strconv"
	"strings"

	"config"

	"net/http"

	"github.com/Sirupsen/logrus"
	"gopkg.in/labstack/echo.v3"
)

type AllData struct {
	Website  []*mr.Website
	Province []*mr.Province
	//Campaign *[]mr.Campaign
	Size map[string]int
	Vast bool
	Data []*mr.AdData
	Len  int
}

func (tc *selectController) allAds(c echo.Context) error {
	w := c.QueryParam("w")
	p := c.QueryParam("p")
	v := c.QueryParam("v")
	cam := c.QueryParam("cam")
	s := c.QueryParam("s")

	//cat := c.QueryParam("cat")
	var campaign int64
	var province mr.Province
	var fltr = []selector.FilterFunc{filter.IsWebNetwork}

	var sizeNumSlice = make(map[string]int)
	var website *mr.Website
	var err error
	var vv bool
	rd := middlewares.MustGetRequestData(c)

	if v != "" || v == "on" {
		fltr = append(fltr, filter.CheckVastSize)
		vv = true
	} else {
		if s != "" {
			ss := strings.Split(s, ",")
			var strin string
			for _, sss := range ss {
				size, err := strconv.Atoi(sss)
				if err == nil {
					strin = fmt.Sprintf("1jhgy%d", rand.Intn(200))
					sizeNumSlice[strin] = size
				}
			}
			if len(sizeNumSlice) > 0 {
				fltr = append(fltr, filter.CheckWebSize)
			}
		}
	}
	if w != "" {
		website, err = mr.NewManager().FetchWebsiteByDomain(w)
		if err == nil {
			fltr = append(fltr, filter.CheckWhiteList, filter.CheckBlackList)
		}
	}
	if cam != "" {

		campaign, err = strconv.ParseInt(s, 10, 0)
		if err == nil {
			fltr = append(fltr, filter.CheckCampaign)
		}
	} else {
		campaign = 0
	}
	if p != "" {
		i64, err := strconv.ParseInt(s, 10, 0)
		if err == nil {
			province, err = mr.NewManager().ConvertProvinceID2Info(i64)
			if err == nil {
				fltr = append(fltr, filter.CheckProvince)
			}
		}

	}
	m := selector.Context{
		RequestData: *rd,
		Website:     website,
		Size:        sizeNumSlice,
		Province:    &province,
		Campaign:    campaign,
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), selector.Mix(fltr...))

	all := make([]*mr.AdData, 0)
	for i := range filteredAds {
		all = append(all, filteredAds[i]...)
	}
	al := allDate()
	al.Vast = vv
	al.Data = all
	al.Len = len(all)

	buf := &bytes.Buffer{}
	err = allAdTemplate.Execute(buf, al)
	logrus.Info(err)
	return c.HTML(http.StatusOK, buf.String())

	return c.JSON(200, struct {
		Count int
		All   []*mr.AdData
	}{
		Count: len(all),
		All:   all,
	})
}
func allDate() AllData {
	/*c, err := mr.NewManager().FetchCampaignAll()
	if err != nil {
		c = nil
	}*/
	p, err := mr.NewManager().FetchProvinceAll()
	if err != nil {
		p = nil
	}
	w, err := mr.NewManager().FetchWebsiteAll()
	if err != nil {
		w = nil
	}
	s := config.GetAllSize()
	al := AllData{
		//Campaign: c,
		Province: p,
		Website:  w,
		Size:     s,
	}
	return al
}
