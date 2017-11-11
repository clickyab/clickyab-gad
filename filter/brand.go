package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// CheckAppBrand return boolean
func CheckAppBrand(c *selector.Context, in mr.AdData) bool {
	return in.Campaign.CampaignAppBrand.Has(true, c.PhoneData.BrandID)
}
