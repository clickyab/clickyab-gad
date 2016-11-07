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
	"middlewares"
)

// SingleAd is the single ad id
type SingleAd struct {
	Link   string
	Width  string
	Height string
	Src    string
	Tiny   bool
}

func (tc *selectController) show(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	var imp transport.Impression
	mega := c.Param("mega")
	ad := c.Param("ad")
	WID,err:=strconv.ParseInt(c.Param("wid"),10,64)
	if err!=nil{
		// TODO : check error
		imp.Suspicious=true
	}
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
		Tiny:   false,
	}

	buf := &bytes.Buffer{}
	err = singleAdTemplate.Execute(buf, sa)
	if err != nil {
		return err
	}

	go func() {
		//validate
		res,err:=aredis.HGetAllString(fmt.Sprintf("%s%s%s",transport.MEGA,transport.DELIMITER,mega),true,2*time.Hour)
		if err!=nil{

		}

		//check ip
		w_id,err:=strconv.ParseInt(res["WS"],10,64)
		if res["IP"]!=rd.IP.String() || res["UA"]!=rd.UserAgent || w_id!=WID{
			imp.Suspicious=true
		}

		//set mega ip in redis
		tmp := []interface{}{
			"IP",
			rd.IP,
			"UA",
			rd.UserAgent,
			"WS",
			WID,
			"T",
			time.Now().Unix(),
		}
		aredis.HMSet(fmt.Sprintf("%s%s%s%s%d", transport.IMP, transport.DELIMITER, mega, transport.DELIMITER, ad), true, 2*time.Hour,tmp...)
		err = rabbit.Publish("cy.imp", imp)
		if err != nil {
			logrus.WithField("cy.imp", imp).Error("error in imp worker ", err)
		}

	}()
	return c.HTML(200, buf.String())
}
