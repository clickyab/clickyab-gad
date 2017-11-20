package filter

import (
	"clickyab.com/gad/models"
	"clickyab.com/gad/models/selector"
)

// IsNativeAd tells if an ad is native
func IsNativeAd(c *selector.Context, in models.AdData) bool {
	return in.AdType == models.NativeAdType
}
