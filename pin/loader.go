package selector

import (
	"fmt"
	"sync"
	"time"

	"clickyab.com/gad/config"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/mr"
	"github.com/clickyab/services/assert"

	"github.com/clickyab/services/mysql"
	"github.com/sirupsen/logrus"
)

var (
	loaded []mr.SlotPinData
	lock   = &sync.RWMutex{}
	once   = &sync.Once{}

	lastTime time.Time
)

type myModel struct {
}

func interval() {
	manager := mr.NewManager()
	fail := 0
	for {
		<-time.After(time.Minute)

		l, err := manager.LoadSlotPin()
		if err != nil {
			//oh crap, failed. can we tolerate this?
			if fail > config.Config.Clickyab.MaxLoadFail {
				assert.Nil(err, fmt.Sprintf("more than %s time failed to load data", fail))
			}
			fail++
			continue
		}
		fail = 0
		lock.Lock()
		loaded = l
		lastTime = time.Now()
		lock.Unlock()
	}

}

// GetPinAdData return the current slot pin data
func GetPinAdData() []mr.SlotPinData {
	lock.RLock()
	defer lock.RUnlock()

	fail := time.Since(lastTime) > 5*time.Minute
	//assert.False(fail, "[BUG] the loader is not called for so long!")

	if fail {
		logrus.Fatal("failed! restart me please")
	}
	return loaded
}

// Initialize initialize the models
func (m *myModel) Initialize() {
	once.Do(func() {
		var err error
		manager := mr.NewManager()
		loaded, err = manager.LoadSlotPin()
		assert.Nil(err)
		lastTime = time.Now()
		middlewares.SafeGO(nil, true, false, interval)
	})
}

func init() {
	mysql.Register(&myModel{})
}
