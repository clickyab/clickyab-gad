package filter

import (
	"selector"
	"mr"
)

func CheckWhiteList(c *selector.Context, in mr.AdData) bool  {
	if len(in.CpWfilter)==0{
		return true;
	}
	for v := range in.CpWfilter{
		if int64(v) == c.WID {
			return true
		}
	}
	return false
}
