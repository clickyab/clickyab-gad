package mr

// SlotData type is record of one slot
type SlotData struct {
	AdID         int64   `json:"ad_id" db:"ad_id"`
	SlotSize     int     `json:"slot_size" db:"slot_size"`
	SLAClicks    int64   `json:"sla_clicks" db:"sla_clicks"`
	SLAImps      int64   `json:"sla_imps" db:"sla_imps"`
	SlotPublicID int64   `json:"slot_pubilc_id" db:"slot_pubilc_id"`
	SlotFloorCPM int     `json:"slot_floor_cpm" db:"slot_floor_cpm"`
	SlotCtr      float64 `json:"slot_ctr" `
}

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
