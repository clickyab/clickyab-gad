package store

import (
	"assert"
	"time"
)

// Store is the blocking store interface
type Store interface {
	// Push data in the store
	Push(string, string, time.Duration)
	// Pop and remove data from store, its blocking pop
	Pop(string, time.Duration) (string, bool)
}

type StoreFactory func() Store

var (
	factory StoreFactory
)

func Register(s StoreFactory) {
	factory = s
}

// Push data in the store
func Push(key, val string, t time.Duration) {
	assert.NotNil(factory, "[BUG] factory is not registered")
	factory().Push(key, val, t)
}

// Pop and remove data from store, its blocking pop
func Pop(key string, t time.Duration) (string, bool) {
	assert.NotNil(factory, "[BUG] factory is not registered")
	return factory().Pop(key, t)
}
