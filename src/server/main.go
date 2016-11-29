package main

import (
	"config"
	"models"
	"modules"
	"os"
	"os/exec"
	"os/signal"
	"rabbit"
	"redis"
	"syscall"
	"utils"
	"version"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo/engine/standard"
	"github.com/pkg/profile"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	defer profile.Start(profile.CPUProfile, profile.NoShutdownHook, profile.ProfilePath("./tmp/"+<-utils.ID)).Stop()

	version.PrintVersion().Info("Application started")
	models.Initialize()
	aredis.Initialize()
	rabbit.Initialize()

	server := modules.Initialize(config.Config.MountPoint)
	go func() {
		_ = server.Run(standard.New(config.Config.Port))
	}()

	sig := make(chan os.Signal, 6)

	signal.Notify(sig, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP)
	s := <-sig
	if s == syscall.SIGHUP {
		_ = server.Stop()
		var args []string
		if len(os.Args) > 1 {
			args = os.Args[1:]
		}
		cmd := exec.Command(os.Args[0], args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		if err := cmd.Start(); err != nil {
			logrus.Fatalf("Restart: Failed to launch, error: %v", err)
		}
	}
}
