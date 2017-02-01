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

// GetClusterStore return an in cluster sync
func GetSyncStore() Store {
	assert.NotNil(factory, "[BUG] cluster factory is not set")
	return factory()
}
