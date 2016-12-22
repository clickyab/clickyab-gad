package filter

import (
	"mr"
	"selector"
)

// IsWebNetwork filter network for campaigns
func IsWebNetwork(c *selector.Context, in mr.MinAdData) bool {
	return in.CampaignNetwork == 0
}
