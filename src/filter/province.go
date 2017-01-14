package filter

import (
	"mr"
	"selector"
)

//CheckProvince find province client in campaign
func CheckProvince(c *selector.Context, in mr.AdData) bool {
	if c.Province == nil {
		return len(in.CampaignGeos) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignGeos.Has(true, c.Province.ID, 1)
}
