package selector

import (
	"assert"
	"config"
	"fmt"
	"middlewares"
	"models"
	"mr"
	"sync"
	"time"
)

var (
	loaded []mr.AdData
	lock   = &sync.RWMutex{}
	once   = &sync.Once{}
)

type myModel struct {
}

func interval() {
	manager := mr.NewManager()
	ticker := time.NewTicker(time.Minute)
	fail := 0
	for range ticker.C {
		l, err := manager.LoadAds()
		if err != nil {
			//oh crap, failed. can we tolerate this?
			if fail > config.Config.Clickyab.MaxLoadFail {
				assert.Nil(err, fmt.Sprintf("more than %s time failed to load data", fail))
			}
			fail++
			break
		}
		fail = 0
		lock.Lock()
		loaded = l
		lock.Unlock()
	}

	ticker.Stop()
}

// GetAdData return the current stored ad data
func GetAdData() []mr.AdData {
	lock.RLock()
	defer lock.RUnlock()

	return loaded
}

// Initialize initialize the models
func (m *myModel) Initialize() {
	once.Do(func() {
		var err error
		manager := mr.NewManager()
		loaded, err = manager.LoadAds()
		assert.Nil(err)
		middlewares.SafeGO(nil, true, interval)
	})
}

func init() {
	models.Register(&myModel{})
}
