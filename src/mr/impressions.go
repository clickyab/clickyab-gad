package mr

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

// Impression is the single impression record
type Impression struct {
	ID              int64          `json:"imp_id" db:"imp_id"`
	WebsiteID       sql.NullInt64  `json:"w_id" db:"w_id"`
	WP              sql.NullInt64  `json:"wp_id" db:"wp_id"`
	AppID           sql.NullInt64  `json:"app_id" db:"app_id"`
	AdID            sql.NullInt64  `json:"ad_id" db:"ad_id"`
	CopID           sql.NullInt64  `json:"cop_id" db:"cop_id"`
	CaID            sql.NullInt64  `json:"ca_id" db:"ca_id"`
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
func (m *Manager) InsertImpression(imp Impression) (int64, error) {
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
							%s,%s
							)`, time.Now().Format("20060102"), strconv.Itoa(int(time.Now().Unix())), time.Now().Format("20060102"))

	res, err := m.GetDbMap().Exec(query,
		imp.WebsiteID, imp.WP, imp.AppID,
		imp.AdID, imp.CopID, imp.CaID,
		imp.IP, imp.ReferralAddress, imp.ParentURL,
		imp.URL, imp.WinnerBid, imp.Status,
		imp.Cookie, imp.Alexa, imp.Flash,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
