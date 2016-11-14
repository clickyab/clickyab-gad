package mr

// InsertSlotAd try to insert slot ad
func (m *Manager) InsertSlotAd(slotID, adID int64) (int64, error) {
	query := `insert into slots_ads (slot_id, ad_id) values (?, ?) ON DUPLICATE KEY UPDATE sla_id=LAST_INSERT_ID(sla_id)`
	res, err := m.GetWDbMap().Exec(
		query,
		slotID,
		adID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FetchSlotAd return the slot ad id by select
func (m *Manager) FetchSlotAd(slotID, adID int64) (int64, error) {
	query := `SELECT sla_id FROM slots WHERE slot_id = ? AND ad_id = ?`
	res, err := m.GetProperDBMap().SelectInt(query, slotID, adID)
	if err != nil {
		return m.InsertSlotAd(slotID, adID)
	}
	return res, nil
}
