package routes

import (
	"fmt"
	"strconv"
	"time"

	"clickyab.com/gad/middlewares"
	"clickyab.com/gad/models"
	"clickyab.com/gad/rabbit"
	"clickyab.com/gad/redis"
	"clickyab.com/gad/transport"
	"github.com/clickyab/services/assert"

	"github.com/sirupsen/logrus"
)

const message = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII="

// Publisher is a publisher object, app, website
type Publisher interface {
	// GetID return the id of object
	GetID() int64
	// GetName return the name of object
	GetName() string

	// FloorCPM is the floor value for this site
	FloorCPM() int64

	// GetActive return if the app is acive
	GetActive() bool

	// GetType return the type of this publisher
	GetType() string
}

func (selectController) fillImp(rd *middlewares.RequestData, sus bool, ads *models.Ad, winnerBid int64, pub Publisher, slotID int64) transport.Impression {
	imp := transport.Impression{
		Suspicious:   sus,
		IP:           rd.IP,
		AdID:         ads.AdID,
		CopID:        rd.CopID,
		CampaignAdID: ads.CampaignAdID.Int64,
		SlotID:       slotID,
		URL:          ads.AdURL.String,
		CampaignID:   ads.CampaignID.Int64,
		UserAgent:    rd.UserAgent,
		Time:         time.Now(),
		WinnerBID:    winnerBid,
		Status:       0,
	}
	if pub.GetType() == "web" {
		imp.Web = &transport.WebSiteImp{
			Referrer:  rd.Referrer,
			ParentURL: rd.Parent,
			SlotID:    slotID,
			WebsiteID: pub.GetID(),
		}
	} else if pub.GetType() == "app" {
		imp.App = &transport.AppImp{
			AppID:  pub.GetID(),
			SlotID: slotID,
		}
	}

	return imp
}

func (selectController) fillNativeImp(rd *middlewares.RequestData, sus bool, ads *models.AdData, winnerBid int64, pub Publisher, slotID int64) transport.Impression {
	imp := transport.Impression{
		Suspicious:   sus,
		IP:           rd.IP,
		AdID:         ads.AdID,
		CopID:        rd.CopID,
		CampaignAdID: ads.CampaignAdID,
		SlotID:       slotID,
		URL:          ads.AdURL.String,
		CampaignID:   ads.CampaignID,
		UserAgent:    rd.UserAgent,
		Time:         time.Now(),
		WinnerBID:    winnerBid,
		Status:       0,
	}

	imp.Web = &transport.WebSiteImp{
		Referrer:  rd.Referrer,
		ParentURL: rd.Parent,
		SlotID:    slotID,
		WebsiteID: pub.GetID(),
	}

	return imp
}

func (selectController) callWebWorker(pub Publisher, slotID int64, adID int64, mega string, rand string, imp transport.Impression, rd *middlewares.RequestData) {
	m := models.NewManager()
	var err error
	imp.SLAID, err = m.InsertSlotAd(slotID, adID)
	if err != nil {
		// not important error
		logrus.Debug(err)
	}
	assert.Nil(m.InsertImpression(&imp))
	//validate
	res, err := aredis.HGetAllString(fmt.Sprintf("%s%s%s", transport.MegaKey, transport.Delimiter, mega), true, 2*time.Hour)
	assert.Nil(err)

	//check ip
	wID, _ := strconv.ParseInt(res["WS"], 10, 64)
	if res["IP"] != rd.IP.String() || res["UA"] != rd.UserAgent || wID != pub.GetID() {
		imp.Suspicious = true
	}

	// TODO : Use constant not strings
	//set mega ip in redis
	tmp := map[string]string{
		"IP":     rd.IP.String(),
		"UA":     rd.UserAgent,
		"AID":    rd.AndroidID,     //android id
		"DID":    rd.AndroidDevice, //android device id
		"GID":    rd.GoogleID,      //google id
		"WS":     strconv.FormatInt(pub.GetID(), 10),
		"T":      strconv.FormatInt(time.Now().Unix(), 10),
		"S":      strconv.FormatInt(slotID, 10),
		"IMPR":   strconv.FormatInt(imp.ID, 10),
		"RND":    rand,
		"WIN":    strconv.FormatInt(imp.WinnerBID, 10),
		"CPADID": strconv.FormatInt(imp.CampaignAdID, 10),
		"SLAID":  strconv.FormatInt(imp.SLAID, 10),
	}
	err = aredis.HMSet(
		fmt.Sprintf(
			"%s%s%s%s%d",
			transport.ImpKey,
			transport.Delimiter,
			mega,
			transport.Delimiter,
			adID),
		megaImpExpire.Duration(),
		tmp)
	if err != nil {
		logrus.WithField("cy.imp", imp).Error("error in hmset", err)
	}
	err = rabbit.Publish(imp)
	if err != nil {
		logrus.WithField("cy.imp", imp).Error("error in  publishing job", err)
	}
}
