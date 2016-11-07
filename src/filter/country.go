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
	return in.CampaignCountry.Has(c.Country2Info.ID)
}
