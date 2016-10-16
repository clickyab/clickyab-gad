package mr

import (
	"database/sql"
)

//Country2Info struct country info
type Country2Info struct {
	ID        int64          `id:"id" db:"id"`
	Iso       string         `json:"iso" db:"iso"`
	Name      string         `json:"name" db:"name"`
	NiceName  string         `json:"nicename" db:"nicename"`
	Iso3      sql.NullString `json:"iso3" db:"iso3"`
	NumCode   sql.NullString `json:"numcode" db:"numcode"`
	Phonecode sql.NullString `json:"phonecode" db:"phonecode"`
}

//ConvertCountry2Info get data country from string
func (m *Manager) ConvertCountry2Info(name string) (Country2Info, error) {
	var country Country2Info
	query := `SELECT * FROM country WHERE nicename = ? LIMIT 1`
	err := m.GetDbMap().SelectOne(
		&country,
		query,
		name,
	)
	if err != nil {
		return country, err
	}

	return country, nil
}
