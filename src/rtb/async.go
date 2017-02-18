package rtb

import (
	"entity"
	"store"
)

// AsyncList is the async procedure to return the list of ads url
func AsyncListCTR(
	imp entity.Impression,
	ads map[int][]entity.Advertise,
) []entity.Slot {

	go func() {
		s := store.GetSyncStore()
		selectCTR(s, imp, ads)
	}()

	return imp.Slots()
}
