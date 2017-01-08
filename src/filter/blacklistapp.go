package filter

import (
	"mr"
	"selector"
)

// CheckAppBlackList filter blacklist
func CheckAppBlackList(c *selector.Context, in mr.AdData) bool {
	return !in.CampaignAppFilter.Has(false, c.App.ID)
}
