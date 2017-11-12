package main

import (
	"clickyab.com/gad/rabbit"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
	"clickyab.com/gad/version"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Initialize("clickyab", "gad", "GAD")

	version.PrintVersion().Info("Application started")
	defer initializer.Initialize()()

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
