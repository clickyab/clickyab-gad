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
	for _, v := range in.CpPlacement {
		if v == c.WID {
			return true
		}
	}
	return false
}
