package filter

import (
	"mr"
	"selector"
)

// CheckWebCategory is the filter for category
func CheckWebCategory(c *selector.Context, in mr.AdData) bool {
	return in.CampaignCat.Match(true, c.Website.WCategories)
}

// CheckAppCategory is the filter for category
func CheckAppCategory(c *selector.Context, in mr.AdData) bool {
	return in.CampaignCat.Match(true, c.App.Appcat)
}
