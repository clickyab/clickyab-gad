package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

//CheckISP find isp
func CheckISP(c *selector.Context, in mr.AdData) bool {
	if c.ISP == 0 {
		return len(in.CampaignISP) == 0
	}
	// The 1 means iran. watch for it please!
	return in.CampaignISP.Has(true, c.ISP)
}
