package mr

import (
	"database/sql"
	"transport"
)

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

func (m *Manager) InsertClick(imp *transport.Click) error {
	query := `INSERT INTO clicks (
							c_winnerbid,w_id,app_id,
							wp_id,cp_id,ca_id,
							slot_id,sla_id,ad_id,cop_id,imp_id,c_status,c_ip,c_referaddress,c_parenturl,
							c_fast,imp_winnerbid,imp_status,
							imp_cookie,imp_alexa,imp_flash,
							imp_time,imp_date
							) VALUES (
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?
							)`
	wid := sql.NullInt64{}
	refer := sql.NullString{}
	parent := sql.NullString{}
	if imp.Web != nil {
		wid.Valid = true
		wid.Int64 = imp.Web.WebsiteID

		refer.Valid = imp.Web.Referrer != ""
		refer.String = imp.Web.Referrer

		parent.Valid = imp.Web.ParentURL != ""
		parent.String = imp.Web.ParentURL
	}
	appID := sql.NullInt64{}
	if imp.App != nil {
		appID.Valid = true
		appID.Int64 = imp.App.AppID
	}

	res, err := m.GetWDbMap().Exec(query,
		wid, 0, appID,
		imp.AdID, imp.CopID, imp.CampaignAdID,
		imp.IP.String(), refer, parent,
		imp.URL, imp.WinnerBID, imp.Status,
		0, 0, 0,
		imp.Time.Unix(), imp.Time.Format("20060102"),
	)
	if err != nil {
		return err
	}
	imp.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
