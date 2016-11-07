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

	"middlewares"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
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
	var suspicious bool
	mega := c.Param("mega")
	ad := c.Param("ad")
	WID, err := strconv.ParseInt(c.Param("wid"), 10, 64)
	assert.Nil(err)
	if err != nil {
		// TODO : check error
		suspicious = true
	}
	megaImp, err := aredis.HGetAllString("mega_"+mega, true, 2*time.Hour)
	assert.Nil(err)
	var winnerBid string
	var winnerFinalBid int64
	var ok bool
	if winnerBid, ok = megaImp[fmt.Sprintf("ad_%s", ad)]; !ok {
		return errors.New("ad not found " + ad)
	}
	winnerFinalBid, err = strconv.ParseInt(winnerBid, 10, 64)
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
	imp := tc.fillImp(c, suspicious, ads, winnerFinalBid)

	go func() {
		mr.NewManager().InsertImpression(imp)
		//validate
		res, err := aredis.HGetAllString(fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, mega), true, 2*time.Hour)
		if err != nil {

		}

		//check ip
		wID, _ := strconv.ParseInt(res["WS"], 10, 64)
		if res["IP"] != rd.IP.String() || res["UA"] != rd.UserAgent || wID != WID {
			imp.Suspicious = true
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
		err = aredis.HMSet(fmt.Sprintf("%s%s%s%s%d", transport.IMP, transport.DELIMITER, mega, transport.DELIMITER, ad), true, 2*time.Hour, tmp...)
		if err != nil {
			logrus.WithField("cy.imp", imp).Error("error in imp worker ", err)
		}
		err = rabbit.Publish("cy.imp", imp)
		if err != nil {
			logrus.WithField("cy.imp", imp).Error("error in imp worker ", err)
		}

	}()
	return c.HTML(200, buf.String())
}

func (selectController) fillImp(ctx echo.Context, sus bool, ads mr.Ad, winnerBid int64) transport.Impression {
	rd := middlewares.MustGetRequestData(ctx)
	return transport.Impression{
		Suspicious:   sus,
		IP:           rd.IP,
		AdID:         ctx.Param("ad"),
		CopID:        rd.CopID,
		CampaignAdID: ads.CampaignAdID,

		URL:        rd.URL,
		CampaignID: ads.CampaignID,
		UserAgent:  rd.UserAgent,
		Time:       time.Now(),
		WinnerBID:  winnerBid,
		Status:     0,
		Web: &transport.WebSiteImp{
			Referrer:  ctx.Request().Referer(),
			ParentURL: rd.Parent,
			SlotID:    ctx.Request().URL().QueryParam("s"),
		},
	}
}
