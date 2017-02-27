package model

import (
	"database/sql"
	"entity"
)

// SingleAdType constant
const SingleAdType = 0

// HtmlAdType constant
const HtmlAdType = 1

// VideoAdType constant
const VideoAdType = 3

// DynamicAdType constant
const DynamicAdType = 2

type Ad struct {
	AdID            int64          `db:"ad_id"`
	AdSize          int            `db:"ad_size"`
	UserID          int64          `db:"u_id"`
	AdName          sql.NullString `db:"ad_name"`
	AdURL           sql.NullString `db:"ad_url"`
	AdCode          sql.NullString `db:"ad_code"`
	AdTitle         sql.NullString `db:"ad_title"`
	AdBody          sql.NullString `db:"ad_body"`
	AdImg           sql.NullString `db:"ad_img"`
	AdStatus        int            `db:"ad_status"`
	AdRejectReason  sql.NullString `db:"ad_reject_reason"`
	AdCtr           float64        `db:"ad_ctr"`
	AdConv          int            `db:"ad_conv"`
	AdTime          int            `db:"ad_time"`
	AdType          int            `db:"ad_type"`
	AdMainText      sql.NullString `db:"ad_mainText"`
	AdDefineText    sql.NullString `db:"ad_defineText"`
	AdTextColor     sql.NullString `db:"ad_textColor"`
	AdTarget        sql.NullString `db:"ad_target"`
	AdAttribute     *Dynamic       `db:"ad_attribute"`
	AdHashAttribute sql.NullString `db:"ad_hash_attribute"`
	CreatedAt       sql.NullString `db:"created_at"`
	UpdatedAt       sql.NullString `db:"updated_at"`
	CampaignAdID    sql.NullInt64  `db:"ca_id"`
	campaign        Campaign       `db:"-"`
	winnerBid       int64          `db:""`
	capping         entity.Capping `db:"-"`
	cpm             int64          `db:"-"`
	ctr             float64        `db:"-"`
}

// Dynamic ad struct
type Dynamic struct {
	Link                            string `db:"-"`
	BannerTitleTextType             string `db:"-"`
	TemplateID                      string `db:"-"`
	CtaTitleTextType                string `db:"-"`
	Logo                            string `db:"-"`
	Product                         string `db:"-"`
	BannerDescriptionTextType       string `db:"-"`
	PriceTextType                   string `db:"-"`
	OffPriceTextType                string `db:"-"`
	BackgroundBannerPickerSelector  string `db:"-"`
	CtaBannerPickerSelector         string `db:"-"`
	TitleBannerPickerSelector       string `db:"-"`
	DescriptionBannerPickerSelector string `db:"-"`
	IconCtaBannerPickerSelector     string `db:"-"`
	TextCtaBannerPickerSelector     string `db:"-"`
	OffPriceBannerPickerSelector    string `db:"-"`
	LivePriceBannerPickerSelector   string `db:"-"`
}

func (ad *Ad) ID() int64 {
	return ad.AdID
}

func (ad *Ad) Type() entity.AdType {
	var m entity.AdType
	switch ad.AdType {
	case SingleAdType:
		m = entity.AdTypeBanner
	case DynamicAdType:
		m = entity.AdTypeDynamic
	case VideoAdType:
		m = entity.AdTypeVideo
	case HtmlAdType:
		m = entity.AdTypeHTML
	}
	return m
}

func (ad *Ad) Campaign() Campaign {
	return ad.campaign
}

func (ad *Ad) Capping() entity.Capping {
	return ad.capping
}

func (ad *Ad) SetCapping(c entity.Capping) {
	ad.capping = c
	ad.campaign.SetCapping(c)

}

func (ad *Ad) SetCPM(a int64) {
	ad.cpm = a
}

func (ad *Ad) CPM() int64 {
	return ad.cpm
}

func (ad *Ad) SetWinnerBID(a int64) {
	ad.winnerBid = a
}

func (ad *Ad) WinnerBID() int64 {
	return ad.winnerBid
}

func (ad *Ad) AdCTR() float64 {
	return ad.AdCtr
}

func (ad *Ad) SetCTR(a float64) {
	ad.ctr = a
}

func (ad *Ad) CTR() float64 {
	return ad.ctr
}
