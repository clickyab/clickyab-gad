package mr

// ByMulti sort by multi sort
type ByMulti struct {
	Video bool
	Ads   []*AdData
}

func (a ByMulti) Len() int {
	return len(a.Ads)
}
func (a ByMulti) Swap(i, j int) {
	a.Ads[i], a.Ads[j] = a.Ads[j], a.Ads[i]
}
func (a ByMulti) Less(i, j int) bool {

	if a.Ads[i].Capping.GetSelected() != a.Ads[j].Capping.GetSelected() {
		return !a.Ads[i].Capping.GetSelected()
	}

	if a.Ads[i].Capping.GetAdCapping(a.Ads[i].AdID) != a.Ads[j].Capping.GetAdCapping(a.Ads[j].AdID) {
		return a.Ads[i].Capping.GetAdCapping(a.Ads[i].AdID) < a.Ads[j].Capping.GetAdCapping(a.Ads[j].AdID)
	}

	if a.Video {
		if a.Ads[i].AdType != a.Ads[j].AdType {
			return a.Ads[i].AdType == VideoAdType
		}
	}

	//if a[i].Capping.GetCapping() != a[j].Capping.GetCapping() {
	//	return a[i].Capping.GetCapping() < a[j].Capping.GetCapping()
	//}
	return a.Ads[i].CPM > a.Ads[j].CPM
}
