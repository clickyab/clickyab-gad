package routes

import (
	"clickyab.com/gad/config"
	"fmt"
	"clickyab.com/gad/mr"
	aredis "clickyab.com/gad/redis"
	"sort"
	"time"
	"clickyab.com/gad/transport"

	"github.com/sirupsen/logrus"
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

func emptyCapping(filteredAds map[int][]*mr.AdData) map[int][]*mr.AdData {
	c := make(mr.CappingContext)
	for i := range filteredAds {
		for j := range filteredAds[i] {
			capp := c.NewCapping(
				filteredAds[i][j].CampaignID,
				0,
				filteredAds[i][j].CampaignFrequency,
			)
			filteredAds[i][j].Capping = capp

		}
		sortCap := mr.ByCapping(filteredAds[i])
		sort.Sort(sortCap)
		filteredAds[i] = []*mr.AdData(sortCap)
	}

	return filteredAds
}

func getCapping(copID int64, sizeNumSlice map[string]int, filteredAds map[int][]*mr.AdData, eventPage string) map[int][]*mr.AdData {
	var selected = make(map[int64]bool)
	if eventPage != "" {
		for _, v := range aredis.SMembersInt(eventPage) {
			selected[v] = true
		}

	}
	c := make(mr.CappingContext)
	results, _ := aredis.HGetAll(getCappingKey(copID), true, config.Config.Clickyab.DailyCapExpire)
	doneSized := make(map[int]bool)
	for i := range sizeNumSlice {
		if doneSized[sizeNumSlice[i]] {
			continue
		}
		doneSized[sizeNumSlice[i]] = true
		found := false
		sizeCap := map[string]string{}
		for ad := range filteredAds[sizeNumSlice[i]] {
			if filteredAds[sizeNumSlice[i]][ad].CampaignFrequency <= 0 {
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency = config.Config.Clickyab.MinFrequency
			}
			key := fmt.Sprintf(
				"%s%s%d",
				transport.ADVERTISE,
				transport.DELIMITER,
				filteredAds[sizeNumSlice[i]][ad].AdID,
			)
			view := results[key]
			sizeCap[key] = "0"
			n := view / filteredAds[sizeNumSlice[i]][ad].CampaignFrequency
			if n <= 1 {
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
			capp := c.NewCapping(
				filteredAds[sizeNumSlice[i]][ad].CampaignID,
				0,
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency,
			)
			capp.IncView(filteredAds[sizeNumSlice[i]][ad].AdID, view, selected[filteredAds[sizeNumSlice[i]][ad].AdID])
			filteredAds[sizeNumSlice[i]][ad].Capping = capp
		}
		//sortCap := mr.ByCapping(filteredAds[sizeNumSlice[i]])
		//sort.Sort(sortCap)
		//filteredAds[sizeNumSlice[i]] = []*mr.AdData(sortCap)
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
