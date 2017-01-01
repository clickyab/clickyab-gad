package mr

//Province type is a location structure
type Province struct {
	ID     int64  `json:"location_id" db:"location_id"`
	Name   string `json:"location_name" db:"location_name"`
	NameFa string `json:"location_name_persian" db:"location_name_persian"`
	Master bool   `json:"location_master" db:"location_master"`
	Select bool   `json:"location_select" db:"location_select"`
}

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
