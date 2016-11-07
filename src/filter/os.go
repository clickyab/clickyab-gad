package filter

import (
	"mr"
	"selector"
)

// CheckOS is the filter function that check for OS in system
func CheckOS(c *selector.Context, in mr.AdData) bool {
	if len(in.CampaignPlatforms) == 0 {
		return true
	}
	return in.CampaignPlatforms.Has(c.RequestData.PlatformID)
}
