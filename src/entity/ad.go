package entity

// Advertise is the single advertise interface
type Advertise interface {
	// GetID return the id of advertise
	ID() int64

	Campaign() Campaign

	SetCapping(Capping)
}
