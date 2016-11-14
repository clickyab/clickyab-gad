package main

import (
	"assert"
	"fmt"
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
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.USER, strconv.FormatInt(in.CopID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	// increment click to campaign
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN, strconv.FormatInt(in.CampaignID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	// increment click to ad
	_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.ADVERTISE, strconv.FormatInt(in.AdID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	if in.Web != nil {

		// increment click to slot
		_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.Web.SlotID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment click to website
		_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.WEBSITE, strconv.FormatInt(in.Web.WebsiteID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the campaign-slot
		_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN_SLOT, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.Web.SlotID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the campaign-website
		_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.CAMPAIGN_WEBSITE, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the ad-slot
		_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.AD_SLOT, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.Web.SlotID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the ad-website
		_, err = utils.IncKeyDaily(utils.KeyGenDaily(transport.AD_WEBSITE, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)
	}

	return true, nil
}
