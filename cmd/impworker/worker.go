package main

import (
	"strconv"

	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"
)

func impWorker(in *transport.Impression) (bool, error) {
	prefix := ""
	if in.Suspicious {
		prefix = transport.FraudPrefix
	}
	var err error
	if in.Web != nil {
		// increment click to slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.Slot, strconv.FormatInt(in.Web.SlotID, 10)), prefix+transport.ImpSubKey, 1)
		assert.Nil(err)
	}

	if in.App != nil {
		// increment click to slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.Slot, strconv.FormatInt(in.App.SlotID, 10)), prefix+transport.ImpSubKey, 1)
		assert.Nil(err)
	}

	return true, nil
}
