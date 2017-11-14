package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

//CheckCampaign find campaign
func CheckCampaign(c *selector.Context, in models.AdData) bool {
	return c.Campaign == 0 || in.CampaignID == c.Campaign
}
