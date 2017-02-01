package rtb

import (
	"entity"
	"store"
	"time"
)

// SyncList is the sync procedure to return the list of ads url
func SyncListCTR(
	pub entity.Publisher,
	imp entity.Impression,
	ads map[int][]entity.Advertise,
	slots []entity.Slot,
	multiVideo bool,
	minCPC int64) []entity.Slot {

	s := store.GetSyncStore()
	SelectCTR(s, pub, imp, ads, slots, multiVideo, minCPC)
	for i := range slots {
		// wait to make sure the ad is ready, since we are sync
		s.Pop(slots[i].StateID(), time.Second*10)
	}

	return slots
}
