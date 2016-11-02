package selectroute

import (
	"fmt"

	"redis"
	"time"

	"errors"
	"mr"

	"assert"

	"github.com/labstack/echo"
)

func (tc *selectController) Show(c echo.Context) error {
	mega := c.Param("mega")
	ad := c.Param("ad")
	megaImp, err := aredis.HGetAllString("mega_"+mega, true, 2*time.Hour)
	assert.Nil(err)
	if _, ok := megaImp[fmt.Sprintf("ad_%s", ad)]; !ok {
		return errors.New("ad not found " + ad)
	}
	ads, err := mr.NewManager().GetAd(ad)

	result := string(ads.AdImg)
	return c.HTML(200, result)

}
