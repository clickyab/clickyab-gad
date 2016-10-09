package filter

import (
	_ "middlewares"
	"mr"
	"selector"
)

// CheckForSize check if the banner size exists in the request
func CheckForSize(c *selector.Context, in mr.AdData) bool {
	validSizes := c.Size
	for size := range validSizes {
		if size == in.AdSize {
			return true
		}
	}
	return false
}
