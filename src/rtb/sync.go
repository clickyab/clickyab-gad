package rtb

import (
	"entity"
	"store"
	"time"
)

// SyncList is the sync procedure to return the list of ads url
func SyncListCTR(
	imp entity.Impression,
	ads map[int][]entity.Advertise,
) []entity.Slot {

	s := store.GetSyncStore()
	selectCTR(s, imp, ads)
	for _, sl := range imp.Slots() {
		// wait to make sure the ad is ready, since we are sync
		s.Pop(sl.StateID(), time.Second*10)
	}

	return imp.Slots()
}
