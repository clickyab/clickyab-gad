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
	User       string    `json:"u"`
	ImpID      string    `json:"i"`
	IP         net.IP    `json:"ip"`
	AdID       int64     `json:"adid"`
	CampaignID int64     `json:"cid"`
	UserAgent  string    `json:"ua"`
	Time       time.Time `json:"it"`
	WinnerBID  int64     `json:"wb"`
	// TODO : better status
	Status     int64 `json:"s"`
	Cookie     bool  `json:"c"`
	Suspicious bool  `json:"sp"`

	Web *WebSiteImp `json:"web"`
}
