package mr

func (m *Manager) LoadAds() ([]AdData, error) {
	var res []AdData
	_, err := m.GetDbMap().Select(&res, "")
	if err != nil {
		return nil, err
	}

	return res, nil
}
