package filter

import (
	"selector"
	"mr"
)

// CheckAppBrand return boolean
func CheckAppHood(c *selector.Context, in mr.AdData)bool{
	return in.Campaign.CampaignHoods.Has(true,c.CellLocation.NeighborhoodsID)
}
