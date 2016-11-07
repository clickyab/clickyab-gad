package filter

import (
	"mr"
	"selector"
)

//CheckCountry find country client in campaign
func CheckCountry(c *selector.Context, in mr.AdData) bool {
	if len(in.CampaignCountry) == 0 {
		return true
	}
	if c.Country == nil {
		return false
	}
	return in.CampaignCountry.Has(c.Country.ID)
}
