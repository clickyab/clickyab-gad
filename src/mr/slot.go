package mr

// SlotData type @todo
type SlotData struct {
	AdID         int64 `json:"ad_id" db:"ad_id"`
	SlotSize     int   `json:"slot_size" db:"slot_size"`
	SLAClicks    int64 `json:"sla_clicks" db:"sla_clicks"`
	SLAImps      int64 `json:"sla_imps" db:"sla_imps"`
	SlotPublicID int64 `json:"slot_pubilc_id" db:"slot_pubilc_id"`
	SlotFloorCPM int   `json:"slot_floor_cpm" db:"slot_floor_cpm"`
}
