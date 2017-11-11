package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

//CheckCampaign find campaign
func CheckCampaign(c *selector.Context, in mr.AdData) bool {
	return c.Campaign == 0 || in.CampaignID == c.Campaign
}
