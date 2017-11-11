package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// CheckAppBlackList filter blacklist
func CheckAppBlackList(c *selector.Context, in mr.AdData) bool {
	return !in.CampaignAppFilter.Has(false, c.App.ID)
}
