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

// IsWebMobile return if the campaign is ok for web mobile
func IsWebMobile(c *selector.Context, in mr.AdData) bool {
	if c.Mobile {
		return in.CampaignWebMobile == 1
	}

	return true
}
