package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

//CheckProvince find province client in campaign
func CheckProvince(c *selector.Context, in mr.AdData) bool {
	if c.Province == 0 {
		return len(in.CampaignGeos) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignGeos.Has(true, c.Province, 1)
}
