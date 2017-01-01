package routes

import (
	"selector"

	"filter"

	"middlewares"
	"mr"

	"strconv"

	"strings"

	"math/rand"

	"fmt"

	"gopkg.in/labstack/echo.v3"
)

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
	rd := middlewares.MustGetRequestData(c)

	if v != "" || v == "true" {
		fltr = append(fltr, filter.CheckVastSize)
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
			fltr = append(fltr, filter.CheckWebSize)
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
		province, err = mr.NewManager().ConvertProvince2Info(p)
		if err == nil {
			fltr = append(fltr, filter.CheckProvince)
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

	return c.JSON(200, struct {
		Count int
		All   []*mr.AdData
	}{
		Count: len(all),
		All:   all,
	})
}
