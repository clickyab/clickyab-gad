package filter

import "entity"

func FilterPublisherWhiteList(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	blacklist := ad.BlackListPublisher()
	elem := impression.Source().ID()

	return hasOneInMany(true, blacklist, elem)
}

func FilterWebCategoty(impression entity.Impression, advertise entity.Advertise) bool {
	if impression.Source().Type() == entity.PublisherTypeWeb {
		return false
	}

	ok, imp, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	whitelist := ad.WebCategory()
	elems := imp.Category()

	if len(whitelist) < len(elems) {
		return false
	}

	return hasManyInMany(true, whitelist, elems)
}

func FilterAppCategoty(impression entity.Impression, advertise entity.Advertise) bool {
	if impression.Source().Type() == entity.PublisherTypeApp {
		return false
	}

	ok, imp, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	whitelist := ad.WebCategory()
	elems := imp.Category()

	if len(whitelist) < len(elems) {
		return false
	}

	return hasManyInMany(true, whitelist, elems)
}

func FilterAppSize(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	if impression.Source().Type() == entity.PublisherTypeApp {
		return false
	}

	var sizes []int
	for _, i := range impression.Slots() {
		sizes = append(sizes, i.Size())
	}

	return hasOneInMany(true, sizes, ad.Size())
}

func FilterWebSize(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	if impression.Source().Type() == entity.PublisherTypeWeb {
		return false
	}

	var sizes []int
	for _, i := range impression.Slots() {
		sizes = append(sizes, i.Size())
	}

	return hasOneInMany(true, sizes, ad.Size())
}

func FilterVastSize(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	if impression.Source().Type() == entity.PublisherTypeVast {
		return false
	}

	var sizes []int
	for _, i := range impression.Slots() {
		sizes = append(sizes, i.Size())
	}

	return hasOneInMany(true, sizes, ad.Size())
}

func FilterHood(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	elems := impression.Location().Hood().ID
	whitelist := ad.Hood()

	return hasOneInMany(true, whitelist, elems)
}

func FilterCountry(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	country := impression.Location().Country().ID
	elem := ad.Country()

	if elem == country {
		return true
	}
	return false
}

func FilterProvince(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	province := impression.Location().Province().ID
	elem := ad.Country()

	if elem == province {
		return true
	}
	return false
}

func FilterCity(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	country := impression.Location().Country().ID
	elem := ad.Country()

	if elem == country {
		return true
	}
	return false
}

/*
func FilterLatLon(impression entity.Impression, advertise entity.Advertise) bool {
	ok, _, ad := clickImp(impression, advertise)
	if !ok {
		return false
	}

	implat := impression.Location().LatLon().Lat
	implon := impression.Location().LatLon().Lon

	adlat, adlon := ad.LanLon()

	if (){
		return false
	}

	return true
}*/
