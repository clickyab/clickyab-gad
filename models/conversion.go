package models

import (
	"errors"
	"time"
)

// InsertConversion try to insert conversion base on click
func (m *Manager) InsertConversion(actionID string, impression *Impression) error {
	q := `INSERT INTO clicks_conv
			(
			w_id,
			app_id,
			wp_id,
			ca_id,
			ad_id,
			cop_id,
			cp_id,
			slot_id,
			imp_id,
			c_time,
			c_date,
			c_action)
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`

	a, err := m.GetWDbMap().Exec(q,
		impression.WebsiteID,
		impression.AppID,
		impression.WP,
		impression.CaID,
		impression.AdID,
		impression.CopID,
		impression.CPID,
		impression.SlotID,
		impression.ID,
		time.Now().Unix(),
		time.Now().Format("20060102"),
		actionID)
	if err != nil {
		return err
	}

	if cnt, err := a.RowsAffected(); err != nil || cnt != 1 {
		return errors.New("No row affected, the clickid is wrong?")
	}

	return nil
}
