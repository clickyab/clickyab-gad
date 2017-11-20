package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

// CheckWebCategory is the filter for category
func CheckWebCategory(c *selector.Context, in models.AdData) bool {
	return in.CampaignCat.Match(true, c.Website.WCategories)
}

// CheckAppCategory is the filter for category
func CheckAppCategory(c *selector.Context, in models.AdData) bool {
	return in.CampaignCat.Match(true, c.App.Appcat)
}
