package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

// CheckAppBrand return boolean
func CheckAppBrand(c *selector.Context, in models.AdData) bool {
	return in.Campaign.CampaignAppBrand.Has(true, c.PhoneData.BrandID)
}
