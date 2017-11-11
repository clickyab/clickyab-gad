package main

import (
	"clickyab.com/gad/assert"
	"clickyab.com/gad/config"
	"clickyab.com/gad/models"
	"clickyab.com/gad/rabbit"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
	"clickyab.com/gad/version"
	_ "github.com/go-sql-driver/mysql"

	"clickyab.com/gad/redis"

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
			&transport.Click{},
			clickWorker,
			10,
		)
		assert.Nil(err)
	}()

	utils.WaitSignal(nil)
	logrus.Info("goodbye")
}
