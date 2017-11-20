package store

import (
	"time"

	"clickyab.com/gad/redis"
	"github.com/clickyab/services/assert"
)

// TODO : {IMPORTANT} create key on selecting ad, so there is a chance to detect no ad request before 10 secound

// Set the data in already reserved key, unlock the key after that
func Set(key string, v string) {
	assert.Nil(aredis.LPush(key, time.Hour, v))
}

// Get the key from the system
func Get(key string) (string, bool) {
	return aredis.BRPopSingle(key, 10*time.Second)
}
