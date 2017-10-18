package filter

import (
	"mr"
	"selector"
)

// CheckAppCarrier return boolean
func CheckAppCarrier(c *selector.Context, in mr.AdData) bool {
	return in.Campaign.CampaignAppsCarriers.Has(true, c.PhoneData.CarrierID)
}
