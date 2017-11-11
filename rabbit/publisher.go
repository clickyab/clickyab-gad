package rabbit

import (
	"container/ring"
	"encoding/json"
	"errors"
	"sync"

	"clickyab.com/gad/config"
	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"

	"github.com/streadway/amqp"
)

var (
	rng *ring.Ring
)

type chnlLock struct {
	chn    Channel
	lock   *sync.Mutex
	rtrn   chan amqp.Confirmation
	wg     *sync.WaitGroup
	closed bool
}

// Publish try to publish an event
func Publish(in Job) (err error) {
	rng = rng.Next()
	v := rng.Value.(*chnlLock)
	v.lock.Lock()
	defer v.lock.Unlock()
	if v.closed {
		return errors.New("waiting for finalize, can not publish")
	}

	msg, err := json.Marshal(in)
	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		CorrelationId: <-utils.ID,
		Body:          msg,
	}

	v.wg.Add(1)
	defer func() {
		// If the result is error, release the lock, there is no message to confirm!
		if err != nil {
			v.wg.Done()
		}
	}()
	topic := in.GetTopic()
	if config.Config.AMQP.Debug {
		topic = "debug." + topic
	}
	return v.chn.Publish(config.Config.AMQP.Exchange, topic, true, false, pub)
}

// MustPublish publish an event with force
func MustPublish(ei Job) {
	assert.Nil(Publish(ei))
}

// FinalizeWait is a function to wait for all publication to finish. after calling this,
// must not call the PublishEvent
func finalizeWait() {
	for i := 0; i < config.Config.AMQP.Publisher; i++ {
		rng = rng.Next()
		v := rng.Value.(*chnlLock)
		v.lock.Lock()
		// I know this lock release at the end, not after for, and this is ok
		defer v.lock.Unlock()

		v.closed = true
		v.wg.Wait()
		_ = v.chn.Close()
	}
}

func publishConfirm(cl *chnlLock) {
	for range cl.rtrn {
		cl.wg.Done()
	}
}
