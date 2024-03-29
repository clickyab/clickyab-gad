package models

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"github.com/clickyab/services/mysql"
)

// SharpArray type is the hack to handle # splited text in our database
type SharpArray string

// AdData min data
type AdData struct {
	AdID              int64            `json:"ad_id" db:"ad_id"`
	CampaignFrequency int              `json:"cp_frequency" db:"cp_frequency"`
	CTR               float64          `json:"ctr" db:"ctr"`
	CaCTR             sql.NullFloat64  `json:"-" db:"ca_ctr"`
	AdCTR             float64          `json:"ad_ctr" db:"-"`
	CPM               int64            `json:"cpm" db:"cpm"`
	Capping           CappingInterface `json:"capping" db:"-"`
	WinnerBid         int64            `json:"winner_bid" db:"-"`
	CampaignMaxBid    int64            `json:"cp_maxbid" db:"cp_maxbid"`
	CampaignID        int64            `json:"cp_id" db:"cp_id"`
	CampaignName      sql.NullString   `json:"cp_name" db:"cp_name"`
	AdType            int              `json:"-" db:"ad_type"`
	SlotID            int64            `json:"-" db:"-"`
	SlotPublicID      string           `json:"-" db:"-"`
	Campaign          `json:"-"`
	AdSize            int            `json:"-" db:"ad_size"`
	UserID            int64          `json:"-" db:"u_id"`
	AdName            sql.NullString `json:"-" db:"ad_name"`
	AdURL             sql.NullString `json:"-" db:"ad_url"`
	AdCode            sql.NullString `json:"-" db:"ad_code"`
	AdTitle           sql.NullString `json:"-" db:"ad_title"`
	AdBody            sql.NullString `json:"-" db:"ad_body"`
	AdImg             sql.NullString `json:"-" db:"ad_img"`
	AdStatus          int            `json:"-" db:"ad_status"`
	AdRejectReason    sql.NullString `json:"-" db:"ad_reject_reason"`
	AdConversion      int            `json:"-" db:"ad_conv"`
	AdTime            int            `json:"-" db:"ad_time"`

	AdMainText      sql.NullString         `json:"-" db:"ad_mainText"`
	AdDefineText    sql.NullString         `json:"-" db:"ad_defineText"`
	AdTextColor     sql.NullString         `json:"-" db:"ad_textColor"`
	AdTarget        sql.NullString         `json:"-" db:"ad_target"`
	AdAttribute     mysql.GenericJSONField `json:"-" db:"ad_attribute"`
	AdHashAttribute sql.NullString         `json:"-" db:"ad_hash_attribute"`
	CreatedAt       sql.NullString         `json:"-" db:"created_at"`
	UpdatedAt       sql.NullString         `json:"-" db:"updated_at"`
	UserEmail       string                 `json:"-" db:"u_email"`
	UserBalance     string                 `json:"-" db:"u_balance"`
	IsCrm           int                    `json:"-" db:"is_crm"`
	CpLock          int                    `json:"-" db:"cp_lock"`
	CampaignAdID    int64                  `json:"-" db:"ca_id"`

	Extra string `db:"-"`
}

// Campaign is a single campaign data
type Campaign struct {
	CampaignType            int             `json:"cp_type" db:"cp_type"`
	CampaignBillingType     sql.NullString  `json:"cp_billing_type" db:"cp_billing_type"`
	CampaignNetwork         int             `json:"cp_network" db:"cp_network"`
	CampaignPlacement       SharpArray      `json:"cp_placement" db:"cp_placement"`
	CampaignWebsiteFilter   SharpArray      `json:"cp_wfilter" db:"cp_wfilter"`
	CampaignRetargeting     sql.NullString  `json:"cp_retargeting" db:"cp_retargeting"`
	CampaignSegmentID       sql.NullInt64   `json:"cp_segment_id" db:"cp_segment_id"`
	CampaignNetProvider     SharpArray      `json:"cp_net_provider" db:"cp_net_provider"`
	CampaignAppBrand        SharpArray      `json:"cp_app_brand" db:"cp_app_brand"`
	CampaignAppLang         sql.NullString  `json:"cp_app_lang" db:"cp_app_lang"`
	CampaignAppMarket       sql.NullInt64   `json:"cp_app_market" db:"cp_app_market"`
	CampaignWebMobile       int             `json:"cp_web_mobile" db:"cp_web_mobile"`
	CampaignWeb             int             `json:"cp_web" db:"cp_web"`
	CampaignApplication     int             `json:"cp_application" db:"cp_application"`
	CampaignVideo           int             `json:"cp_video" db:"cp_video"`
	CampaignAppsCarriers    SharpArray      `json:"cp_apps_carriers" db:"cp_apps_carriers"`
	CampaignLongMap         sql.NullFloat64 `json:"cp_longmap" db:"cp_longmap"`
	CampaignLatMap          sql.NullFloat64 `json:"cp_latmap" db:"cp_latmap"`
	CampaignRadius          sql.NullFloat64 `json:"cp_radius" db:"cp_radius"`
	CampaignOptCTR          int             `json:"cp_opt_ctr" db:"cp_opt_ctr"`
	CampaignOptConv         int             `json:"cp_opt_conv" db:"cp_opt_conv"`
	CampaignOptBr           int             `json:"cp_opt_br" db:"cp_opt_br"`
	CampaignGender          int             `json:"cp_gender" db:"cp_gender"`
	CampaignAlexa           int             `json:"cp_alexa" db:"cp_alexa"`
	CampaignFatfinger       int             `json:"cp_fatfinger" db:"cp_fatfinger"`
	CampaignUnder           int             `json:"cp_under" db:"cp_under"`
	CampaignGeos            SharpArray      `json:"cp_geos" db:"cp_geos"`
	CampaignISP             SharpArray      `json:"cp_isp" db:"cp_isp"`
	CampaignRegion          SharpArray      `json:"cp_region" db:"cp_region"`
	CampaignCountry         SharpArray      `json:"cp_country" db:"cp_country"`
	CampaignHoods           SharpArray      `json:"cp_hoods" db:"cp_hoods"`
	CampaignIspBlacklist    SharpArray      `json:"cp_isp_blacklist" db:"cp_isp_blacklist"`
	CampaignCat             SharpArray      `json:"cp_cat" db:"cp_cat"`
	CampaignLikeApp         SharpArray      `json:"cp_like_app" db:"cp_like_app"`
	CampaignApp             SharpArray      `json:"cp_app" db:"cp_app"`
	CampaignAppFilter       SharpArray      `json:"cp_app_filter" db:"cp_app_filter"`
	CampaignKeywords        SharpArray      `json:"cp_keywords" db:"cp_keywords"`
	CampaignPlatforms       SharpArray      `json:"cp_platforms" db:"cp_platforms"`
	CampaignPlatformVersion SharpArray      `json:"cp_platform_version" db:"cp_platform_version"`
	CampaignWeeklyBudget    int             `json:"cp_weekly_budget" db:"cp_weekly_budget"`
	CampaignDailyBudget     int             `json:"cp_daily_budget" db:"cp_daily_budget"`
	CampaignTotalBudget     int             `json:"cp_total_budget" db:"cp_total_budget"`
	CampaignWeeklySpend     int             `json:"cp_weekly_spend" db:"cp_weekly_spend"`
	CampaignTotalSpend      int             `json:"cp_total_spend" db:"cp_total_spend"`
	CampaignTodaySpend      int             `json:"cp_today_spend" db:"cp_today_spend"`
	CampaignClicks          int             `json:"cp_clicks" db:"cp_clicks"`
	CampaignCTR             float64         `json:"cp_ctr" db:"cp_ctr"`
	CampaignImps            int             `json:"cp_imps" db:"cp_imps"`
	CampaignCPM             int             `json:"cp_cpm" db:"cp_cpm"`
	CampaignCPA             int             `json:"cp_cpa" db:"cp_cpa"`
	CampaignCPC             int             `json:"cp_cpc" db:"cp_cpc"`
	CampaignConv            int             `json:"cp_conv" db:"cp_conv"`
	CampaignConvRate        float64         `json:"cp_conv_rate" db:"cp_conv_rate"`
	CampaignRevenue         int             `json:"cp_revenue" db:"cp_revenue"`
	CampaignRoi             int             `json:"cp_roi" db:"cp_roi"`
	CampaignStart           int             `json:"cp_start" db:"cp_start"`
	CampaignEnd             int             `json:"cp_end" db:"cp_end"`
	CampaignStatus          int             `json:"cp_status" db:"cp_status"`
	CampaignLastupdate      int             `json:"cp_lastupdate" db:"cp_lastupdate"`
	CampaignHourStart       int             `json:"cp_hour_start" db:"cp_hour_start"`
	CampaignHourEnd         int             `json:"cp_hour_end" db:"cp_hour_end"`
}

// ByCapping sort by Capping
type ByCapping []*AdData

// StringSS formats data
func (ad AdData) StringSS() string {
	return fmt.Sprintf("%d => %d * %f (%s)", ad.CPM, ad.CampaignMaxBid, ad.AdCTR, ad.Extra)
}

func (a ByCapping) Len() int {
	return len(a)
}
func (a ByCapping) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a ByCapping) Less(i, j int) bool {

	// This is a multisort function.
	if a[i].Capping.GetSelected() != a[j].Capping.GetSelected() {
		return !a[i].Capping.GetSelected()
	}
	if a[i].Capping.GetAdCapping(a[i].AdID) != a[j].Capping.GetAdCapping(a[j].AdID) {
		return a[i].Capping.GetAdCapping(a[i].AdID) < a[j].Capping.GetAdCapping(a[j].AdID)
	}
	if a[i].Capping.GetCapping() != a[j].Capping.GetCapping() {
		return a[i].Capping.GetCapping() < a[j].Capping.GetCapping()
	}

	return a[i].CPM > a[j].CPM
}

// Scan convert the json array ino string slice
func (pa *SharpArray) Scan(src interface{}) error {
	s := &sql.NullString{}
	err := s.Scan(src)
	if err != nil {
		return err
	}

	if s.Valid {
		*pa = SharpArray(s.String)
	} else {
		*pa = ""
	}
	return nil

}

// Value try to get the string slice representation in database
func (pa SharpArray) Value() (driver.Value, error) {
	s := sql.NullString{}
	s.Valid = pa != ""
	s.String = string(pa)

	return s.Value()
}

// Has check exist value in sharpArray
func (pa SharpArray) Has(empty bool, in ...int64) (x bool) {
	if len(in) == 0 || len(pa) == 0 {
		return empty
	}
	if len(in) == 1 {
		return strings.Contains(string(pa), fmt.Sprintf("#%d#", in[0]))
	}

	reg := []string{}
	for i := range in {
		reg = append(reg, fmt.Sprintf("#%d#", in[i]))
	}

	return regexp.MustCompile("(" + strings.Join(reg, "|") + ")").MatchString(string(pa))
}

// Match check for at least one match
func (pa SharpArray) Match(empty bool, in SharpArray) bool {
	if len(in) == 0 || len(pa) == 0 {
		return empty
	}
	inTrim := strings.Trim(string(in), "# \n\t")
	return regexp.MustCompile("(" + strings.Replace(inTrim, "#", "#|#", -1) + ")").MatchString(string(pa))
}
