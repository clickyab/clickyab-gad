package rtb

import (
	"assert"
	"config"
	"eav"
	"eav/redis"
	"entity"
	"fmt"
	"time"
)

const (
	Mega          string = "MEGA_"
	MegaIP        string = "IP"
	MegaUserAgent string = "UA"
	MegaPubID     string = "PID"
	MegaTimeUnix  string = "TU"
)

func createMegaStore(imp entity.Impression) eav.Kiwi {
	kiwi := redis.NewRedisEAVStore(Mega + imp.MegaIMP())
	assert.Nil(kiwi.SubKey(MegaIP, imp.IP().String()).
		SubKey(MegaUserAgent, imp.UserAgent()).
		SubKey(MegaPubID, fmt.Sprint(imp.Source().GetID())).
		SubKey(MegaTimeUnix, fmt.Sprint(time.Now().Unix())).
		Save(config.Config.Clickyab.MegaImpExpire))
	return kiwi
}

// Select is the key function to select an ad for an imp base on real time biding
func Select(pub entity.Publisher, imp entity.Impression, ads map[int64][]entity.Advertise, slots []entity.Slot) {
	
}
