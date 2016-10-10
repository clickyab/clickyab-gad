package platform

import (
	"errors"
	"strings"

	"middlewares"
)

type Platform map[string]int

var Platforms = Platform{
	"windows":    3,
	"macintosh":  1,
	"x11":        4,
	"android":    6,
	"tablet":     4,
	"iPhone":     5,
	"like Mac":   5,
	"iPod":       2,
	"iPad":       5,
	"blackberry": 2,
	"symbian":    2,
	"linux":      6,
	"en-us":      2,
	"mobile":     4,
	"weboS":      2,
}

//find ID OS.
func FindIdOs(c middlewares.RequestData) (int, error) {
	c.Platform = strings.ToLower(c.Platform)
	for OSName, ID := range Platforms {
		if strings.Contains(OSName, c.Platform) {
			return ID, nil
		}
		//if OSName == c.Platform{
		//	return ID,nil
		//}
	}
	err := errors.New("OS not found")
	return 0, err
}
