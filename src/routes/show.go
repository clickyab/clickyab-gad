package routes

import (
	"assert"
	"bytes"
	"config"
	"fmt"
	"html/template"
	"math/rand"
	"middlewares"
	"mr"
	"net"
	"net/http"
	"net/url"
	"redis"
	"store"
	"strconv"
	"strings"
	"transport"
	"utils"

	"gopkg.in/labstack/echo.v3"
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
	rd := middlewares.MustGetRequestData(c)
	var suspicious bool
	mega := c.Param("mega")
	typ := c.Param("type")
	long := c.Request().URL.Query().Get("l")
	pos := c.Request().URL.Query().Get("pos")
	if typ != "web" && typ != "vast" {
		return c.String(http.StatusNotFound, "not found")
	}

	ad, _ := store.Get(c.Param("ad"))
	adID, err := strconv.ParseInt(ad, 10, 64)
	if err != nil {
		// Can not find ad
		return c.String(http.StatusNoContent, "")
	}
	websiteID, err := strconv.ParseInt(c.Param("wid"), 10, 64)
	website, err := mr.NewManager().FetchWebsite(websiteID)
	assert.Nil(err)

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

	ads, err := mr.NewManager().GetAd(adID, false)
	if err != nil {
		return c.String(http.StatusNotFound, "not found")
	}

	rnd := <-utils.ID
	u := url.URL{
		Scheme: rd.Scheme,
		Host:   rd.Host,
		Path:   fmt.Sprintf("/click/%s/%d/%s/%d/%s", typ, websiteID, mega, adID, rnd),
	}
	v := url.Values{}
	v.Set("tid", rd.TID)
	v.Set("ref", rd.Referrer)
	v.Set("parent", rd.Parent)
	u.RawQuery = v.Encode()

	res, err := tc.makeAdData(c, typ, ads, u.String(), long, pos, rd.Scheme == "https")
	if err != nil {
		return err
	}
	slotID, err := strconv.ParseInt(megaImp[fmt.Sprintf("%s%s%d", transport.SLOT, transport.DELIMITER, adID)], 10, 64)
	assert.Nil(err)
	imp := tc.fillImp(rd, suspicious, ads, winnerFinalBid, website, slotID)

	go tc.callWebWorker(website, slotID, adID, mega, rnd, imp, rd)
	if typ == "vast" {
		return c.XMLBlob(http.StatusOK, []byte(res))
	}
	return c.HTML(http.StatusOK, res)
}

func (tc *selectController) makeWebTemplate(c echo.Context, typ string, ads *mr.Ad, url string, long string, pos string) (string, error) {
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
func (tc *selectController) makeAdData(c echo.Context, typ string, ads *mr.Ad, url string, long string, pos string, https bool) (string, error) {
	if typ == "web" {
		return tc.makeWebTemplate(c, typ, ads, url, long, pos)
	}

	buf := &bytes.Buffer{}
	if !config.NonLinearVastSize(ads.AdSize) {
		res := tc.makeVastAdData(ads, url, long, pos, https)
		if err := linear.Execute(buf, res); err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	res := tc.makeVastAdData(ads, url, long, pos, https)
	if err := nonlinear.Execute(buf, res); err != nil {
		return "", err
	}

	return buf.String(), nil

}

func (tc *selectController) makeVideoAdData(ad *mr.Ad, url string) VideoAd {
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

func (tc *selectController) makeSingleAdData(ad *mr.Ad, url string) SingleAd {
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
func (tc *selectController) makeVastAdData(ad *mr.Ad, urll string, long string, pos string, https bool) vastTemplate {
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
	src := ad.AdImg.String
	if https {
		src = strings.Replace(src, "http://", "https://", 1)
	}
	sa := vastTemplate{
		Link:        template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", urll)),
		Tracking:    template.HTML(fmt.Sprintf("<![CDATA[\n%s?tv=1\n]]>", urll)),
		Height:      h,
		Width:       w,
		Src:         template.HTML(fmt.Sprintf("<![CDATA[\n%s\n]]>", strings.Replace(src, " ", "%20", -1))),
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
