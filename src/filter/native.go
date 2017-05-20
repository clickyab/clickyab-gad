package filter

import (
	"mr"
	"selector"
)

func IsNativeAd(c *selector.Context, in mr.AdData) bool {
	return in.AdType == mr.NativeAdType
}
