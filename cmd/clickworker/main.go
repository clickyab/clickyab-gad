package main

import (
	"clickyab.com/gad/rabbit"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/version"
	_ "github.com/clickyab/services/mysql/connection/mysql"

	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

func main() {
	config.Initialize("clickyab", "gad", "GAD")

	version.PrintVersion().Info("Application started")
	defer initializer.Initialize()()

	go func() {
		err := rabbit.RunWorker(
			&transport.Click{},
			clickWorker,
			10,
		)
		assert.Nil(err)
	}()

	sig := shell.WaitExitSignal()
	logrus.Infof("goodbye (%s received)", sig)
}
