package filter

import (
	"config"
	"mr"
	"selector"
)

// CheckForVideo check if the banner size exists in the request
func CheckForVideo(c *selector.Context, in mr.AdData) bool {
	return config.CheckIfBannerIsVideo(in.AdSize)
}
