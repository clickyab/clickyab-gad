package mr

import (
	"time"
	"config"
)

func (m *Manager) LoadAds() ([]AdData, error) {
	var res []AdData
	//t:= strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	t:= time.Now()
	u:= t.Unix() //return date in unixtimestamp
	h:= t.Round(time.Minute).Format("15") //round time in minute scale

	query:= `select A.ad_id, C.u_id, ad_name, ad_url,ad_code, ad_title, ad_body, ad_img, ad_status,
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
	 is_crm, cp_lock from ads as A left join campaigns_ads as CA on A.ad_id=CA.ad_id
	left join campaigns as C on CA.cp_id=C.cp_id left join users as U on C.u_id=U.u_id
		where A.ad_status=1 and C.cp_status=1 and C.cp_start <= ? and C.cp_end >= ?
		and cp_hour_start <= ? and cp_hour_end >= ?
		and U.u_balance >= ?`

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
