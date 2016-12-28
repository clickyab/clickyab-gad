package main

import (
	"config"
	"errors"
	"fmt"
	"mr"
	"redis"
	"strconv"
	"time"
	"transport"
)

// error means Ack/Nack the boolean maens only when error is not nil, and means re-queue
func convWorker(in *transport.Conversion) (bool, error) {

	_, err := strconv.ParseInt(in.ConvID, 10, 0)
	if err == nil { // the conversion is the integer value, th old syste,
		return false, mr.NewManager().InsertConversion(in.ConvID, in.ActionID)
	}

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
