package ip2location

import (
	"net"
	"regexp"
)

var ispConst map[int64]*regexp.Regexp = map[int64]*regexp.Regexp{
	1: regexp.MustCompile(`(?i)iran\s?cell`),
	2: regexp.MustCompile(`(?i)Mobile Communication Company of Iran PLC`),
}

var m map[string]int64 = map[string]int64{
	"IR": 1,
	"Azarbayjan-e Sharqi":         2,
	"Ostan-e Azarbayjan-e Gharbi": 3,
	"Ardabil":                     4,
	"Esfahan":                     6,
	"Alborz":                      7,
	"Ilam":                        8,
	"Bushehr":                     9,
	"Tehran":                      10,
	"Chahar Mahall va Bakhtiari":  11,
	"Khorasan-e Janubi":           13,
	"Khorasan-e Razavi":           14,
	"Khorasan-e Shemali":          15,
	"Khuzestan":                   16,
	"Zanjan":                      17,
	"Semnan":                      18,
	"Sistan va Baluchestan":       19,
	"Fars":                        21,
	"Qazvin":                      22,
	"Qom":                         23,
	"Kordestan":                   24,
	"Kerman":                      25,
	"Kermanshah":                  26,
	"Kohkiluyeh va Buyer Ahmadi":  27,
	"Golestan":                    29,
	"Gilan":                       30,
	"Lorestan":                    31,
	"Mazandaran":                  32,
	"Markazi":                     33,
	"Hormozgan":                   34,
	"Hamadan":                     35,
	"Yazd":                        36,
	//"Hamadan":37,
}

var isoMap map[string]int64 = map[string]int64{
	"IR-01": 2,
	"IR-02": 3,
	"IR-03": 4,
	"IR-04": 6,
	"IR-32": 7,
	"IR-05": 8,
	"IR-06": 9,
	"IR-07": 10,
	"IR-08": 11,
	"IR-29": 13,
	"IR-30": 14,
	"IR-31": 15,
	"IR-10": 16,
	"IR-11": 17,
	"IR-12": 18,
	"IR-13": 19,
	"IR-14": 21,
	"IR-28": 22,
	"IR-26": 23,
	"IR-16": 24,
	"IR-15": 25,
	"IR-17": 26,
	"IR-18": 27,
	"IR-27": 29,
	"IR-19": 30,
	"IR-20": 31,
	"IR-21": 32,
	"IR-22": 33,
	"IR-23": 34,
	"IR-24": 35,
	"IR-25": 36,
}

// GetProvinceIDByIP get province id by ip
func GetProvinceIDByIP(ip net.IP) int64 {
	rec := IP2Location(ip.String())
	if i, ok := m[rec.Region]; ok {
		return i
	}
	return 0
}

// GetProvinceIDByIP get province id by ip
func GetProvinceISPByIP(ip net.IP) (int64, int64) {
	var province int64
	var uISP int64
	rec := IP2Location(ip.String())
	if i, ok := m[rec.Region]; ok {
		province = i
	}
	if rec.ISP != "" {
		//check isp
		for j := range ispConst {
			if ispConst[j].Match([]byte(rec.ISP)) {
				uISP = j
				break
			}
		}
	}

	return province, uISP
}

// GetProvinceIDByName get province id by name
func GetProvinceIDByName(s string) int64 {
	if i, ok := m[s]; ok {
		return i
	}
	return 0
}

// GetProvinceIDByISOName get province id by iso 3166-2
func GetProvinceIDByISOName(s string) int64 {
	if i, ok := isoMap[s]; ok {
		return i
	}
	return 0
}
