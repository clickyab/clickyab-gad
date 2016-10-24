package main

import (
	"assert"
	"config"
	"models"
	"rabbit"
	"time"
	"transport"
	"utils"
	"version"

	"github.com/Sirupsen/logrus"
)

func main() {
	config.Initialize()
	ver := version.GetVersion()
	//logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(
		logrus.Fields{
			"Commit hash":       ver.Hash,
			"Commit short hash": ver.Short,
			"Commit date":       ver.Date.Format(time.RFC3339),
			"Build date":        ver.BuildDate.Format(time.RFC3339),
		},
	).Infof("Application started")
	models.Initialize()
	rabbit.Initialize()

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
		assert.Nil(err)
	}()

	utils.WaitSignal(exit)
	logrus.Info("goodbye")
}
