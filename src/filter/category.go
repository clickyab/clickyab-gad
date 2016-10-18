package filter

import (
	"mr"
	"selector"
)

// CheckCategory is the filter for category
func CheckCategory(c *selector.Context, in mr.AdData) bool {

	lenC := len(in.CpCat)
	lenW := len(c.WCategories)

	if lenC == 0 {
		return true
	}

	//compare two slice
	if lenC >= lenW {
		for _, WCat := range c.WCategories {
			for _, cCat := range in.CpCat {
				if cCat == WCat {
					return true
				}
			}
		}
	} else {
		for _, cCat := range in.CpCat {
			for _, WCat := range c.WCategories {
				if cCat == WCat {
					return true
				}
			}
		}
	}

	return false
}
