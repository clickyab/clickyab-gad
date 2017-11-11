// Package main
// Fraud is for checking fraud based on this rules
//    ***Historical rules***
//1.   duplicate click: one user by same imps, get two click on one ads
//2.  unknown reference: when imp not have reference address
//9.  fast clicks: under 4 second click
//3.  Extra Cookie Active OR 3 clicks in month
//5.  One Month Cookie Block: after active "Extra Cookie Active" all click fault
//16.  There is no ad; ad id is not valid
//4.  total click 4 per day
//17.  one person before clicked on ads of same campaigns on day
//
//    ***new rules impression***
//    mega impression same in select & show
//
//
//    ***new rules click***
//    same impersion_id in impression & click
//    same ip in impression & click
//
package main

import (
	"github.com/clickyab/services/assert"
	"clickyab.com/gad/mr"
	"strconv"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"
)

// error means Ack/Nack the boolean maens only when error is not nil, and means re-queue
func clickWorker(in *transport.Click) (bool, error) {

	prefix := ""
	if in.FraudReason != 0 {
		prefix = transport.FRAUD_PREFIX
	}
	var err error
	// increment click to slot
	_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.SlotID, 10)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)
	//insert click in db
	err = mr.NewManager().InsertClick(in)
	return false, err
}
