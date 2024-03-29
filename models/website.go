package models

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"clickyab.com/gad/utils"
	"github.com/clickyab/services/config"
)

// Website type for website
type Website struct {
	WID                int64          `json:"w_id" db:"w_id"`
	UserID             int64          `json:"u_id" db:"u_id"`
	WPubID             int64          `json:"w_pub_id" db:"w_pub_id"`
	WDomain            sql.NullString `json:"w_domain" db:"w_domain"`
	WSupplier          string         `json:"w_supplier" db:"w_supplier"`
	WName              sql.NullString `json:"w_name" db:"w_name"`
	WCategories        SharpArray     `json:"w_categories" db:"w_categories"`
	WMinBid            int64          `json:"w_minbid" db:"w_minbid"`
	WFloorCpm          sql.NullInt64  `json:"w_floor_cpm" db:"w_floor_cpm"`
	WProfileType       sql.NullInt64  `json:"w_profile_type" db:"w_profile_type"`
	WStatus            int            `json:"w_status" db:"w_status"`
	WReview            int            `json:"w_review" db:"w_review"`
	WAlexaRank         int64          `json:"w_alexarank" db:"w_alexarank"`
	WDiv               float64        `json:"w_div" db:"w_div"`
	WMobad             int            `json:"w_mobad" db:"w_mobad"`
	WNativeAd          int            `json:"w_nativead" db:"w_nativead"`
	WFatFinger         int            `json:"w_fatfinger" db:"w_fatfinger"`
	WPublishStart      int            `json:"w_publish_start" db:"w_publish_start"`
	WPublishEnd        int            `json:"w_publish_end" db:"w_publish_end"`
	WPublishCost       int            `json:"w_publish_cost" db:"w_publish_cost"`
	WPrePayment        int            `json:"w_prepayment" db:"w_prepayment"`
	WTodayCtr          float64        `json:"w_today_ctr" db:"w_today_ctr"`
	WTodayImps         int64          `json:"w_today_imps" db:"w_today_imps"`
	WTodayClicks       int64          `json:"w_today_clicks" db:"w_today_clicks"`
	WDate              int            `json:"w_date" db:"w_date"`
	WNotApprovedReason SharpArray     `json:"w_notapprovedreason" db:"w_notapprovedreason"`
	CreatedAt          sql.NullString `json:"created_at" db:"created_at"`
	UpdatedAt          sql.NullString `json:"updated_at" db:"updated_at"`
}

// GetID return the id of app
func (w *Website) GetID() int64 {
	return w.WID
}

// GetName return the name of object
func (w *Website) GetName() string {
	return w.WDomain.String
}

var (
	floorWeb = config.RegisterInt64("clickyab.min_cpm_floor_web", 1000, "minimum web florcpm")
)

// FloorCPM is the floor value for this site
func (w *Website) FloorCPM() int64 {
	if w.WFloorCpm.Int64 < floorWeb.Int64() {
		w.WFloorCpm.Int64 = floorWeb.Int64()
		w.WFloorCpm.Valid = true
	}
	return w.WFloorCpm.Int64 / 3
}

// GetActive return if app is active or not
func (w *Website) GetActive() bool {
	return w.WStatus == 0 || w.WStatus == 1
}

// GetType of this object
func (w *Website) GetType() string {
	return "web"
}

// FetchWebsiteByPublicID function @todo
func (m *Manager) FetchWebsiteByPublicID(publicID int64) (*Website, error) {
	var res = Website{}
	key := utils.Hash(fmt.Sprintf("Website_%d", publicID))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	query := `SELECT * FROM websites WHERE w_pub_id = ?  LIMIT 1`

	err = m.GetRDbMap().SelectOne(
		&res,
		query,
		publicID,
	)
	if err != nil {
		return nil, err
	}
	_ = store(key, &res, time.Hour)
	return &res, nil
}

// FetchWebsite function @todo
func (m *Manager) FetchWebsite(ID int64) (*Website, error) {
	var res = Website{}
	key := utils.Hash(fmt.Sprintf("WebsiteID_%d", ID))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	query := `SELECT * FROM websites WHERE w_id = ?  LIMIT 1`

	err = m.GetRDbMap().SelectOne(
		&res,
		query,
		ID,
	)
	if err != nil {
		return nil, err
	}
	_ = store(key, &res, time.Hour)
	return &res, nil
}

// FetchWebsiteByDomain return a function based on its domain
func (m *Manager) FetchWebsiteByDomain(domain, supplier string) (*Website, error) {
	var res = Website{}
	key := utils.Hash(fmt.Sprintf("WebsiteDomainSupplier_%s_%s", domain, supplier))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	query := `SELECT * FROM websites WHERE w_domain = ? AND w_status NOT IN (2,3) AND w_supplier= ? LIMIT 1`

	err = m.GetRDbMap().SelectOne(
		&res,
		query,
		domain,
		supplier,
	)
	if err != nil {
		return nil, err
	}
	_ = store(key, &res, time.Hour)
	return &res, nil
}

// FindWebsiteByDomain return a function based on its domain
func (m *Manager) FindWebsiteByDomain(domain string) (*Website, error) {
	var res = Website{}
	query := `SELECT * FROM websites WHERE w_domain = ? AND w_status IN (0, 1) ORDER BY w_today_imps DESC LIMIT 1`

	err := m.GetRDbMap().SelectOne(
		&res,
		query,
		domain,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// InsertWebsite adds a website if its new in the request
func (m *Manager) InsertWebsite(domain, supplier string, userID int64) (*Website, error) {
	if supplier == "clickyab" {
		// we are not allow to register sites from clickyab
		return nil, fmt.Errorf("the clickyab supplier is not allowed to register website on the fly")
	}
	ins := Website{
		UserID:    userID,
		WDomain:   sql.NullString{String: domain, Valid: true},
		WSupplier: supplier,
		CreatedAt: sql.NullString{String: time.Now().String(), Valid: true},
		UpdatedAt: sql.NullString{String: time.Now().String(), Valid: true},
		WStatus:   1,
		WDate:     int(time.Now().Unix()),
		WPubID:    int64(rand.Intn(899) + 100),
	}
	err := m.GetWDbMap().Insert(&ins)
	if err != nil {
		return nil, err
	}
	return &ins, nil
}
