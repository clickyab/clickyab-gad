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

	"net/http"

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

// VideoAd the video add
type VideoAd struct {
	Link   string
	Src    string
	Tiny   bool
	Width  string
	Height string
}

func (tc *selectController) show(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	var suspicious bool
	mega := c.Param("mega")
	ad := c.Param("ad")
	adID, err := strconv.ParseInt(ad, 10, 64)
	assert.Nil(err)
	WID, err := strconv.ParseInt(c.Param("wid"), 10, 64)
	assert.Nil(err)
	if err != nil {
		// TODO : check error
		suspicious = true
	}
	megaImp, err := aredis.HGetAllString(transport.MEGA+transport.DELIMITER+mega, false, 0)
	assert.Nil(err)
	var winnerBid string
	var winnerFinalBid int64
	var ok bool
	if winnerBid, ok = megaImp[fmt.Sprintf("%s%s%s", transport.ADVERTISE, transport.DELIMITER, ad)]; !ok {
		return errors.New("ad not found " + ad)
	}
	winnerFinalBid, err = strconv.ParseInt(winnerBid, 10, 64)

	ads, err := mr.NewManager().GetAd(adID)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}

	res, err := tc.makeAdData(c, ads)
	if err != nil {
		return err
	}
	imp := tc.fillImp(c, suspicious, ads, winnerFinalBid)

	go func() {
		assert.Nil(mr.NewManager().InsertImpression(&imp))
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
	return c.HTML(200, res)
}

func (selectController) fillImp(ctx echo.Context, sus bool, ads mr.Ad, winnerBid int64) transport.Impression {
	rd := middlewares.MustGetRequestData(ctx)
	adID, err := strconv.ParseInt(ctx.Param("ad"), 10, 0)
	assert.Nil(err)
	// TODO : Slot is not a big deal. but check this later
	slot, _ := strconv.ParseInt(ctx.Request().URL().QueryParam("s"), 10, 0)

	return transport.Impression{
		Suspicious:   sus,
		IP:           rd.IP,
		AdID:         adID,
		CopID:        rd.CopID,
		CampaignAdID: ads.CampaignAdID.Int64,

		URL:        rd.URL,
		CampaignID: ads.CampaignID.Int64,
		UserAgent:  rd.UserAgent,
		Time:       time.Now(),
		WinnerBID:  winnerBid,
		Status:     0,
		Web: &transport.WebSiteImp{
			Referrer:  ctx.Request().Referer(),
			ParentURL: rd.Parent,
			SlotID:    slot,
		},
	}
}

// makeAdData
func (tc *selectController) makeAdData(c echo.Context, ads mr.Ad) (string, error) {
	buf := &bytes.Buffer{}

	res := tc.makeVideoAdData(ads)
	if err := videoAdTemplate.Execute(buf, res); err != nil {
		return "", err
	}

	return buf.String(), nil

}

func (tc *selectController) makeVideoAdData(ad mr.Ad) VideoAd {
	w, h := config.GetSizeByNum(ad.AdSize)
	sa := VideoAd{
		Link:   ad.AdURL.String,
		Height: h,
		Width:  w,
		Src:    ad.AdImg.String,
		Tiny:   true,
	}
	return sa
}
