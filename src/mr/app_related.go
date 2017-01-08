package mr

import (
	"errors"
	"fmt"
	"gmaps"
	"time"
	"utils"
)

// UnknownNetwork is the default network
const UnknownNetwork = 2

// App is the applications structure
type App struct {
	ID                   int64      `db:"app_id"`
	UserID               int64      `db:"u_id"`
	AppToken             string     `db:"app_token"`
	AppName              string     `db:"app_name"`
	EnAppName            string     `db:"en_app_name"`
	AppPackage           string     `db:"app_package"`
	AmID                 int        `db:"am_id"`
	MinBID               int64      `db:"app_minbid"`
	AppFloorCPM          int64      `db:"app_floor_cpm"`
	AppDIV               float64    `db:"app_div"`
	AppStatus            int        `db:"app_status"`
	AppReview            int        `db:"app_review"`
	AppTodayCTR          int64      `db:"app_today_ctr"`
	AppTodayIMPs         int64      `db:"app_today_imps"`
	AppTodayClicks       int64      `db:"app_today_clicks"`
	AppDate              int        `db:"app_date"`
	Appcat               SharpArray `db:"app_cat"`
	AppNotApprovedReason string     `db:"app_notapprovedreason"`
	AppFatFinger         string     `db:"app_fatfinger"`
	CreatedAt            time.Time  `db:"created_at"`
	UpdatedAt            time.Time  `db:"updated_at"`
}

// CellLocation is the location of the cell
type CellLocation struct {
	ID              int64  `db:"id"`
	CellID          int64  `db:"cell_id"`
	Location        string `db:"location"`
	NeighborhoodsID int64  `db:"neighborhoods_id"`
}

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
	res := tmpData{}
	key := utils.Sha1(fmt.Sprintf("caca_%s", p))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	err = m.GetRDbMap().SelectOne(&res, q, p)
	if err == nil {
		return nil, errors.New("not found")
	}
	// Found one
	_ = store(key, &res, 720*time.Hour)
	return &res, nil
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

// GetApp try to get application from the system
func (m *Manager) GetApp(token string) (*App, error) {
	res := App{}
	key := utils.Sha1(fmt.Sprintf("app_%s", token))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	q := "SELECT * FROM `apps` WHERE `app_token`=?"
	err = m.GetRDbMap().SelectOne(&res, q, token)
	if err != nil {
		return nil, err
	}

	_ = store(key, &res, 72*time.Hour)
	return &res, nil
}

// IsUserActive return if the user is active
func (m *Manager) IsUserActive(u int64) bool {
	q := "SELECT u_close FROM users WHERE u_id = ?"
	res, err := m.GetRDbMap().SelectInt(q, u)
	if err != nil || res != 0 {
		return true
	}

	return false
}

func (m *Manager) findCell(lat, lon float64) (int64, int64, error) {
	tmp := struct {
		ID  int64 `db:"id"`
		NID int64 `db:"neighborhoods_id"`
	}{}
	q := "SELECT id,neighborhoods_id FROM finder_cells WHERE top_left_lat >= ? AND bottom_left_lat <= ? AND top_left_long <= ? AND bottom_right_long >=? LIMIT 1"
	err := m.GetRDbMap().SelectOne(&tmp, q, lat, lat, lon, lon)
	if err != nil {
		return 0, 0, err
	}

	return tmp.ID, tmp.NID, nil
}

// GetCellLocation try to get cell location from mmap database
func (m *Manager) GetCellLocation(mcc, mnc, lac, cid int, carrier string) (*CellLocation, error) {
	res := CellLocation{}
	key := utils.Sha1(fmt.Sprintf("loc_%d%d%d%d", mcc, mnc, lac, cid))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	q := "SELECT id, cell_id, locations, neighborhoods_id FROM `finder_logs_sdk_true` WHERE `mcc`=? AND `mnc`=? AND `lac`=? AND `cid`=? LIMIT 1"
	err = m.GetRDbMap().SelectOne(&res, q, mcc, mnc, lac, cid)
	if err == nil {
		_ = store(key, &res, 72*time.Hour)
		return &res, nil
	}

	lat, lon, err := gmaps.LockUp(mcc, mnc, lac, cid)
	if err != nil {
		return nil, err
	}

	cellID, NeighborhoodID, err := m.findCell(lat, lon)
	if err != nil {
		return nil, err
	}

	q = "INSERT INTO `finder_logs_sdk_true` (`cell_id`,`neighborhoods_id`, `carrier`, `mcc`, `mnc`, `lac`, `cid`, `locations`, `time`) VALUES (?,?,?,?,?,?,?,?)"
	t, err := m.GetWDbMap().Exec(q, cellID, NeighborhoodID, carrier, mcc, mnc, lac, cid, fmt.Sprintf("%f,%f", lat, lon), time.Now().Unix())
	if err != nil {
		return nil, err
	}
	res.ID, err = t.LastInsertId()
	if err != nil {
		return nil, err
	}
	res.CellID = cellID
	res.NeighborhoodsID = NeighborhoodID
	res.Location = fmt.Sprintf("%f,%f", lat, lon)

	return &res, nil
}