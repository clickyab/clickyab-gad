package utils

import (
	"config"
	"crypto/sha1"
	"net"
)

//Ctr calculate ctr
func Ctr(imps int64, clicks int64) float64 {
	if imps == 0 || clicks == 0 {
		return config.Config.DefaultCTR
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

// CreateCopID create COP ID
func CreateCopID(useragent string, ip net.IP) string {
	h := sha1.New()
	h.Write([]byte(useragent))
	bs := h.Sum(ip)
	return string(bs)
}