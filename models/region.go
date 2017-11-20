package models

// FetchProvinceAll get the list of all region in database
func (m *Manager) FetchProvinceAll() ([]*Province, error) {
	var res = []*Province{}

	query := `SELECT * FROM list_locations`

	_, err := m.GetRDbMap().Select(
		&res,
		query,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}
