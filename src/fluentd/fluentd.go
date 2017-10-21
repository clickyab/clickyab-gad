package fluentd

import (
	"context"

	"config"
	"github.com/evalphobia/logrus_fluent"
	"github.com/sirupsen/logrus"
)

func Initialize(ctx context.Context) {

	var (
		active     = config.Config.Fluentd.Enable
		host       = config.Config.Fluentd.Host
		allLevels  = config.Config.Fluentd.All_levels
		defaultTag = config.Config.Fluentd.Tag
		port       = config.Config.Fluentd.Port
	)

	if !active {
		return
	}

	hook, err := logrus_fluent.New(host, int(port))
	if err != nil {
		logrus.Error("fluentd logger failed, if this is in production check for the problem")
	}

	// set custom fire level
	l := []logrus.Level{
		logrus.PanicLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
	if allLevels {
		l = append(l, logrus.DebugLevel)
	}
	hook.SetLevels(l)

	// set static tag
	hook.SetTag(defaultTag)
	// filter func
	// TODO : write more filter to handle more type (for clarification on logger side)
	hook.AddFilter("error", logrus_fluent.FilterError)
	logrus.AddHook(hook)

	go func() {
		<-ctx.Done()
		// somehow it can make a race, but there is no objection, its dying anyway
		tmp := hook.Fluent
		hook.Fluent = nil
		err := tmp.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
}
