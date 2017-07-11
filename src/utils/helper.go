package utils

import (
	"config"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"redis"
	"regexp"
	"strings"
	"syscall"
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

// Int64InArray check for a string in other strings
func Int64InArray(q int64, arr ...int64) bool {
	for i := range arr {
		if arr[i] == q {
			return true
		}
	}

	return false
}

// Long2IP change the integer to ip
func Long2IP(ipLong uint32) net.IP {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip
}

// IP2long change ip to integer
func IP2long(ip net.IP) (uint32, error) {
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip2 := ip.To4()
	if ip2 == nil {
		return 0, fmt.Errorf("ipv6? the input was %s", ip.String())
	}
	return binary.BigEndian.Uint32(ip2), nil
}

// WaitSignal get os signal
func WaitSignal(exit chan chan struct{}) {
	quit := make(chan os.Signal, 6)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGABRT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGQUIT,
	)

	<-quit
	if exit != nil {
		tmp := make(chan struct{})
		exit <- tmp

		<-tmp
	}
}

//IncKeyDaily function increase redis daily key
func IncKeyDaily(key, subKey string, count int64) (int64, error) {
	res, err := aredis.IncHash(
		key,
		subKey,
		count,
		config.Config.Clickyab.DailyImpExpire)
	return res, err
}

// InSlice check if the value exists in the slice
func InSlice(a interface{}, list []interface{}) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Sha1 is the sha1 generation func
func Sha1(k string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(k))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// LimitCharacter limit character
func LimitCharacter(s string, c int) string {
	l := len([]rune(s))
	var res = ""
	if l > c {
		strArr := strings.Split(s, " ")
		count := 0
		for i := range strArr {
			count = count + len([]rune(strArr[i]))
			count++ //space
			if count > c-3 {
				break
			}
			res = res + " " + strArr[i]
		}
	} else {
		return s
	}
	return res + "..."
}
