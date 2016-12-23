package filter

import (
	"mr"
	"selector"
)

// CheckCategory is the filter for category
func CheckCategory(c *selector.Context, in mr.AdData) bool {
	return in.CampaignCat.Match(true, c.Website.WCategories)
}
