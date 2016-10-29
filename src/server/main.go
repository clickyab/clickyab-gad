package main

import (
	"config"
	"models"
	"modules"
	"rabbit"
	"version"

	"github.com/labstack/echo/engine/fasthttp"
	"redis"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	version.PrintVersion().Info("Application started")
	models.Initialize()
	aredis.Initialize()
	rabbit.Initialize()

	_ = modules.Initialize(config.Config.MountPoint).Run(fasthttp.New(config.Config.Port))
}
