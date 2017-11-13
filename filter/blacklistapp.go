package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

// CheckAppBlackList filter blacklist
func CheckAppBlackList(c *selector.Context, in models.AdData) bool {
	return !in.CampaignAppFilter.Has(false, c.App.ID)
}
