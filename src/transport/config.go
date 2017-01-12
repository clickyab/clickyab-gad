package transport

const (
	// DELIMITER is the delimiter in all redis key
	DELIMITER = "_"
	// USER is for user keys
	USER = "U"
	// CAMPAIGN for campaign keys
	CAMPAIGN = "CP"
	// ADVERTISE for advertise keys
	ADVERTISE = "AD"
	// SLOT for slot keys
	SLOT = "S"
	// WEBSITE is for website keys
	WEBSITE = "W"
	// APP is for app keys
	APP = "A"

	// FRAUD_PREFIX for all fraud key
	FRAUD_PREFIX = "F"
	// SUBKEY_IMP is for impression key inside each master hash
	SUBKEY_IMP = "I"
	// SUBKEY_Cl is for click count
	SUBKEY_Cl = "C"
	// CAMPAIGN_SLOT is the campaign related stuff
	CAMPAIGN_SLOT = "CPS"
	// CAMPAIGN_WEBSITE is for campaign related to website
	CAMPAIGN_WEBSITE = "CPW"
	// CAMPAIGN_APP campaign relation to app keys
	CAMPAIGN_APP = "CPA"
	// AD_SLOT ad related to slots
	AD_SLOT = "ADS"
	// AD_WEBSITE ads related to websites
	AD_WEBSITE = "ADW"
	// AD_APP is ad related to apps
	AD_APP = "ADA"
	// USER_WEBSITE is user realted to website
	USER_WEBSITE = "UW"
	// USER_APP is user related to apps
	USER_APP = "UA"

	USER_CAPPING     = "CAP"
	USER_RETARGETING = "RET"
	IMP              = "IMP"
	CLICK            = "CLK"
	MEGA             = "MGA"
	CONV             = "CNV"
)
