package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

// CheckDesktopNetwork filter network for desktop
func CheckDesktopNetwork(c *selector.Context, in models.AdData) bool {
	if in.CampaignWeb == 1 {
		if in.CampaignWebMobile == 0 {
			return !c.Mobile
		}
	} else if in.CampaignWeb == 0 {
		if in.CampaignWebMobile == 1 {
			return c.Mobile
		}
	}
	return true
}
