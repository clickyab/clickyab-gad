package entity

import "net"

// ImpressionAttributes is the imp attr key
type ImpressionAttributes string

// Impression is the single impression object
type Impression interface {
	// MegaIMP return the random id of this imp object
	MegaIMP() string
	// ClientID is the key to identify client
	ClientID() string
	// IP return the client ip
	IP() net.IP
	// UserAgent return the client user agent
	UserAgent() string
	// Source return the publisher that this client is going into system from that
	Source() Publisher
	// Location of the request
	Location() Location
	// Attributes is the generic attribute system
	Attributes(ImpressionAttributes) interface{}
}
