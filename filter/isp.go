package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

//CheckISP find isp
func CheckISP(c *selector.Context, in models.AdData) bool {
	if c.ISP == 0 {
		return len(in.CampaignISP) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignISP.Has(true, c.ISP)
}
