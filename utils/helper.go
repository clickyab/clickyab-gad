package utils

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"time"

	"clickyab.com/gad/redis"
	"github.com/clickyab/services/config"
)

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

var (
	impExpireDaily = config.RegisterDuration("clickyab.daily_imp_expire", 7*24*time.Hour, "daily impression expiration")
)

//IncKeyDaily function increase redis daily key
func IncKeyDaily(key, subKey string, count int64) (int64, error) {
	res, err := aredis.IncHash(
		key,
		subKey,
		count,
		impExpireDaily.Duration())
	return res, err
}

// Hash is the hash generation func for keys, md5 normally
func Hash(k string) string {
	h := md5.New()
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
			if count > c-1 {
				break
			}
			res = res + " " + strArr[i]
		}
	} else {
		return s
	}
	return res + "…"
}
