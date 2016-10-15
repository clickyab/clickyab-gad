package filter

import (
	"mr"
	"selector"
)

// CheckNetwork filter network for campaigns
func CheckNetwork(c *selector.Context, in mr.AdData) bool {
	if c.Mobile == true {
		return in.CpNetwork == 1
	}
	if c.Mobile == false {
		return in.CpNetwork == 0
	}
	return false
}
