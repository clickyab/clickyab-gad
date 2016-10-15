package config

import "strings"

const (
	osMac     = 1
	osUnknown = 2
	osWindows = 3
	osLinux   = 4
	osIOS     = 5
	osAndroid = 6
)

var platforms = map[string]int{
	"windows":   osWindows,
	"macintosh": osMac,
	"x11":       osLinux,
	"android":   osAndroid,
	"tablet":    osAndroid,
	"iPhone":    osIOS,
	"like Mac":  osIOS,
	"iPod":      osIOS,
	"iPad":      osIOS,
	"linux":     osAndroid,
	"mobile":    osAndroid,
}

// FindOsID try to find os ID base on old id of system
func FindOsID(platform string) int {
	platform = strings.ToLower(platform)
	for OSName, ID := range platforms {
		if strings.Contains(OSName, platform) {
			return ID
		}
	}
	return osUnknown
}
