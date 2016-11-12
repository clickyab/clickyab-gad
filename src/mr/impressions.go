package mr

import (
	"database/sql"
	"fmt"
	"time"
	"transport"
)

//// Impression is the single impression record
//type Impression struct {
//	ID              int64          `json:"imp_id" db:"imp_id"`
//	WebsiteID       sql.NullInt64  `json:"w_id" db:"w_id"`
//	WP              sql.NullInt64  `json:"wp_id" db:"wp_id"`
//	AppID           sql.NullInt64  `json:"app_id" db:"app_id"`
//	AdID            sql.NullInt64  `json:"ad_id" db:"ad_id"`
//	CopID           sql.NullInt64  `json:"cop_id" db:"cop_id"`
//	CaID            sql.NullInt64  `json:"ca_id" db:"ca_id"`
//	IP              sql.NullString `json:"imp_ipaddress" db:"imp_ipaddress"`
//	ReferralAddress sql.NullString `json:"imp_referaddress" db:"imp_referaddress"`
//	ParentURL       sql.NullString `json:"imp_parenturl" db:"imp_parenturl"`
//	URL             sql.NullString `json:"imp_url" db:"imp_url"`
//	WinnerBid       sql.NullInt64  `json:"imp_winnerbid" db:"imp_winnerbid"`
//	Status          sql.NullInt64  `json:"imp_status" db:"imp_status"`
//	Cookie          sql.NullInt64  `json:"imp_cookie" db:"imp_cookie"`
//	Alexa           sql.NullInt64  `json:"imp_alexa" db:"imp_alexa"`
//	Flash           sql.NullInt64  `json:"imp_flash" db:"imp_flash"`
//	Time            time.Time      `json:"imp_time" db:"imp_time"`
//	Date            int            `json:"imp_date" db:"imp_date"`
//}

// InsertImpression insert into impression table
func (m *Manager) InsertImpression(imp *transport.Impression) error {
	query := fmt.Sprintf(`INSERT INTO impressions%s (
							w_id,wp_id,app_id,
							ad_id,cop_id,ca_id,
							imp_ipaddress,imp_referaddress,imp_parenturl,
							imp_url,imp_winnerbid,imp_status,
							imp_cookie,imp_alexa,imp_flash,
							imp_time,imp_date
							) VALUES (
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?,?,
							?,?
							)`, time.Now().Format("20060102"))
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
