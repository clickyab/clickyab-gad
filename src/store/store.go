package store

import (
	"assert"
	"redis"
	"time"
)

// Set the data in already reserved key, unlock the key after that
func Set(key string, v string) {
	assert.Nil(aredis.LPush(key, time.Hour, v))
}

// Get the key from the system
func Get(key string) (string, bool) {
	return aredis.BRPopSingle(key, 10*time.Second)
}
