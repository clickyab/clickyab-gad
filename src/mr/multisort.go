package mr

// ByMulti sort by multi sort
type ByMulti []*AdData

func (a ByMulti) Len() int {
	return len(a)
}
func (a ByMulti) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByMulti) Less(i, j int) bool {
	if a[i].Capping.GetSelected() != a[j].Capping.GetSelected() {
		return !a[i].Capping.GetSelected()
	}
	if a[i].Capping.GetAdCapping(a[i].AdID) != a[j].Capping.GetAdCapping(a[j].AdID) {
		return a[i].Capping.GetAdCapping(a[i].AdID) < a[j].Capping.GetAdCapping(a[j].AdID)
	}
	if a[i].Capping.GetCapping() != a[j].Capping.GetCapping() {
		return a[i].Capping.GetCapping() < a[j].Capping.GetCapping()
	}
	if a[i].CampaignNetwork != a[j].CampaignNetwork {
		return a[i].CampaignNetwork == 2
	}
	return a[i].CPM < a[j].CPM
}
