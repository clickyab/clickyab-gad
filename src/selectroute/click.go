package selectroute

import (
	"assert"
	"mr"
	"strconv"
	"transport"

	"fmt"

	"middlewares"

	"time"

	"config"
	"redis"

	"net/http"

	"rabbit"

	"github.com/labstack/echo"
)

const (
	SUSP_DUPLICATE_CLICK = 1
	SUSP_FAST_CLICK      = 9
	SUSP_WS_MISMATCH     = 1024
	SUSP_RND_MISMATCH    = 1025
	SUSP_IP_MISMATCH     = 1026
	SUSP_UA_MISMATCH     = 1027
	SUSP_SLOW_CLICK      = 1028
)

func assertNil(status bool, err error) {
	if !status {
		return
	}

	assert.Nil(err)
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
	result, err := aredis.HGetAllString(fmt.Sprintf("%s%s%s%s%d",
		transport.IMP,
		transport.DELIMITER,
		mega,
		transport.DELIMITER,
		adID,
	), true, config.Config.Clickyab.DailyImpExpire)
	if err != nil {
		status = SUSP_SLOW_CLICK
		noRedisKey = true
	}
	ads, err := mr.NewManager().GetAd(adID)
	if err != nil {
		return c.String(http.StatusNotFound, "Not found")
	}

	go func() {

		winnerBid, err := strconv.ParseInt(result["WIN"], 10, 0)
		assertNil(noRedisKey, err)

		if noRedisKey || wIDStr != result["WS"] {
			status = SUSP_WS_MISMATCH
		}

		if noRedisKey || rand != result["RND"] {
			status = SUSP_RND_MISMATCH
		}

		if noRedisKey || rd.IP.String() != result["IP"] {
			status = SUSP_IP_MISMATCH
		}

		if noRedisKey || rd.UserAgent != result["UA"] {
			status = SUSP_UA_MISMATCH
		}

		wID, err := strconv.ParseInt(result["WS"], 10, 0)
		assertNil(noRedisKey, err)

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
			status = SUSP_FAST_CLICK
		}

		clickRedis := fmt.Sprintf("%s%s%s%s%s", transport.CLICK, transport.DELIMITER, mega, transport.DELIMITER, transport.ADVERTISE)
		count, err := aredis.IncHash(clickRedis, "CLICK", 1, config.Config.Clickyab.DailyImpExpire)
		assert.Nil(err)

		if count != 1 {
			status = SUSP_DUPLICATE_CLICK
		}

		click := tc.fillClick(c, ads, winnerBid, wID, slotID, inTime, outTime, slaID, impID, cpAdID, status)
		err = rabbit.Publish("cy.click", click)
		assert.Nil(err)
	}()

	return c.Redirect(http.StatusFound, ads.AdURL.String)

}

func (selectController) fillClick(ctx echo.Context, ads mr.Ad, winnerBid int64, websiteID int64, slotID int64, inTime, outTime time.Time, slaID int64, impID int64, campaignAdID int64, status int64) *transport.Click {
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
		Web: &transport.WebSiteImp{
			WebsiteID: websiteID,
			SlotID:    slotID,
			Referrer:  rd.Referrer,
			ParentURL: rd.Parent,
		},
	}
}
