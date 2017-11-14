package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
	"clickyab.com/gad/utils"
)

// CheckWebSize check if the banner size exists in the request
func CheckWebSize(c *selector.Context, in models.AdData) bool {
	if in.AdType == utils.AdTypeVideo {
		for _, size := range c.Size {
			if utils.InVideoSize(size) {
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
func CheckVastSize(_ *selector.Context, in models.AdData) bool {
	if in.AdType == utils.AdTypeDynamic {
		return false
	}

	return in.AdType == utils.AdTypeVideo || utils.InVastSize(in.AdSize)
}

// CheckAppSize check if the banner size exists in the request
func CheckAppSize(c *selector.Context, in models.AdData) bool {
	if in.AdType == utils.AdTypeVideo || in.AdType == utils.AdTypeDynamic {
		return false
	}

	for _, size := range c.Size {
		if size == in.AdSize {
			return true
		}
	}
	return false
}

// CheckWebMobileSize check if the banner size exists in the request
func CheckWebMobileSize(c *selector.Context, in models.AdData) bool {
	return in.AdSize == 8
}
