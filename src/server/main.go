package main

import (
	"config"
	"models"
	"modules"
	"rabbit"
	"version"

	"redis"

	"github.com/labstack/echo/engine/fasthttp"
	"github.com/pkg/profile"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	if config.Config.DevelMode {
		defer profile.Start().Stop()
	}

	version.PrintVersion().Info("Application started")
	models.Initialize()
	aredis.Initialize()
	rabbit.Initialize()

	_ = modules.Initialize(config.Config.MountPoint).Run(fasthttp.New(config.Config.Port))
}
