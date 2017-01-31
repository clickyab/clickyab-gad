package mock

import (
	"eav"
	"time"
)

type MockKiwi struct {
	MasterKey string
	Data      map[string]string
	Duration  time.Duration
}

var (
	pre   = make(map[string]map[string]string)
	store = make(map[string]*MockKiwi)
)

// Key return the parent key
func (m *MockKiwi) Key() string {
	return m.MasterKey
}

// SubKey for adding a sub key
func (m *MockKiwi) SubKey(key, value string) eav.Kiwi {
	m.Data[key] = value
	return m
}

// GetKey return a key
func (m *MockKiwi) GetKey(key string) string {
	return m.Data[key]
}

// GetAllKeys from the store
func (m *MockKiwi) GetAllKeys() map[string]string {
	return m.Data
}

// Save the entire keys (mostly first time)
func (m *MockKiwi) Save(t time.Duration) error {
	m.Duration = t
	return nil
}

func NewMockStore(key string) eav.Kiwi {
	var (
		data map[string]string
		ok   bool
	)
	if data, ok = pre[key]; !ok {
		data = make(map[string]string)
	}
	m := &MockKiwi{
		MasterKey: key,
		Data:      data,
	}

	store[key] = m
	return m
}

func SetMockData(key string, data map[string]string) {
	pre[key] = data
}
