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

// Impression is the record for the single impression
type Impression struct {
	User          string    `json:"u"`
	MegaUser      string    `json:"mu"`
	ImpID         string    `json:"i"`
	IP            net.IP    `json:"ip"`
	MegaIP        net.IP    `json:"mip"`
	AdID          int64     `json:"adid"`
	CampaignID    int64     `json:"cid"`
	UserAgent     string    `json:"ua"`
	MegaUserAgent string    `json:"mua"`
	InTime        time.Time `json:"it"`
	OutTime       time.Time `json:"ot"`
	WinnerBID     int64     `json:"wb"`
	// TODO : better status
	Status int64 `json:"s"`
	Cookie bool  `json:"c"`

	Web *WebSiteImp `json:"web"`
}
