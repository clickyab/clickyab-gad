package mr

import (
	"database/sql"
	"fmt"
	"time"
	"utils"
)

// WebsiteData type @todo
type WebsiteData struct {
	WID                int64          `json:"w_id" db:"w_id"`
	UserID             int64          `json:"u_id" db:"u_id"`
	WPubID             int64          `json:"w_pub_id" db:"w_pub_id"`
	WDomain            sql.NullString `json:"w_domain" db:"w_domain"`
	WName              sql.NullString `json:"w_name" db:"w_name"`
	WCategories        SharpArray     `json:"w_categories" db:"w_categories"`
	WMinBid            sql.NullInt64  `json:"w_minbid" db:"w_minbid"`
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

// FetchWebsiteByPublicID function @todo
func (m *Manager) FetchWebsiteByPublicID(publicID int64) (*WebsiteData, error) {
	var res = WebsiteData{}
	key := utils.Sha1(fmt.Sprintf("Website_%d", publicID))
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
func (m *Manager) FetchWebsite(ID int64) (*WebsiteData, error) {
	var res = WebsiteData{}
	key := utils.Sha1(fmt.Sprintf("WebsiteID_%d", ID))
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
func (m *Manager) FetchWebsiteByDomain(domain string) (*WebsiteData, error) {
	var res = WebsiteData{}
	key := utils.Sha1(fmt.Sprintf("WebsiteDomain_%s", domain))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	query := `SELECT * FROM websites WHERE w_domain = ? AND w_status NOT IN (2,3)  LIMIT 1`

	err = m.GetRDbMap().SelectOne(
		&res,
		query,
		domain,
	)
	if err != nil {
		return nil, err
	}

	_ = store(key, &res, time.Hour)
	return &res, nil
}
