package filter

import (
	"mr"
	"selector"
)

// CheckBlackList function @todo
func CheckBlackList(c *selector.Context, in mr.AdData) bool {
	if len(in.CpWfilter) == 0 {
		return true
	}
	for _, v := range in.CpWfilter {
		if v == c.WID {
			return false
		}
	}
	return false
}
