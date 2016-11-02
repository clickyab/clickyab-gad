package main

import (
	"config"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"models"
	"mr"
	"rabbit"
	"redis"
	"selector"
	"sort"
	"time"
	"version"

	"github.com/labstack/echo"
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
	e := echo.New().NewContext(nil, nil)
	//CopID:=""
	time.Sleep(time.Second * 15)
	x := selector.GetAdData()
	var fff []*mr.MinAdData
	for i := range x {
		if i > 10 {
			break
		}
		x[i].Capping = mr.NewCapping(e, x[i].CpID, 0, 2)
		x[i].CPM = rand.Int63n(1000)
		if rand.Intn(100) < 50 {
			x[i].MinAdData.Capping.IncView(1)
		}
		fff = append(fff, &x[i].MinAdData)
	}
	j, _ := json.MarshalIndent(fff, "\t", "\t")
	fmt.Println(string(j))
	fmt.Println("===========================================================")
	sort.Sort(mr.ByCPM(fff))
	j, _ = json.MarshalIndent(fff, "\t", "\t")
	fmt.Println(string(j))

}
