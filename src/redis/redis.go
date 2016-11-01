// Package aredis is for initializing the redis requirements, the pool and connection
package aredis

import (
	"assert"
	"config"
	"sync"
	"time"

	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

var (
	// Pool the actual pool to use with redis
	Pool *redis.Pool
	once = &sync.Once{}
)

// Initialize try to create a redis pool
func Initialize() {
	once.Do(func() {
		Pool = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial(config.Config.Redis.Network, config.Config.Redis.Address)
				if err != nil {
					return nil, err
				}
				if config.Config.Redis.Password != "" {
					if _, err := c.Do("AUTH", config.Config.Redis.Password); err != nil {
						_ = c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}

		// PING the server to make sure every thing is fine
		conn := Pool.Get()
		defer func() { _ = conn.Close() }()

		_, err := conn.Do("PING")
		assert.Nil(err)

		logrus.Info("redis is reday.")
	})
}

// StoreKey is a simple key value store with timeout
func StoreKey(key, data string, expire time.Duration) error {
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()
	_, err := r.Do("SET", key, data)
	if err != nil {
		return err
	}
	_, err = r.Do("EXPIRE", key, int64(expire.Seconds()))

	return err
}

// GetKey Get a key from redis
func GetKey(key string, touch bool, expire time.Duration) (string, error) {
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()
	res, err := r.Do("GET", key)
	if err != nil {
		return "", err
	}

	if touch {
		_, err = r.Do("EXPIRE", key, int64(expire.Seconds()))
		assert.Nil(err)
	}
	data, err := redis.Bytes(res, err)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// HGetAll Get a key and value from redis
func HGetAll(key string, touch bool, expire time.Duration) (map[string]int, error) {
	var res map[string]int
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()
	n, err := r.Do("HGETALL", key)

	res, err = redis.IntMap(n, err)
	if err != nil {
		return res, err
	}

	if touch {
		for k := range res {
			_, err = r.Do("EXPIRE", k, int64(expire.Seconds()))
			assert.Nil(err)
		}

	}
	return res, nil
}

// RemoveKey for removing a key in redis
func RemoveKey(key string) error {
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()
	_, err := r.Do("DEL", key)

	return err
}

// IncHash try to inc hash
func IncHash(key string, hash string, value int, touch bool, expire time.Duration) (int64, error) {
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()

	res, err := r.Do("HINCRBY", key, hash, value)
	data, err := redis.Int64(res, err)
	if err != nil {
		return 0, err
	}
	if touch {
		_, err = r.Do("EXPIRE", key, int64(expire.Seconds()))
		assert.Nil(err)
	}
	return data, nil
}

// HGetAll Get a key and value from redis
func HGetByField(key string, field ...string) (map[string]int64, error) {
	var res map[string]int64
	final := map[string]int64{
		"c":  0,
		"i":  0,
		"fc": 0,
		"fi": 0,
	}
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()
	n, err := r.Do("HGETALL", key)

	res, err = redis.Int64Map(n, err)
	if err != nil {
		return final, err
	}
	for f := range field {
		final[field[f]] = res[field[f]]
	}
	return final, nil
}
func SumHMGetField(prefix string, days int, field ...string) (map[string]int64, error) {
	now := time.Now()
	var res map[string]int64
	final := map[string]int64{
		"c":  0,
		"i":  0,
		"fc": 0,
		"fi": 0,
	}

	for i := 0; i <= days; i++ {
		res, _ = HGetByField(fmt.Sprintf("%s%s", prefix, now.AddDate(0, 0, -1*i).Format("060102")), field...)
		for f := range field {
			final[field[f]] += res[field[f]]
		}
	}
	return final, nil
}

func HMSet(key string, touch bool, expire time.Duration, fields ...interface{}) error {
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()

	nf := make([]interface{}, len(fields)+1)
	nf[0] = key
	for i := range fields {
		nf[i+1] = fields[i]
	}
	fmt.Println(nf)
	_, err := redis.String(r.Do("HMSET", nf...))
	if err != nil {
		return err
	}
	if touch {
		_, err = r.Do("EXPIRE", key, int64(expire.Seconds()))
		assert.Nil(err)
	}
	return nil
}
