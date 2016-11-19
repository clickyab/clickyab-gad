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
	"assert"
	"config"
	"fmt"
	"mr"
	"redis"
	"strconv"
	"transport"
	"utils"
)

// error means Ack/Nack the boolean maens only when error is not nil, and means re-queue
func clickWorker(in *transport.Click) (bool, error) {

	//Validation Click TODO : should be changed using redis

	// increment click to user
	prefix := ""
	if in.FraudReason != 0 {
		prefix = transport.FRAUD_PREFIX
	}
	var err error
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.USER, fmt.Sprintf("%d", in.CopID)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment click to campaign
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN, strconv.FormatInt(in.CampaignID, 10)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment click to ad
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.ADVERTISE, strconv.FormatInt(in.AdID, 10)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment click to slot
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.SlotID, 10)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment click to website
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.WEBSITE, strconv.FormatInt(in.Web.WebsiteID, 10)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment the campaign-slot
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN_SLOT, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.SlotID)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment the campaign-website
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN_WEBSITE, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment the ad-slot
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.AD_SLOT, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.SlotID)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment the ad-website
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.AD_WEBSITE, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	// increment the user website kry in redis
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.USER_WEBSITE, fmt.Sprintf("%d%s%d", in.CopID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_Cl, 1)
	assert.Nil(err)

	//insert click in db
	err = mr.NewManager().InsertClick(in)
	if err != nil {
		return false, err
	}
	err = aredis.HMSet(fmt.Sprintf("%s%s%s", transport.CONV, transport.DELIMITER, in.Rand), config.Config.Clickyab.DailyClickExpire, "CLICK", in.ID)
	return false, err
}
