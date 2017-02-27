package filter

import (
	"entity"
	"reflect"
)

type filterFunc func(impression entity.Impression, advertise entity.Advertise) bool

func clickImp(impression entity.Impression, advertise entity.Advertise) (bool, entity.ClickyabImp, entity.ClickyabAd) {
	ad, ok := advertise.(entity.ClickyabAd)
	imp, ok2 := impression.(entity.ClickyabImp)
	if !ok || !ok2 {
		return false, nil, nil
	}
	return true, imp, ad
}

func hasOneInMany(con bool, slice interface{}, elem interface{}) bool {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		panic("not a slice kind")
	}
	if reflect.TypeOf(slice).Elem() == reflect.TypeOf(elem) {
		panic("types doesnt match")
	}

	c := reflect.ValueOf(slice)
	for i := 0; i < c.Len(); i++ {
		if c.Index(i) == elem {
			return con
		}
	}
	return !con
}

func hasManyInMany(con bool, s1 interface{}, s2 interface{}) bool {
	if reflect.TypeOf(s1) != reflect.TypeOf(s2) {
		panic("difrent types passed")
	}
	first := reflect.ValueOf(s1)
	second := reflect.ValueOf(s2)

	var a, b reflect.Value
	if first.Len() > second.Len() {
		a = first
		b = second
	} else {
		b = first
		a = second
	}

outer:
	for i := 0; i < a.Len(); i++ {
		for j := 0; j < b.Len(); j++ {
			if a.Index(i) == b.Index(j) {
				continue outer
			}
		}
		return !con
	}

	return !con
}
