package rabbit

import (
	"assert"
	"config"
	"container/ring"
	"encoding/json"
	"errors"
	"sync"

	"utils"

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
func Publish(topic string, in interface{}) (err error) {
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

	err = v.chn.Publish(config.Config.AMQP.Exchange, topic, true, false, pub)

	return err
}

// MustPublish publish an event with force
func MustPublish(topic string, ei interface{}) {
	assert.Nil(Publish(topic, ei))
}

// FinalizeWait is a function to wait for all publication to finish. after calling this,
// must not call the PublishEvent
func FinalizeWait() {
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
