package filter

import (
	"mr"
	"selector"
	"strings"
)

// ExcludeCPM filter to exclude cpm campaigns
func ExcludeCPM(c *selector.Context, in mr.AdData) bool {
	if in.Campaign.CampaignBillingType.Valid{
		return strings.ToLower(in.Campaign.CampaignBillingType.String)=="cpc"
	}
	return false
}
