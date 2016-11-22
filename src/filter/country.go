package filter

import (
	"mr"
	"selector"
)

//CheckCountry find country client in campaign
func CheckCountry(c *selector.Context, in mr.AdData) bool {
	if c.Country == nil {
		return len(in.CampaignCountry)==0
	}
	return in.CampaignCountry.Has(true, c.Country.ID)
}
