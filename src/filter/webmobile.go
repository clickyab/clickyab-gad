package filter

import (
	"mr"
	"selector"
)

// IsNotWebMobile filter for webmobile
func IsNotWebMobile(c *selector.Context, in mr.AdData) bool {
	if c.Mobile {
		return true
	}
	return in.CampaignWebMobile == 0
}

// IsWebMobile return if the campaign is ok for web mobile
func IsWebMobile(c *selector.Context, in mr.AdData) bool {
	if c.Mobile {
		return in.CampaignWebMobile == 1
	}

	return true
}
