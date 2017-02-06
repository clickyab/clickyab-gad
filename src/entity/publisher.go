package entity

type (
	// PublisherType is the publisher type
	PublisherType string
	// PublisherAttributes is the publisher attributes
	PublisherAttributes string
)

const (
	// PublisherTypeApp is the app
	PublisherTypeApp PublisherType = "app"
	// PublisherTypeWeb is the web
	PublisherTypeWeb PublisherType = "web"
	// PublisherTypeVast is the vast
	PublisherTypeVast PublisherType = "vast"
)

// Publisher is the publisher interface
type Publisher interface {
	// GetID return the publisher id
	ID() int64
	// FloorCPM is the floor cpm for publisher
	FloorCPM() int64
	// Name of publisher
	Name() string
	// Active is the publisher active?
	Active() bool
	// Type return the publisher type
	Type() PublisherType
	// Attributes is the generic attribute system
	Attributes(PublisherAttributes) interface{}
	// BIDType return this publisher bid type
	BIDType() BIDType
	// MinCPC is the minimum CPC requested for this requets
	MinCPC() int64
}
