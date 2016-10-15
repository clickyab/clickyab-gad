package modules

import (
	"io"

	"github.com/Sirupsen/logrus"
	ll "github.com/labstack/echo/log"
	"github.com/labstack/gommon/log"
)

type myLogger struct {
	logger *logrus.Logger
}

func (m *myLogger) SetOutput(o io.Writer) {
	m.logger.Out = o
}
func (m *myLogger) SetLevel(l log.Lvl) {
	switch l {
	case log.DEBUG:
		m.logger.Level = logrus.DebugLevel
	case log.INFO:
		m.logger.Level = logrus.InfoLevel
	case log.WARN:
		m.logger.Level = logrus.WarnLevel
	case log.ERROR:
		m.logger.Level = logrus.ErrorLevel
	case log.FATAL:
		m.logger.Level = logrus.FatalLevel
	case log.OFF:
		m.logger.Level = logrus.PanicLevel
	}
}
func (m *myLogger) Print(args ...interface{}) {
	m.logger.Print(args...)
}
func (m *myLogger) Printf(s string, a ...interface{}) {
	m.logger.Panicf(s, a...)
}
func (m *myLogger) Printj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Print("-")
}

func (m *myLogger) Debug(args ...interface{}) {
	m.logger.Debug(args...)
}
func (m *myLogger) Debugf(s string, args ...interface{}) {
	m.logger.Debugf(s, args...)
}
func (m *myLogger) Debugj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Debug("-")
}
func (m *myLogger) Info(args ...interface{}) {
	m.logger.Info(args...)
}
func (m *myLogger) Infof(s string, args ...interface{}) {
	m.logger.Infof(s, args...)
}
func (m *myLogger) Infoj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Info("-")
}
func (m *myLogger) Warn(args ...interface{}) {
	m.logger.Warn(args)
}
func (m *myLogger) Warnf(s string, args ...interface{}) {
	m.logger.Warnf(s, args...)
}
func (m *myLogger) Warnj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Warn("-")
}
func (m *myLogger) Error(args ...interface{}) {
	m.logger.Error(args...)
}
func (m *myLogger) Errorf(s string, args ...interface{}) {
	m.logger.Errorf(s, args...)
}
func (m *myLogger) Errorj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Error("-")
}
func (m *myLogger) Fatal(args ...interface{}) {
	m.logger.Fatal(args...)
}
func (m *myLogger) Fatalj(j log.JSON) {
	l := logrus.Fields{}
	for i := range j {
		l[i] = j[i]
	}

	m.logger.WithFields(l).Fatal("-")
}
func (m *myLogger) Fatalf(s string, args ...interface{}) {
	m.logger.Fatalf(s, args...)
}

// NewLogger function that  logs events on show ads
func NewLogger() ll.Logger {
	return &myLogger{logger: logrus.New()}
}
