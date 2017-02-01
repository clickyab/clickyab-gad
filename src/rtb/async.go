package rtb

import (
	"entity"
	"store"
)

// AsyncList is the async procedure to return the list of ads url
func AsyncListCTR(
	u entity.URLProvider,
	pub entity.Publisher,
	imp entity.Impression,
	ads map[int][]entity.Advertise,
	slots []entity.Slot,
	multiVideo bool,
	minCPC int64) map[string]string {
	res := make(map[string]string)
	for i := range slots {
		res[slots[i].PublicID()] = u.ShowURL(slots[i], imp, pub)
	}

	go func() {
		s := store.GetClusterStore()
		SelectCTR(s, pub, imp, ads, slots, multiVideo, minCPC)
	}()

	return res
}
