package routes

import (
	"assert"
	"config"
	"encoding/base64"
	"fmt"
	"middlewares"
	"redis"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"

	"gopkg.in/labstack/echo.v3"
)

const invalidRedirect = "__retarget__"

var lock = &sync.Mutex{}

func cleanHashKey(key string, exp time.Duration) map[string]string {
	k, err := aredis.HGetAllString(key, false, 0)
	if err != nil {
		return make(map[string]string)
	}
	for i := range k {
		t, err := strconv.ParseInt(k[i], 10, 0)
		if err != nil {
			delete(k, i)
			continue
		}
		if s := time.Since(time.Unix(t, 0)); s > exp || s < 0 {
			delete(k, i)
		}
	}

	return k
}

func (tc *selectController) retarget(c echo.Context) error {

	middlewares.SafeGO(c, false, func() {
		cpID := strings.Trim(c.Param("cpid"), " \n\t")
		_, err := strconv.ParseInt(cpID, 10, 0)
		if c.Get(invalidRedirect) == nil && err == nil {
			rd := middlewares.MustGetRequestData(c)
			// set the retargetting key
			key := retargetingKey(rd.CopID)
			exp := time.Duration(config.Config.Clickyab.RetargettingHour) * time.Hour
			keys := cleanHashKey(key, exp)
			keys[cpID] = fmt.Sprint(time.Now().Unix())
			_ = aredis.HMSet(key, exp, keys)
		}
	})
	data, err := base64.StdEncoding.DecodeString(message)
	assert.Nil(err)
	c.Response().Header().Set("Content-Type", "image/png")
	return c.String(http.StatusOK, string(data))
}
