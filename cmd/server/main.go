package main

import (
	"clickyab.com/gad/version"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"

	_ "github.com/clickyab/services/fluentd"

	"os"

	"clickyab.com/gad/modules"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	"github.com/clickyab/services/shell"
	"github.com/sirupsen/logrus"
	"gopkg.in/fzerorubigd/onion.v2"
)

var port = config.RegisterInt("port", 80, "port of app")

func main() {
	envLayer := onion.NewEnvLayer("PORT")
	config.Initialize("clickyab", "gad", "GAD", envLayer)
	config.DumpConfig(os.Stdout)
	version.PrintVersion().Info("Application started")
	defer initializer.Initialize()()

	server := modules.Initialize("/")
	go func() {
		_ = server.Start(fmt.Sprintf(":%d", port.Int()))
	}()

	sig := shell.WaitExitSignal()
	logrus.Infof("getting %s signal, bye...", sig)

}
