package main

import (
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
			"cy.imp",
			"cy_imp_queue",
			&transport.Impression{},
			impWorker,
			10,
			exit,
		)
		if err != nil {
			// Fatal is only allowed in main
			logrus.Fatal(err)
		}
	}()

	utils.WaitSignal(exit)
	logrus.Info("goodbye")
}
