package mr

import (
	"fmt"
	"time"
)

// LoadSlotPin try to load the slot pin into system
func (m *Manager) LoadSlotPin() (res []SlotPinData, err error) {
	t := time.Now()
	r := t.Format(time.RFC3339)
	//TODO u_balance to 0
	q := fmt.Sprintf(`SELECT SP.bid,SP.chance,SP.start,SP.end,SP.direct,SP.slot_id,S.slot_size,
		A.ad_id, C.u_id, A.ad_name, A.ad_url,A.ad_code, A.ad_title, A.ad_body, A.ad_img, A.ad_status,A.ad_size,
	 A.ad_reject_reason, CA.ca_ctr , A.ad_conv, A.ad_time, A.ad_type, A.ad_mainText, A.ad_defineText,
	 A.ad_textColor, A.ad_target, A.ad_attribute, A.ad_hash_attribute, A.created_at, A.updated_at,
	 U.u_email,0 AS u_balance, C.cp_id, cp_type, cp_billing_type, cp_name, cp_network, cp_placement,
	 cp_wfilter, cp_retargeting, cp_frequency, cp_segment_id, cp_app_brand, cp_net_provider,
	 cp_app_lang, cp_app_market, cp_web_mobile, cp_web, cp_application, cp_video, cp_apps_carriers,
	 cp_longmap, cp_latmap, cp_radius, cp_opt_ctr, cp_opt_conv, cp_opt_br, cp_gender, cp_alexa,
	 cp_fatfinger, cp_under, cp_geos, cp_region, cp_country, cp_hoods, cp_isp_blacklist, cp_cat,
	 cp_like_app, cp_app, cp_app_filter, cp_keywords, cp_platforms, cp_platform_version, cp_maxbid,
	 cp_weekly_budget, cp_daily_budget, cp_total_budget, cp_weekly_spend, cp_total_spend,cp_isp,
	 cp_today_spend, cp_clicks, cp_ctr, cp_imps, cp_cpm, cp_cpa, cp_cpc, cp_conv, cp_conv_rate,
	 cp_revenue, cp_roi, cp_start, cp_end, cp_status, cp_lastupdate, cp_hour_start, cp_hour_end,
	 is_crm, cp_lock,CA.ca_id
	 	FROM slot_pin AS SP
	 	INNER JOIN slots AS S ON S.slot_id=SP.slot_id
	 	INNER JOIN ads AS A ON A.ad_id=SP.ad_id
	 	INNER JOIN campaigns_ads AS CA ON A.ad_id=CA.ad_id
	 	INNER JOIN campaigns AS C ON C.cp_id=CA.cp_id
	 	INNER JOIN users AS U ON C.u_id=U.u_id WHERE SP.start <= ? AND SP.end >=?`)

	_, err = m.GetRDbMap().Select(&res, q, r, r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
