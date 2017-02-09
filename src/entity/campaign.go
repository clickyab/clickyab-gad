package entity

// Campaign is the single campaign in ssytem
type Campaign interface {
	// ID return the campaign id
	ID() int64
	// Name is the campaign name
	Name() string
	// MaxBID get the campaign max bid
	MaxBID() int64
	// Capping return the campaign capping object
	Capping() Capping
	// Make sure the result is >= 1
	Frequency() int
}
