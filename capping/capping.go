package capping

import (
	"fmt"
	"sort"
	"time"

	"clickyab.com/gad/models"
	"clickyab.com/gad/redis"
	"clickyab.com/gad/transport"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/sirupsen/logrus"
)

var (
	minFrequency   = config.RegisterInt("clickyab.min_frequency", 2, "")
	dailyCapExpire = config.RegisterDuration("clickyab.daily_cap_expire", 72*time.Hour, "")
)

func getCappingKey(copID int64) string {
	return fmt.Sprintf(
		"%s%s%d%s%s",
		transport.CappingKey,
		transport.Delimiter,
		copID,
		transport.Delimiter,
		time.Now().Format("060102"),
	)
}

// EmptyCapping is a hack to handle no capping situation
func EmptyCapping(filteredAds map[int][]*models.AdData) map[int][]*models.AdData {
	c := make(models.CappingContext)
	for i := range filteredAds {
		for j := range filteredAds[i] {
			capp := c.NewCapping(
				filteredAds[i][j].CampaignID,
				0,
				filteredAds[i][j].CampaignFrequency,
			)
			filteredAds[i][j].Capping = capp

		}
		sortCap := models.ByCapping(filteredAds[i])
		sort.Sort(sortCap)
		filteredAds[i] = []*models.AdData(sortCap)
	}

	return filteredAds
}

// GetCapping try to get capping for current ad
func GetCapping(copID int64, sizeNumSlice map[string]int, filteredAds map[int][]*models.AdData, eventPage string) map[int][]*models.AdData {
	var selected = make(map[int64]bool)
	if eventPage != "" {
		for _, v := range aredis.SMembersInt(eventPage) {
			selected[v] = true
		}

	}
	c := make(models.CappingContext)
	results, _ := aredis.HGetAll(getCappingKey(copID), true, dailyCapExpire.Duration())
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
				filteredAds[sizeNumSlice[i]][ad].CampaignFrequency = minFrequency.Int()
			}
			key := fmt.Sprintf(
				"%s%s%d",
				transport.Advertise,
				transport.Delimiter,
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
			if len(sizeCap) != 0 {
				assert.Nil(aredis.HMSet(getCappingKey(copID), dailyCapExpire.Duration(), sizeCap))
			}
			for i := range sizeCap {
				results[i] = 0
			}
		}
		for ad := range filteredAds[sizeNumSlice[i]] {
			view := 0
			if found {
				view = results[fmt.Sprintf(
					"%s%s%d",
					transport.Advertise,
					transport.Delimiter,
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
		//sortCap := models.ByCapping(filteredAds[sizeNumSlice[i]])
		//sort.Sort(sortCap)
		//filteredAds[sizeNumSlice[i]] = []*models.AdData(sortCap)
	}
	return filteredAds
}

// StoreCapping try to store a capping object
func StoreCapping(copID int64, cpID int64) error {
	_, err := aredis.IncHash(
		getCappingKey(copID),
		fmt.Sprintf("%s%s%d", transport.Advertise, transport.Delimiter, cpID),
		1,
		dailyCapExpire.Duration(),
	)
	return err
}
