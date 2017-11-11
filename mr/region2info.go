package mr

import (
	"errors"
	"fmt"
	"time"

	"clickyab.com/gad/utils"
)

//Province struct province info
type Province struct {
	ID     int64  `id:"location_id" db:"location_id"`
	Name   string `json:"location_name" db:"location_name"`
	NameFa string `json:"location_name_persian" db:"location_name_persian"`
	Master int    `json:"location_master" db:"location_master"`
	Select int    `json:"location_select" db:"location_select"`
}

//ConvertProvince2Info get data province from string
func (m *Manager) ConvertProvince2Info(name string) (Province, error) {
	var province Province
	if len(name) < 2 {
		return province, errors.New("invalid province name")
	}
	key := utils.Sha1("Province_" + name)
	err := fetch(key, &province)
	if err == nil {
		return province, nil
	}

	query := `SELECT * FROM list_locations WHERE location_name = ? LIMIT 1`
	err = m.GetRDbMap().SelectOne(
		&province,
		query,
		name,
	)
	if err != nil {
		return province, err
	}

	_ = store(key, &province, 72*time.Hour)
	return province, nil
}

//ConvertProvinceID2Info get data province from id
func (m *Manager) ConvertProvinceID2Info(id int64) (Province, error) {
	var province Province
	if id < 1 {
		return province, errors.New("invalid province name")
	}
	key := utils.Sha1(fmt.Sprintf("Province_%d", id))
	err := fetch(key, &province)
	if err == nil {
		return province, nil
	}

	query := `SELECT * FROM list_locations WHERE location_id = ? LIMIT 1`
	err = m.GetRDbMap().SelectOne(
		&province,
		query,
		id,
	)
	if err != nil {
		return province, err
	}

	_ = store(key, &province, 72*time.Hour)
	return province, nil
}
