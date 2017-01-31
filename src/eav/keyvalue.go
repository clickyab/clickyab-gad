package eav

import (
	"time"
	"assert"
)

// Kiwi is the key value storage system in a parent key
type Kiwi interface {
	// Key return the parent key
	Key() string
	// SubKey for adding a sub key
	SubKey(key, value string) Kiwi
	// GetKey return a key
	GetKey(key string) string
	// GetAllKeys from the store
	GetAllKeys() map[string]string
	// Save the entire keys (mostly first time)
	Save(time.Duration) error
}

type StoreFactory func(string) Kiwi

var (
	factory StoreFactory
)

func Register(s StoreFactory) {
	factory = s
}

// NewEavStore return a new eav store
func NewEavStore(key string) Kiwi {
	assert.NotNil(factory, "[BUG] factory is not registered")
	return factory(key)
}
