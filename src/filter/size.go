package filter

import (
	"mr"
	"selector"
)

// CheckForSize check if the banner size exists in the request
func CheckForSize(c *selector.Context, in mr.AdData) bool {

	for _, size := range c.Size {
		if size == in.AdSize {
			return true
		}
	}
	return false
}
