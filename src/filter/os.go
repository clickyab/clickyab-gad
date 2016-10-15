package filter

import (
	"config"
	"mr"
	"selector"
)

// CheckOS is the filter function that check for OS in system
func CheckOS(c *selector.Context, in mr.AdData) bool {

	if len(in.CpPlatforms) == 0 {
		return true
	}
	os := config.FindOsID(c.RequestData.Platform)
	for _, v := range in.CpPlatforms {
		if v == int64(os) {
			return true
		}
	}
	return false
}
