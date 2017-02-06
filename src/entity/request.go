package entity

import (
	"net"
	"net/http"
)

// ImpressionAttributes is the imp attr key
type ImpressionAttributes string

// BIDType is the bid type for this imp cpc or cpm
type BIDType string

const (
	// BIDTypeCPC is the cost per click type
	BIDTypeCPC = "CPC"
	//BIDTypeCPM is the cost per view type
	BIDTypeCPM = "CPM"
)

// Impression is the single impression object
type Impression interface {
	Request() *http.Request
	// MegaIMP return the random id of this imp object
	MegaIMP() string
	// ClientID is the key to identify client
	ClientID() int64
	// IP return the client ip
	IP() net.IP
	// UserAgent return the client user agent
	UserAgent() string
	// Source return the publisher that this client is going into system from that
	Source() Publisher
	// Location of the request
	Location() Location
	// OS the os of requester if available
	OS() OS
	// Attributes is the generic attribute system
	Attributes(ImpressionAttributes) interface{}
	// AcceptedTypes is the type accepted by this impression
	AcceptedTypes() []AdType
	// Slots is the slot for this request
	Slots() []Slot
}
