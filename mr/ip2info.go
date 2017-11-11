package mr

import (
	"database/sql"
	"errors"
	"time"

	"clickyab.com/gad/utils"
)

//CountryInfo struct country info
type CountryInfo struct {
	ID        int64          `id:"id" db:"id"`
	Iso       string         `json:"iso" db:"iso"`
	Name      string         `json:"name" db:"name"`
	NiceName  string         `json:"nicename" db:"nicename"`
	Iso3      sql.NullString `json:"iso3" db:"iso3"`
	NumCode   sql.NullString `json:"numcode" db:"numcode"`
	Phonecode sql.NullString `json:"phonecode" db:"phonecode"`
}

//ConvertCountry2Info get data country from string
// @DEPRICATED
func (m *Manager) ConvertCountry2Info(name string) (CountryInfo, error) {
	var country CountryInfo
	if len(name) < 2 {
		return country, errors.New("invalid country name")
	}

	key := utils.Hash("Country_" + name)
	err := fetch(key, &country)
	if err == nil {
		return country, nil
	}

	query := `SELECT * FROM country WHERE iso = ? LIMIT 1`
	err = m.GetRDbMap().SelectOne(
		&country,
		query,
		name,
	)
	if err != nil {
		return country, err
	}

	_ = store(key, &country, 72*time.Hour)
	return country, nil
}
