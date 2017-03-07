package mr

import (
	"assert"
	"config"
	"database/sql"
	"fmt"
	"redis"
	"sort"
	"strconv"
	"strings"
	"time"
	"transport"
	"utils"
)

var last time.Time = time.Now()

// LoadAds load all ads at once and return them
func (m *Manager) LoadAds() ([]AdData, error) {
	var res []AdData
	//t:= strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	t := time.Now()
	u := t.Unix()                          //return date in unixtimestamp
	h := t.Round(time.Minute).Format("15") //round time in minute scale

	query := fmt.Sprintf(`SELECT
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
	 	INNER JOIN users AS U ON C.u_id=U.u_id
		INNER JOIN campaigns_ads AS CA ON C.cp_id=CA.cp_id
		INNER JOIN ads AS A ON A.ad_id=CA.ad_id
		WHERE A.ad_status=1
				AND C.cp_status=1
				AND CA.ca_status = 1
				AND (C.cp_start <= %d OR C.cp_start=0)
				AND (C.cp_end >= %d OR C.cp_end=0)
				AND (cp_time_duration IS NULL OR cp_time_duration LIKE "%%#%s#%%")
				AND C.cp_daily_budget > C.cp_today_spend
				AND C.cp_total_budget > C.cp_total_spend
				AND U.u_balance > U.u_today_spend AND
				U.u_balance > 5000`, u, u, h)

	_, err := m.GetRDbMap().Select(
		&res,
		query,
	)
	if err != nil {
		return nil, err
	}

	for i := range res {
		//get redis key for ad
		result, err := aredis.SumHMGetField(
			transport.KeyGenDaily(transport.ADVERTISE, strconv.FormatInt(res[i].AdID, 10)),
			config.Config.Redis.Days,
			"i",
			"c",
		)
		if err != nil || result["c"] == 0 || result["i"] < config.Config.Clickyab.MinImp {
			res[i].AdCTR = config.Config.Clickyab.DefaultCTR
		} else {
			res[i].AdCTR = utils.Ctr(result["i"], result["c"])
		}
		// if res[i].CampaignNetwork != 1 {
		// 	// Web and vast
		// 	if res[i].CampaignMaxBid < config.Config.Clickyab.WebMinBid {
		// 		res[i].CampaignMaxBid = config.Config.Clickyab.WebMinBid
		// 	}
		// } else {
		// 	// app
		// 	if res[i].CampaignMaxBid < config.Config.Clickyab.AppMinBid {
		// 		res[i].CampaignMaxBid = config.Config.Clickyab.AppMinBid
		// 	}
		// }
	}

	assert.False(time.Since(last) > 5*time.Minute, "[BUG] the loader is not called for so long!")
	last = time.Now()
	return res, nil
}

// Build implode slice of string with ,
func Build(slot []string) string {
	sort.Strings(slot)
	return strings.Join(slot, ",")
}

//ToNullString invalidates a sql.NullString if empty, validates if not empty
func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

//ToNullInt64 validates a sql.NullInt64
func toNullInt64(s int64) sql.NullInt64 {
	return sql.NullInt64{Int64: s, Valid: true}
}
