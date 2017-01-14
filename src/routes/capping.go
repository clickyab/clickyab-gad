package routes

import (
	"config"
	"fmt"
	"mr"
	aredis "redis"
	"sort"
	"time"
	"transport"

	"strconv"

	"github.com/Sirupsen/logrus"
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

func retargetingKey(copID int64) string {
	return fmt.Sprintf(
		"%s%s%d",
		transport.USER_RETARGETING,
		transport.DELIMITER,
		copID,
	)
}

func getCapping(copID int64, sizeNumSlice map[string]int, filteredAds map[int][]*mr.AdData) map[int][]*mr.AdData {
	// Retargeting structure is like this :
	/*
		map[string]int {
			"campaign_id" : unix_time_of_retargeting
			"cp_2" : utime_2
		}

	*/
	c := make(mr.CappingContext)
	retargetings, _ := aredis.HGetAll(retargetingKey(copID), false, 0)
	if retargetings == nil {
		retargetings = make(map[string]int)
	}
	results, _ := aredis.HGetAll(getCappingKey(copID), true, config.Config.Clickyab.DailyCapExpire)
	for i := range sizeNumSlice {
		found := false
		sizeCap := map[string]string{}
		for ad := range filteredAds[sizeNumSlice[i]] {
			key := fmt.Sprintf(
				"%s%s%d",
				transport.ADVERTISE,
				transport.DELIMITER,
				filteredAds[sizeNumSlice[i]][ad].AdID,
			)
			view := results[key]
			sizeCap[key] = "0"
			n := view / filteredAds[sizeNumSlice[i]][ad].CampaignFrequency
			if n < 1 {
				found = true
				break // there is still one campaign
			}
		}
		// if not found then reset all capping
		if !found {
			logrus.Debugf("Removing key for size %d", sizeNumSlice[i])
			aredis.HMSet(getCappingKey(copID), config.Config.Clickyab.DailyCapExpire, sizeCap)
			for i := range sizeCap {
				results[i] = 0
			}
		}
		for ad := range filteredAds[sizeNumSlice[i]] {
			view := 0
			if found {
				view = results[fmt.Sprintf(
					"%s%s%d",
					transport.ADVERTISE,
					transport.DELIMITER,
					filteredAds[sizeNumSlice[i]][ad].AdID,
				)]
			}
			if filteredAds[sizeNumSlice[i]][ad].CampaignFrequency <= 0 {
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency = config.Config.Clickyab.MinFrequency
			}
			retarget := false
			if v, ok := retargetings[strconv.FormatInt(filteredAds[sizeNumSlice[i]][ad].CampaignID, 10)]; ok {
				if time.Since(time.Unix(int64(v), 0)) < time.Duration(config.Config.Clickyab.RetargettingHour)*time.Hour {
					retarget = true
				}
			}
			capp := c.NewCapping(
				filteredAds[sizeNumSlice[i]][ad].CampaignID,
				0,
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency,
				retarget,
			)
			capp.IncView(filteredAds[sizeNumSlice[i]][ad].AdID, view, false)
			filteredAds[sizeNumSlice[i]][ad].Capping = capp
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
		fmt.Sprintf("%s%s%d", transport.ADVERTISE, transport.DELIMITER, cpID),
		1,
		config.Config.Clickyab.DailyCapExpire,
	)
	return err
}
