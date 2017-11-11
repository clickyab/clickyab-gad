package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// CheckAppHood return boolean
func CheckAppHood(c *selector.Context, in mr.AdData) bool {
	if c.CellLocation == nil {
		return in.Campaign.CampaignHoods == ""
	}
	return in.Campaign.CampaignHoods.Has(true, c.CellLocation.NeighborhoodsID)
}
