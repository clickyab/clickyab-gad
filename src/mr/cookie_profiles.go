package mr

import (
	"database/sql"
	"time"
)

// cookie_profiles struct table
type CookieProfiles struct {
	ID      int64          `json:"cop_id" db:"cop_id"`
	Key     string         `json:"cop_key" db:"cop_key"`
	Email   sql.NullString `json:"cop_email" db:"cop_email"`
	IP      sql.NullString `json:"cop_last_ip" db:"cop_last_ip"`
	Gender  sql.NullInt64  `json:"cop_gender" db:"cop_gender"`
	Alexa   sql.NullInt64  `json:"cop_alexa" db:"cop_alexa"`
	OS      sql.NullInt64  `json:"cop_os" db:"cop_os"`
	Browser sql.NullInt64  `json:"cop_browser" db:"cop_browser"`
	City    sql.NullInt64  `json:"cop_city" db:"cop_city"`
	Age     sql.NullInt64  `json:"cop_age" db:"cop_age"`
	KeyWord sql.NullString `json:"cop_keywords" db:"cop_keywords"`
	Date    sql.NullInt64  `json:"cop_active_date" db:"cop_active_date"`
}

// FetchCookieProfile get data from table cookie
func (m *Manager) FetchCookieProfile(key string) (*CookieProfiles, error) {
	var res = CookieProfiles{}
	query := `SELECT * FROM cookie_profiles WHERE cop_key = ?  LIMIT 1`
	err := m.GetDbMap().SelectOne(
		&res,
		query,
		key,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (m *Manager) InsertCookieProfile(cop, ip string) (*CookieProfiles, error) {

	ipNullString := ToNullString(ip)
	date := ToNullInt64(time.Now().Unix())
	co := &CookieProfiles{
		Key:  cop,
		IP:   ipNullString,
		Date: date,
	}
	err := m.GetDbMap().Insert(co)
	if err != nil {
		return nil, err
	}
	return co, nil
}
