package mr

import (
	"redis"
	"time"
	"utils"
)

func store(key string, in interface{}, d time.Duration) error {
	t, err := utils.InterfaceToByte(in)
	if err != nil {
		return err
	}

	return aredis.StoreKey("CACHE_"+key, string(t), d)
}

func fetch(key string, in interface{}) error {
	s, err := aredis.GetKey("CACHE_"+key, false, 0)
	if err != nil {
		return err
	}

	return utils.ByteToInterface([]byte(s), in)
}

func fetchTouch(key string, in interface{}, expire time.Duration) error {
	s, err := aredis.GetKey("CACHE_"+key, true, expire)
	if err != nil {
		return err
	}

	return utils.ByteToInterface([]byte(s), in)
}
