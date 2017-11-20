package main

import (
	"clickyab.com/gad/version"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"

	_ "github.com/clickyab/services/fluentd"

	"os"

	"clickyab.com/gad/modules"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
)

var port = config.RegisterInt("port.echo", 80, "port of app")

func main() {
	config.Initialize("clickyab", "gad", "GAD")
	config.DumpConfig(os.Stdout)
	version.PrintVersion().Info("Application started")
	defer initializer.Initialize()()

	server := modules.Initialize("/")
	go func() {
		assert.Nil(server.Start(fmt.Sprintf(":%d", port.Int())))
	}()

	sig := shell.WaitExitSignal()
	logrus.Infof("goodbye (%s received)", sig)
}
