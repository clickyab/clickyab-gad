package aredis

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/healthz"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/safe"
	"github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var (
	// Client the actual pool to use with redis
	Client *redis.Client
	all    []initializer.Simple
	lock   sync.RWMutex
)

type initRedis struct {
}

// Healthy return true if the databases are ok and ready for ping
func (initRedis) Healthy(context.Context) error {
	ping, err := Client.Ping().Result()
	if err != nil || strings.ToUpper(ping) != "PONG" {
		return fmt.Errorf("Redis PING failed. result was '%s' and the error was %s", ping, err)
	}

	return nil
}

// Initialize try to create a redis pool
func (i *initRedis) Initialize(ctx context.Context) {
	Client = redis.NewClient(
		&redis.Options{
			Network:  networkType.String(),
			Addr:     address.String(),
			Password: password.String(),
			PoolSize: poolsize.Int(),
			DB:       db.Int(),
		},
	)
	// PING the server to make sure every thing is fine
	safe.Try(func() error { return Client.Ping().Err() }, tryLimit.Duration())

	healthz.Register(i)

	for i := range all {
		all[i].Initialize()
	}
	logrus.Debug("redis is ready.")
	go func() {
		c := ctx.Done()
		assert.NotNil(c, "[BUG] context has no mean to cancel/deadline/timeout")
		<-c
		assert.Nil(Client.Close())
		logrus.Debug("redis finalized.")
	}()
}

// Register a new object to inform it after redis is loaded
func Register(in initializer.Simple) {
	lock.Lock()
	defer lock.Unlock()

	all = append(all, in)
}

func init() {
	initializer.Register(&initRedis{}, 0)
}
