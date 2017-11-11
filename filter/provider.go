package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

//CheckProvder find provider client in campaign
func CheckProvder(c *selector.Context, in mr.AdData) bool {
	if c.PhoneData.NetworkID == 0 {
		return len(in.CampaignNetProvider) == 0
	}
	return in.CampaignNetProvider.Has(true, c.PhoneData.NetworkID)
}
