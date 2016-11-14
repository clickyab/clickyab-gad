package selectroute

import (
	"assert"
	"strconv"

	"fmt"

	"github.com/labstack/echo"
)

func (tc *selectController) click(c echo.Context) error {
	//rd := middlewares.MustGetRequestData(c)
	var suspicious bool
	adID, err := strconv.ParseInt(c.Param("ad"), 10, 64)
	assert.Nil(err)
	wID := c.Param("wid")
	impID := c.Param("imp")
	rand := c.Param("rand")

	//validate in imp id exists in redis
	//megaImp, err := aredis.HGetAllString(transport.MEGA+transport.DELIMITER+mega, false, 0)
	fmt.Println(wID, suspicious, adID, impID, rand)
	assert.Nil(err)
	return nil

}
