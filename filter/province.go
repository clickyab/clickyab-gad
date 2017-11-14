package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

//CheckProvince find province client in campaign
func CheckProvince(c *selector.Context, in models.AdData) bool {
	if c.Province == 0 {
		return len(in.CampaignGeos) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignGeos.Has(true, c.Province, 1)
}
