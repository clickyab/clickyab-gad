package filter

import (
	"math"
	"mr"
	"selector"
)

//CheckMinBid find isp
func CheckMinBid(c *selector.Context, in mr.AdData) bool {
	if c.Website == nil {
		return true
	}
	t := c.Website.WMinBid
	if c.MinBidPercentage > 0 {
		t = int64(math.Ceil(c.MinBidPercentage * float64(t)))
	}

	return in.CampaignMaxBid >= t
}
