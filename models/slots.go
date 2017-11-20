package models

import (
	"fmt"

	"database/sql"

	"clickyab.com/gad/utils"

	"time"

	"github.com/clickyab/services/assert"

	"github.com/go-sql-driver/mysql"
)

// Slot is the slots table
type Slot struct {
	ID               int64          `json:"slot_id" db:"slot_id"`
	PublicID         int64          `json:"slot_pubilc_id" db:"slot_pubilc_id"`
	Name             sql.NullString `json:"slot_name" db:"slot_name"`
	Size             sql.NullString `json:"slot_size" db:"slot_size"`
	WID              int64          `json:"w_id" db:"w_id"`
	AppID            int64          `json:"app_id" db:"app_id"`
	AvgDailyImps     int64          `json:"slot_avg_daily_imps" db:"slot_avg_daily_imps"`
	AvgDailyClicks   int64          `json:"slot_avg_daily_clicks" db:"slot_avg_daily_clicks"`
	FloorCPM         int            `json:"slot_floor_cpm" db:"slot_floor_cpm"`
	TotalMonthlyCost int64          `json:"slot_total_monthly_cost" db:"slot_total_monthly_cost"`
	LastUpdate       int64          `json:"slot_lastupdate" db:"slot_lastupdate"`
	CreatedAt        mysql.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt        mysql.NullTime `json:"updated_at" db:"updated_at"`
}

// FetchWebSlots fetch all slots
func (m *Manager) FetchWebSlots(publicID string, wID int64) ([]Slot, error) {
	var res []Slot
	// TODO : this is dangerous to cache this one
	//key := utils.Hash(fmt.Sprintf("slot_%s_%d", publicID, wID))
	//err := fetch(key, &res)
	//if err == nil {
	//	return res, nil
	//}

	query := fmt.Sprintf(`SELECT * FROM slots WHERE slot_pubilc_id IN (%s) AND w_id = ?`, publicID)

	_, err := m.GetProperDBMap().Select(
		&res,
		query,
		wID,
	)
	if err != nil {
		return nil, err
	}

	//_ = store(key, &res, time.Hour)
	return res, nil
}

// FetchAppSlots fetch all slots
func (m *Manager) FetchAppSlots(publicID string, appID int64) ([]Slot, error) {
	var res []Slot
	// TODO : this is dangerous to cache this one
	//key := utils.Hash(fmt.Sprintf("slot_%s_%d", publicID, wID))
	//err := fetch(key, &res)
	//if err == nil {
	//	return res, nil
	//}

	query := fmt.Sprintf(`SELECT * FROM slots WHERE slot_pubilc_id IN (%s) AND app_id = ?`, publicID)

	_, err := m.GetProperDBMap().Select(
		&res,
		query,
		appID,
	)
	if err != nil {
		return nil, err
	}

	//_ = store(key, &res, time.Hour)
	return res, nil
}

// InsertSlots create as many slots you want
func (m *Manager) InsertSlots(wID int64, appID int64, slotsPublic int64, size int) (*Slot, error) {
	assert.True((appID == 0 && wID > 0) || (appID > 0 && wID == 0), "[BUG] invalid input")
	s := &Slot{PublicID: slotsPublic, Size: sql.NullString{String: fmt.Sprint(size), Valid: size != 0}}
	if wID > 0 {
		s.WID = wID
	} else {
		s.AppID = appID
	}

	err := m.GetWDbMap().Insert(s)
	if err != nil {
		return nil, err
	}

	return s, nil

}

// FetchAppSlot is the app version of fetch slot
func (m *Manager) FetchAppSlot(appID int64, slotID int64) (*Slot, error) {
	var res Slot

	key := utils.Hash(fmt.Sprintf("slotapp_%d_%d", slotID, appID))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	err = m.GetProperDBMap().SelectOne(
		&res,
		`SELECT * FROM slots WHERE slot_pubilc_id = ? AND app_id = ?`,
		slotID,
		appID,
	)
	if err != nil {
		return nil, err
	}

	_ = store(key, &res, time.Hour)
	return &res, nil
}
