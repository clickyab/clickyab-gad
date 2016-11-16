package selectroute

import (
	"assert"
	"bytes"
	"config"
	"fmt"
	"middlewares"
	"mr"
	"net/http"
	"rabbit"
	"redis"
	"strconv"
	"time"
	"transport"
	"utils"

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
	typ:=c.Param("type")
	if typ!="web" && typ!="vast"{
		return c.String(http.StatusNotFound,"not found")
	}
	ad := c.Param("ad")
	adID, err := strconv.ParseInt(ad, 10, 64)
	assert.Nil(err)
	websiteID, err := strconv.ParseInt(c.Param("wid"), 10, 64)

	//TODO :validate Wid compare to us

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
		return c.String(http.StatusNotFound, "ad not found")
	}
	winnerFinalBid, err = strconv.ParseInt(winnerBid, 10, 64)
	ads, err := mr.NewManager().GetAd(adID)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}
	rand := <-utils.ID
	url := fmt.Sprintf("%s://%s/click/%d/%s/%d/%s", rd.Proto, rd.URL, websiteID, mega, adID, rand)
	res, err := tc.makeAdData(c,typ,ads, url)
	if err != nil {
		return err
	}
	slotID, err := strconv.ParseInt(megaImp[fmt.Sprintf("%s%s%d", transport.SLOT, transport.DELIMITER, adID)], 10, 64)
	assert.Nil(err)
	imp := tc.fillImp(c, suspicious, ads, winnerFinalBid, websiteID, slotID)

	go tc.callWorker(websiteID, slotID, adID, mega, rand, imp, rd)
	return c.HTML(200, res)
}

func (selectController) callWorker(WID int64, slotID int64, adID int64, mega string, rand string, imp transport.Impression, rd *middlewares.RequestData) {
	m := mr.NewManager()
	var err error
	imp.SlaID, err = m.InsertSlotAd(slotID, adID)
	if err != nil {
		// not important error
		logrus.Warn(err)
	}
	assert.Nil(m.InsertImpression(&imp))
	//validate
	res, err := aredis.HGetAllString(fmt.Sprintf("%s%s%s", transport.MEGA, transport.DELIMITER, mega), true, 2*time.Hour)
	assert.Nil(err)

	//check ip
	wID, _ := strconv.ParseInt(res["WS"], 10, 64)
	if res["IP"] != rd.IP.String() || res["UA"] != rd.UserAgent || wID != WID {
		imp.Suspicious = true
	}

	// TODO : Use constant not strings
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
		"S",
		slotID,
		"IMPR",
		imp.ID,
		"RND",
		rand,
		"WIN",
		imp.WinnerBID,
		"CPADID",
		imp.CampaignAdID,
		"SLAID",
		imp.SlaID,
	}
	// TODO : Config time
	err = aredis.HMSet(fmt.Sprintf("%s%s%s%s%d", transport.IMP, transport.DELIMITER, mega, transport.DELIMITER, adID), 2*time.Hour, tmp...)
	if err != nil {
		logrus.WithField("cy.imp", imp).Error("error in hmset", err)
	}
	err = rabbit.Publish("cy.imp", imp)
	if err != nil {
		logrus.WithField("cy.imp", imp).Error("error in  publishing job", err)
	}
}

func (selectController) fillImp(ctx echo.Context, sus bool, ads mr.Ad, winnerBid int64, websiteID int64, slotID int64) transport.Impression {
	rd := middlewares.MustGetRequestData(ctx)
	adID, err := strconv.ParseInt(ctx.Param("ad"), 10, 0)
	assert.Nil(err)

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
			SlotID:    slotID,
			WebsiteID: websiteID,
		},
	}
}

// makeAdData
func (tc *selectController) makeAdData(c echo.Context,typ string, ads mr.Ad, url string) (string, error) {
	buf := &bytes.Buffer{}
	if typ == "web"{
		switch ads.AdType {
		case mr.SingleAdType:
			res := tc.makeSingleAdData(ads, url)
			if err := singleAdTemplate.Execute(buf, res); err != nil {
				return "", err
			}
		case mr.VideoAdType:
			res := tc.makeVideoAdData(ads, url)
			if err := videoAdTemplate.Execute(buf, res); err != nil {
				return "", err
			}
		case mr.DynamicAdType:
			res := getTemplate("", ads.AdSize)
			ads.AdAttribute.Link = url
			if err := res.Execute(buf, ads.AdAttribute); err != nil {
				return "", err
			}

		}
	}else{

	}


	return buf.String(), nil

}

func (tc *selectController) makeVideoAdData(ad mr.Ad, url string) VideoAd {
	w, h := config.GetSizeByNum(ad.AdSize)
	sa := VideoAd{
		Link:   url,
		Height: h,
		Width:  w,
		Src:    ad.AdImg.String,
		Tiny:   true,
	}
	return sa
}

func (tc *selectController) makeSingleAdData(ad mr.Ad, url string) SingleAd {
	w, h := config.GetSizeByNum(ad.AdSize)
	sa := SingleAd{
		Link:   url,
		Height: h,
		Width:  w,
		Src:    ad.AdImg.String,
		Tiny:   true,
	}
	return sa
}
