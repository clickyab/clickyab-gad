package rtb

import (
	"entity"
	"store"
)

// AsyncList is the async procedure to return the list of ads url
func AsyncListCTR(
	pub entity.Publisher,
	imp entity.Impression,
	ads map[int][]entity.Advertise,
	slots []entity.Slot,
	multiVideo bool,
	minCPC int64) []entity.Slot {

	go func() {
		s := store.GetSyncStore()
		SelectCTR(s, pub, imp, ads, slots, multiVideo, minCPC)
	}()

	return slots
}
