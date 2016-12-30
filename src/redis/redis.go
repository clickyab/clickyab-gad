// Package aredis is for initializing the redis requirements, the pool and connection
package aredis

import (
	"assert"
	"config"
	"sync"
	"time"

	"fmt"

	"strconv"

	"github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var (
	// Pool the actual pool to use with redis
	Client *redis.Client
	once   = &sync.Once{}
)

// Initialize try to create a redis pool
func Initialize() {
	once.Do(func() {
		Client = redis.NewClient(
			&redis.Options{
				Network:  config.Config.Redis.Network,
				Addr:     config.Config.Redis.Address,
				Password: config.Config.Redis.Password,
				PoolSize: config.Config.Redis.Size,
				DB:       config.Config.Redis.Databse,
			},
		)
		// PING the server to make sure every thing is fine
		assert.Nil(Client.Ping().Err())
		logrus.Debug("redis is ready.")
	})
}

// StoreKey is a simple key value store with timeout
func StoreKey(key, data string, expire time.Duration) error {
	return Client.Set(key, data, expire).Err()
}

// GetKey Get a key from redis
func GetKey(key string, touch bool, expire time.Duration) (string, error) {
	cmd := Client.Get(key)
	if err := cmd.Err(); err != nil {
		return "", err
	}

	if touch {
		bCmd := Client.Expire(key, expire)
		if err := bCmd.Err(); err != nil {
			return "", err
		}
	}
	return cmd.Result()
}

// RemoveKey for removing a key in redis
func RemoveKey(key string) error {
	bCmd := Client.Del(key)
	return bCmd.Err()
}

// HGetAll Get a key and value from redis
// @deprecated use the HGetAllString
func HGetAll(key string, touch bool, expire time.Duration) (map[string]int, error) {
	res, err := HGetAllString(key, touch, expire)
	if err != nil {
		return nil, err
	}

	newRes := make(map[string]int)
	for i := range res {
		ii, _ := strconv.ParseInt(res[i], 10, 0)
		newRes[i] = int(ii)
	}
	return newRes, nil
}

// HGetAll Get a key and value from redis
func HGetAllString(key string, touch bool, expire time.Duration) (map[string]string, error) {
	cmd := Client.HGetAll(key)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	if touch {
		Client.Expire(key, expire)
	}

	return cmd.Result()
}

// IncHash try to inc hash
func IncHash(key string, hash string, value int64, expire time.Duration) (int64, error) {
	cmd := Client.HIncrBy(key, hash, value)

	Client.Expire(key, expire)
	return cmd.Result()
}

// HGetAll Get a key and value from redis
func HGetByField(key string, field ...string) (map[string]int64, error) {

	final := map[string]int64{}

	res, err := HGetAllString(key, false, 0)
	if err != nil {
		for i := range field {
			final[field[i]] = 0
		}
		return final, err
	}
	for f := range field {
		final[field[f]], _ = strconv.ParseInt(res[field[f]], 10, 0)
	}
	return final, nil
}
func SumHMGetField(prefix string, days int, field ...string) (map[string]int64, error) {
	now := time.Now()
	var (
		res   map[string]int64
		final = make(map[string]int64)
	)

	for i := 0; i <= days; i++ {
		res, _ = HGetByField(fmt.Sprintf("%s%s", prefix, now.AddDate(0, 0, -1*i).Format("060102")), field...)
		for f := range field {
			final[field[f]] += res[field[f]]
		}
	}
	return final, nil
}

func HMSet(key string, expire time.Duration, fields map[string]string) error {

	cmd := Client.HMSet(key, fields)
	if err := cmd.Err(); err != nil {
		return err
	}
	Client.Expire(key, expire)
	return nil
}

// StoreHashKey is a simple function to set hash key
func StoreHashKey(key, subkey, data string, expire time.Duration) error {
	err := Client.HSet(key, subkey, data).Err()
	if err == nil {
		err = Client.Expire(key, expire).Err()
	}

	return err
}

// RPush perform an rpush command
func LPush(key string, value ...interface{}) error {
	return Client.LPush(key, value...).Err()
}

// BRPopSingle is the function to pop a value from a single list
func BRPopSingle(key string, t time.Duration) (string, bool) {
	res := Client.BRPop(t, key)

	v := res.Val()
	if len(v) == 0 {
		return "", false
	}

	if len(v) == 2 && v[0] == key {
		return v[1], true
	}

	return "", false
}
