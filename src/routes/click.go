package routes

import (
	"assert"
	"config"
	"fmt"
	"middlewares"
	"mr"
	"net/http"
	"rabbit"
	"redis"
	"strconv"
	"strings"
	"time"
	"transport"
	"utils"

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
	tv := c.QueryParam("tv") != ""

	ads, err := mr.NewManager().GetAd(adID)
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

	wID, err := strconv.ParseInt(result["WS"], 10, 0)
	assertNil(noRedisKey, err)

	webSite, err := mr.NewManager().FetchWebsite(wID)
	status = changeStatus(status, suspInvalidWebsite, err != nil || (webSite.WStatus != 0 && webSite.WStatus != 1))

	clickID := <-utils.ID

	go func() {

		winnerBid, err := strconv.ParseInt(result["WIN"], 10, 0)
		assertNil(noRedisKey, err)

		status = changeStatus(status, suspWSMismatch, wIDStr != result["WS"])

		status = changeStatus(status, suspIPMismatch, rd.IP.String() != result["IP"])

		status = changeStatus(status, suspUAMismatch, rd.UserAgent != result["UA"])

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
		count, err := aredis.IncHash(clickRedis, "CLICK", 1, config.Config.Clickyab.DailyImpExpire)
		assert.Nil(err)

		if count != 1 {
			status = suspDuplicateClick
		}

		click := tc.fillClick(c, ads, winnerBid, wID, slotID, inTime, outTime, slaID, impID, cpAdID, status, clickID, tv)

		rabbit.MustPublish("cy.click", click)
	}()

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
	if webSite != nil {
		domain = webSite.WDomain.String
	}
	body := tc.replaceParameters(url, domain, cpName, rand, result["IMPR"])
	return c.HTML(http.StatusOK, body)
}

func (selectController) fillClick(
	ctx echo.Context,
	ads *mr.Ad,
	winnerBid int64,
	websiteID int64,
	slotID int64,
	inTime, outTime time.Time,
	slaID int64,
	impID int64,
	campaignAdID int64,
	status int64,
	rand string,
	tv bool) *transport.Click {
	rd := middlewares.MustGetRequestData(ctx)
	adID, err := strconv.ParseInt(ctx.Param("ad"), 10, 0)
	assert.Nil(err)

	return &transport.Click{
		CopID:        rd.CopID,
		IP:           rd.IP,
		AdID:         adID,
		SlotID:       slotID,
		CampaignID:   ads.CampaignAdID.Int64,
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
		Web: &transport.WebSiteImp{
			WebsiteID: websiteID,
			SlotID:    slotID,
			Referrer:  rd.Referrer,
			ParentURL: rd.Parent,
		},
	}
}

func (selectController) replaceParameters(url, domain, campaign, clickID, impID string) string {
	r := strings.NewReplacer(
		"[domain]",
		domain,
		"[campaign]",
		campaign,
		"[click_id]",
		clickID,
		"{domain}",
		domain,
		"{campaign}",
		campaign,
		"{imp_id}",
		impID,
		"{click_id}",
		clickID,
	)

	url = r.Replace(url)
	return `<html><head><title>$imp_url</title><meta name="robots" content="nofollow"/></head>
			<body><script>window.setTimeout( function() { window.location.href = '` + url + `' }, 500 );</script></body>
			</html>`
}
