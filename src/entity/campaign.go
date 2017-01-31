package entity

// Campaign is the single campaign in ssytem
type Campaign interface {
	ID() int64

	Name() string

	MaxBID() int64

	Capping() Capping
	// Make sure the result is >= 1
	Frequency() int
}
