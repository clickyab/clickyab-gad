package transport

import (
	"net"
	"time"
)

// WebSiteImp is detail of website related to imp
type WebSiteImp struct {
	WebsiteID    int64  `json:"wid"`
	SlotID       int64  `json:"sid"`
	Referrer     string `json:"r"`
	MegaReferrer string `json:"mr"`
}

// AppImp is detail of app
type AppImp struct {
	AppID int64 `json:"aid"`
}

// Impression is the record for the single impression
type Impression struct {
	User            string    `json:"u"`
	ID              string    `json:"id"`
	IP              net.IP    `json:"ip"`
	AdID            int64     `json:"adid"`
	CopID           int64     `json:"adid"`
	CampaignAdID    int64     `json:"caid"`
	ReferralAddress string    `json:"ref"`
	ParentURL       string    `json:"par"`
	URL             string    `json:"url"`
	CampaignID      int64     `json:"cid"`
	UserAgent       string    `json:"ua"`
	MegaUserAgent   string    `json:"mua"`
	Time            time.Time `json:"it"`
	WinnerBID       int64     `json:"wb"`
	// TODO : better status
	Status     int64       `json:"s"`
	Cookie     int         `json:"c"`
	Suspicious bool        `json:"sp"`
	Web        *WebSiteImp `json:"web"`
	App        *AppImp     `json:"app"`
}
