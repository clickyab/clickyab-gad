package mr

import (
	"database/sql"
	"utils"
)

type IP2Location struct {
	IPFrom      int64          `json:"ip_from" db:"ip_from"`
	IPTO        int64          `json:"ip_to" db:"ip_to"`
	CountryCode sql.NullString `json:"country_code" db:"country_code"`
	CountryName sql.NullString `json:"country_name" db:"country_name"`
	RegionName  sql.NullString `json:"region_name" db:"region_name"`
	CityName    sql.NullString `json:"city_name" db:"city_name"`
}

//
func (m *Manager) GetLocation(ip string) (*IP2Location, error) {
	var res IP2Location
	long,err:=utils.IP2long(ip)
	if err!=nil{
		return nil,err
	}
	query := `SELECT * FROM ip2location WHERE ip_from >= ? AND ip_to <= ? LIMIT 1`
	err = m.GetDbMap().SelectOne(
		&res,
		query,
		long,
		long,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
