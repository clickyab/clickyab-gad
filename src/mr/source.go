package mr

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"models/common"
	"strconv"
	"strings"
)

// SharpArray type @todo
type SharpArray []int64

// AdData type @todo
type AdData struct {
	AdID              int64                   `json:"ad_id" db:"ad_id"`
	AdSize            int                     `json:"ad_size" db:"ad_size"`
	UserID            int64                   `json:"u_id" db:"u_id"`
	AdName            sql.NullString          `json:"ad_name" db:"ad_name"`
	AdURL             sql.NullString          `json:"ad_url" db:"ad_url"`
	AdCode            sql.NullString          `json:"ad_code" db:"ad_code"`
	AdTitle           sql.NullString          `json:"ad_title" db:"ad_title"`
	AdBody            sql.NullString          `json:"ad_body" db:"ad_body"`
	AdImg             sql.NullString          `json:"ad_img" db:"ad_img"`
	AdStatus          int                     `json:"ad_status" db:"ad_status"`
	AdRejectReason    sql.NullString          `json:"ad_reject_reason" db:"ad_reject_reason"`
	AdCtr             float64                 `json:"ad_ctr" db:"ad_ctr"`
	AdConv            int                     `json:"ad_conv" db:"ad_conv"`
	AdTime            int                     `json:"ad_time" db:"ad_time"`
	AdType            int                     `json:"ad_type" db:"ad_type"`
	AdMainText        sql.NullString          `json:"ad_mainText" db:"ad_mainText"`
	AdDefineText      sql.NullString          `json:"ad_defineText" db:"ad_defineText"`
	AdTextColor       sql.NullString          `json:"ad_textColor" db:"ad_textColor"`
	AdTarget          sql.NullString          `json:"ad_target" db:"ad_target"`
	AdAttribute       common.GenericJSONField `json:"ad_attribute" db:"ad_attribute"`
	AdHashAttribute   sql.NullString          `json:"ad_hash_attribute" db:"ad_hash_attribute"`
	CreatedAt         sql.NullString          `json:"created_at" db:"created_at"`
	UpdatedAt         sql.NullString          `json:"updated_at" db:"updated_at"`
	UEmail            string                  `json:"u_email" db:"u_email"`
	UBalance          string                  `json:"u_balance" db:"u_balance"`
	CpID              int64                   `json:"cp_id" db:"cp_id"`
	CpType            int                     `json:"cp_type" db:"cp_type"`
	CpBillingType     sql.NullString          `json:"cp_billing_type" db:"cp_billing_type"`
	CpName            sql.NullString          `json:"cp_name" db:"cp_name"`
	CpNetwork         int                     `json:"cp_network" db:"cp_network"`
	CpPlacement       SharpArray              `json:"cp_placement" db:"cp_placement"`
	CpWfilter         SharpArray              `json:"cp_wfilter" db:"cp_wfilter"`
	CpRetargeting     sql.NullString          `json:"cp_retargeting" db:"cp_retargeting"`
	CpFrequency       int                     `json:"cp_frequency" db:"cp_frequency"`
	CpSegmentID       sql.NullInt64           `json:"cp_segment_id" db:"cp_segment_id"`
	CpAppBrand        sql.NullString          `json:"cp_app_brand" db:"cp_app_brand"`
	CpNetProvider     sql.NullString          `json:"cp_net_provider" db:"cp_net_provider"`
	CpAppLang         sql.NullString          `json:"cp_app_lang" db:"cp_app_lang"`
	CpAppMarket       sql.NullInt64           `json:"cp_app_market" db:"cp_app_market"`
	CpWebMobile       int                     `json:"cp_web_mobile" db:"cp_web_mobile"`
	CpWeb             int                     `json:"cp_web" db:"cp_web"`
	CpApplication     int                     `json:"cp_application" db:"cp_application"`
	CpVideo           int                     `json:"cp_video" db:"cp_video"`
	CpAppsCarriers    sql.NullString          `json:"cp_apps_carriers" db:"cp_apps_carriers"`
	CpLongmap         sql.NullString          `json:"cp_longmap" db:"cp_longmap"`
	CpLatmap          sql.NullString          `json:"cp_latmap" db:"cp_latmap"`
	CpRadius          int                     `json:"cp_radius" db:"cp_radius"`
	CpOptCtr          int                     `json:"cp_opt_ctr" db:"cp_opt_ctr"`
	CpOptConv         int                     `json:"cp_opt_conv" db:"cp_opt_conv"`
	CpOptBr           int                     `json:"cp_opt_br" db:"cp_opt_br"`
	CpGender          int                     `json:"cp_gender" db:"cp_gender"`
	CpAlexa           int                     `json:"cp_alexa" db:"cp_alexa"`
	CpFatfinger       int                     `json:"cp_fatfinger" db:"cp_fatfinger"`
	CpUnder           int                     `json:"cp_under" db:"cp_under"`
	CpGeos            SharpArray              `json:"cp_geos" db:"cp_geos"`
	CpRegion          SharpArray              `json:"cp_region" db:"cp_region"`
	CpCountry         SharpArray              `json:"cp_country" db:"cp_country"`
	CpHoods           SharpArray              `json:"cp_hoods" db:"cp_hoods"`
	CpIspBlacklist    SharpArray              `json:"cp_isp_blacklist" db:"cp_isp_blacklist"`
	CpCat             SharpArray              `json:"cp_cat" db:"cp_cat"`
	CpLikeApp         SharpArray              `json:"cp_like_app" db:"cp_like_app"`
	CpApp             SharpArray              `json:"cp_app" db:"cp_app"`
	CpAppFilter       SharpArray              `json:"cp_app_filter" db:"cp_app_filter"`
	CpKeywords        SharpArray              `json:"cp_keywords" db:"cp_keywords"`
	CpPlatforms       SharpArray              `json:"cp_platforms" db:"cp_platforms"`
	CpPlatformVersion SharpArray              `json:"cp_platform_version" db:"cp_platform_version"`
	CpMaxbid          int                     `json:"cp_maxbid" db:"cp_maxbid"`
	CpWeeklyBudget    int                     `json:"cp_weekly_budget" db:"cp_weekly_budget"`
	CpDailyBudget     int                     `json:"cp_daily_budget" db:"cp_daily_budget"`
	CpTotalBudget     int                     `json:"cp_total_budget" db:"cp_total_budget"`
	CpWeeklySpend     int                     `json:"cp_weekly_spend" db:"cp_weekly_spend"`
	CpTotalSpend      int                     `json:"cp_total_spend" db:"cp_total_spend"`
	CpTodaySpend      int                     `json:"cp_today_spend" db:"cp_today_spend"`
	CpClicks          int                     `json:"cp_clicks" db:"cp_clicks"`
	CpCtr             float64                 `json:"cp_ctr" db:"cp_ctr"`
	CpImps            int                     `json:"cp_imps" db:"cp_imps"`
	CpCpm             int                     `json:"cp_cpm" db:"cp_cpm"`
	CpCpa             int                     `json:"cp_cpa" db:"cp_cpa"`
	CpCpc             int                     `json:"cp_cpc" db:"cp_cpc"`
	CpConv            int                     `json:"cp_conv" db:"cp_conv"`
	CpConvRate        float64                 `json:"cp_conv_rate" db:"cp_conv_rate"`
	CpRevenue         int                     `json:"cp_revenue" db:"cp_revenue"`
	CpRoi             int                     `json:"cp_roi" db:"cp_roi"`
	CpStart           int                     `json:"cp_start" db:"cp_start"`
	CpEnd             int                     `json:"cp_end" db:"cp_end"`
	CpStatus          int                     `json:"cp_status" db:"cp_status"`
	CpLastupdate      int                     `json:"cp_lastupdate" db:"cp_lastupdate"`
	CpHourStart       int                     `json:"cp_hour_start" db:"cp_hour_start"`
	CpHourEnd         int                     `json:"cp_hour_end" db:"cp_hour_end"`
	IsCrm             int                     `json:"is_crm" db:"is_crm"`
	CpLock            int                     `json:"cp_lock" db:"cp_lock"`
}

// Scan convert the json array ino string slice
func (pa *SharpArray) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}

	s := strings.Split(string(b), "#")
	for i := range s {
		v, err := strconv.ParseInt(s[i], 10, 0)
		if err == nil {
			*pa = append(*pa, v)
		}
	}

	return nil

}

// Value try to get the string slice representation in database
func (pa SharpArray) Value() (driver.Value, error) {
	tmp := "#"
	arr := make([]interface{}, len(pa))
	for i := range pa {
		tmp += "%d#"
		arr[i] = pa[i]
	}
	res := fmt.Sprintf(tmp, arr...)
	return []byte(res), nil
}

// Has check exist value in sharpArray
func (pa SharpArray) Has(in int64) bool {
	for i := range pa {
		if pa[i] == in {
			return true
		}
	}
	return false
}
