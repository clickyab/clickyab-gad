package main

import (
	"config"
	"models"
	"os"
	"os/signal"
	"rabbit"
	"syscall"
	"time"
	"transport"
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

	rabbit.RunWorker(
		config.Config.AMQP.Exchange,
		"cy.imp",
		"cy_imp_queue",
		&transport.Impression{},
		impWorker,
		10,
		exit,
	)

	quit := make(chan os.Signal, 5)
	signal.Notify(quit, syscall.SIGABRT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	<-quit

	tmp := make(chan struct{})
	exit <- tmp

	<-tmp
	logrus.Info("goodbye")
}
