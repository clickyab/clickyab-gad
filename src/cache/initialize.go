package cache

import (
	"config"
	"errors"
	"store"
	"strings"

	"net"

	"github.com/golang/groupcache"
	"gopkg.in/labstack/echo.v3"
)

var (
	me   *groupcache.Group
	pool *HTTPPool
)

// Initialize the cache system
func Initialize(ip net.IP, port int, server *echo.Echo) {
	me = groupcache.NewGroup("ads", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			data := strings.Split(key, "-")
			if len(data) != 2 {
				return errors.New("invalid key")
			}

			if data[0] != config.Config.MachineName {
				return errors.New("this is not generated here")
			}

			v, ok := store.Get(key)
			if !ok {
				return errors.New("cache expired")
			}
			dest.SetString(v)
			return nil
		},
	))
	pool = NewHTTPPool(ip, port, server)
}

// Get the key from cache
func Get(key string) (string, error) {
	var b []byte
	err := me.Get(nil, key, groupcache.AllocatingByteSliceSink(&b))

	if err != nil {
		return "", err
	}

	return string(b), nil
}
