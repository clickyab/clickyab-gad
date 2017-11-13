package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"clickyab.com/gad/utils"
)

// SingleAdType constant
const SingleAdType = 0

// VideoAdType constant
const VideoAdType = 3

// DynamicAdType constant
const DynamicAdType = 2

// NativeAdType constant
const NativeAdType = 4

// Dynamic ad struct
type Dynamic struct {
	Link                            string `json:"-"`
	BannerTitleTextType             string `json:"banner_title_text_type"`
	TemplateID                      string `json:"template_id"`
	CtaTitleTextType                string `json:"cta_title_text_type"`
	Logo                            string `json:"logo"`
	Product                         string `json:"product"`
	BannerDescriptionTextType       string `json:"banner_description_text_type"`
	PriceTextType                   string `json:"price_text_type"`
	OffPriceTextType                string `json:"off_price_text_type"`
	BackgroundBannerPickerSelector  string `json:"background_banner_picker_selector"`
	CtaBannerPickerSelector         string `json:"cta_banner_picker_selector"`
	TitleBannerPickerSelector       string `json:"title_banner_picker_selector"`
	DescriptionBannerPickerSelector string `json:"description_banner_picker_selector"`
	IconCtaBannerPickerSelector     string `json:"icon_cta_banner_picker_selector"`
	TextCtaBannerPickerSelector     string `json:"text_cta_banner_picker_selector"`
	OffPriceBannerPickerSelector    string `json:"off_price_banner_picker_selector"`
	LivePriceBannerPickerSelector   string `json:"live_price_banner_picker_selector"`
}

// Scan convert the json array ino string slice
func (gjf *Dynamic) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
		return nil
	default:
		return errors.New("unsupported type")
	}

	return json.Unmarshal(b, gjf)
}

// Value try to get the string slice representation in database
func (gjf Dynamic) Value() (driver.Value, error) {
	return json.Marshal(gjf)
}

//Ad struct ad info
type Ad struct {
	AdID            int64              `json:"ad_id" db:"ad_id"`
	AdSize          int                `json:"ad_size" db:"ad_size"`
	UserID          int64              `json:"u_id" db:"u_id"`
	AdName          sql.NullString     `json:"ad_name" db:"ad_name"`
	AdURL           sql.NullString     `json:"ad_url" db:"ad_url"`
	AdCode          sql.NullString     `json:"ad_code" db:"ad_code"`
	AdTitle         sql.NullString     `json:"ad_title" db:"ad_title"`
	AdBody          sql.NullString     `json:"ad_body" db:"ad_body"`
	AdImg           sql.NullString     `json:"ad_img" db:"ad_img"`
	AdStatus        int                `json:"ad_status" db:"ad_status"`
	AdRejectReason  sql.NullString     `json:"ad_reject_reason" db:"ad_reject_reason"`
	AdCtr           float64            `json:"ad_ctr" db:"ad_ctr"`
	AdConv          int                `json:"ad_conv" db:"ad_conv"`
	AdTime          int                `json:"ad_time" db:"ad_time"`
	AdType          int                `json:"ad_type" db:"ad_type"`
	AdMainText      sql.NullString     `json:"ad_mainText" db:"ad_mainText"`
	AdDefineText    sql.NullString     `json:"ad_defineText" db:"ad_defineText"`
	AdTextColor     sql.NullString     `json:"ad_textColor" db:"ad_textColor"`
	AdTarget        sql.NullString     `json:"ad_target" db:"ad_target"`
	AdAttribute     *Dynamic           `json:"ad_attribute" db:"ad_attribute"`
	AdHashAttribute sql.NullString     `json:"ad_hash_attribute" db:"ad_hash_attribute"`
	CreatedAt       sql.NullString     `json:"created_at" db:"created_at"`
	UpdatedAt       sql.NullString     `json:"updated_at" db:"updated_at"`
	CampaignAdID    sql.NullInt64      `db:"ca_id" json:"ca_id"`
	CampaignID      sql.NullInt64      `db:"cp_id" json:"cp_id"`
	CampaignName    sql.NullString     `db:"cp_name" json:"cp_name"`
	RawSlotSize     *RawSlotDimensions `db:"-" json:"-"`
}

// RawSlotDimensions is the raw information of size
type RawSlotDimensions struct {
	Width  string
	Height string
}

//GetAd get data ad from id
func (m *Manager) GetAd(id int64, withCPName bool) (*Ad, error) {
	var ad Ad
	key := utils.Hash(fmt.Sprintf("adData_%d", id))
	err := fetch(key, &ad)
	if err == nil {
		return &ad, nil
	}
	var query string
	if withCPName {
		query = `SELECT ads.*,campaigns_ads.ca_id,campaigns_ads.cp_id, campaigns.cp_name FROM ads
	LEFT JOIN campaigns_ads ON ads.ad_id = campaigns_ads.ad_id
	LEFT JOIN campaigns ON campaigns_ads.cp_id = campaigns.cp_id WHERE ads.ad_id = ? LIMIT 1`
	} else {
		query = `SELECT ads.*,campaigns_ads.ca_id,campaigns_ads.cp_id, "CPNAME" as cp_name FROM ads
	LEFT JOIN campaigns_ads ON ads.ad_id = campaigns_ads.ad_id
	WHERE ads.ad_id = ? LIMIT 1`

	}
	err = m.GetRDbMap().SelectOne(
		&ad,
		query,
		id,
	)
	if err != nil {
		return &ad, err
	}
	_ = store(key, &ad, time.Minute)
	return &ad, nil
}
