package main

import (
	"config"
	"models"
	"modules"
	"rabbit"
	"version"

	"redis"

	"os"

	"os/signal"
	"syscall"

	"utils"

	"github.com/labstack/echo/engine/fasthttp"
	"github.com/pkg/profile"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	defer profile.Start(profile.CPUProfile, profile.NoShutdownHook, profile.ProfilePath("./"+<-utils.ID))

	version.PrintVersion().Info("Application started")
	models.Initialize()
	aredis.Initialize()
	rabbit.Initialize()

	go func() {
		_ = modules.Initialize(config.Config.MountPoint).Run(fasthttp.New(config.Config.Port))
	}()

	sig := make(chan os.Signal, 5)

	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-sig
}
