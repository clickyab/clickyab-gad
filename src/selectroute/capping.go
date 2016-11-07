package selectroute

import (
	"config"
	"fmt"
	"mr"
	"sort"
	"time"
	"transport"

	"redis"

	"github.com/labstack/echo"
)

func getCappingKey(copID string) string {
	return fmt.Sprintf(
		"%s%s%s%s%s",
		transport.USER_CAPPING,
		transport.DELIMITER,
		copID,
		transport.DELIMITER,
		time.Now().Format("060102"),
	)
}

func getCapping(c echo.Context, copID string, sizeNumSlice []int, filteredAds map[int][]*mr.MinAdData) map[int][]*mr.MinAdData {
	var userMinView int

	results, _ := aredis.HGetAll(getCappingKey(copID), true, 72*time.Hour)
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
			if userMinView == 0 {
				userMinView = view
			} else if view > 0 && userMinView > view {
				userMinView = view
			}
		}
		sortCap := mr.ByCapping(filteredAds[sizeNumSlice[i]])
		sort.Sort(sortCap)
		filteredAds[sizeNumSlice[i]] = []*mr.MinAdData(sortCap)
	}
	return filteredAds
}

func storeCapping(copID string, cpID int64) error {
	_, err := aredis.IncHash(
		getCappingKey(copID),
		fmt.Sprintf("%s%s%d", transport.CAMPAIGN, transport.DELIMITER, cpID),
		1,
		true,
		config.Config.Clickyab.DailyCapExpireTime,
	)
	return err
}
