package routes

import (
	"sync"

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

	if len(errorMessage) > 0 {
		return c.JSON(http.StatusInternalServerError, errResponse{
			Res: errorMessage,
		})
	}
	return c.JSON(http.StatusOK, nil)
}
