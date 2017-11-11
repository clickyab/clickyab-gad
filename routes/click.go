package routes

import (
	"clickyab.com/gad/assert"
	"clickyab.com/gad/config"
	"fmt"
	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/mr"
	"net/http"
	"clickyab.com/gad/rabbit"
	"clickyab.com/gad/redis"
	"strconv"
	"strings"
	"time"
	"clickyab.com/gad/transport"
	"clickyab.com/gad/utils"

	"gopkg.in/labstack/echo.v3"
)

const (
	suspDuplicateClick = 1
	suspFastClick      = 9
	suspSlowClick      = 8
	suspNoAdFound      = 16
	suspWSMismatch     = 94
	suspRndMismatch    = 95
	suspIPMismatch     = 96
	suspUAMismatch     = 97
	suspInvalidWebsite = 98
	suspInvalidApp     = 99
)

func assertNil(status bool, err error) {
	if !status {
		return
	}

	assert.Nil(err)
}

func changeStatus(oldStatus int64, newStatus int64, conditions bool) int64 {
	if oldStatus != 0 {
		return oldStatus
	}
	if conditions {
		return newStatus
	}

	return 0
}

func (tc *selectController) click(c echo.Context) error {
	rd := middlewares.MustGetRequestData(c)
	var (
		status     int64
		noRedisKey bool
	)
	adID, _ := strconv.ParseInt(c.Param("ad"), 10, 64)
	wIDStr := c.Param("wid")
	rand := c.Param("rand")
	mega := c.Param("mega")
	typ := c.Param("typ")
	tv := c.QueryParam("tv") != ""

	ads, err := mr.NewManager().GetAd(adID, true)
	status = changeStatus(0, suspNoAdFound, err != nil)

	result, err := aredis.HGetAllString(fmt.Sprintf("%s%s%s%s%d",
		transport.IMP,
		transport.DELIMITER,
		mega,
		transport.DELIMITER,
		adID,
	), true, config.Config.Clickyab.DailyImpExpire)
	if err != nil {
		status = changeStatus(status, suspSlowClick, true)
		noRedisKey = true
		result = make(map[string]string)
	}

	wID, err := strconv.ParseInt(wIDStr, 10, 0)
	assertNil(noRedisKey, err)

	var pub Publisher
	if typ != "app" {
		pub, err = mr.NewManager().FetchWebsite(wID)
		status = changeStatus(status, suspInvalidWebsite, err != nil || !pub.GetActive())
	} else {
		pub, err = mr.NewManager().GetAppByID(wID)
		status = changeStatus(status, suspInvalidApp, err != nil || !pub.GetActive())
	}
	clickID := <-utils.ID

	middlewares.SafeGO(c, false, false, func() {

		winnerBid, err := strconv.ParseInt(result["WIN"], 10, 0)
		assertNil(noRedisKey, err)
		status = changeStatus(status, suspWSMismatch, wIDStr != result["WS"])

		status = changeStatus(status, suspIPMismatch, rd.IP.String() != result["IP"])

		// App is special case, since the app is clicked via browser and the UA is changed
		if typ != "app" {
			status = changeStatus(status, suspUAMismatch, rd.UserAgent != result["UA"])
		}
		status = changeStatus(status, suspRndMismatch, rand != result["RND"])

		slotID, err := strconv.ParseInt(result["S"], 10, 0)
		assertNil(noRedisKey, err)

		in, err := strconv.ParseInt(result["T"], 10, 0)
		assertNil(noRedisKey, err)

		inTime := time.Unix(in, 0)

		slaID, err := strconv.ParseInt(result["SLAID"], 10, 0)
		assertNil(noRedisKey, err)

		impID, err := strconv.ParseInt(result["IMPR"], 10, 0)
		assertNil(noRedisKey, err)

		cpAdID, err := strconv.ParseInt(result["CPADID"], 10, 0)
		assertNil(noRedisKey, err)

		outTime := time.Now()
		if noRedisKey || outTime.Unix()-inTime.Unix() < config.Config.Clickyab.FastClick {
			status = suspFastClick
		}

		clickRedis := fmt.Sprintf("%s%s%s%s%s", transport.CLICK, transport.DELIMITER, mega, transport.DELIMITER, transport.ADVERTISE)
		count, err := aredis.IncHash(clickRedis, fmt.Sprintf("CLICK_%d",adID), 1, config.Config.Clickyab.DailyImpExpire)
		assert.Nil(err)

		if count != 1 {
			status = suspDuplicateClick
		}

		click := tc.fillClick(rd, ads, winnerBid, pub, slotID, inTime, outTime, slaID, impID, cpAdID, status, clickID, tv)

		rabbit.MustPublish(click)
	})

	if status == suspNoAdFound {
		return c.String(http.StatusNotFound, "Not found")
	}
	// TODO : better handling
	_, _ = aredis.IncHash(fmt.Sprintf("%s%s%s", transport.CONV, transport.DELIMITER, clickID), "OK", 1, config.Config.Clickyab.DailyClickExpire)
	url := ""
	cpName := ""
	if ads != nil {
		url = ads.AdURL.String
		cpName = ads.CampaignName.String
	}
	domain := ""
	if pub != nil {
		domain = pub.GetName()
	}
	body := tc.replaceParameters(url, domain, cpName, rand, result["IMPR"], rd.IP.String(), result["GID"], result["AID"], result["DID"])
	return c.HTML(http.StatusOK, body)
}

func (selectController) fillClick(
	rd *middlewares.RequestData,
	ads *mr.Ad,
	winnerBid int64,
	pub Publisher,
	slotID int64,
	inTime, outTime time.Time,
	slaID int64,
	impID int64,
	campaignAdID int64,
	status int64,
	rand string,
	tv bool) *transport.Click {

	var (
		web *transport.WebSiteImp
		app *transport.AppImp
	)
	var id int64
	// in some case (forged request) its possible to pub to be empty. just ignore it
	if pub != nil {
		id = pub.GetID()
	}
	if pub.GetType() == "web" {
		web = &transport.WebSiteImp{
			WebsiteID: id,
			SlotID:    slotID,
			Referrer:  rd.Referrer,
			ParentURL: rd.Parent,
		}
	} else {
		app = &transport.AppImp{
			AppID:  id,
			SlotID: slotID,
		}
	}
	return &transport.Click{
		CopID:        rd.CopID,
		IP:           rd.IP,
		AdID:         ads.AdID,
		SlotID:       slotID,
		CampaignID:   ads.CampaignID.Int64,
		UserAgent:    rd.UserAgent,
		WinnerBid:    winnerBid,
		InTime:       inTime,
		OutTime:      outTime,
		SlaID:        slaID,
		ImpID:        impID,
		OS:           rd.PlatformID,
		Status:       status,
		CampaignAdID: campaignAdID,
		Rand:         rand,
		TrueView:     tv,
		Web:          web,
		App:          app,
	}
}

func (selectController) replaceParameters(url, domain, campaign, clickID, impID, ip, googlead_id, android_id, android_device string) string {
	r := strings.NewReplacer(
		"[app]",
		domain,
		"[domain]",
		domain,
		"[campaign]",
		campaign,
		"[click_id]",
		clickID,
		"{app}",
		domain,
		"{domain}",
		domain,
		"{campaign}",
		campaign,
		"{imp_id}",
		impID,
		"{click_id}",
		clickID,

		"{ip}",
		ip,
		"[ip]",
		ip,
		"{googlead_id}",
		googlead_id,
		"[googlead_id]",
		googlead_id,
		"{android_id}",
		android_id,
		"[android_id]",
		android_id,
		"{android_device}",
		android_device,
		"[android_device]",
		android_device,
	)

	url = r.Replace(url)
	return `<html><head><title>$imp_url</title><meta name="robots" content="nofollow"/></head>
			<body><script>window.setTimeout( function() { window.location.href = '` + url + `' }, 500 );</script></body>
			</html>`
}
