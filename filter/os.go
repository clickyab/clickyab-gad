package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// CheckOS is the filter function that check for OS in system
func CheckOS(c *selector.Context, in mr.AdData) bool {
	return in.CampaignPlatforms.Has(true, c.RequestData.PlatformID)
}
