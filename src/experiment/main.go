package main

import (
	"config"
	"errors"
	"fmt"
	"models"
	"rabbit"
	"redis"
	"sync"
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

	wg := sync.WaitGroup{}
	for i := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		wg.Add(1)
		go func(ij int) {
			defer wg.Done()
			fmt.Print(ij)
		}(i)
	}

	wg.Wait()
}
