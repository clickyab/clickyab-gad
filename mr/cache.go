package mr

import (
	"time"

	"clickyab.com/gad/redis"
	"clickyab.com/gad/utils"
)

func store(key string, in interface{}, d time.Duration) error {
	t, err := utils.InterfaceToByte(in)
	if err != nil {
		return err
	}

	return aredis.StoreKey("C_"+key, string(t), d)
}

func fetch(key string, in interface{}) error {
	s, err := aredis.GetKey("C_"+key, false, 0)
	if err != nil {
		return err
	}

	return utils.ByteToInterface([]byte(s), in)
}

func fetchTouch(key string, in interface{}, expire time.Duration) error {
	s, err := aredis.GetKey("C_"+key, true, expire)
	if err != nil {
		return err
	}

	return utils.ByteToInterface([]byte(s), in)
}
