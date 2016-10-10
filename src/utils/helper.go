package utils

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"regexp"
	"strings"
)

var spaceMatch = regexp.MustCompile(`\s+`)

// PrefixMatch return the matched items in array
func PrefixMatch(in string, data ...string) []string {
	var res []string
	for i := range data {
		if len(in) < len(data[i]) {
			if data[i][:len(in)] == in {
				res = append(res, data[i][len(in):])
			}
		}
	}

	return res
}

// CleanSplit replace all multiple strings with one and then split them using the space as delimiter
func CleanSplit(line string) []string {
	str := strings.Trim(spaceMatch.ReplaceAllString(line, " "), " ")
	return strings.Split(str, " ")
}

// Exists returns whether the given file or directory exists or not
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// StringInArray check for a string in other strings
func StringInArray(q string, arr ...string) bool {
	for i := range arr {
		if arr[i] == q {
			return true
		}
	}

	return false
}

// StringInArray check for a string in other strings
func Int64InArray(q int64, arr ...int64) bool {
	for i := range arr {
		if arr[i] == q {
			return true
		}
	}

	return false
}

func Long2IP(ipLong uint32) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip.String()
}
func IP2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}
