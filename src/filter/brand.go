package filter

import (
	"selector"
	"mr"
)

// CheckAppBrand return boolean
func CheckAppBrand(c *selector.Context, in mr.AdData)bool{
	return in.Campaign.CampaignAppBrand.Has(true,c.PhoneData.BrandID)
}