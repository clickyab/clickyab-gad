package filter

import (
	"config"
	"mr"
	"selector"
)

// CheckForSize check if the banner size exists in the request
func CheckForSize(c *selector.Context, in mr.AdData) bool {
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
