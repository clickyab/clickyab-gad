package mr

// FetchWebsiteAll return all domain
func (m *Manager) FetchWebsiteAll() ([]*Website, error) {
	var res = []*Website{}
	//key := utils.Hash("WebsiteAll")
	query := `SELECT w_id,
	u_id,
	w_pub_id,
	w_domain,
	w_supplier,
	w_name,
	w_categories,
	w_minbid,
	w_floor_cpm,
	w_profile_type,
	w_status,
	w_review,
	w_alexarank,
	w_div,
	w_mobad,
	w_nativead,
	w_fatfinger,
	w_publish_start,
	w_publish_end,
	w_publish_cost,
	w_prepayment,
	w_today_ctr,
	w_today_imps,
	w_today_clicks,
	w_date,
	w_notapprovedreason,
	created_at,
	updated_at FROM websites ORDER BY w_today_clicks DESC LIMIT 100`
	_, err := m.GetRDbMap().Select(
		&res,
		query,
	)
	if err != nil {
		return nil, err
	}
	//_ = store(key, &res, time.Hour)
	return res, nil
}

// FetchCampaignAll return all campaigns
func (m *Manager) FetchCampaignAll() ([]*Campaign, error) {
	var res = []*Campaign{}
	//key := utils.Hash("CampaignAll")
	query := `SELECT * from campaigns LEFT JOIN statistics_campaigns on campaigns.cp_id = statistics_campaigns.cp_id ORDER BY statistics_campaigns.sc_clicks DESC LIMIT 100`
	_, err := m.GetRDbMap().Select(
		&res,
		query,
	)
	if err != nil {
		return nil, err
	}
	//_ = store(key, &res, time.Hour)
	return res, nil
}
