package main

import (
	"assert"
	"cache"
	_ "cache/redis"
	"config"
	"fmt"
	"services/redis"
	"time"
)

type Test struct {
	X int
	Y string
	Z float64
}

func main() {
	config.Initialize()
	aredis.Initialize()

	t := Test{
	//X: 11111,
	//Y: "s;wlswijdoiw",
	//Z: 1000.999,
	}

	cc := cache.CreateCacheWrapper("SSS", &t)
	err := cache.Hit("SSS", cc)
	if err == nil {
		fmt.Println("cache hit")
		fmt.Printf("%+v", t)
	}

	err = cache.Cache(cc, time.Hour, nil)
	assert.Nil(err)

}
