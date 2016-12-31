package filter

import (
	"config"
	"mr"
	"selector"
)

// CheckWebSize check if the banner size exists in the request
func CheckWebSize(c *selector.Context, in mr.AdData) bool {
	if in.AdType == config.AdTypeVideo {
		for _, size := range c.Size {
			if config.InVideoSize(size) {
				return true
			}
		}
		return false
	}

	for _, size := range c.Size {
		if size == in.AdSize {
			return true
		}
	}
	return false
}

// CheckVastSize check if the banner size fits for Vast Template
func CheckVastSize(_ *selector.Context, in mr.AdData) bool {
	if in.AdTarget == config.AdTypeDynamic {
		return false
	}
	
	return in.AdType == config.AdTypeVideo || config.InVastSize(in.AdSize)
}
