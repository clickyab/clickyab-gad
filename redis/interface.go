package aredis

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisClient interface {
	Ping() *redis.StatusCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(key string) *redis.StringCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	Del(keys ...string) *redis.IntCmd
	HGetAll(key string) *redis.StringStringMapCmd
	HIncrBy(key, field string, incr int64) *redis.IntCmd
	HMSet(key string, fields map[string]interface{}) *redis.StatusCmd
	HSet(key, field string, value interface{}) *redis.BoolCmd
	LPush(key string, values ...interface{}) *redis.IntCmd
	SAdd(key string, members ...interface{}) *redis.IntCmd
	SMembers(key string) *redis.StringSliceCmd
	BRPop(timeout time.Duration, keys ...string) *redis.StringSliceCmd
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
}

