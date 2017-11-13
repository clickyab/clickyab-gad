package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

//CheckProvder find provider client in campaign
func CheckProvder(c *selector.Context, in models.AdData) bool {
	if c.PhoneData.NetworkID == 0 {
		return len(in.CampaignNetProvider) == 0
	}
	return in.CampaignNetProvider.Has(true, c.PhoneData.NetworkID)
}
