package main

import (
	"config"
	"errors"
	"fmt"
	"models"
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
	m:=selec

}
