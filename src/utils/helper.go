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
	"time"
	"transport"
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

// Long2IP function @todo
func Long2IP(ipLong uint32) net.IP {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip := net.IP(ipByte)
	return ip
}

// IP2long function @todo
func IP2long(ip net.IP) (uint32, error) {
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip), nil
}

// WaitSignal get os signal
func WaitSignal(exit chan chan struct{}) {
	quit := make(chan os.Signal, 5)
	signal.Notify(quit, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	<-quit

	tmp := make(chan struct{})
	exit <- tmp

	<-tmp
}

//KeyGenDaily function  generate a redis key
func KeyGenDaily(prefix, value string) string {
	date := time.Now().Format("060102")
	return prefix + transport.DELIMITER + value + transport.DELIMITER + date
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
