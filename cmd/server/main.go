package main

import (
	"clickyab.com/gad/config"
	"clickyab.com/gad/models"
	"clickyab.com/gad/modules"
	"os"
	"os/signal"
	"clickyab.com/gad/rabbit"
	"clickyab.com/gad/redis"
	"syscall"
	"clickyab.com/gad/utils"
	"clickyab.com/gad/version"

	"fmt"

	"github.com/pkg/profile"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	defer profile.Start(profile.CPUProfile, profile.NoShutdownHook, profile.ProfilePath("./tmp/"+<-utils.ID)).Stop()

	version.PrintVersion().Info("Application started")
	aredis.Initialize()
	rabbit.Initialize()
	models.Initialize()

	server := modules.Initialize(config.Config.MountPoint)
	go func() {
		_ = server.Start(fmt.Sprintf(":%d", config.Config.Port))
	}()

	sig := make(chan os.Signal, 6)

	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP)
	<-sig
}
