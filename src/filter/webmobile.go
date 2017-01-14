package filter

import (
	"mr"
	"selector"
)

// IsWebMobile filter for webmobile
func IsWebMobile(c *selector.Context, in mr.AdData) bool {
	return in.CampaignWebMobile == 1
}

// IsNotWebMobile filter for webmobile
func IsNotWebMobile(c *selector.Context, in mr.AdData) bool {
	return in.CampaignWebMobile == 0
}

