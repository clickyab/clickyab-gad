package models

import (
	"database/sql"

	"clickyab.com/gad/transport"
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

// InsertClick try to insert a click in database
func (m *Manager) InsertClick(click *transport.Click) error {
	query := `INSERT INTO clicks (c_winnerbid,
	w_id,
	app_id,
	wp_id,
	cp_id,
	ca_id,
	slot_id,
	sla_id,
	ad_id,
	cop_id,
	imp_id,
	c_status,
	c_ip,
	c_referaddress,
	c_parenturl,
	c_fast,
	c_os,
	c_time,
	c_date) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	web := sql.NullInt64{}
	parent := sql.NullString{}
	ref := sql.NullString{}
	app := sql.NullInt64{}
	if click.Web != nil {
		web.Valid = true
		web.Int64 = click.Web.WebsiteID
		parent.Valid = true
		parent.String = click.Web.ParentURL
		ref.Valid = true
		ref.String = click.Web.Referrer
	} else if click.App != nil {
		app.Valid = true
		app.Int64 = click.App.AppID
	}

	res, err := m.GetWDbMap().Exec(query,
		click.WinnerBid,
		web,
		app, //app id
		0,   //wp_id
		click.CampaignID,
		click.CampaignAdID, //default in mysql ca_id
		click.SlotID,
		//fetch slot AD
		click.SLAID,
		click.AdID,
		click.CopID,
		click.ImpID,
		click.Status, // c_status
		click.IP.String(),
		ref,
		parent,
		click.OutTime.Unix()-click.InTime.Unix(),
		click.OS,
		click.OutTime.Unix(),
		click.OutTime.Format("20060102"),
	)
	if err != nil {
		return err
	}
	click.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	if click.TrueView {
		return m.InsertTrueView(click.ID)
	}
	return nil
}

// InsertTrueView is for vast true view
func (m *Manager) InsertTrueView(clickID int64) error {
	query := `INSERT INTO trueview (tv_click_id) VALUES (?)`
	_, err := m.GetProperDBMap().Exec(query, clickID)
	if err != nil {
		return err
	}
	return nil
}
