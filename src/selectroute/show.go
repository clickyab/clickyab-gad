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

	"html/template"
	"math/rand"

	"net"

	"net/url"

	"store"

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

type vastTemplate struct {
	Link        template.HTML
	Tracking    template.HTML
	Width       string
	Height      string
	Src         template.HTML
	Tiny        bool
	Linear      bool
	RND         int64
	RND2        int64
	SkipOffset  string
	Duration    string
	Video       bool
	Title       template.HTML
	Description template.HTML
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
	reserve := c.Param("reserve")
	fmt.Println(reserve)
	k, ok := store.Get(reserve)
	if !ok {
		return c.String(http.StatusNotFound, "no key")
	}
	return c.HTML(http.StatusOK, k)
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
	tmp := map[string]string{
		"IP":     rd.IP.String(),
		"UA":     rd.UserAgent,
		"WS":     strconv.FormatInt(WID, 10),
		"T":      strconv.FormatInt(time.Now().Unix(), 10),
		"S":      strconv.FormatInt(slotID, 10),
		"IMPR":   strconv.FormatInt(imp.ID, 10),
		"RND":    rand,
		"WIN":    strconv.FormatInt(imp.WinnerBID, 10),
		"CPADID": strconv.FormatInt(imp.CampaignAdID, 10),
		"SLAID":  strconv.FormatInt(imp.SlaID, 10),
	}
	// TODO : Config time
	err = aredis.HMSet(fmt.Sprintf("%s%s%s%s%d", transport.IMP, transport.DELIMITER, mega, transport.DELIMITER, adID), 2*time.Hour, tmp)
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

func (tc *selectController) makeWebTemplate(c echo.Context, typ string, ads mr.Ad, url string, long string, pos string) (string, error) {
	buf := &bytes.Buffer{}

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
	return buf.String(), nil
}

// makeAdData
func (tc *selectController) makeAdData(c echo.Context, typ string, ads mr.Ad, url string, long string, pos string) (string, error) {
	if typ == "web" {
		return tc.makeWebTemplate(c, typ, ads, url, long, pos)
	}

	buf := &bytes.Buffer{}
	if !config.NonLinearVastSize(ads.AdSize) {
		res := tc.makeVastAdData(ads, url, long, pos)
		if err := linear.Execute(buf, res); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	res := tc.makeVastAdData(ads, url, long, pos)
	if err := nonlinear.Execute(buf, res); err != nil {
		return "", err
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
func (tc *selectController) makeVastAdData(ad mr.Ad, urll string, long string, pos string) vastTemplate {
	w, h := config.GetSizeByNum(ad.AdSize)
	_, ma := config.MakeVastLen(long)

	skipOffset := config.Config.Clickyab.Vast.DefaultSkipOff
	duration := config.Config.Clickyab.Vast.DefaultDuration
	if k, ok := ma[pos]; ok {
		duration = k[2]
		if len(k) == 4 {
			skipOffset = k[3]
		}
	}
	var v = ad.AdType == 3
	r := rand.Int63n(99999)
	r2 := rand.Int63n(99999)
	u, _ := url.Parse(ad.AdURL.String)
	host, _, _ := net.SplitHostPort(u.Host)

	sa := vastTemplate{
		Link:        template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", urll)),
		Tracking:    template.HTML(fmt.Sprintf("<![CDATA[\n%s?tv=1\n]]>", urll)),
		Height:      h,
		Width:       w,
		Src:         template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", ad.AdImg.String)),
		Tiny:        true,
		RND:         r,
		RND2:        r2,
		Video:       v,
		Duration:    duration,
		SkipOffset:  skipOffset,
		Title:       template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", host)),
		Description: template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", ad.AdBody.String)),
	}
	return sa
}
