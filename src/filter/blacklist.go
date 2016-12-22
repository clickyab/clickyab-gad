package filter

import (
	"mr"
	"selector"
)

// CheckBlackList filter blacklist
func CheckBlackList(c *selector.Context, in mr.MinAdData) bool {
	return !in.CampaignWebsiteFilter.Has(false, c.Website.WID)
}
