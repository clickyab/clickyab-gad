package filter

import "entity"

func FilterPublisherBlackList(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	blacklist := ad.BlackListPublisher()
	elem := impression.Source().ID()

	return hasOneInMany(false, blacklist, elem)
}

func FilterOS(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	blacklist := ad.BlackListOS()
	elem := impression.OS().ID

	return hasOneInMany(false, blacklist, elem)
}

func FilterPublisherType(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	blacklist := ad.BlackListPublisherType()
	elem := impression.Source().Type()

	return hasOneInMany(false, blacklist, elem)
}

/*
func FilterAppBrand(impression entity.Impression, advertise entity.Advertise) bool {
	if !impression.OS().Mobile {
		return false
	}

	blacklist := entity.Advertise().Campaign().AllowedBrands()
	elem := impression.Source().
}*/
