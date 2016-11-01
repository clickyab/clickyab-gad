package main

import (
	"config"
	"errors"
	"fmt"
	"rabbit"
	"redis"
	"selectroute"
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
	//models.Initialize()
	rabbit.Initialize()
	aredis.Initialize()
	//CopID:=""

	fmt.Println(selector.CalculateCtr(12, 23, 343, "54"))
}
