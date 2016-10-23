package main

import (
	"fmt"
	"transport"
)

func impWorker(in *transport.Impression) bool {
	fmt.Println("IM HERE", in)
	return true
}
