package main

import (
	"assert"
	"config"
	"models"
	"rabbit"
	"transport"
	"utils"
	"version"

	"redis"

	"github.com/sirupsen/logrus"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	config.Config.AMQP.Publisher = 1 // Do not waste many publisher channel
	version.PrintVersion().Info("Application started")
	models.Initialize()
	rabbit.Initialize()
	defer rabbit.Finalize()
	aredis.Initialize()

	go func() {
		err := rabbit.RunWorker(
			&transport.Warning{},
			warnWorker,
			10,
		)
		assert.Nil(err)
	}()

	utils.WaitSignal(nil)
	logrus.Info("goodbye")
}
