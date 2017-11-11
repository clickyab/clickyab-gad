package routes

import (
	"github.com/clickyab/services/assert"
	"bytes"
	"clickyab.com/gad/config"
	"encoding/base64"
	"fmt"
	"html/template"
	"math/rand"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/mr"
	"net/http"
	"net/url"
	"clickyab.com/gad/redis"
	"clickyab.com/gad/store"
	"strconv"
	"strings"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"

	"gopkg.in/labstack/echo.v3"
)

// SingleAd is the single ad id
type SingleAd struct {
	Link   string
	Width  string
	Height string
	Src    string
	Tiny   bool

	ShowT bool
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
	Rand   int
}

func (tc *selectController) show(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	var suspicious bool
	mega := c.Param("mega")

	typ := c.Param("type")
	long := c.Request().URL.Query().Get("l")
	pos := c.Request().URL.Query().Get("pos")
	if typ == "sync" {
		typ = "web"
	}
	if typ != "web" && typ != "vast" && typ != "app" {
		return c.String(http.StatusNotFound, "not found")
	}

	ad, _ := store.Get(c.Param("ad"))
	adID, err := strconv.ParseInt(ad, 10, 64)
	if err != nil {
		// Can not find ad
		return c.String(http.StatusNoContent, "")
	}
	websiteID, err := strconv.ParseInt(c.Param("wid"), 10, 64)
	var publisher Publisher
	if typ == "web" || typ == "vast" {
		website, err := mr.NewManager().FetchWebsite(websiteID)
		assert.Nil(err)
		publisher = website
	} else {
		app, err := mr.NewManager().FetchValidAppByID(websiteID)
		publisher = app
		assert.Nil(err)
	}

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
	slotID, err := strconv.ParseInt(megaImp[fmt.Sprintf("%s%s%d", transport.SLOT, transport.DELIMITER, adID)], 10, 64)
	assert.Nil(err)

	ccu, ccuok := megaImp[fmt.Sprintf(
		"%s%s%d",
		transport.CUSTOM_CLICK_URL,
		transport.DELIMITER,
		slotID)]
	ccp, ccpok := megaImp[fmt.Sprintf(
		"%s%s%d",
		transport.CUSTOM_CLICK_PARAM,
		transport.DELIMITER,
		slotID)]
	cct := megaImp[fmt.Sprintf(
		"%s%s%d",
		transport.CUSTOM_CLICK_TYPE,
		transport.DELIMITER,
		slotID)]
	if ccuok && ccpok {
		cu, e := url.Parse(ccu)
		assert.Nil(e)
		b := base64.URLEncoding.WithPadding(rune('.')).EncodeToString([]byte(u.String()))
		if cct == "replace" {
			tu, e := url.Parse(strings.Replace(cu.String(), ccp, b, -1))
			assert.Nil(e)
			u = *tu
		} else {

			qu := cu.Query()
			qu.Set(ccp, b)
			cu.RawQuery = qu.Encode()
			u = *cu
		}

	}
	res, err := tc.makeAdData(c, typ, ads, u.String(), long, pos, rd.Scheme != "http")
	if err != nil {
		return err
	}

	imp := tc.fillImp(rd, suspicious, ads, winnerFinalBid, publisher, slotID)

	go tc.callWebWorker(publisher, slotID, adID, mega, rnd, imp, rd)
	if typ == "vast" {
		return c.XMLBlob(http.StatusOK, []byte(res))
	}
	return c.HTML(http.StatusOK, res)
}

func (tc *selectController) makeWebTemplate(c echo.Context, typ string, ads *mr.Ad, url string, long string, pos string, https bool, showT bool) (string, error) {
	buf := &bytes.Buffer{}
	switch ads.AdType {
	case mr.SingleAdType:
		res := tc.makeSingleAdData(ads, url, https, showT)
		if err := singleAdTemplate.Execute(buf, res); err != nil {
			return "", err
		}
	case mr.VideoAdType:
		res := tc.makeVideoAdData(ads, url, https)
		if err := videoAdTemplate.Execute(buf, res); err != nil {
			return "", err
		}
	case mr.DynamicAdType:
		if https {
			ads.AdAttribute.Product = strings.Replace(ads.AdAttribute.Product, "http://", "https://", -1)
			ads.AdAttribute.Logo = strings.Replace(ads.AdAttribute.Logo, "http://", "https://", -1)
		}
		res := getTemplate(ads.AdSize)
		ads.AdAttribute.Link = url
		if err := res.Execute(buf, ads.AdAttribute); err != nil {
			return "", err
		}

	}
	return buf.String(), nil
}

// makeAdData
func (tc *selectController) makeAdData(c echo.Context, typ string, ads *mr.Ad, url string, long string, pos string, https bool) (string, error) {
	if typ == "web" || typ == "app" {
		return tc.makeWebTemplate(c, typ, ads, url, long, pos, https, false)
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

func (tc *selectController) makeVideoAdData(ad *mr.Ad, url string, https bool) VideoAd {
	w, h := config.GetSizeByNum(ad.AdSize)
	src := ad.AdImg.String
	if https {
		src = strings.Replace(src, "http://", "https://", -1)
	}
	if ad.RawSlotSize != nil {
		w = ad.RawSlotSize.Width
		h = ad.RawSlotSize.Height
	}
	sa := VideoAd{
		Link:   url,
		Height: h,
		Width:  w,
		Src:    src,
		Tiny:   true,
		Rand:   rand.Intn(100),
	}
	return sa
}

func (tc *selectController) makeSingleAdData(ad *mr.Ad, url string, https, showT bool) SingleAd {
	w, h := config.GetSizeByNum(ad.AdSize)
	src := ad.AdImg.String
	if https {
		src = strings.Replace(src, "http://", "https://", -1)
	}
	sa := SingleAd{
		Link:   url,
		Height: h,
		Width:  w,
		Src:    src,
		Tiny:   true,
		ShowT:  showT,
	}
	return sa
}
func (tc *selectController) makeVastAdData(ad *mr.Ad, urll string, long string, pos string, https bool) vastTemplate {
	w, h := config.GetSizeByNum(ad.AdSize)
	_, ma := config.MakeVastLen(long, "", "", "")

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
	u, err := url.Parse(ad.AdURL.String)
	var host string
	if err == nil {
		host = u.Host
	}
	src := ad.AdImg.String
	if https {
		src = strings.Replace(src, "http://", "https://", -1)
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
