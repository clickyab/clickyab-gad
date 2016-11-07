package main

import (
	"config"
	"errors"
	"fmt"
	"models"
	"mr"
	"rabbit"
	"redis"
	"time"
	"version"
)

func findMinCap(userKey string) (int, error) {
	//get user capping data
	result, err := aredis.HGetAll(userKey, true, 72*time.Hour)
	if err != nil {
		return 0, errors.New("error")
	}
	fmt.Println(result)
	return 0, nil
}

func main() {
	config.Initialize()
	config.SetConfigParameter()
	version.PrintVersion().Info("Application started")
	models.Initialize()
	rabbit.Initialize()
	aredis.Initialize()
	//e := echo.New().NewContext(nil, nil)

	i := mr.ToNullInt64
	s := mr.ToNullString
	imp := mr.Impression{
		AdID:            i(123),
		Alexa:           i(1),
		AppID:           i(12),
		Cookie:          i(1),
		CopID:           i(345345),
		Flash:           i(0),
		CaID:            i(345),
		WebsiteID:       i(234),
		IP:              s("23dfg4"),
		ParentURL:       s("asdasdasd"),
		ReferralAddress: s("sdfwr4e"),
		Status:          i(0),
		URL:             s("asdas"),
		WinnerBid:       i(32),
		WP:              i(1),
	}
	res, err := mr.NewManager().InsertImpression(imp)
	fmt.Println(res, err)

}
