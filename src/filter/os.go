package filter

import (
	"selector"
	"mr"
	"platform"

)

func CheckOS(c *selector.Context,in mr.AdData) bool  {

	if in.CpPlatforms == nil{
		return true
	}
	Os,err := platform.FindIdOs(c.RequestData)
	if err != nil{
		return false
	}
	for _,v := range in.CpPlatforms{
		if v == int64(Os){
			return true
		}
	}
	return false
}