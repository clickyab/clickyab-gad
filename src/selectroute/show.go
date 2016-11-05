package selectroute

import (
	"fmt"

	"redis"
	"time"

	"errors"
	"mr"

	"assert"

	"strconv"

	"bytes"
	"config"

	"github.com/labstack/echo"
)

type SingleAd struct {
	Link   string
	Width  string
	Height string
	Src    string
}

func (tc *selectController) Show(c echo.Context) error {
	mega := c.Param("mega")
	ad := c.Param("ad")
	megaImp, err := aredis.HGetAllString("mega_"+mega, true, 2*time.Hour)
	assert.Nil(err)
	if _, ok := megaImp[fmt.Sprintf("ad_%s", ad)]; !ok {
		return errors.New("ad not found " + ad)
	}
	adId, _ := strconv.ParseInt(ad, 10, 64)
	ads, err := mr.NewManager().GetAd(adId)

	w, h := config.GetSizeByNum(ads.AdSize)
	fmt.Println(h)
	sa := SingleAd{
		Link:   ads.AdURL.String,
		Height: h,
		Width:  w,
		Src:    ads.AdImg.String,
	}

	buf := &bytes.Buffer{}
	err = SingleAdTemplate.Execute(buf, sa)
	if err != nil {
		return err
	}

	return c.HTML(200, buf.String())
}
