package mr

import "database/sql"

// Click data table clicks
type Click struct {
	ID              int64          `json:"c_id" db:"c_id"`
	WinnerBid       sql.NullInt64  `json:"imp_winnerbid" db:"imp_winnerbid"`
	WebsiteID       sql.NullInt64  `json:"w_id" db:"w_id"`
	AppID           sql.NullInt64  `json:"app_id" db:"app_id"`
	WP              sql.NullInt64  `json:"wp_id" db:"wp_id"`
	CampaignID      sql.NullInt64  `json:"cp_id" db:"cp_id"`
	CaID            sql.NullInt64  `json:"ca_id" db:"ca_id"`
	SlotID          sql.NullInt64  `json:"slot_id" db:"slot_id"`
	SlotADID        sql.NullInt64  `json:"sla_id" db:"sla_id"`
	AdID            sql.NullInt64  `json:"ad_id" db:"ad_id"`
	CopID           sql.NullInt64  `json:"cop_id" db:"cop_id"`
	ImpressionID    int64          `json:"imp_id" db:"imp_id"`
	Status          sql.NullInt64  `json:"imp_status" db:"imp_status"`
	IP              sql.NullString `json:"c_ip" db:"c_ip"`
	ReferralAddress sql.NullString `json:"imp_referaddress" db:"imp_referaddress"`
	ParentURL       sql.NullString `json:"imp_parenturl" db:"imp_parenturl"`
	Fast            sql.NullInt64  `json:"c_fast" db:"c_fast"`
	OS              sql.NullInt64  `json:"c_os" db:"c_os"`
	Time            sql.NullInt64  `json:"imp_time" db:"imp_time"`
	Date            sql.NullInt64  `json:"imp_date" db:"imp_date"`
}
