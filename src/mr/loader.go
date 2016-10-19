package mr

import (
	"config"
	"strings"
	"time"
)

// LoadAds function @todo
func (m *Manager) LoadAds() ([]AdData, error) {
	var res []AdData
	//t:= strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	t := time.Now()
	u := t.Unix()                          //return date in unixtimestamp
	h := t.Round(time.Minute).Format("15") //round time in minute scale

	query := `SELECT
		A.ad_id, C.u_id, ad_name, ad_url,ad_code, ad_title, ad_body, ad_img, ad_status,ad_size,
	 ad_reject_reason, ad_ctr, ad_conv, ad_time, ad_type, ad_mainText, ad_defineText,
	 ad_textColor, ad_target, ad_attribute, ad_hash_attribute, A.created_at, A.updated_at,
	 U.u_email, U.u_balance, C.cp_id, cp_type, cp_billing_type, cp_name, cp_network, cp_placement,
	 cp_wfilter, cp_retargeting, cp_frequency, cp_segment_id, cp_app_brand, cp_net_provider,
	 cp_app_lang, cp_app_market, cp_web_mobile, cp_web, cp_application, cp_video, cp_apps_carriers,
	 cp_longmap, cp_latmap, cp_radius, cp_opt_ctr, cp_opt_conv, cp_opt_br, cp_gender, cp_alexa,
	 cp_fatfinger, cp_under, cp_geos, cp_region, cp_country, cp_hoods, cp_isp_blacklist, cp_cat,
	 cp_like_app, cp_app, cp_app_filter, cp_keywords, cp_platforms, cp_platform_version, cp_maxbid,
	 cp_weekly_budget, cp_daily_budget, cp_total_budget, cp_weekly_spend, cp_total_spend,
	 cp_today_spend, cp_clicks, cp_ctr, cp_imps, cp_cpm, cp_cpa, cp_cpc, cp_conv, cp_conv_rate,
	 cp_revenue, cp_roi, cp_start, cp_end, cp_status, cp_lastupdate, cp_hour_start, cp_hour_end,
	 is_crm, cp_lock
	 	FROM campaigns AS C
	 	LEFT JOIN campaigns_ads AS CA ON C.cp_id=CA.cp_id
		LEFT JOIN ads AS A ON A.ad_id=CA.ad_id
		LEFT JOIN users AS U ON C.u_id=U.u_id
		WHERE A.ad_status=1 AND C.cp_status=1 AND (C.cp_start <= ? OR C.cp_start=0)
				AND (C.cp_end >= ? OR C.cp_end=0)
				AND cp_hour_start <= ? AND cp_hour_end >= ?
				AND U.u_balance >= ?`

	_, err := m.GetDbMap().Select(
		&res,
		query,
		u,
		u,
		h,
		h,
		config.Config.Select.Balance,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FetchWebsite function @todo
func (m *Manager) FetchWebsite(publicID int) (*WebsiteData, error) {
	var res = WebsiteData{}

	query := `SELECT * FROM websites WHERE w_pub_id = ?  LIMIT 1`

	err := m.GetDbMap().SelectOne(
		&res,
		query,
		publicID,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FetchRegion function @todo
func (m *Manager) FetchRegion() (*RegionData, error) {
	var res = RegionData{}

	query := `SELECT * FROM list_locations`

	_, err := m.GetDbMap().Select(
		&res,
		query,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// FetchSlotAd fetch slot ad
func (m *Manager) FetchSlotAd(slotString string, adIDString string) ([]SlotData, error) {
	var res []SlotData
	query := `SELECT slots.slot_pubilc_id,
		slots.slot_size,
		slots_ads.sla_clicks,
		slots_ads.sla_imps,
		slots.slot_floor_cpm,
		slots_ads.ad_id
	FROM slots INNER JOIN slots_ads ON slots_ads.slot_id=slots.slot_id WHERE slots.slot_pubilc_id IN (?) AND slots.slot_lastupdate=? AND slots_ads.ad_id IN (?)`
	_, err := m.GetDbMap().Select(
		&res,
		query,
		slotString,
		time.Now().AddDate(0, 0, -1).Format("20060102"),
		adIDString,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Build imlode slice of string with ,
func Build(slot []string) string {
	return strings.Join(slot, ",")
}
