package main

import (
	"config"
	"fmt"
	"redis"
	"time"
)

func main() {
	config.Initialize()

	aredis.Initialize()

	fmt.Println(aredis.BRPopSingle("x", time.Second))
}
