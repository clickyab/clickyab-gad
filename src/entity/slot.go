package entity

// Slot is the slot of the app
type Slot interface {
	// ID of slot
	ID() int64
	// PublicID of slot
	PublicID() int64
	// Size of slot
	Size() int
	// StateID is an string for this slot, its a random at first but the value is not changed at all other calls
	StateID() string
}
