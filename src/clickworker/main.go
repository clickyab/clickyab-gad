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
	version.PrintVersion().Info("Application started")
	models.Initialize()
	rabbit.Initialize()
	aredis.Initialize()

	exit := make(chan chan struct{})

	go func() {
		err := rabbit.RunWorker(
			config.Config.AMQP.Exchange,
			"cy.click",
			"cy_click_queue",
			&transport.Click{},
			clickWorker,
			10,
			exit,
		)
		assert.Nil(err)
	}()

	utils.WaitSignal(exit)
	logrus.Info("goodbye")
}
