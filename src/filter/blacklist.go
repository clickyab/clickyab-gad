package filter

import (
	"mr"
	"selector"
)

// CheckBlackList filter blacklist
func CheckBlackList(c *selector.Context, in mr.AdData) bool {
	if len(in.CampaignWebsiteFilter) == 0 {
		return true
	}
	for _, v := range in.CampaignWebsiteFilter {
		if v == c.Website.WID {
			return false
		}
	}
	return false
}
