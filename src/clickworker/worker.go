package main

import (
	"assert"
	"strconv"
	"transport"
	"utils"
)

func clickWorker(in *transport.Click) bool {

	// increment click to user
	prefix := ""
	if in.Suspicious {
		prefix = transport.FRAUD_PREFIX
	}
	var err error
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.USER, in.User), prefix+transport.SUBKEY_Cl, 1)
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

	// persist in mysql database
	return true
}
