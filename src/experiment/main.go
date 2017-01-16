package main

import (
	"config"
	"version"

	"assert"
	"net"

	"cache"
	"fmt"
	"time"
)

func main() {
	config.Initialize()
	config.SetConfigParameter()
	version.PrintVersion().Info("Application started")

	p, err := cache.NewPeerSelector(net.IPv4(4, 4, 4, 4), 1000)
	assert.Nil(err)
	for range time.NewTicker(time.Second).C {
		fmt.Println("---")
		fmt.Println(p.AllPeers())
		fmt.Println("===")
	}
}
