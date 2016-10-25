package main

import (
	"assert"
	"strconv"
	"transport"
	"utils"
)

func impWorker(in *transport.Impression) (bool, error) {
	prefix := ""
	if in.Suspicious {
		prefix = transport.FRAUD_PREFIX
	}
	var err error
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.USER, in.User), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN, strconv.FormatInt(in.CampaignID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.ADVERTISE, strconv.FormatInt(in.AdID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.Web.SlotID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.WEBSITE, strconv.FormatInt(in.Web.WebsiteID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	return true, nil
}
