package store

import (
	"assert"
	"fmt"
	"redis"
	"time"
)

// Set the data in already reserved key, unlock the key after that
func Set(key string, v string) {
	fmt.Println("set ->", key)
	defer fmt.Println("set <-", key)
	assert.Nil(aredis.LPush(key, v))
}

// Get the key from the system
func Get(key string) (string, bool) {
	fmt.Println("get ->", key)
	defer fmt.Println("get <-", key)
	return aredis.BRPopSingle(key, 10*time.Second)
}
