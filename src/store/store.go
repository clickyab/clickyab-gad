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
	clusterFactory StoreFactory
	syncFactory    StoreFactory
)

func RegisterCluster(s StoreFactory) {
	clusterFactory = s
}

func RegisterSync(s StoreFactory) {
	syncFactory = s
}

// GetSyncStore return an inner process sync
func GetSyncStore() Store {
	assert.NotNil(syncFactory, "[BUG] sync factory is not set")
	return syncFactory()
}

// GetClusterStore return an in cluster sync
func GetClusterStore() Store {
	assert.NotNil(clusterFactory, "[BUG] cluster factory is not set")
	return clusterFactory()
}