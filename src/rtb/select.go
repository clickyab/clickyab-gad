package rtb

import (
	"assert"
	"config"
	"eav"
	"entity"
	"fmt"
	"sort"
	"store"
	"time"
)

const (
	// Mega is the mega store prefix
	Mega string = "MEGA_"
	// MegaIP is the ip subkey
	MegaIP string = "IP"
	// MegaUserAgent is the user agent subkey
	MegaUserAgent string = "UA"
	// MegaPubID is the publisher id subkey
	MegaPubID string = "PID"
	// MegaTimeUnix is the impression timestamp subkey
	MegaTimeUnix string = "TU"
	// MegaAdvertise is the selected ad subkey
	MegaAdvertise string = "AD"
	// MegaSlot is the slot subkey
	MegaSlot string = "SLOT"
)

func createMegaStore(imp entity.Impression) eav.Kiwi {
	kiwi := eav.NewEavStore(Mega + imp.MegaIMP())
	assert.Nil(kiwi.SetSubKey(MegaIP, imp.IP().String()).
		SetSubKey(MegaUserAgent, imp.UserAgent()).
		SetSubKey(MegaPubID, fmt.Sprint(imp.Source().ID())).
		SetSubKey(MegaTimeUnix, fmt.Sprint(time.Now().Unix())).
		Save(config.Config.Clickyab.MegaImpExpire))
	return kiwi
}

// SelectCTR is the key function to select an ad for an imp base on real time biding
func SelectCTR(
	store store.Store,
	imp entity.Impression,
	ads map[int][]entity.Advertise) {

	// TODO : better implementation
	multiVideo := imp.Source().Type() == entity.PublisherTypeVast

	// Get the capping
	slots := imp.Slots()
	pub := imp.Source()
	ads = getCapping(imp.ClientID(), ads, slots)
	kiwi := createMegaStore(imp)
	for i := range slots {
		var (
			exceedFloor []entity.Advertise
			underFloor  []entity.Advertise
			size        = slots[i].Size()
			noVideo     bool
		)
		for _, ad := range ads[size] {
			if ad.Type() == entity.AdTypeVideo && noVideo {
				continue
			}
			if ad.WinnerBID() == 0 && doBid(ad, pub, slots[i]) {
				exceedFloor = append(exceedFloor, ad)
			} else if ad.WinnerBID() == 0 {
				underFloor = append(underFloor, ad)
			}
		}
		var sorted []entity.Advertise
		var (
			ef     entity.SortByCap
			secBid bool
		)
		// order is to get data from exceed flor, then capping passed and if the config allowed,
		// use the under floor. for under floor there is no second biding pricing
		if len(exceedFloor) > 0 {
			ef = entity.SortByCap(exceedFloor)
			secBid = true
		} else if config.Config.Clickyab.UnderFloor && len(underFloor) > 0 {
			ef = entity.SortByCap(underFloor)
			secBid = false
		}

		if len(ef) == 0 {
			// TODO : Warnings
			store.Push(slots[i].StateID(), "", time.Hour)
			continue
		}

		sort.Sort(ef)
		sorted = []entity.Advertise(ef)
		// Do not do second biding pricing on this ads, they can not pass CPMFloor
		if secBid {
			secondCPM := getSecondCPM(pub.FloorCPM(), sorted)
			sorted[0].SetWinnerBID(winnerBid(secondCPM, sorted[0].CTR()))
		} else {
			sorted[0].SetWinnerBID(sorted[0].Campaign().MaxBID())
		}

		// Force price on min CPC
		if sorted[0].WinnerBID() < imp.Source().MinCPC() {
			sorted[0].SetWinnerBID(imp.Source().MinCPC())
		}

		sorted[0].Capping().IncView(sorted[0].ID(), 1, true)
		slots[i].SetWinnerAdvertise(sorted[0])

		if !multiVideo {
			noVideo = noVideo || sorted[0].Type() == entity.AdTypeVideo
		}

		kiwi.SetSubKey(fmt.Sprintf("%s_%d", MegaAdvertise, sorted[0].ID()), fmt.Sprint(sorted[0].WinnerBID()))
		kiwi.SetSubKey(fmt.Sprintf("%s_%d", MegaSlot, sorted[0].ID()), fmt.Sprint(slots[i].ID()))
		assert.Nil(kiwi.Save(config.Config.Clickyab.MegaImpExpire))

		store.Push(slots[i].StateID(), fmt.Sprint(sorted[0].ID()), time.Hour)
	}
}

func doBid(ad entity.Advertise, pub entity.Publisher, slot entity.Slot) bool {
	ad.SetCTR(calculateCTR(
		ad,
		slot,
	))
	ad.SetCPM(cpm(ad.Campaign().MaxBID(), ad.CTR()))
	//exceed cpm floor
	return ad.CPM() >= pub.FloorCPM()
}

// CalculateCtr calculate ctr
func calculateCTR(ad entity.Advertise, slot entity.Slot) float64 {
	return (ad.AdCTR()*float64(config.Config.Clickyab.AdCTREffect) + slot.SlotCTR()*float64(config.Config.Clickyab.SlotCTREffect)) / float64(100)
}

//Cpm calculate cpm
func cpm(bid int64, ctr float64) int64 {
	return int64(float64(bid) * ctr * 10.0)
}

func getSecondCPM(floorCPM int64, exceedFloor []entity.Advertise) int64 {
	var secondCPM = floorCPM
	if len(exceedFloor) > 1 && exceedFloor[0].Capping().Selected() == exceedFloor[1].Capping().Selected() {
		secondCPM = exceedFloor[1].CPM()
	}

	return secondCPM
}

// winnerBid calculate winner bid
func winnerBid(cpm int64, ctr float64) int64 {
	return int64(float64(cpm)/(ctr*10)) + 1
}
