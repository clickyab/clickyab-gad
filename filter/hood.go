package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

// CheckAppHood return boolean
func CheckAppHood(c *selector.Context, in models.AdData) bool {
	if c.CellLocation == nil {
		return in.Campaign.CampaignHoods == ""
	}
	return in.Campaign.CampaignHoods.Has(true, c.CellLocation.NeighborhoodsID)
}
