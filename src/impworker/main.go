package main

import (
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
			&transport.Impression{},
			impWorker,
			10,
		)
		if err != nil {
			// Fatal is only allowed in main
			logrus.Fatal(err)
		}
	}()

	utils.WaitSignal(nil)
	logrus.Info("goodbye")
}
