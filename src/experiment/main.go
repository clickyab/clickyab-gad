package main

import (
	"assert"
	"config"
	"errors"
	"fmt"
	"models"
	"rabbit"
	"redis"
	"time"
	"transport"
	"utils"
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
	fmt.Println(utils.KeyGenDaily(transport.USER, "5"))
	aredis.IncHash("ab","b",4,false,0)
	_, err := utils.IncKeyDaily(utils.KeyGenDaily(transport.USER, "5"), "fc",1)
	assert.Nil(err)
	time.Sleep(2*time.Second)

}
