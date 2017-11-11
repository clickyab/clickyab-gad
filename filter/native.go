package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// IsNativeAd tells if an ad is native
func IsNativeAd(c *selector.Context, in mr.AdData) bool {
	return in.AdType == mr.NativeAdType
}
