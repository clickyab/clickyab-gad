package redis

import (
	"eav"
	"services/redis"
	"sync"
	"time"
)

type kiwiRedis struct {
	key  string
	v    map[string]string
	lock sync.Mutex
}

// Key return the parent key
func (kr *kiwiRedis) Key() string {
	return kr.key
}

// SubKey for adding a sub key
func (kr *kiwiRedis) SubKey(key, value string) eav.Kiwi {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	kr.v[key] = value

	return kr
}

// GetKey return a key
func (kr *kiwiRedis) GetKey(key string) string {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	if v, ok := kr.v[key]; ok {
		return v
	}
	res := aredis.Client.HGet(kr.key, key)
	if res.Err() != nil {
		return ""
	}

	r, err := res.Result()
	if err != nil {
		return ""
	}

	return r
}

// GetAllKeys from the store
func (kr *kiwiRedis) GetAllKeys() map[string]string {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	kr.v = map[string]string{}
	res := aredis.Client.HGetAll(kr.key)
	if res.Err() != nil {
		return kr.v
	}

	r, err := res.Result()
	if err != nil {
		return kr.v
	}

	kr.v = r
	return kr.v
}

// Save the entire keys (mostly first time)
func (kr *kiwiRedis) Save(t time.Duration) error {
	kr.lock.Lock()
	defer kr.lock.Unlock()

	res := aredis.Client.HMSet(kr.key, kr.v)
	if res.Err() != nil {
		return res.Err()
	}

	b := aredis.Client.Expire(kr.key, t)
	return b.Err()
}

// NewRedisEAVStore return a redis store for eav
func NewRedisEAVStore(key string) eav.Kiwi {
	return &kiwiRedis{
		key:  key,
		v:    make(map[string]string),
		lock: sync.Mutex{},
	}
}
