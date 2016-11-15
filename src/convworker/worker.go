package main

import (
	"config"
	"errors"
	"fmt"
	"mr"
	"redis"
	"time"
	"transport"
)

// error means Ack/Nack the boolean maens only when error is not nil, and means re-queue
func convWorker(in *transport.Conversion) (bool, error) {
	m, err := aredis.HGetAllString(fmt.Sprintf("%s%s%s", transport.CONV, transport.DELIMITER, in.ConvID), true, config.Config.Clickyab.DailyClickExpire)
	if err != nil {
		return false, err
	}

	if c, ok := m["CLICK"]; ok {
		// TODO : can we use imp for fraud conversion :))))
		err = mr.NewManager().InsertConversion(c, in.ActionID)
		if err != nil {
			return false, err
		}
		_, _ = aredis.IncHash(fmt.Sprintf("%s%s%s", transport.CONV, transport.DELIMITER, in.ConvID), "DONE", 1, 0)
		return false, nil
	}

	count, _ := aredis.IncHash(fmt.Sprintf("%s%s%s", transport.CONV, transport.DELIMITER, in.ConvID), "OK", 1, config.Config.Clickyab.DailyClickExpire)
	if count > config.Config.Clickyab.ConvRetry {
		return false, errors.New("limit is done")
	}
	time.Sleep(config.Config.Clickyab.ConvDelay)
	return true, errors.New("no data yet")
}
