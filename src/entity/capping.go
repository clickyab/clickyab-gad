package entity

// CappingInterface interface capping
type Capping interface {
	// View return the view of this campaign for this user
	View() int
	// View return the view of this campaign for this user
	AdView(int64) int
	// Frequency return the frequency for this object
	Frequency() int
	// Capping return the frequency capping value, the view/frequency
	Capping() int
	// AdCapping return the frequency capping value, the view/frequency
	AdCapping(int64) int
	// IncView increase the vie
	IncView(int64, int, bool)
	// GetSelected return if this campaign is already selected in this batch
	Selected() bool
	//// IsTargeted return if the current campaign is targeted for this user?
	//IsTargeted() bool
}
