package main

import (
	"github.com/clickyab/services/assert"
	"strconv"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
)

func impWorker(in *transport.Impression) (bool, error) {
	prefix := ""
	if in.Suspicious {
		prefix = transport.FRAUD_PREFIX
	}
	var err error
	if in.Web != nil {
		// increment click to slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.Web.SlotID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)
	}

	if in.App != nil {
		// increment click to slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.App.SlotID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)
	}

	return true, nil
}
