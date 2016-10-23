package main

import (
	"fmt"
	"transport"
)

func impWorker(in *transport.Impression) (bool, error) {
	fmt.Println("IM HERE", in)
	return true, nil
}
