package redlock

import (
	"time"

	"clickyab.com/gad/utils"
	"github.com/clickyab/services/aredis"

	"github.com/sirupsen/logrus"
)

const unlockScript = `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
        `

var tryCoolDown time.Duration = 200 * time.Microsecond

type mux struct {
	ttl      time.Duration
	now      time.Time
	resource string

	value   string
	retries int
}

// LockInterface is the interface to use as distributed lock
type LockInterface interface {
	Lock()
	Unlock()
	Resource() string
	TTL() time.Duration
}

// Lock is used to set a record in redis and tries until it gets its goal
func (m *mux) Lock() {
	m.now = time.Now()

	m.value = <-utils.ID

	for i := 0; i < m.retries; i++ {
		res := aredis.Client.SetNX(m.resource, m.value, m.TTL())
		if ok, err := res.Result(); ok == false || err != nil {
			time.Sleep(tryCoolDown)
			continue
		}
		break
	}
}

// Unlock tries to get the record from redis and tries until it can
func (m *mux) Unlock() {
	h := unlockScript

	cmd := aredis.Client.Eval(h, []string{m.resource}, m.value)
	if cmd.Err() != nil {
		logrus.Warn("unlock failed with error :", cmd.Err())
	}
}

// Resource returns resource for no reason
func (m *mux) Resource() string {
	return m.resource
}

// TTL returns the duration of a lock
func (m *mux) TTL() time.Duration {
	return m.ttl
}

// NewRedisDistributedLock returns interface of a redlock
func NewRedisDistributedLock(resource string, ttl time.Duration) LockInterface {
	return &mux{
		retries:  int((ttl / tryCoolDown) + 1),
		resource: resource,
		ttl:      ttl,
	}
}
