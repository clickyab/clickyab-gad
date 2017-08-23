package ip2location

import "net"

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

// GetProvinceIDByIP get province id by ip
func GetProvinceIDByIP(ip net.IP) int64 {
	rec := IP2Location(ip.String())
	if i, ok := m[rec.Region]; ok {
		return i
	}
	return 0
}

// GetProvinceIDByName get province id by name
func GetProvinceIDByName(s string) int64 {
	if i, ok := m[s]; ok {
		return i
	}
	return 0
}