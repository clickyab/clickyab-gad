package routes

import (
	"config"
	"fmt"
	"mr"
	"sort"
	"time"
	"transport"

	"redis"

	"gopkg.in/labstack/echo.v3"
)

func getCappingKey(copID int64) string {
	return fmt.Sprintf(
		"%s%s%d%s%s",
		transport.USER_CAPPING,
		transport.DELIMITER,
		copID,
		transport.DELIMITER,
		time.Now().Format("060102"),
	)
}

func getCapping(c echo.Context, copID int64, sizeNumSlice map[string]int, filteredAds map[int][]*mr.AdData) map[int][]*mr.AdData {

	results, _ := aredis.HGetAll(getCappingKey(copID), true, config.Config.Clickyab.DailyCapExpire)
	for i := range sizeNumSlice {
		for ad := range filteredAds[sizeNumSlice[i]] {
			view := results[fmt.Sprintf(
				"%s%s%d",
				transport.CAMPAIGN,
				transport.DELIMITER,
				filteredAds[sizeNumSlice[i]][ad].CampaignID,
			)]

			if filteredAds[sizeNumSlice[i]][ad].CampaignFrequency <= 0 {
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency = config.Config.Clickyab.MinFrequency
			}
			filteredAds[sizeNumSlice[i]][ad].Capping = mr.NewCapping(
				c,
				filteredAds[sizeNumSlice[i]][ad].CampaignID,
				view,
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency,
			)
		}
		sortCap := mr.ByCapping(filteredAds[sizeNumSlice[i]])
		sort.Sort(sortCap)
		filteredAds[sizeNumSlice[i]] = []*mr.AdData(sortCap)
	}
	return filteredAds
}

func storeCapping(copID int64, cpID int64) error {
	_, err := aredis.IncHash(
		getCappingKey(copID),
		fmt.Sprintf("%s%s%d", transport.CAMPAIGN, transport.DELIMITER, cpID),
		1,
		config.Config.Clickyab.DailyCapExpire,
	)
	return err
}
