package main

import (
	"config"
	"modules"
	"time"
	"version"

	"models"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo/engine/fasthttp"
)

func main() {

	config.Initialize()
	ver := version.GetVersion()
	//logrus.SetLevel(logrus.PanicLevel)
	logrus.WithFields(
		logrus.Fields{
			"Commit hash":       ver.Hash,
			"Commit short hash": ver.Short,
			"Commit date":       ver.Date.Format(time.RFC3339),
			"Build date":        ver.BuildDate.Format(time.RFC3339),
		},
	).Infof("Application started")

	models.Initialize()
	_ = modules.Initialize(config.Config.MountPoint).Run(fasthttp.New(config.Config.Port))
}
