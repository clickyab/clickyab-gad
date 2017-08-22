package mr

import (
	"assert"
	"database/sql"
	"fmt"
	"net"
	"strconv"
	"time"
	"utils"
)

// CookieProfile cookie_profiles struct table
type CookieProfile struct {
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

func randInt64() int64 {
	x := <-utils.ID
	i, err := strconv.ParseInt(x[:8], 16, 64)
	assert.Nil(err)
	return i
}

// fetchCookieProfile get data from table cookie
func (m *Manager) fetchCookieProfile(key string) (*CookieProfile, error) {
	var res = CookieProfile{}
	hash := utils.Sha1(fmt.Sprintf("a_new_cookie_%s", key))
	err := fetchTouch(hash, &res, 3*24*time.Hour)
	if err == nil {
		return &res, nil
	}
	return nil, err
}

// insertCookieProfile create a new cookie profile and return it
func (m *Manager) insertCookieProfile(key string, ip net.IP) (*CookieProfile, error) {
	hash := utils.Sha1(fmt.Sprintf("a_new_cookie_%s", key))
	ipNullString := toNullString(ip.String())
	date := toNullInt64(time.Now().Unix())
	co := &CookieProfile{
		Key:  key,
		IP:   ipNullString,
		Date: date,
	}
	co.ID = randInt64()
	_ = store(hash, co, 3*24*time.Hour)

	return co, nil
}

// CreateCookieProfile try to select/create a cookie profile
func (m *Manager) CreateCookieProfile(key string, ip net.IP) *CookieProfile {
	//key = key[:config.Config.Clickyab.CopLen]
	res, err := m.fetchCookieProfile(key)
	if err != nil {
		res, err = m.insertCookieProfile(key, ip)
		assert.Nil(err)
	}
	return res
}
