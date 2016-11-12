package filter

import (
	"mr"
	"selector"
)

// CheckNetwork filter network for campaigns
func CheckNetwork(c *selector.Context, in mr.AdData) bool {
	if c.Mobile {
		return in.CampaignNetwork == 1
	}
	return in.CampaignNetwork == 0
}
