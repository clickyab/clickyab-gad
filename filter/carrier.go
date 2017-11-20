package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

// CheckAppCarrier return boolean
func CheckAppCarrier(c *selector.Context, in models.AdData) bool {
	return in.Campaign.CampaignAppsCarriers.Has(true, c.PhoneData.CarrierID)
}
