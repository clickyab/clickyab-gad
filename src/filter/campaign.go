package filter

import (
	"mr"
	"selector"
)

//CheckCampaign find campaign
func CheckCampaign(c *selector.Context, in mr.AdData) bool {
	return c.Campaign == 0 || in.CampaignID == c.Campaign
}
