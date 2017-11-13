package models

import (
	"database/sql"
	"fmt"
	"time"

	"clickyab.com/gad/transport"
)

// Impression is the single impression record
type Impression struct {
	ID              int64          `json:"imp_id" db:"imp_id"`
	WP              sql.NullInt64  `json:"wp_id" db:"wp_id"`
	CPID            sql.NullInt64  `json:"cp_id" db:"cp_id"`
	WebsiteID       sql.NullInt64  `json:"w_id" db:"w_id"`
	AppID           sql.NullInt64  `json:"app_id" db:"app_id"`
	AdID            sql.NullInt64  `json:"ad_id" db:"ad_id"`
	CopID           sql.NullInt64  `json:"cop_id" db:"cop_id"`
	CaID            sql.NullInt64  `json:"ca_id" db:"ca_id"`
	SlotID          sql.NullInt64  `json:"slot_id" db:"slot_id"`
	SLAID           sql.NullInt64  `json:"sla_id" db:"sla_id"`
	CellID          sql.NullInt64  `json:"cell_id" db:"cell_id"`
	HoodID          sql.NullInt64  `json:"hood_id" db:"hood_id"`
	IP              sql.NullString `json:"imp_ipaddress" db:"imp_ipaddress"`
	ReferralAddress sql.NullString `json:"imp_referaddress" db:"imp_referaddress"`
	ParentURL       sql.NullString `json:"imp_parenturl" db:"imp_parenturl"`
	URL             sql.NullString `json:"imp_url" db:"imp_url"`
	WinnerBid       sql.NullInt64  `json:"imp_winnerbid" db:"imp_winnerbid"`
	Status          sql.NullInt64  `json:"imp_status" db:"imp_status"`
	Cookie          sql.NullInt64  `json:"imp_cookie" db:"imp_cookie"`
	Alexa           sql.NullInt64  `json:"imp_alexa" db:"imp_alexa"`
	Flash           sql.NullInt64  `json:"imp_flash" db:"imp_flash"`
	Time            sql.NullInt64  `json:"imp_time" db:"imp_time"`
	Date            sql.NullInt64  `json:"imp_date" db:"imp_date"`
}

// InsertImpression insert into impression table
func (m *Manager) InsertImpression(imp *transport.Impression) error {
	query := fmt.Sprintf(`INSERT INTO impressions%s (
							cp_id,
							w_id,wp_id,app_id,
							ad_id,cop_id,ca_id,
							imp_ipaddress,imp_referaddress,imp_parenturl,
							imp_url,imp_winnerbid,imp_status,
							imp_cookie,imp_alexa,imp_flash,
							imp_time,imp_date,sla_id,slot_id
							) VALUES (
							?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,?
							)`, time.Now().Format("20060102"))
	wid := sql.NullInt64{}
	refer := sql.NullString{}
	parent := sql.NullString{}
	var slot int64
	if imp.Web != nil {
		wid.Valid = true
		wid.Int64 = imp.Web.WebsiteID

		refer.Valid = imp.Web.Referrer != ""
		refer.String = imp.Web.Referrer

		parent.Valid = imp.Web.ParentURL != ""
		parent.String = imp.Web.ParentURL
		slot = imp.Web.SlotID
	}
	appID := sql.NullInt64{}
	if imp.App != nil {
		appID.Valid = true
		appID.Int64 = imp.App.AppID
		slot = imp.App.SlotID
	}

	res, err := m.GetWDbMap().Exec(query,
		imp.CampaignID,
		wid, 0, appID,
		imp.AdID, imp.CopID, imp.CampaignAdID,
		imp.IP.String(), refer, parent,
		imp.URL, imp.WinnerBID, imp.Status,
		0, 0, 0,
		imp.Time.Unix(), imp.Time.Format("20060102"),
		imp.SLAID,
		slot,
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

// FindImpByIDDate find imp by id and date
func (m *Manager) FindImpByIDDate(ID int64, date string) (*Impression, error) {
	res := Impression{}
	q := fmt.Sprintf("SELECT * FROM impressions%s WHERE imp_id=? LIMIT 1", date)
	err := m.GetRDbMap().SelectOne(&res, q, ID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
