package entity

import "context"

// AdProvider is the interface to handle ad in system base on impression
type AdProvider interface {
	// Provide is the function to handle the request, provider shoud response
	// to this function and return all eligible ads
	Provide(context.Context, Impression) <-chan map[int][]Advertise
}
