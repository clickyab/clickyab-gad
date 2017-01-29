package config

import "strconv"

const (
	short = "short"
	def   = "default"
	long  = "long"
)

var vastConfig = map[string]map[string][]string{
	"short": map[string][]string{
		"start": []string{"linear", "11", "00:00:10", "00:00:03"},
		"end":   []string{"linear", "12", "00:00:10", "00:00:03"},
	},
	"default": map[string][]string{
		"start":    []string{"linear", "11", "00:00:10", "00:00:03"},
		"00:00:10": []string{"non-linear", "22", "00:00:12"},
		"end":      []string{"linear", "13", "00:00:10", "00:00:03"},
	},
	"long": map[string][]string{
		"start":    []string{"linear", "11", "00:00:10", "00:00:03"},
		"00:00:20": []string{"non-linear", "22", "00:00:12"},
		"00:01:20": []string{"non-linear", "23", "00:00:12"},
		"00:03:20": []string{"non-linear", "24", "00:00:12"},
		"00:05:20": []string{"non-linear", "25", "00:00:12"},
	},
}

// MakeVastLen return vast len
func MakeVastLen(len string) (string, map[string][]string) {
	if m, found := vastConfig[len]; found {
		return def, m
	}

	if m, err := strconv.ParseInt(len, 10, 64); err == nil {
		if m < 30 {
			return short, vastConfig[short]
		} else if m < 90 {
			return def, vastConfig[def]
		} else {
			return long, vastConfig[long]
		}
	}
	return def, vastConfig[def]
}
