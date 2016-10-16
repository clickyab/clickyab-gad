package filter

import (
	"mr"
	"selector"
)

// CheckWhiteList return boolean
func CheckWhiteList(c *selector.Context, in mr.AdData) bool {
	if len(in.CpPlacement) == 0 {
		return true
	}
	return in.CpPlacement.Has(c.WID)
}
