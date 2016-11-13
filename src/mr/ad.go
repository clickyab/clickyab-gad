package mr

import (
	"database/sql"
	"models/common"
)

const SingleAdType  = 0
const VideoAdType  = 3

//Ad struct ad info
type Ad struct {
	AdID            int64                   `json:"ad_id" db:"ad_id"`
	AdSize          int                     `json:"ad_size" db:"ad_size"`
	UserID          int64                   `json:"u_id" db:"u_id"`
	AdName          sql.NullString          `json:"ad_name" db:"ad_name"`
	AdURL           sql.NullString          `json:"ad_url" db:"ad_url"`
	AdCode          sql.NullString          `json:"ad_code" db:"ad_code"`
	AdTitle         sql.NullString          `json:"ad_title" db:"ad_title"`
	AdBody          sql.NullString          `json:"ad_body" db:"ad_body"`
	AdImg           sql.NullString          `json:"ad_img" db:"ad_img"`
	AdStatus        int                     `json:"ad_status" db:"ad_status"`
	AdRejectReason  sql.NullString          `json:"ad_reject_reason" db:"ad_reject_reason"`
	AdCtr           float64                 `json:"ad_ctr" db:"ad_ctr"`
	AdConv          int                     `json:"ad_conv" db:"ad_conv"`
	AdTime          int                     `json:"ad_time" db:"ad_time"`
	AdType          int                     `json:"ad_type" db:"ad_type"`
	AdMainText      sql.NullString          `json:"ad_mainText" db:"ad_mainText"`
	AdDefineText    sql.NullString          `json:"ad_defineText" db:"ad_defineText"`
	AdTextColor     sql.NullString          `json:"ad_textColor" db:"ad_textColor"`
	AdTarget        sql.NullString          `json:"ad_target" db:"ad_target"`
	AdAttribute     common.GenericJSONField `json:"ad_attribute" db:"ad_attribute"`
	AdHashAttribute sql.NullString          `json:"ad_hash_attribute" db:"ad_hash_attribute"`
	CreatedAt       sql.NullString          `json:"created_at" db:"created_at"`
	UpdatedAt       sql.NullString          `json:"updated_at" db:"updated_at"`
	CampaignAdID    sql.NullInt64           `db:"ca_id" json:"ca_id"`
	CampaignID      sql.NullInt64           `db:"cp_id" json:"cp_id"`
}

//GetAd get data ad from id
func (m *Manager) GetAd(id int64) (Ad, error) {
	var Ads Ad
	query := `SELECT ads.*,campaigns_ads.ca_id,campaigns_ads.cp_id FROM ads LEFT JOIN campaigns_ads ON ads.ad_id = campaigns_ads.ad_id WHERE ads.ad_id = ? LIMIT 1`
	err := m.GetRDbMap().SelectOne(
		&Ads,
		query,
		id,
	)
	if err != nil {
		return Ads, err
	}

	return Ads, nil
}
