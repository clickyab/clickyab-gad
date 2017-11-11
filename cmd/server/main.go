package main

import (
	"os"
	"os/signal"
	"syscall"

	"clickyab.com/gad/config"
	"clickyab.com/gad/modules"
	"clickyab.com/gad/utils"
	"clickyab.com/gad/version"
	_ "github.com/clickyab/services/mysql/connection/mysql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"

	"context"

	"clickyab.com/gad/fluentd"

	"github.com/clickyab/services/initializer"
	"github.com/pkg/profile"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	defer profile.Start(profile.CPUProfile, profile.NoShutdownHook, profile.ProfilePath("./tmp/"+<-utils.ID)).Stop()

	ctx, cancel := context.WithCancel(context.Background())

	version.PrintVersion().Info("Application started")
	fluentd.Initialize(ctx)
	defer initializer.Initialize()()

	server := modules.Initialize(config.Config.MountPoint)
	go func() {
		_ = server.Start(fmt.Sprintf(":%d", config.Config.Port))
	}()

	sig := make(chan os.Signal, 6)

	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP)
	<-sig
	cancel()
}
