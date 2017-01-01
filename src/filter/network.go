package filter

import (
	"mr"
	"selector"
)

// IsWebNetwork filter network for campaigns
func IsWebNetwork(c *selector.Context, in mr.AdData) bool {
	return in.CampaignNetwork == 0 || in.CampaignNetwork == 2
}