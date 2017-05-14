package utils

import (
	"config"
	"crypto/sha1"
	"fmt"
	"math"
)

//Ctr calculate ctr
func Ctr(imps int64, clicks int64) float64 {
	if imps == 0 || clicks == 0 {
		return config.Config.Clickyab.DefaultCTR
	}
	return (float64(clicks) / float64(imps)) * 100
}

//Cpc calculate cpc
func Cpc(spend int64, clicks int64) int64 {
	return spend / clicks
}

//Cpm calculate cpm
func Cpm(bid int64, ctr float64) int64 {
	return int64(float64(bid) * ctr * 10.0)
}

// WinnerBid calculate winner bid
func WinnerBid(cpm int64, ctr float64) int64 {
	return int64(float64(cpm)/(ctr*10)) + 1
}

// CreateHash is used to handle the cop key
func CreateHash(l int, items ...[]byte) string {
	h := sha1.New()
	for i := range items {
		_, _ = h.Write(items[i])
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	if l >= len(sum) {
		l = len(sum)
	}
	if l < 1 {
		l = 1
	}
	return sum[:l]
}

// AreaInGlob is a helper fuunction to handle check point in a globe
func AreaInGlob(lat, lon, centerLat, centerLon, radius float64) bool {
	var ky = 40000.0 / 360.0
	var kx = math.Cos(math.Pi*centerLat/180.0) * ky
	dx := math.Abs(centerLon-lon) * kx
	dy := math.Abs(centerLat-lat) * ky
	return (math.Sqrt(dx*dx+dy*dy) <= radius)
}
