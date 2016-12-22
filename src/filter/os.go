package filter

import (
	"mr"
	"selector"
)

// CheckOS is the filter function that check for OS in system
func CheckOS(c *selector.Context, in mr.MinAdData) bool {
	return in.CampaignPlatforms.Has(true, c.RequestData.PlatformID)
}
