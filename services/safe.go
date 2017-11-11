package services

import (
	"runtime/debug"

	"time"

	"github.com/sirupsen/logrus"
)

func actual(a func() error) (res error) {
	defer func() {
		if e := recover(); e != nil {
			stack := debug.Stack()
			logrus.Debug(string(stack))

		}
	}()
	res = a()
	return
}

// GoRoutine is a safe go routine system with recovery and a way to inform finish of the routine
func GoRoutine(f func(), extra ...interface{}) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				stack := debug.Stack()
				logrus.Debug(string(stack))
			}
		}()

		f()
	}()
}

// Try retry by fibonacci way the given function
func Try(a func() error, max time.Duration, extra ...interface{}) {
	x, y := 0, 1
	for {
		err := actual(a)
		if err == nil {
			return
		}
		logrus.Error(err)
		t := time.Duration(x) * time.Second
		if t < max {
			x, y = y, x+y
		}
		time.Sleep(t)
	}

}
