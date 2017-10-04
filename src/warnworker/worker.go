// Package main
package main

import (
	"transport"

	"github.com/sirupsen/logrus"
)

// error means Ack/Nack the boolean maens only when error is not nil, and means re-queue
func warnWorker(in *transport.Warning) (bool, error) {
	// Simply log the error
	logrus.WithFields(
		logrus.Fields{
			"Level":   in.Level,
			"When":    in.When,
			"Where":   in.Where,
			"Message": in.Message,
		},
	).Warn(string(in.Request))

	return false, nil
}
