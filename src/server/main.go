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

	"cache"
	"fmt"
	"net"

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
		cache.Initialize(net.ParseIP(config.Config.IPString), config.Config.Port, server)
		_ = server.Start(fmt.Sprintf(":%d", config.Config.Port))
	}()

	sig := make(chan os.Signal, 6)

	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-sig
}
