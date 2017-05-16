package config

import (
	"strconv"

	"github.com/Sirupsen/logrus"
)

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
func MakeVastLen(len string, first string, mid string, end string) (string, map[string][]string) {
	//apply vast customization
	preRes := vastConfig
	for i := range preRes {
		if first != "" {
			delete(preRes[i], "start")
		}
		if end != "" {
			delete(preRes[i], "end")
		}
		if mid != "" {
			for j := range preRes[i] {
				if j != "start" && j != "end" {
					logrus.Info(j)
					delete(preRes[i], j)
				}
			}
		}
	}
	if m, found := preRes[len]; found {
		return def, m
	}
	if m, err := strconv.ParseInt(len, 10, 64); err == nil {
		if m < 30 {
			logrus.Warn("<30")
			logrus.Warn(preRes[short])
			return short, preRes[short]
		} else if m < 90 {
			return def, preRes[def]
		} else {
			return long, preRes[long]
		}
	}
	return def, preRes[def]
}
