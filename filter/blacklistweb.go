package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

// CheckWebBlackList filter blacklist
func CheckWebBlackList(c *selector.Context, in models.AdData) bool {
	return !in.CampaignWebsiteFilter.Has(false, c.Website.WID)
}
