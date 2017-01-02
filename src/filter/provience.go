package filter

import (
	"mr"
	"selector"
)

//CheckProvince find province client in campaign
func CheckProvince(c *selector.Context, in mr.AdData) bool {
	if c.Province == nil {
		return len(in.CampaignRegion) == 0
	}
	return in.CampaignRegion.Has(true, c.Province.ID)
}
