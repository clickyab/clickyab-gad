package filter

import (
	"mr"
	"selector"
)

// CheckWhiteList return boolean
func CheckWhiteList(c *selector.Context, in mr.AdData) bool {
	return in.CampaignPlacement.Has(true, c.Website.WID)
}
