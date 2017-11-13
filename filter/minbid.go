package filter

import (
	"math"

	"clickyab.com/gad/models"
	"clickyab.com/gad/selector"
)

//CheckMinBid find isp
func CheckMinBid(c *selector.Context, in models.AdData) bool {
	if c.Website == nil {
		return true
	}
	t := c.Website.WMinBid
	if c.MinBidPercentage > 0 {
		t = int64(math.Ceil(c.MinBidPercentage * float64(t)))
	}

	return in.CampaignMaxBid >= t
}
