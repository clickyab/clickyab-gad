package routes

import (
	"models"
	"redis"
	"sync"

	"mr"
	"net/http"

	"gopkg.in/labstack/echo.v3"
)

var (
	lock1 sync.RWMutex
)

type errResponse struct {
	Res []string `json:"res"`
}

// healthz halth route check

func (tc *selectController) healthz(c echo.Context) error {
	lock1.Lock()
	defer lock1.Unlock()
	var errorMessage = []string{}
	//check mysql
	errSlice := models.MysqlHealth()
	for i := range errSlice {
		errorMessage = append(errorMessage, errSlice[i].Error())
	}

	//check redis
	errSlice = aredis.RedisHealth()
	for i := range errSlice {
		errorMessage = append(errorMessage, errSlice[i].Error())
	}

	// check Ad pool
	errSlice = mr.LoaderHealth()
	for i := range errSlice {
		errorMessage = append(errorMessage, errSlice[i].Error())
	}

	if len(errorMessage) > 0 {
		return c.JSON(http.StatusInternalServerError, errResponse{
			Res: errorMessage,
		})
	}
	return c.JSON(http.StatusOK, nil)
}
