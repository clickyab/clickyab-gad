package pin

import (
	"fmt"
	"sync"
	"time"

	"clickyab.com/gad/models"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"

	"context"

	"github.com/clickyab/services/mysql"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
)

var (
	loaded []models.SlotPinData
	lock   = &sync.RWMutex{}
	once   = &sync.Once{}

	lastTime time.Time

	maxLoadFail = config.RegisterInt("clickyab.max_load_fail", 3, "maximum fail allowed in sequence")
)

type myModel struct {
}

func interval(_ context.CancelFunc) {
	manager := models.NewManager()
	fail := 0
	for {
		<-time.After(time.Minute)

		l, err := manager.LoadSlotPin()
		if err != nil {
			//oh crap, failed. can we tolerate this?
			if fail > maxLoadFail.Int() {
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
func GetPinAdData() []models.SlotPinData {
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
		manager := models.NewManager()
		loaded, err = manager.LoadSlotPin()
		assert.Nil(err)
		lastTime = time.Now()
		safe.ContinuesGoRoutine(interval, time.Second)
	})
}

func init() {
	mysql.Register(&myModel{})
}
