package mr

import (
	"database/sql"
	"net"
	"time"
	"utils"
)

// IP2Location struct table ip2location
type IP2Location struct {
	IPFrom      int64          `json:"ip_from" db:"ip_from"`
	IPTo        int64          `json:"ip_to" db:"ip_to"`
	CountryCode sql.NullString `json:"country_code" db:"country_code"`
	CountryName sql.NullString `json:"country_name" db:"country_name"`
	RegionName  sql.NullString `json:"region_name" db:"region_name"`
	CityName    sql.NullString `json:"city_name" db:"city_name"`
}

//GetLocation @todo
func (m *Manager) GetLocation(ip net.IP) (*IP2Location, error) {
	var res IP2Location
	long, err := utils.IP2long(ip)
	if err != nil {
		return nil, err
	}
	key := utils.Sha1("IP_" + ip.String())
	err = fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	query := `SELECT * FROM ip2location WHERE ip_from <= ? AND ip_to >= ? LIMIT 1`
	err = m.GetRDbMap().SelectOne(
		&res,
		query,
		long,
		long,
	)

	if err != nil {
		return nil, err
	}
	_ = store(key, &res, time.Hour)
	return &res, nil
}
