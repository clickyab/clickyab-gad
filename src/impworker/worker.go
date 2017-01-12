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
	_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.USER, strconv.FormatInt(in.CopID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	// increment click to campaign
	_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.CAMPAIGN, strconv.FormatInt(in.CampaignID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	// increment click to ad
	_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.ADVERTISE, strconv.FormatInt(in.AdID, 10)), prefix+transport.SUBKEY_IMP, 1)
	assert.Nil(err)

	if in.Web != nil {

		// increment click to slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.Web.SlotID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment click to website
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.WEBSITE, strconv.FormatInt(in.Web.WebsiteID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the campaign-slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.CAMPAIGN_SLOT, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.Web.SlotID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the campaign-website
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.CAMPAIGN_WEBSITE, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the ad-slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.AD_SLOT, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.Web.SlotID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the ad-website
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.AD_WEBSITE, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the user website kry in redis
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.USER_WEBSITE, fmt.Sprintf("%d%s%d", in.CopID, transport.DELIMITER, in.Web.WebsiteID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)
	}

	if in.App != nil {

		// increment click to slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.SLOT, strconv.FormatInt(in.App.SlotID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment click to website
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.APP, strconv.FormatInt(in.App.AppID, 10)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the campaign-slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.CAMPAIGN_SLOT, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.App.SlotID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the campaign-website
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.CAMPAIGN_APP, fmt.Sprintf("%d%s%d", in.CampaignID, transport.DELIMITER, in.App.AppID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the ad-slot
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.AD_SLOT, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.App.SlotID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the ad-website
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.AD_APP, fmt.Sprintf("%d%s%d", in.AdID, transport.DELIMITER, in.App.AppID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)

		// increment the user website kry in redis
		_, err = utils.IncKeyDaily(transport.KeyGenDaily(transport.USER_APP, fmt.Sprintf("%d%s%d", in.CopID, transport.DELIMITER, in.App.AppID)), prefix+transport.SUBKEY_IMP, 1)
		assert.Nil(err)
	}

	return true, nil
}
