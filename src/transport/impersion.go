package transport

import (
	"net"
	"time"
)

// WebSiteImp is detail of website related to imp
type WebSiteImp struct {
	WebsiteID int64  `json:"wid"`
	SlotID    int64  `json:"sid"`
	Referrer  string `json:"r"`
	ParentURL string `json:"par"`
}

// AppImp is detail of app
type AppImp struct {
	AppID  int64 `json:"aid"`
	SlotID int64 `json:"sid"`
}

// Impression is the record for the single impression
type Impression struct {
	ID           int64     `json:"id"`
	IP           net.IP    `json:"ip"`
	AdID         int64     `json:"adid"`
	CopID        int64     `json:"copid"`
	CampaignAdID int64     `json:"caid"`
	SLAID        int64     `json:"sla"`
	SlotID       int64     `json:"slot"`
	URL          string    `json:"url"`
	CampaignID   int64     `json:"cid"`
	UserAgent    string    `json:"ua"`
	Time         time.Time `json:"it"`
	WinnerBID    int64     `json:"wb"`
	// TODO : better status
	Status     int64       `json:"s"`
	Suspicious bool        `json:"sp"`
	Web        *WebSiteImp `json:"web"`
	App        *AppImp     `json:"app"`
}
