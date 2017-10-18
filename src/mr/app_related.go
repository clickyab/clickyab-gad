package mr

import (
	"assert"
	"config"
	"crypto/md5"
	"database/sql"
	"fmt"
	"gmaps"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"utils"
)

// UnknownNetwork is the default network
const UnknownNetwork = 2

// App is the applications structure
type App struct {
	ID                   int64          `db:"app_id"`
	UserID               int64          `db:"u_id"`
	AppToken             string         `db:"app_token"`
	AppName              string         `db:"app_name"`
	EnAppName            string         `db:"en_app_name"`
	AppPackage           string         `db:"app_package"`
	AppSupplier          string         `db:"app_supplier"`
	AmID                 int            `db:"am_id"`
	MinBID               int64          `db:"app_minbid"`
	AppFloorCPM          sql.NullInt64  `db:"app_floor_cpm"`
	AppDIV               float64        `db:"app_div"`
	AppStatus            int            `db:"app_status"`
	AppReview            int            `db:"app_review"`
	AppTodayCTR          int64          `db:"app_today_ctr"`
	AppTodayIMPs         int64          `db:"app_today_imps"`
	AppTodayClicks       int64          `db:"app_today_clicks"`
	AppDate              int            `db:"app_date"`
	Appcat               SharpArray     `db:"app_cat"`
	AppNotApprovedReason sql.NullString `db:"app_notapprovedreason"`
	AppFatFinger         sql.NullBool   `db:"app_fatfinger"`
	CreatedAt            time.Time      `db:"created_at"`
	UpdatedAt            time.Time      `db:"updated_at"`

	AppPrepayment  int `db:"app_prepayment"`
	AppPublishCost int `db:"app_publish_cost"`
}

// CellLocation is the location of the cell
type CellLocation struct {
	ID              int64   `db:"id"`
	CellID          int64   `db:"cell_id"`
	Location        string  `db:"location"`
	Lat             float64 `db:"-"`
	Lon             float64 `db:"-"`
	NeighborhoodsID int64   `db:"neighborhoods_id"`
}

// PhoneData is the phone data united in one structure for filtering
type PhoneData struct {
	Brand   string
	BrandID int64
	// Model     string
	// ModelID   int64
	Carrier   string
	CarrierID int64
	// Lang      string
	// LangID    int64
	Network   string
	NetworkID int64
}

type tmpData struct {
	ID   int64  `db:"id"`
	Text string `db:"string"`
	Show int    `db:"show"`
}

// GetID return the id of app
func (w *App) GetID() int64 {
	return w.ID
}

// GetName return the name of object
func (w *App) GetName() string {
	return w.AppPackage
}

// FloorCPM is the floor value for this site
func (w *App) FloorCPM() int64 {
	if w.AppFloorCPM.Int64 < config.Config.Clickyab.MinCPMFloorApp {
		w.AppFloorCPM.Int64 = config.Config.Clickyab.MinCPMFloorApp
		w.AppFloorCPM.Valid = true
	}
	return w.AppFloorCPM.Int64
}

// GetActive return if app is active or not
func (w *App) GetActive() bool {
	return w.AppStatus == 0 || w.AppStatus == 1
}

// GetType of this object
func (w *App) GetType() string {
	return "app"
}

func (m *Manager) doCacheQuery(q string, p string) (*tmpData, error) {
	res := tmpData{}
	key := utils.Hash(fmt.Sprintf("caca_%s", p))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	err = m.GetRDbMap().SelectOne(&res, q, p)
	if err != nil {
		return nil, err
	}
	// Found one
	_ = store(key, &res, 720*time.Hour)
	return &res, nil
}

// GetPhoneData try to insert/retrieve brand for phone
func (m *Manager) GetPhoneData(brand, carrier, network string) *PhoneData {
	result := PhoneData{
		Brand: brand,
		// Model:   model,
		Carrier:   carrier,
		Network:   network,
		NetworkID: UnknownNetwork,
	}
	q := "SELECT ab_id as id, ab_brand as string, ab_show as `show` FROM apps_brands WHERE ab_brand = ? LIMIT 1"
	t, err := m.doCacheQuery(q, result.Brand)
	if err == nil && t.Show > 0 {
		// Found one
		result.BrandID = t.ID
	}
	q = "SELECT ac_id as id, ac_carrier as string , ac_show as `show` FROM  apps_carriers WHERE ac_carrier = ? LIMIT 1"
	t, err = m.doCacheQuery(q, result.Carrier)
	if err == nil && t.Show > 0 {
		// Found one
		result.CarrierID = t.ID
	}

	q = "SELECT an_id as id, an_network as string, an_show AS `show` FROM `apps_networks` WHERE `an_network` = ? LIMIT 1;"
	t, err = m.doCacheQuery(q, result.Network)
	if err == nil && t.Show > 0 {
		// Found one
		result.NetworkID = t.ID
	}
	return &result
}

// GetApp try to get application from the system
func (m *Manager) GetApp(token string) (*App, error) {
	res := App{}
	key := utils.Hash(fmt.Sprintf("app_%s", token))
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

// GetAppByID try to get application from the system
func (m *Manager) GetAppByID(id int64) (*App, error) {
	res := App{}
	key := utils.Hash(fmt.Sprintf("app_%d", id))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	q := "SELECT * FROM `apps` WHERE `app_id`=?"
	err = m.GetRDbMap().SelectOne(&res, q, id)
	if err != nil {
		return nil, err
	}

	_ = store(key, &res, 72*time.Hour)
	return &res, nil
}

// IsUserActive return if the user is active
func (m *Manager) IsUserActive(u int64) bool {
	key := utils.Hash(fmt.Sprintf("user_active_%d", u))
	var act bool
	err := fetch(key, &act)
	if err == nil {
		return act
	}
	q := "SELECT u_close FROM users WHERE u_id = ?"
	res, err := m.GetRDbMap().SelectInt(q, u)
	if err != nil || res != 0 {
		return false
	}
	x := true
	_ = store(key, &x, 72*time.Hour)
	return true
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
func (m *Manager) GetCellLocation(mcc, mnc, lac, cid int64, carrier string) (*CellLocation, error) {
	res := CellLocation{}
	key := utils.Hash(fmt.Sprintf("loc_%d%d%d%d", mcc, mnc, lac, cid))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}

	q := "SELECT id, cell_id, locations, neighborhoods_id FROM `finder_logs_sdk_true` WHERE `mcc`=? AND `mnc`=? AND `lac`=? AND `cid`=? LIMIT 1"
	err = m.GetRDbMap().SelectOne(&res, q, mcc, mnc, lac, cid)
	if err == nil {
		arr := strings.Split(res.Location, ",")
		assert.True(len(arr) == 2, fmt.Sprintf("[DATA-BUG] finder_logs_sdk_true location for %d is invalid", res.ID))
		res.Lat, err = strconv.ParseFloat(arr[0], 64)
		assert.Nil(err, fmt.Sprintf("[DATA-BUG] finder_logs_sdk_true location for %d is invalid", res.ID))
		res.Lon, err = strconv.ParseFloat(arr[1], 64)
		assert.Nil(err, fmt.Sprintf("[DATA-BUG] finder_logs_sdk_true location for %d is invalid", res.ID))
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

	q = "INSERT INTO `finder_logs_sdk_true` (`cell_id`,`neighborhoods_id`, `carrier`, `mcc`, `mnc`, `lac`, `cid`, `locations`, `time`) VALUES (?,?,?,?,?,?,?,?,?)"
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
	res.Lat = lat
	res.Lon = lon

	return &res, nil
}

// FetchAppByPack fetch app by package and supplier name
func (m *Manager) FetchAppByPack(pack, supplier string) (*App, error) {
	res := App{}
	key := utils.Hash(fmt.Sprintf("AppPackageSupplier_%s_%s", pack, supplier))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	q := "SELECT * FROM apps WHERE app_supplier=? AND app_package=? AND app_status NOT IN (2,3) LIMIT 1"
	err = m.GetRDbMap().SelectOne(&res, q, supplier, pack)
	if err != nil {
		return nil, err
	}
	_ = store(key, &res, time.Hour)
	return &res, nil
}

// FetchValidAppByID find app by ID
func (m *Manager) FetchValidAppByID(ID int64) (*App, error) {
	res := App{}
	key := utils.Hash(fmt.Sprintf("App_%d", ID))
	err := fetch(key, &res)
	if err == nil {
		return &res, nil
	}
	q := "SELECT * FROM apps WHERE app_id=? AND app_status NOT IN (2,3) LIMIT 1"
	err = m.GetProperDBMap().SelectOne(&res, q, ID)
	if err != nil {
		return nil, err
	}
	_ = store(key, &res, time.Hour)
	return &res, nil
}

// InsertApp insert application
func (m *Manager) InsertApp(pack, supplier string, userID int64) (*App, error) {
	if supplier == "clickyab" {
		// we are not allow to register sites from clickyab
		return nil, fmt.Errorf("the clickyab supplier is not allowed to register website on the fly")
	}
	ins := App{
		UserID:      userID,
		AppPackage:  pack,
		AppSupplier: supplier,
		AppToken:    fmt.Sprintf("%x", md5.Sum([]byte(pack+fmt.Sprintf("%d", rand.Intn(899999999999)+100000000000)))),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AppStatus:   1,
		AppDate:     int(time.Now().Unix()),
	}
	err := m.GetWDbMap().Insert(&ins)
	if err != nil {
		return nil, err
	}
	return &ins, nil
}
