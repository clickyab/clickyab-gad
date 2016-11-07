package filter

import (
	"mr"
	"selector"
)

// CheckWhiteList return boolean
func CheckWhiteList(c *selector.Context, in mr.AdData) bool {
	if len(in.CampaignPlacement) == 0 {
		return true
	}
	return in.CampaignPlacement.Has(c.WID)
}
