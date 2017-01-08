package mr

import (
	"errors"
	"fmt"
	"time"
)

// UnknownNetwork is the default network
const UnknownNetwork = 2

// PhoneData is the phone data united in one structure for filtering
type PhoneData struct {
	Brand   string
	BrandID int64
	// Model     string
	// ModelID   int64
	Carrier   string
	CarrierID int64
	Lang      string
	LangID    int64
	Network   string
	NetworkID int64
}

type tmpData struct {
	ID   int64  `db:"id"`
	Text string `db:"string"`
	Show int    `db:"show"`
}

func (m *Manager) doCacheQuery(q string, p string) (*tmpData, error) {
	res := tmpData
	key := utils.Sha1(fmt.Sprintf("caca_%s", p))
	err = fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	_, err := m.GetRDbMap().SelectOne(&res, q, p)
	if err == nil {
		return nil, errors.New("not found")
	}
	// Found one
	_ = store(key, &res, 720*time.Hour)
	return &res /**/, nil
}

// GetPhoneData try to insert/retrieve brand for phone
func (m *Manager) GetPhoneData(brand, carrier, network string) (*PhoneData, error) {
	result := PhoneData{
		Brand: brand,
		// Model:   model,
		Carrier:   carrier,
		Network:   network,
		NetworkID: UnknownNetwork,
	}
	q := "SELECT ab_id as id, ab_brand as string, ab_show as show FROM apps_brands WHERE ab_brand = ? LIMIT 1"
	t, err := m.doCacheQuery(q, result.Brand)
	if err == nil && t.Show > 0 {
		// Found one
		result.BrandID = t.ID
	}

	// q = "INSERT INTO apps_brand_models (`abm_model`,`ab_id` ) VALUES (?, ?) ON DUPLICATE KEY UPDATE ab_id=ab_id"
	// d, err = m.GetWDbMap().Exec(q, result.Model, result.BrandID)
	// if err != nil {
	// 	return nil, err
	// }
	// result.ModelID, err = d.LastInsertId()
	// if err != nil {
	// 	return nil, err
	// }
	q = "SELECT ac_id as id, ac_carrier as string , ac_show as show FROM  apps_carriers WHERE ac_carrier = ? LIMIT 1"
	t, err = m.doCacheQuery(q, result.Carrier)
	if err == nil && t.Show > 0 {
		// Found one
		result.CarrierID = t.ID
	}

	q = "SELECT an_id as id, an_network as string, an_show FROM `apps_networks` WHERE `an_network` = ? LIMIT 1;"
	t, err = m.doCacheQuery(q, result.Network)
	if err == nil && t.Show > 0 {
		// Found one
		result.NetworkID = t.ID
	}
	return &result, nil
}
