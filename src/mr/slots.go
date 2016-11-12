package mr

import (
	"fmt"

	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type Slots struct {
	ID               int64          `json:"slot_id" db:"slot_id"`
	PublicID         int64          `json:"slot_pubilc_id" db:"slot_pubilc_id"`
	Name             sql.NullString `json:"slot_name" db:"slot_name"`
	Size             sql.NullString `json:"slot_size" db:"slot_size"`
	WID              int            `json:"w_id" db:"w_id"`
	AppID            int64          `json:"app_id" db:"app_id"`
	AvgDailyImps     int64          `json:"slot_avg_daily_imps" db:"slot_avg_daily_imps"`
	AvgDailyClicks   int64          `json:"slot_avg_daily_clicks" db:"slot_avg_daily_clicks"`
	FloorCPM         int            `json:"slot_floor_cpm" db:"slot_floor_cpm"`
	TotalMonthlyCost int64          `json:"slot_total_monthly_cost" db:"slot_total_monthly_cost"`
	LastUpdate       int64          `json:"slot_lastupdate" db:"slot_lastupdate"`
	CreatedAt        mysql.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt        mysql.NullTime `json:"updated_at" db:"updated_at"`
}

func (m *Manager) FetchSlots(publicID string, wID int64) ([]Slots, error) {
	var res []Slots

	query := fmt.Sprintf(`SELECT * FROM slots WHERE slot_pubilc_id IN (%s) AND w_id = ?`, publicID)

	_, err := m.GetDbMap().Select(
		&res,
		query,
		wID,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *Manager) InsertSlots(slotsPublic []int64) error {
	var slot []interface{}
	for s:= range slotsPublic{
		slot=append(slot,&Slots{PublicID:slotsPublic[s]})
	}
	err := m.GetDbMap().Insert(slot...)
	if err != nil {
		return err
	}
	return nil

}
