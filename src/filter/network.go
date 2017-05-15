package filter

import (
	"mr"
	"selector"
)

// IsWebNetwork filter network for campaigns
func IsWebNetwork(c *selector.Context, in mr.AdData) bool {
	return in.CampaignNetwork == 0 || in.CampaignNetwork == 2
}

// IsAppNetwork filter network for campaigns
func IsAppNetwork(c *selector.Context, in mr.AdData) bool {
	return in.CampaignNetwork == 1
}

// IsNativeNetwork filter network for native
func IsNativeNetwork(c *selector.Context, in mr.AdData) bool {
	return in.CampaignNetwork == 3
}
