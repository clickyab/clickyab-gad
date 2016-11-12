package filter

import (
	"mr"
	"selector"
)

// CheckCategory is the filter for category
func CheckCategory(c *selector.Context, in mr.AdData) bool {

	lenC := len(in.CampaignCat)
	lenW := len(c.Website.WCategories)

	if lenC == 0 {
		return true
	}

	//compare two slice
	if lenC >= lenW {
		for _, WCat := range c.Website.WCategories {
			for _, cCat := range in.CampaignCat {
				if cCat == WCat {
					return true
				}
			}
		}
	} else {
		for _, cCat := range in.CampaignCat {
			for _, WCat := range c.Website.WCategories {
				if cCat == WCat {
					return true
				}
			}
		}
	}

	return false
}
