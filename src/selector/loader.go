package selector

import (
	"assert"
	"models"
	"mr"
	"sync"
	"time"

	"fmt"
)

var (
	loaded []mr.AdData
	lock   = &sync.RWMutex{}
	once   = &sync.Once{}
)

type myModel struct {
}

func interval() {
	var err error
	manager := mr.NewManager()
	loaded, err = manager.LoadAds()
	assert.Nil(err)
	ticker := time.NewTicker(time.Minute)
	fail := 0
	for {
		select {
		case <-ticker.C:
			{
				l, err := manager.LoadAds()
				if err != nil {
					// TODO : handle this
					//oh crap, failed. can we tolerate this?
					if fail > 3 { // TODO Read from config
						assert.Nil(err, fmt.Sprintf("more than %s time failed to load data", fail))
					}
					fail++
					break
				}
				fail = 0
				lock.Lock()
				copy(loaded, l)
				lock.Unlock()
			}
		}
	}
	ticker.Stop()
}

func GetAdData() []mr.AdData {
	lock.RLock()
	defer lock.RUnlock()

	return loaded
}

func (m *myModel) Initialize() {
	once.Do(func() {
		go interval()
	})
}

func init() {
	models.Register(&myModel{})
}
