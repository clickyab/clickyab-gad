package filter

import (
	"selector"
	"mr"
)

func CheckWhiteList(c *selector.Context, in mr.AdData) bool  {
	if len(in.CpWfilter)==0{
		return true;
	}
	for _, v := range in.CpWfilter{
		if v == c.WID {
			return true
		}
	}
	return false
}
