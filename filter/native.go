package filter

import (
	"clickyab.com/gad/mr"
	"clickyab.com/gad/selector"
)

// IsNativeAd shows if the ad is native or not
func IsNativeAd(c *selector.Context, in mr.AdData) bool {
	return in.AdType == mr.NativeAdType
}
