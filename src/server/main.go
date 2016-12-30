package main

import (
	"config"
	"models"
	"modules"
	"os"
	"os/signal"
	"rabbit"
	"redis"
	"syscall"
	"utils"
	"version"

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
