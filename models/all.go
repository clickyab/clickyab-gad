package models

// FetchWebsiteAll return all domain
func (m *Manager) FetchWebsiteAll() ([]*Website, error) {
	var res = []*Website{}
	//key := utils.Hash("WebsiteAll")
	query := `SELECT * FROM websites ORDER BY w_today_clicks DESC LIMIT 100`
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
