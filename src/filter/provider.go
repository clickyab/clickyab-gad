package filter

import (
	"mr"
	"selector"
)

//CheckProvder find provider client in campaign
func CheckProvder(c *selector.Context, in mr.AdData) bool {
	if c.PhoneData.NetworkID == nil {
		return len(in.CampaignNetProvider) == 0
	}
	return in.CampaignNetProvider.Has(true, c.PhoneData.NetworkID)
}
