package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// CheckWhiteList return boolean
func CheckWhiteList(c *selector.Context, in mr.AdData) bool {
	return in.CampaignPlacement.Has(true, c.Website.WID)
}

// CheckAppWhiteList return boolean
func CheckAppWhiteList(c *selector.Context, in mr.AdData) bool {
	return in.CampaignApp.Has(true, c.App.ID)
}
