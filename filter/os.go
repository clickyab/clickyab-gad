package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

// CheckOS is the filter function that check for OS in system
func CheckOS(c *selector.Context, in models.AdData) bool {
	return in.CampaignPlatforms.Has(true, c.RequestData.PlatformID)
}
