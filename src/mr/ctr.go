package mr

import (
	"config"
	"fmt"
	"redis"
	"time"
	"transport"
)

func CalculateCTR(ad AdData) float64 {
	maxBid := ad.CpMaxbid
	fmt.Println(maxBid)
	ctrCalculateConst := config.Config.CtrConst
	for m := range ctrCalculateConst {
		switch ctrCalculateConst[m] {
		case transport.AD_SLOT:
			{
				mapVal, _ := aredis.HGetAll(transport.AD_SLOT, true, 72*time.Hour)
				fmt.Println(mapVal)
			}
		case transport.AD_WEBSITE:
			{
				mapVal, _ := aredis.HGetAll(transport.AD_WEBSITE, true, 72*time.Hour)
				fmt.Println(mapVal)
			}
		case transport.ADVERTISE:
			{
				mapVal, _ := aredis.HGetAll(transport.ADVERTISE, true, 72*time.Hour)
				fmt.Println(mapVal)
			}
		case transport.CAMPAIGN_WEBSITE:
			{
				mapVal, _ := aredis.HGetAll(transport.CAMPAIGN_WEBSITE, true, 72*time.Hour)
				fmt.Println(mapVal)
			}
		case transport.SLOT:
			{
				mapVal, _ := aredis.HGetAll(transport.SLOT, true, 72*time.Hour)
				fmt.Println(mapVal)
			}
		case transport.CAMPAIGN_SLOT:
			{
				mapVal, _ := aredis.HGetAll(transport.CAMPAIGN_SLOT, true, 72*time.Hour)
				fmt.Println(mapVal)
			}
		default:
			{
				CTR := .1
				fmt.Println(CTR)
			}
		}
	}
	return .3243
}
