package mr

import (
	"errors"
	"time"
)

// InsertConversion try to insert conversion base on click
func (m *Manager) InsertConversion(clickID, action string) error {
	q := `INSERT INTO click_conv (c_id,
		w_id,
		app_id,
		wp_id,
		ca_id,
		ad_id,
		cop_id,
		cp_id,
		slot_id,
		sla_id,
		imp_id,
		c_ip,
		c_status,
		c_referaddress,
		c_parenturl,
		c_fast,
		c_os,
		c_winnerbid,
		c_time,
		c_date,
		c_action)
		SELECT c_id,
		w_id,
		app_id,
		wp_id,
		ca_id,
		ad_id,
		cop_id,
		cp_id,
		slot_id,
		sla_id,
		imp_id,
		c_ip,
		c_status,
		c_referaddress,
		c_parenturl,
		c_fast,
		c_os,
		c_winnerbid, ?, ?, ? FROM clicks WHERE c_id=?`

	a, err := m.GetWDbMap().Exec(q, time.Now().Unix(), time.Now().Format("20060102"), action, clickID)
	if err != nil {
		return err
	}

	if cnt, err := a.RowsAffected(); err != nil || cnt != 1 {
		return errors.New("No row affected, the clickid is wrong?")
	}

	return nil
}
