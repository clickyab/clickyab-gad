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

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	config.Config.AMQP.Publisher = 1 // Do not waste many publisher channel
	version.PrintVersion().Info("Application started")
	models.Initialize()
	rabbit.Initialize()
	aredis.Initialize()

	exit := make(chan chan struct{})

	go func() {
		err := rabbit.RunWorker(
			config.Config.AMQP.Exchange,
			"cy.warn",
			"cy_warn_queue",
			&transport.Warning{},
			warnWorker,
			10,
			exit,
		)
		assert.Nil(err)
	}()

	utils.WaitSignal(exit)
	logrus.Info("goodbye")
}
