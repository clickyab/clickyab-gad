package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

// CheckWhiteList return boolean
func CheckWhiteList(c *selector.Context, in models.AdData) bool {
	return in.CampaignPlacement.Has(true, c.Website.WID)
}

// CheckAppWhiteList return boolean
func CheckAppWhiteList(c *selector.Context, in models.AdData) bool {
	return in.CampaignApp.Has(true, c.App.ID)
}
