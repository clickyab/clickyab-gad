// Package aredis is for initializing the redis requirements, the pool and connection
package aredis

import (
	"assert"
	"config"
	"sync"
	"time"

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

// RemoveKey for removing a key in redis
func RemoveKey(key string) error {
	r := Pool.Get()
	defer func() { assert.Nil(r.Close()) }()
	_, err := r.Do("DEL", key)

	return err
}

// IncHash
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
