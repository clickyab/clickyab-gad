package selectroute

import (
	"assert"
	"bytes"
	"config"
	"errors"
	"fmt"
	"mr"
	"rabbit"
	"redis"
	"strconv"
	"time"
	"transport"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

// SingleAd is the single ad id
type SingleAd struct {
	Link   string
	Width  string
	Height string
	Src    string
}

func (tc *selectController) show(c echo.Context) error {
	var imp transport.Impression
	mega := c.Param("mega")
	ad := c.Param("ad")
	megaImp, err := aredis.HGetAllString("mega_"+mega, true, 2*time.Hour)
	assert.Nil(err)
	if _, ok := megaImp[fmt.Sprintf("ad_%s", ad)]; !ok {
		return errors.New("ad not found " + ad)
	}
	adID, err := strconv.ParseInt(ad, 10, 64)
	assert.Nil(err)
	ads, err := mr.NewManager().GetAd(adID)

	w, h := config.GetSizeByNum(ads.AdSize)
	sa := SingleAd{
		Link:   ads.AdURL.String,
		Height: h,
		Width:  w,
		Src:    ads.AdImg.String,
	}

	buf := &bytes.Buffer{}
	err = singleAdTemplate.Execute(buf, sa)
	if err != nil {
		return err
	}

	go func() {
		err := rabbit.Publish("cy.imp", imp)
		if err != nil {
			logrus.WithField("cy.imp", imp).Error("error in imp worker ", err)
		}
	}()
	return c.HTML(200, buf.String())
}
