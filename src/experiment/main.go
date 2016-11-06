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

	"github.com/Sirupsen/logrus"
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

	cop := "123"

	CookieProfiles, err := mr.NewManager().FetchCookieProfile(cop)

	if CookieProfiles == nil {
		//copData := mr.CookieProfiles{
		//	Key:  cop,
		//	IP:   net.ParseIP(e.Request().RealIP()),
		//	Date: int64(time.Now().Unix()),
		//}
		CookieProfiles, err = mr.NewManager().InsertCookieProfile(cop, "192.168.1.1")
		if err != nil {
			logrus.Error("can not insert cop id , ", err)
		}
	}
	fmt.Println(CookieProfiles)

}
