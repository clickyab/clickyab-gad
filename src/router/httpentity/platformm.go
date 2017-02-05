package httpentity

import (
	"entity"
	"strings"

	"github.com/mssola/user_agent"
)

const (
	osMac     int64 = 1
	osUnknown       = 2
	osWindows       = 3
	osLinux         = 4
	osIOS           = 5
	osAndroid       = 6
)

var platforms = map[string]int64{
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

// findHTTPOS try to find os ID base on old id of system
func findHTTPOS(ua *user_agent.UserAgent) entity.OS {
	res := entity.OS{}
	res.Mobile = ua.Mobile()
	res.Name = ua.Platform()
	res.ID = osUnknown

	p, ok := platforms[strings.ToLower(res.Name)]
	if ok {
		res.ID = p
	}
	return res
}
