package aredis

import (
	"assert"
	"config"
	"sync"

	"github.com/Sirupsen/logrus"
	redis "gopkg.in/redis.v5"
)

var (
	// Pool the actual pool to use with redis
	Client *redis.Client
	once   = &sync.Once{}
)

// Initialize try to create a redis pool
func Initialize() {
	once.Do(func() {
		Client = redis.NewClient(
			&redis.Options{
				Network:  config.Config.Redis.Network,
				Addr:     config.Config.Redis.Address,
				Password: config.Config.Redis.Password,
				PoolSize: config.Config.Redis.Size,
				DB:       config.Config.Redis.Databse,
			},
		)
		// PING the server to make sure every thing is fine
		assert.Nil(Client.Ping().Err())
		logrus.Debug("redis is ready.")
	})
}
