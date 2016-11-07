package config

import (
	"fmt"
	"strings"
)

var sizes = map[string]int{
	"120x600":  1,
	"160x600":  2,
	"300x250":  3,
	"336x280":  4,
	"468x60":   5,
	"728x90":   6,
	"120x240":  7,
	"320x50":   8,
	"800x440":  9,
	"300x600":  11,
	"970x90":   12,
	"970x250":  13,
	"250x250":  14,
	"300x1050": 15,
	"320x480":  16,
	"48x320":   17,
	"128x128":  18,
}

var videoSize = []int{3,4,9,16,14,17}

// GetSize return the size of a banner in clickyab std
func GetSize(size string) (int, error) {
	s, ok := sizes[size]
	if ok {
		return s, nil
	}

	return 0, fmt.Errorf("size %s is not valid", size)
}

// GetSizeByNum return the size
func GetSizeByNum(num int) (string, string) {
	// TODO : better way / no loop please
	for i, s := range sizes {
		if s == num {
			a := strings.Split(i, "x")
			return a[0], a[1]
		}
	}
	return "", ""
}

// CheckIfBannerIsVideo check if the banner suitable for video
func CheckIfBannerIsVideo(size int) bool{
	for i := range videoSize{
		if videoSize[i]==size{
			return true
		}
	}
	return false;
}
