package routes

import (
	"selector"

	"filter"

	"config"
	"middlewares"
	"mr"

	"gopkg.in/labstack/echo.v3"
)

var (
	currentSelector = selector.Mix(
		filter.CheckWhiteList,
		filter.CheckBlackList,
		filter.IsWebNetwork,
		filter.CheckCategory,
	)

	currentVastSelector = selector.Mix(
		filter.CheckVastSize,
		filter.CheckWhiteList,
		filter.CheckBlackList,
		filter.IsWebNetwork,
		filter.CheckCategory,
	)
)

func (tc *selectController) allAds(c echo.Context) error {
	ws := c.Request().URL.Query().Get("ws")
	if ws == "" {
		return c.JSON(200, selector.GetAdData())
	}

	vast := c.Request().URL.Query().Get("vast") != ""
	website, err := mr.NewManager().FetchWebsiteByDomain(ws)
	if err != nil {
		// return all of them
		return c.String(400, "no website")
	}
	m := selector.Context{
		RequestData: *middlewares.MustGetRequestData(c),
		Website:     website,
		Size:        config.GetAllSize(),
		Country:     nil,
	}
	fltr := currentSelector
	if vast {
		fltr = currentVastSelector
	}
	filteredAds := selector.Apply(&m, selector.GetAdData(), fltr)

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
