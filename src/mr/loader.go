package mr

func (m *Manager) LoadAds() ([]AdData, error) {
	var res []AdData
	m.GetDbMap().Select(&res, "")
}
