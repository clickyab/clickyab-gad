package redis

import (
	"bytes"
	"cache"
	"crypto/sha1"
	"entity"
	"fmt"
	"services/redis"
	"time"
)

type redisCache struct {
}

// Sha1 is the sha1 generation func
func sha(k string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(k))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Cache is called to store the cache
func (redisCache) Cache(e entity.Cacheable, t time.Duration) error {
	name := "CACHE_" + sha(e.String())
	target := &bytes.Buffer{}
	err := e.Decode(target)
	if err != nil {
		return err
	}

	res := aredis.Client.Set(name, target.String(), t)
	return res.Err()
}

// Hit called when we need to load the cache
func (redisCache) Hit(key string, e entity.Cacheable) error {
	name := "CACHE_" + sha(e.String())
	res := aredis.Client.Get(name)
	if err := res.Err(); err != nil {
		return err
	}
	data, err := res.Result()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString(data)
	return e.Encode(buf)
}

// NewRedisCacheProvider return a new cache storage in redis
func NewRedisCacheProvider() cache.CacheProvider {
	return &redisCache{}
}
