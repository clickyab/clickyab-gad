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

type SortByCap []Advertise

func (a SortByCap) Len() int {
	return len(a)
}
func (a SortByCap) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a SortByCap) Less(i, j int) bool {

	// This is a multisort function.
	if a[i].Capping().Selected() != a[j].Capping().Selected() {
		return !a[i].Capping().Selected()
	}
	if a[i].Capping().AdCapping(a[i].ID()) != a[j].Capping().AdCapping(a[j].ID()) {
		return a[i].Capping().AdCapping(a[i].ID()) < a[j].Capping().AdCapping(a[j].ID())
	}
	if a[i].Capping().Capping() != a[j].Capping().Capping() {
		return a[i].Capping().Capping() < a[j].Capping().Capping()
	}

	return a[i].CPM() < a[j].CPM()
}
