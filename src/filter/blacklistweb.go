package filter

import (
	"mr"
	"selector"
)

// CheckWebBlackList filter blacklist
func CheckWebBlackList(c *selector.Context, in mr.AdData) bool {
	return !in.CampaignWebsiteFilter.Has(false, c.Website.WID)
}
